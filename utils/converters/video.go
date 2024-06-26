package converters

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

func SecondVideo(inputPattern, outputFile string) error {

	// Use ffmpeg-go to create a video from images
	fmt.Println(inputPattern)

	err := ffmpeg_go.
		Input(inputPattern, ffmpeg_go.KwArgs{"r": 24}).
		Output(outputFile, ffmpeg_go.KwArgs{"vcodec": "libx264", "pix_fmt": "yuv420p"}).
		OverWriteOutput().
		Run()

	if err != nil {
		return err
	}

	log.Println("Video created!")
	return nil
}

func LoopVideo(inputFile, outputFile string) error {

	err := ffmpeg_go.
		Input(inputFile, ffmpeg_go.KwArgs{"stream_loop": "50"}).
		Output(outputFile).
		OverWriteOutput().
		Run()

	if err != nil {
		return err
	}

	return nil
}

func AddAudio(audioFile, videoFile, outputFile string) error {
	cmd := exec.Command("ffmpeg", "-i", videoFile, "-i", audioFile, "-c", "copy", "-map", "0:v:0", "-map", "1:a:0", outputFile)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
