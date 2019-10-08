package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
)

var iname *string
var iiname *string
var oname *string
var appmode *string
var isCreateNew *bool
var isEncodeMode *bool
var isDecodeMode *bool
var messageToEncode *string

func init() {
	iname = flag.String("i", "input.png", "input file name the program will use")
	iiname = flag.String("ii", "", "input file name the program will use for decoding mode")
	oname = flag.String("o", "output.png", "output file name the result will be written to")

	// for new no required additional parameters
	// for encode -m is required
	// for decode -iiname is required
	appmodehelp := `mode to use app on accepted values: 
	new - outputs an image with name given in -o flag, size: 10x10, colors: chosen randomly
		example: 'executible -mode=new -o=newimg.png'
	encode - outputs a file (name given in -o) encoded with message in -m over image given by -i
		example: 'executible -mode=encode -i=img.png -o=encodedimg.png -m="some very secret message to be encoded over pixels"'
	decode - outputs a message using two images (one encoded with message of other) given by -i and -ii flags
		example: 'executible -mode=decode -i=img.png -ii=encodedimg.png'
	grayscale - outputs an image (name given in -o) using image given in -i
		example: 'executible -mode=grayscale -o=newimg.png'`
	appmode = flag.String("mode", "", appmodehelp)
	messageToEncode = flag.String("m", "", "message that will be encoded onto input image")
	flag.Parse()
}

func main() {
	switch *appmode {
	case "new":
		fmt.Println("Creating new image")
		createNewPng(oname)
	case "grayscale":
		fmt.Println("Grayscaling image")
		f, err := openGivenFile(iname)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()
		grayscale(f, oname)
	case "encode":
		fmt.Println("Encoding image")
		f, err := openGivenFile(iname)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()
		encode(f, oname, messageToEncode)
	case "decode":
		fmt.Printf("Decoding images using files: %v and %v\n", *iname, *iiname)
		f1, err := openGivenFile(iname)
		if err != nil {
			fmt.Println(err)
		}
		defer f1.Close()
		f2, err := openGivenFile(iiname)
		if err != nil {
			fmt.Println(err)
		}
		defer f2.Close()
		// decoded message
		dm := decode(f1, f2)
		fmt.Println("Decoded message: ", dm)
	default:
		fmt.Println("You need to pass app mode with -mode, possible values are: new / encode / decode / grayscale, see --help for more details")
	}
}

func openGivenFile(fname *string) (*os.File, error) {
	f, err := os.OpenFile(*fname, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func createNewPng(fname *string) {
	f, err := os.Create(*fname)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	r := image.Rect(0, 0, 10, 10)
	ni := image.NewRGBA(r)
	or := ni.Bounds()
	size := or.Size()
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			ni.Set(x, y, color.RGBA{
				R: uint8(rand.Int()),
				G: uint8(rand.Int()),
				B: uint8(rand.Int()),
				A: 255,
			})
		}
	}
	png.Encode(f, ni)
}

func grayscale(f *os.File, fname *string) {
	ri, err := png.Decode(f)
	if err != nil {
		fmt.Println(err)
	}
	or := ri.Bounds()
	size := or.Size()
	ni := image.NewRGBA(or)

	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			r, g, b, _ := ri.At(x, y).RGBA()
			gs := (r + g + b) / 3
			ni.Set(x, y, color.RGBA{
				R: uint8(gs),
				G: uint8(gs),
				B: uint8(gs),
				A: 255,
			})
		}
	}

	nf, err := os.Create(*fname)
	if err != nil {
		fmt.Println(err)
	}
	defer nf.Close()

	err = png.Encode(nf, ni)
	if err != nil {
		fmt.Println(err)
	}
}

func encode(f *os.File, fname *string, m *string) {
	if len(*m) < 1 {
		fmt.Println("Please input a message via -m flag to encode")
		return
	}
	ri, err := png.Decode(f)
	if err != nil {
		fmt.Println(err)
	}
	or := ri.Bounds()
	maxMessageLength := or.Dx() * or.Dy()
	if maxMessageLength < len(*m) {
		fmt.Println("Message is too long to fit onto this image, either provide bigger image or shorter message")
		fmt.Println("MaxLetterCount = (width * height) [of image]")
	}
	size := or.Size()
	ni := image.NewRGBA(or)

	bs := []byte(*m)
	i := 0

	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			r, g, b, _ := ri.At(x, y).RGBA()

			// encoded byte
			var eb byte

			if i < len(*m) {
				// if equal to space turn to 0x1f otherwise remove binary abcxyzmn's abc part to reduce pixel mutation amount
				if bs[i] == 0x20 {
					eb = bs[i] - 1
				} else {
					eb = bs[i] & 0x1f
				}
				i++
			}

			// red green blue (to write)
			rtw := encodeColor(uint8(r), eb)
			gtw := encodeColor(uint8(g), eb)
			btw := encodeColor(uint8(b), eb)

			ni.Set(x, y, color.RGBA{
				R: rtw,
				G: gtw,
				B: btw,
				A: 255,
			})
		}
	}

	nf, err := os.Create(*fname)
	if err != nil {
		fmt.Println(err)
	}
	defer nf.Close()

	err = png.Encode(nf, ni)
	if err != nil {
		fmt.Println(err)
	}
}

func decode(f1, f2 *os.File) string {
	ri1, err := png.Decode(f1)
	if err != nil {
		fmt.Println(err)
	}
	ri2, err := png.Decode(f2)
	if err != nil {
		fmt.Println(err)
	}
	or1 := ri1.Bounds()
	or2 := ri2.Bounds()

	if or1 != or2 {
		fmt.Println("Image bounds do not match, please provide normal and encoded version of same image")
		return ""
	}
	maxMessageLength := or1.Dx() * or1.Dy()
	bs := make([]byte, 0, maxMessageLength)
	size := or1.Size()

	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			r1, g1, b1, _ := ri1.At(x, y).RGBA()
			r2, g2, b2, _ := ri2.At(x, y).RGBA()

			// decoded byte, temp byte
			var db, tb byte

			tb = decodeColor(uint8(r1), uint8(r2))
			db = tb
			tb = decodeColor(uint8(g1), uint8(g2))
			if tb != db {
				fmt.Println("mismatch between red and green byte encodes at index ", len(bs))
			}
			tb = decodeColor(uint8(b1), uint8(b2))
			if tb != db {
				fmt.Println("mismatch between red and blue byte encodes")
			}
			if db != 0x60 {
				bs = append(bs, db)
			}
		}
	}
	return string(bs)
}

// check if added encoded byte is going to exceed uint8 max boundaries and act accordingly
func encodeColor(c, b uint8) uint8 {
	if b == 0 {
		return c
	} else if uint16(c)+uint16(b) > 255 {
		return c - b
	}
	return c + b
}

// returns absolute value of compared color bytes and adds 01100000 to match it to original ascii value
func decodeColor(c1, c2 uint8) byte {
	var v byte
	if c1 > c2 {
		v = c1 - c2
	} else {
		v = (c2 - c1)
	}
	if v == 0x1f {
		return 0x20
	}
	v |= 0x60
	return v
}
