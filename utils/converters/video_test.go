package converters

import (
	"fmt"
	"os"
	"testing"
)

func TestSecondVideo(t *testing.T) {
	inputPattern := "../../users/286946560/%02d.png"
	// inputPattern := "%02d.png"
	outputFile := "../../users/286946560/secondvideo.mp4"

	for i := 1; i <= 24; i++ {
		fileName := "../../users/286946560/" + fmt.Sprintf("%02d", i) + ".png"
		//fileName := fmt.Sprintf("%02d", i) + ".png"
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			t.Error(err)
		}
	}

	err := SecondVideo(inputPattern, outputFile)
	if err != nil {
		t.Error(err)
	}
}

func TestLoopVideo(t *testing.T) {
	inputFile := "../../users/286946560/secondvideo.mp4"
	outputFile := "minutevideo.mp4"
	LoopVideo(inputFile, outputFile)
}
