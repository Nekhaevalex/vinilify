package main

import (
	"errors"
	"fmt"
	"image"
	"log"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

// Structs and interfaces
type Coord struct {
	x int
	y int
}

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

// Errors
var (
	ErrorEmptyArray         = errors.New("empty array of images")
	ErrorImageCoordMismatch = errors.New("len of coords must be equal to len of images")
)

//Functions

// Rotates the image about the center and returns it
func CropAndRotateImage(img image.Image, d float64) image.Image {

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	target := min(width, height)

	cropSize := image.Rect(0, 0, target, target).Add(image.Point{width/2 - target/2, height/2 - target/2})
	img = img.(SubImager).SubImage(cropSize)

	dc := gg.NewContext(img.Bounds().Dx(), img.Bounds().Dy())
	dc.RotateAbout(
		gg.Radians(d),
		float64(img.Bounds().Dx()/2),
		float64(img.Bounds().Dy()/2),
	)
	dc.DrawImage(img, 0, 0)
	return dc.Image()
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

// Loads an image and returns a resized image

func LoadAndResizeImage(path string, width uint, height uint) (image.Image, error) {
	img, err := gg.LoadImage(path)
	if err != nil {
		log.Panic(err)
	}

	img_res := resize.Resize(width, height, img, resize.Lanczos3)

	return img_res, nil
}

func (u User) AssembleImages() (string, error) {

	disk, _ := LoadAndResizeImage("./Assets/Images/Disk.png", 1000, 1000)
	pin, _ := LoadAndResizeImage("./Assets/Images/Pin.png", 1000, 1000)
	path, err := u.getImage()
	if err != nil {
		return "", err
	}
	userpic, _ := LoadAndResizeImage(path, 1000, 1000)

	for i := range 45 {
		dc := gg.NewContext(1000, 1000)
		us := CropAndRotateImage(userpic, float64(i*8))
		dc.DrawImage(us, 0, 0)
		dc.DrawImage(disk, 0, 0)
		dc.DrawImage(pin, 0, 0)
		dc.SavePNG(fmt.Sprintf("./users/%d/image_%d", u.Id, i))
	}

	return fmt.Sprintf("users/%d", u.Id), nil

	//1. Load image assets
	//2. Load image from the user struct
	//3. Generate stacked images with rotated user image
	//4. Return the path to the folder with the iamges

}
