package model

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
)

// Image of the labyrinth
var LabImg *image.RGBA = image.NewRGBA(image.Rect(0, 0, LabWidth, LabHeight))

// Gopher images for each direction, each has zero Min point
var GopherImgs []*image.RGBA = make([]*image.RGBA, DirMax+1)

// Image of the wall block
var WallImg = image.NewUniform(WallCol)

// Image of the empty block
var EmptyImg = image.NewUniform(Black)

// Image of the empty block
var TargetImg = image.NewUniform(color.RGBA{0xff, 0x10, 0x10, 0xff})

func init() {
	// Load gopher images
	for i := Dir(0); i <= DirMax; i++ {
		data, err := ioutil.ReadFile(fmt.Sprintf("w:/gopher-%s.png", i.String()))
		if err != nil {
			panic(err)
		}
		src, err := png.Decode(bytes.NewBuffer(data))
		if err != nil {
			panic(err)
		}

		// Convert to image.RGBA, also make sure result image has zero Min point
		b := src.Bounds()
		if b.Dx() != BlockSize || b.Dy() != BlockSize {
			panic("Invalid image size!")
		}
		GopherImgs[i] = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(GopherImgs[i], src.Bounds(), src, b.Min, draw.Src)
	}
}
