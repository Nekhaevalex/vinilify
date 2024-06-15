package main

import (
	"fmt"
	"image"
	"log"
	"testing"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

func TestRotateImage(t *testing.T) {
	img, err := gg.LoadImage("./t.png")
	if err != nil {
		t.Error(err)
	}

	newWidth := 500
	newHeight := 500

	res_img := resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos2)

	// dc.RotateAbout(1, 1000/2, 1000/2)
	for i := range 45 {
		dc := gg.NewContext(1000, 1000)
		dc.RotateAbout(gg.Radians(float64(8*(i+1))), 500, 500)
		dc.DrawImage(res_img, 250, 250)
		dc.SavePNG(fmt.Sprintf("./tests/test_%d.png", i))
	}
}

func TestStackImage(t *testing.T) {
	img1, err := gg.LoadImage("./t.png")
	if err != nil {
		t.Error(err)
	}

	back1, err := gg.LoadImage("./Assets/Images/Disk.png")

	b1s := resize.Resize(1000, 1000, back1, resize.Lanczos3)
	i1s := resize.Resize(500, 500, img1, resize.Lanczos3)

	dc := gg.NewContext(1000, 1000)
	dc.DrawImage(i1s, 250, 250)
	dc.DrawImage(b1s, 0, 0)
	dc.SavePNG("./test.png")
}

func LoadAndResizeImage(path string, width uint, height uint) (image.Image, error) {
	img, err := gg.LoadImage(path)
	if err != nil {
		log.Panic(err)
	}

	img_res := resize.Resize(width, height, img, resize.Lanczos3)

	return img_res, nil
}

func TestStackImages(t *testing.T) {
	img1, _ := LoadAndResizeImage("./t.png", 500, 500)
	img2, _ := LoadAndResizeImage("./test.png", 250, 250)
	img3, _ := LoadAndResizeImage("./Assets/Images/Disk.png", 1000, 1000)

	imgs := []image.Image{
		img3,
		img1,
		img2,
	}

	coords := []Coord{
		Coord{x: 0, y: 0},
		Coord{x: 250, y: 250},
		Coord{x: 375, y: 375},
	}

	dc := gg.NewContext(1000, 1000)
	dc, err := StackImages(dc, imgs, coords)
	if err != nil {
		t.Error(err)
	}
	dc.SavePNG("./test_stack.png")
}
