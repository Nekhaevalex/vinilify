package main

import (
	"errors"
	"fmt"
	"image"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

// Gets the image by the `filePath`, rotates it by `d` degrees until a complete loop and returns the directory where the rotated images are stored.
func RotateImage(filePath string, d int) (*gg.Context, error) {

	img, err := gg.LoadImage("./t.png")
	if err != nil {
		return gg.NewContext(0, 0), err
	}

	newWidth := 500
	newHeight := 500

	res_img := resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos2)

	// dc.RotateAbout(1, 1000/2, 1000/2)

	dc := gg.NewContext(1000, 1000)
	dc.RotateAbout(gg.Radians(float64(d)), 500, 500)
	dc.DrawImage(res_img, 250, 250)
	dc.SavePNG(fmt.Sprintf("./tests/test_%d.png", d))

	return dc, nil
}

// Takes in gg.Context `dc` and stacks `imgs` one on top of another at `coords` (top left). Returns the same context.
func StackImages(dc *gg.Context, imgs []image.Image, coords []Coord) (*gg.Context, error) {
	if len(imgs) == 0 {
		return gg.NewContext(0, 0), ErrorEmptyArray
	}
	if len(imgs) != len(coords) {
		return gg.NewContext(0, 0), ErrorImageCoordMismatch
	}
	for i := range len(imgs) {
		dc.DrawImage(imgs[i], coords[i].x, coords[i].y)
	}
	return dc, nil
}

var (
	ErrorEmptyArray         = errors.New("empty array of images")
	ErrorImageCoordMismatch = errors.New("len of coords must be equal to len of images")
)

type Coord struct {
	x int
	y int
}
