package converters

import (
	"path/filepath"
	"testing"

	"github.com/Nekhaevalex/vinilify/utils"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

func TestTrim(t *testing.T) {
	testResultsPath := filepath.Join(utils.GetRoot(), "utils", "converters", "test_results", "trimmed.mp3")
	soundAssetPath := filepath.Join(utils.GetAssets(), "Sounds", "Vinyl.mp3")
	err := Trim(ffmpeg_go.Input(soundAssetPath)).Output(filepath.Join(testResultsPath, "trimmed.mp3")).Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestTrimFile(t *testing.T) {
	inputFilePath := filepath.Join(utils.GetAssets(), "sounds", "vinyl.mp3")
	outputFilePath := filepath.Join(utils.GetRoot(), "utils", "converters", "test_results", "test.mp3")
	err := TrimFile(inputFilePath, outputFilePath)
	if err != nil {
		t.Fatal(err)
	}
}

// func TestMix(t *testing.T) {
// 	soundAssetPath := filepath.Join(utils.GetAssets(), "Sounds", "Vinyl.mp3")
// 	effect := Truncate(ffmpeg_go.Input(soundAssetPath))
// 	userAudio := Truncate(ffmpeg_go.Input(soundAssetPath))

// 	result := Mix(userAudio, effect)
// 	if result != nil {
// 		t.Fatal(result)
// 	}
// }

// func TestTrim(t *testing.T) {
// 	soundAssetPath := filepath.Join(utils.GetAssets(), "Sounds", "Vinyl.mp3")
// 	effect := ffmpeg_go.Input(soundAssetPath).Filter("trim", )
// }
