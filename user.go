package main

import (
	"errors"
	"time"

	tg "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

// user record structure
type user struct {
	State      State
	Audiofile  *string
	Image      *string
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
)

// Returns true if user's audio file specified
func (u user) HasAudio() bool {
	return u.Audiofile != nil
}

// Returns true if user's image file specified
func (u user) HasImage() bool {
	return u.Image != nil
}

// Builds keyboard markup according to user state
func (u user) GenerateKeyboard() (*tg.ReplyKeyboardMarkup, error) {
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
