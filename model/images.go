package model

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"image/color"
)

// Image of the labyrinth
var LabImg *image.RGBA = image.NewRGBA(image.Rect(0, 0, LabWidth, LabHeight))

// Gopher image which has zero Min point
var GopherImg *image.RGBA

// Image of the wall block
var WallImg = image.NewUniform(WallCol)

// Image of the empty block
var EmptyImg = image.NewUniform(Black)

// Image of the empty block
var TargetImg = image.NewUniform(color.RGBA{0xff, 0x10, 0x10, 0xff })

func init() {
	// Load gopher image
	data, err := ioutil.ReadFile("w:/gopher-small.png")
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
	GopherImg = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(GopherImg, src.Bounds(), src, b.Min, draw.Src)
}
