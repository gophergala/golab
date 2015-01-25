package model

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
)

// Image of the labyrinth
var LabImg *image.RGBA = image.NewRGBA(image.Rect(0, 0, LabWidth, LabHeight))

// Gopher images for each direction, each has zero Min point
var GopherImgs []*image.RGBA = make([]*image.RGBA, DirLength)

// Dead Gopher image.
var DeadImg *image.RGBA

// Bulldog images for each direction, each has zero Min point
var BulldogImgs []*image.RGBA = make([]*image.RGBA, DirLength)

// Image of the wall block
var WallImg *image.RGBA

// Image of the empty block
var EmptyImg = image.NewUniform(Black)

// Image of the empty block
//var TargetImg = image.NewUniform(color.RGBA{0x00, 0xff, 0x00, 0xff})
var TargetImg *image.RGBA

// Image of a door, this is the exit sign
var ExitImg *image.RGBA

// Image of a congratulation
var WonImg *image.RGBA

func init() {
	for i := Dir(0); i < DirLength; i++ {
		// Load gopher images
		GopherImgs[i] = loadImg(fmt.Sprintf("w:/gopher-%s.png", i), true)
		// Load Bulldog images
		BulldogImgs[i] = loadImg(fmt.Sprintf("w:/bulldog-%s.png", i), true)
	}

	WallImg = loadImg("w:/wall5.png", true)
	DeadImg = loadImg("w:/gopher-dead.png", true)
	ExitImg = loadImg("w:/door.png", true)

	TargetImg = loadImg("w:/marker.png", false)
	WonImg = loadImg("w:/won.png", false)
}

// loadImg loads a PNG image from the specified file, and converts it to image.RGBA and makes sure image has zero Min point.
func loadImg(name string, blockSize bool) *image.RGBA {
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
	if blockSize && (b.Dx() != BlockSize || b.Dy() != BlockSize) {
		panic("Invalid image size!")
	}

	img := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(img, src.Bounds(), src, b.Min, draw.Src)

	return img
}
