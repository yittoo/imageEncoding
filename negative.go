package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

// Negative takes pointer f to image file (PNG encoded) and output inverted PNG with fname
func Negative(f *os.File, fname *string) {
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
			ni.Set(x, y, color.RGBA{
				R: 255 - uint8(r),
				G: 255 - uint8(g),
				B: 255 - uint8(b),
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
