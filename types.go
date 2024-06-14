package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

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
type user struct {
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
func (u user) hasAudioURL() bool {
	return u.AudioURL != ""
}

// Returns true if user's image file specified
func (u user) hasImageURL() bool {
	return u.ImageURL != ""
}

// Builds path to user's local audio file even without file existing
func (u user) getAudioPath() string {
	return filepath.Join("users", fmt.Sprintf("%d", u.Id), "audio.mp3")
}

// Builds path to user's local image file even without file existing
func (u user) getImagePath() string {
	return filepath.Join("users", fmt.Sprintf("%d", u.Id), "video.mp4")
}

// Gets path to user's local audio file creating it's local copy if one does not exists.
// Returns
func (u user) getAudio() (string, error) {
	if !u.hasAudioURL() {
		return "", ErrorNoAudioURL
	}
	_, err := os.Stat(u.getAudioPath())
	if os.IsNotExist(err) {
		return u.getAudioPath(), downloadAttachment(u.getAudioPath(), u.AudioURL)
	}
	return u.getAudioPath(), nil
}

func (u user) getImage() (string, error) {
	if !u.hasImageURL() {
		return "", ErrorNoImageURL
	}
	_, err := os.Stat(u.getImagePath())
	if os.IsNotExist(err) {
		return u.getImagePath(), downloadAttachment(u.getImagePath(), u.ImageURL)
	}
	return "", err
}

// Builds keyboard markup according to user state
func (u user) generateKeyboard() (*tg.ReplyKeyboardMarkup, error) {
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
func (u user) mixAudio() error {
	effect := ffmpeg_go.Input(SoundAssetPath)
	userAudioPath, err := u.getAudio()
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
