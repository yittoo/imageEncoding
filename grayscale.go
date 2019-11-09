package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

// Grayscale takes pointer f to image file (PNG encoded) and output grayscaled PNG with fname
func Grayscale(f *os.File, fname *string) {
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
