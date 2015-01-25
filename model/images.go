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
var GopherImgs []*image.RGBA = make([]*image.RGBA, DirLength)

// Bulldog images for each direction, each has zero Min point
var BulldogImgs []*image.RGBA = make([]*image.RGBA, DirLength)

// Image of the wall block
var WallImg *image.RGBA

// Image of the empty block
var EmptyImg = image.NewUniform(Black)

// Image of the empty block
var TargetImg = image.NewUniform(color.RGBA{0xff, 0x10, 0x10, 0xff})

func init() {
	for i := Dir(0); i < DirLength; i++ {
		// Load gopher images
		GopherImgs[i] = loadImg(fmt.Sprintf("w:/gopher-%s.png", i))
		// Load Bulldog images
		BulldogImgs[i] = loadImg(fmt.Sprintf("w:/bulldog-%s.png", i))
	}

	WallImg = loadImg("w:/wall5.png")
}

// loadImg loads a PNG image from the specified file, and converts it to image.RGBA and makes sure image has zero Min point.
func loadImg(name string) *image.RGBA {
	data, err := ioutil.ReadFile(name)
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

	img := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(img, src.Bounds(), src, b.Min, draw.Src)

	return img
}
