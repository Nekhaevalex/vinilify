package main

import (
	"testing"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

func TestLoadImage(t *testing.T) {
	img, err := gg.LoadImage("./t.png")
	if err != nil {
		t.Error(err)
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	newWidth := 500
	newHeight := 500

	res_img := resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos2)

	width = res_img.Bounds().Dx()
	height = res_img.Bounds().Dy()

	t.Log(width, height)

	dc := gg.NewContext(1000, 1000)
	// dc.RotateAbout(1, 1000/2, 1000/2)
	dc.RotateAbout(1, 500, 500)
	dc.DrawImage(res_img, 250, 250)
	dc.SavePNG("./test.png")
}

// func TestRotateImage(t *testing.T) {
// 	img, err := gg.LoadImage("./Assets/Images/Disk.png")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	img.Rotate
// }
