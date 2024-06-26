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

	// cmd := exec.Command("ffmpeg", "-r", "24", "-i", inputPattern, "-vcodec", "libx265", "-pix_fmt", "yuv420p", "-crf", "30", outputFile, "-y")
	// cmd.Stderr = os.Stderr
	// cmd.Stdout = os.Stdout
	// cmd.Run()

	err := ffmpeg_go.
		Input(inputPattern, ffmpeg_go.KwArgs{"r": 24}).
		Output(outputFile, ffmpeg_go.KwArgs{"vcodec": "libx265", "pix_fmt": "yuv420p", "crf": "28"}).
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
		Input(inputFile, ffmpeg_go.KwArgs{"stream_loop": "43"}).
		Output(outputFile).
		OverWriteOutput().
		Run()

	if err != nil {
		return err
	}

	return nil
}

func AddAudio(audioFile, videoFile, outputFile string) error {
	cmd := exec.Command("ffmpeg", "-i", videoFile, "-i", audioFile, "-c", "copy", "-map", "0:v:0", "-map", "1:a:0", outputFile, "-y")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
