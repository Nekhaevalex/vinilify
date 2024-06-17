package types

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Nekhaevalex/vinilify/utils"
	tg "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

// Type for representing current state of user
type State int

// State constatnts
const (
	Welcome State = iota
	HasImage
	HasAudio
	HasBoth
	Generating
)

const SoundAssetPath = "Assets/Sounds/Vinyl.mp3"

// user record structure
type User struct {
	Id         int64
	State      State
	Generating bool
	Cooldown   time.Time
	ImageURL   string
	AudioURL   string
}

var (
	// Error showing that there's nothing to display
	ErrorNothingToDisplay = errors.New("nothing to display")
	// Error showing that state value is unknown (not one of constants defined)
	ErrorUnknownState = errors.New("unknown state received")
	// No audio URL found for user
	ErrorNoAudioURL = errors.New("no audio URL")
	// No video URL found for user
	ErrorNoImageURL = errors.New("no video URL")
)

// Returns true if user's audio file specified
func (u User) HasAudioURL() bool {
	return u.AudioURL != ""
}

// Returns true if user's image file specified
func (u User) HasImageURL() bool {
	return u.ImageURL != ""
}

// Builds path to user's local audio file even without file existing
func (u User) GetAudioPath() string {
	return filepath.Join("users", fmt.Sprintf("%d", u.Id), "audio.mp3")
}

// Builds path to user's local image file even without file existing
func (u User) GetImagePath() string {
	return filepath.Join("users", fmt.Sprintf("%d", u.Id), "video.mp4")
}

// Gets path to user's local audio file creating it's local copy if one does not exists.
// Returns
func (u User) GetAudio() (string, error) {
	if !u.HasAudioURL() {
		return "", ErrorNoAudioURL
	}
	_, err := os.Stat(u.GetAudioPath())
	if os.IsNotExist(err) {
		return u.GetAudioPath(), utils.DownloadAttachment(u.GetAudioPath(), u.AudioURL)
	}
	return u.GetAudioPath(), nil
}

func (u User) GetImage() (string, error) {
	if !u.HasImageURL() {
		return "", ErrorNoImageURL
	}
	_, err := os.Stat(u.GetImagePath())
	if os.IsNotExist(err) {
		return u.GetImagePath(), utils.DownloadAttachment(u.GetImagePath(), u.ImageURL)
	}
	return "", err
}

// Builds keyboard markup according to user state
func (u User) GenerateKeyboard() (*tg.ReplyKeyboardMarkup, error) {
	var keyboard *tg.ReplyKeyboardMarkup
	var err error = nil
	switch u.State {
	case Welcome:
		err = ErrorNothingToDisplay
	case HasImage:
		keyboard = tu.Keyboard(
			tu.KeyboardRow(
				tu.KeyboardButton("/remove_image"),
			),
		)
	case HasAudio:
		keyboard = tu.Keyboard(
			tu.KeyboardRow(
				tu.KeyboardButton("/remove_audio"),
			),
		)
	case HasBoth:
		keyboard = tu.Keyboard(
			tu.KeyboardRow(
				tu.KeyboardButton("/remove_image"),
				tu.KeyboardButton("/remove_audio"),
			),
		)
	default:
		keyboard = nil
		err = ErrorUnknownState
	}
	return keyboard, err
}

// Mixes specified user audio and vinyl audio asset
func (u User) MixAudio() error {
	effect := ffmpeg_go.Input(SoundAssetPath)
	userAudioPath, err := u.GetAudio()
	if err != nil {
		return err
	}
	userAudio := ffmpeg_go.Input(userAudioPath)
	err = ffmpeg_go.Filter(
		[]*ffmpeg_go.Stream{
			userAudio,
			effect,
		},
		"amix",
		ffmpeg_go.Args{"inputs=2:duration=longest:dropout_transition=2"},
	).Output("mixed_audio.mp3").Run()
	return err
}