package converters

import (
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

// Mixes specified user audio and vinyl audio asset
func Mix(music, effect *ffmpeg_go.Stream) error {
	result := ffmpeg_go.Filter(
		[]*ffmpeg_go.Stream{
			music,
			effect,
		},
		"amix",
		ffmpeg_go.Args{"inputs=2:duration=longest:dropout_transition=2"},
	).Output("mixed_audio.mp3").Run()
	return result
}

func Trim(music *ffmpeg_go.Stream) *ffmpeg_go.Stream {
	return music.Trim(ffmpeg_go.KwArgs{"ss": "0", "t": "60"})
}

func TrimFile(inputFilePath, outputFilePath string) error {
	err := ffmpeg_go.Input(inputFilePath).
		Output(outputFilePath, ffmpeg_go.KwArgs{"ss": "0", "t": "60"}).
		Run()
	return err
}

// getDuration returns the duration of the audio file in seconds
// func getDuration(filePath string) (float64, error) {
// 	out, err := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", filePath).Output()
// 	if err != nil {
// 		return 0, err
// 	}
// 	durationStr := strings.TrimSpace(string(out))
// 	return strconv.ParseFloat(durationStr, 64)
// }

// func main() {
// 	file1 := "input1.mp3"
// 	file2 := "input2.mp3"
// 	outputFile := "output_mixed.mp3"

// 	// Get durations of both files
// 	duration1, err := getDuration(file1)
// 	if err != nil {
// 		log.Fatalf("Error getting duration of file %s: %v", file1, err)
// 	}

// 	duration2, err := getDuration(file2)
// 	if err != nil {
// 		log.Fatalf("Error getting duration of file %s: %v", file2, err)
// 	}

// 	// Get the shortest duration
// 	shortestDuration := duration1
// 	if duration2 < duration1 {
// 		shortestDuration = duration2
// 	}

// 	// Mix the two MP3 files and limit the duration to the shortest file
// 	err = ffmpeg_go.Input(file1).
// 		Input(file2).
// 		FilterComplex("amix=inputs=2:duration=shortest").
// 		Output(outputFile, ffmpeg_go.KwArgs{"t": strconv.Itoa(int(shortestDuration))}).
// 		Run()

// 	if err != nil {
// 		log.Fatalf("Error mixing the MP3 files: %v", err)
// 	}

// 	log.Println("MP3 files mixed successfully!")
// }
