package converters

import (
	"fmt"
	"image"
	"path/filepath"
	"testing"

	"github.com/Nekhaevalex/vinilify/utils"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

func TestCropAndRotateImage(t *testing.T) {
	testAssetsPath := filepath.Join(utils.GetRoot(), "utils", "converters", "test_assets")
	testResultsPath := filepath.Join(utils.GetRoot(), "utils", "converters", "test_results")

	img, err := gg.LoadImage(filepath.Join(testAssetsPath, "t.png"))
	if err != nil {
		t.Error(err)
	}

	img = resize.Resize(1000, 1000, CropAndRotateImage(img, float64(50)), resize.Lanczos3)
	dc := gg.NewContext(1000, 1000)
	dc.DrawImage(img, 0, 0)
	dc.SavePNG(filepath.Join(testResultsPath, "test1.png"))

}

func TestStackImage(t *testing.T) {
	testAssetsPath := filepath.Join(utils.GetRoot(), "utils", "converters", "test_assets")
	testResultsPath := filepath.Join(utils.GetRoot(), "utils", "converters", "test_results")

	img1, err := gg.LoadImage(filepath.Join(testAssetsPath, "t.png"))
	if err != nil {
		t.Error(err)
	}

	back1, _ := gg.LoadImage(filepath.Join(utils.GetAssets(), "Images", "Disk.png"))

	b1s := resize.Resize(1000, 1000, back1, resize.Lanczos3)
	i1s := resize.Resize(500, 500, img1, resize.Lanczos3)

	dc := gg.NewContext(1000, 1000)
	dc.DrawImage(i1s, 250, 250)
	dc.DrawImage(b1s, 0, 0)
	dc.SavePNG(filepath.Join(testResultsPath, "test2.png"))
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
		{x: 0, y: 0},
		{x: 250, y: 250},
		{x: 375, y: 375},
	}

	dc, err := StackImages(gg.NewContext(1000, 1000), imgs, coords)
	if err != nil {
		t.Error(err)
	}
	dc.SavePNG("./test_stack.png")
}

func TestAssembleImage(t *testing.T) {
	testAssetsPath := filepath.Join(utils.GetRoot(), "utils", "converters", "test_assets")
	testResultsPath := filepath.Join(utils.GetRoot(), "utils", "converters", "test_results")

	b, _ := LoadAndResizeImage(filepath.Join(utils.GetAssets(), "Images", "Disk.png"), 1000, 1000)
	f, _ := LoadAndResizeImage(filepath.Join(testAssetsPath, "t.png"), 500, 500)

	for i := range 3 {
		dc := gg.NewContext(1000, 1000)
		f = CropAndRotateImage(f, float64(8*i))
		dc, err := StackImages(
			dc,
			[]image.Image{
				f,
				b,
			},
			[]Coord{
				{x: 250, y: 250},
				{x: 0, y: 0},
			},
		)
		if err != nil {
			t.Error(err)
		}

		dc.SavePNG(filepath.Join(testResultsPath, fmt.Sprintf("t_%d.png", i)))
	}
}
