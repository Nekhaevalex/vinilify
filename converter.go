package main

import (
	"fmt"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

// Gets the image by the `filePath`, rotates it by `d` degrees until a complete loop and returns the directory where the rotated images are stored.
func RotateImage(filePath string, d int) (string, error) {

	img, err := gg.LoadImage("./t.png")
	if err != nil {
		return "", err
	}

	newWidth := 500
	newHeight := 500

	res_img := resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos2)

	// dc.RotateAbout(1, 1000/2, 1000/2)
	for i := range 45 {
		dc := gg.NewContext(1000, 1000)
		dc.RotateAbout(gg.Radians(float64(d*(i+1))), 500, 500)
		dc.DrawImage(res_img, 250, 250)
		dc.SavePNG(fmt.Sprintf("./tests/test_%d.png", i))
	}

	return "./tests/", nil
}

func StackImage() (string, error) {

	return "", nil
}
