package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

// Negative takes pointer f to image file (PNG encoded) and output inverted PNG with fname
func PaintTo(f *os.File, fname *string, to *string) {
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
			if *to == "red" {
				g = 0
				b = 0
			} else if *to == "green" {
				r = 0
				b = 0
			} else if *to == "blue" {
				r = 0
				g = 0
			} else {
				fmt.Println("Please provide one of the following colors to paint with -to flag\n\tred, green, blue")
				return
			}

			ni.Set(x, y, color.RGBA{
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
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
