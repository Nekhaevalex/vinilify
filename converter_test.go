package main

import (
	"testing"

	"github.com/Nekhaevalex/vinilify/utils/converters"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

func TestCropAndRotateImage(t *testing.T) {
	img, err := gg.LoadImage("./t.png")
	if err != nil {
		t.Error(err)
	}

	img = resize.Resize(1000, 1000, converters.CropAndRotateImage(img, float64(50)), resize.Lanczos3)
	dc := gg.NewContext(1000, 1000)
	dc.DrawImage(img, 0, 0)
	dc.SavePNG("./test.png")

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

// func TestStackImages(t *testing.T) {
// 	img1, _ := converters.LoadAndResizeImage("./t.png", 500, 500)
// 	img2, _ := converters.LoadAndResizeImage("./test.png", 250, 250)
// 	img3, _ := converters.LoadAndResizeImage("./Assets/Images/Disk.png", 1000, 1000)

// 	imgs := []image.Image{
// 		img3,
// 		img1,
// 		img2,
// 	}

// 	coords := []converters.Coord{
// 		converters.Coord{x: 0, y: 0},
// 		converters.Coord{x: 250, y: 250},
// 		converters.Coord{x: 375, y: 375},
// 	}

// 	dc, err := converters.StackImages(gg.NewContext(1000, 1000), imgs, coords)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	dc.SavePNG("./test_stack.png")
// }

// func TestAssembleImage(t *testing.T) {
// 	b, _ := converters.LoadAndResizeImage("./Assets/Images/Disk.png", 1000, 1000)
// 	f, _ := converters.LoadAndResizeImage("./t.png", 500, 500)

// 	for i := range 3 {
// 		dc := gg.NewContext(1000, 1000)
// 		f = converters.CropAndRotateImage(f, float64(8*i))
// 		dc, err := converters.StackImages(
// 			dc,
// 			[]image.Image{
// 				f,
// 				b,
// 			},
// 			[]converters.Coord{
// 				converters.Coord{x: 250, y: 250},
// 				converters.Coord{x: 0, y: 0},
// 			},
// 		)
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		dc.SavePNG(fmt.Sprintf("./t_%d.png", i))
// 	}

// }
