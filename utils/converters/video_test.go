package converters

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestSecondVideo(t *testing.T) {
	inputPattern := "../../users/286946560/%02d.png"
	outputFile := "../../users/286946560/secondvideo.mp4"
	// inputPattern := "%02d.png"
	// outputFile := "secondvideo.mp4"

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
	outputFile := "../../users/286946560/minutevideo.mp4"
	LoopVideo(inputFile, outputFile)
}

func TestAddAudioCmd(t *testing.T) {
	videoPath := "../../users/286946560/minutevideo.mp4"
	audioPath := "../../users/286946560/mix.mp3"
	outPath := "../../users/286946560/output.mp4"

	cmd := exec.Command("ffmpeg", "-i", videoPath, "-i", audioPath, "-c", "copy", "-map", "0:v:0", "-map", "1:a:0", outPath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		t.Error(err)
	}

}

func TestCompress(t *testing.T) {
	inPath := "../../users/286946560/output.mp4"
	outPath := "../../users/286946560/output_compressed.mp4"
	cmd := exec.Command("ffmpeg", "-i", inPath, "-vcodec", "h264", "-acodec", "mp2", outPath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		t.Error(err)
	}
}

func TestAddAudio(t *testing.T) {
	videoPath := "../../users/286946560/minutevideo.mp4"
	audioPath := "../../users/286946560/mix.mp3"
	outPath := "../../users/286946560/output.mp4"

	err := AddAudio(audioPath, videoPath, outPath)
	if err != nil {
		t.Error(err)
	}
}
