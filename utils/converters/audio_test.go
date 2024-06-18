package converters

import (
	"path/filepath"
	"testing"

	"github.com/Nekhaevalex/vinilify/utils"
)

func TestMixFull(t *testing.T) {
	effect := filepath.Join(utils.GetAssets(), "sounds", "vinyl.mp3")
	music := filepath.Join(utils.GetAssets(), "sounds", "music.mp3")
	outfile := "./mix.mp3"
	err := Mix(effect, music, outfile)
	if err != nil {
		t.Fatal("error: ", err)
	}
}
