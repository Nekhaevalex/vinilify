package converters

import (
	"os"
	"os/exec"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

// Mixes specified user audio and vinyl audio asset

func Mix(effect, music, outfile string) error {

	//remove outfile if exists
	if fileExists(outfile) {
		if err := deleteFile(outfile); err != nil {
			return err
		}
	}

	s1 := ffmpeg_go.Input(effect)
	s2 := ffmpeg_go.Input(music)
	return ffmpeg_go.Filter(
		[]*ffmpeg_go.Stream{
			s1,
			s2,
		},
		"amix",
		ffmpeg_go.Args{"inputs=2:duration=first:dropout_transition=2"},
	).Output(outfile).Run()
}

func Convert(inputFile, outputFile string) error {
	cmd := exec.Command("ffmpeg", "-i", inputFile, outputFile, "-y")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func AddVinylEffects(inputFile, outputFile string) error {
	cmd := exec.Command("sox", inputFile, outputFile, "sinc", "220-12k")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	return err
}

// Helpers

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func deleteFile(filename string) error {
	err := os.Remove(filename)
	if err != nil {
		return err
	}
	return nil
}

// func Mix(music, effect *ffmpeg_go.Stream) *ffmpeg_go.Stream {
// 	return ffmpeg_go.Filter(
// 		[]*ffmpeg_go.Stream{
// 			music,
// 			effect,
// 		},
// 		"amix",
// 		ffmpeg_go.Args{"inputs=2:duration=first:dropout_transition=2"},
// 	)
// }
// func AssebmleAudio(music, effect string, userId int) error {

// 	outputFilePath := filepath.Join(utils.GetRoot(), "users", fmt.Sprintf("%d", userId), "audio.mp3")

// 	musicStream := Trim(ffmpeg_go.Input(music), 60)
// 	effectStream := Trim(ffmpeg_go.Input(effect), 60)
// 	s := Mix(musicStream, effectStream)
// 	err := s.Output(outputFilePath).Run()
// 	if err != nil {
// 		return err
// 	}
// 	return nil

// }
// func Trim(music *ffmpeg_go.Stream, duration int) *ffmpeg_go.Stream {
// 	return music.Trim(ffmpeg_go.KwArgs{"duration": fmt.Sprintf("%d", duration)})
// }

// func TrimFile(inputFilePath, outputFilePath string) error {
// 	err := ffmpeg_go.Input(inputFilePath).
// 		Output(outputFilePath, ffmpeg_go.KwArgs{"ss": "0", "t": "60"}).
// 		Run()
// 	return err
// }
