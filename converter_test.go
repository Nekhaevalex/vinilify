package main

import (
	"fmt"
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
