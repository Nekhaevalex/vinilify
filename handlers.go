package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tg "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func handleStart(bot *tg.Bot, update tg.Update) {

	userID := update.Message.From.ID

	//check if user is in the map, if not - add
	_, ok := users[userID]
	if !ok {
		users[userID] = user{
			Id:       userID,
			State:    0,
			Cooldown: time.Now(),
			AudioURL: "",
			ImageURL: "",
		}
	}
	usr := users[userID]

	// Checking if "users" directory exists
	usersExists, _ := dirExists("./users")
	if !usersExists {
		// if not -- create
		os.Mkdir("./users", os.ModeDir)
	}

	// Checking if user directory exists
	userDirExists, _ := dirExists(fmt.Sprintf("./users/%d", userID))
	if !userDirExists {
		os.Mkdir(fmt.Sprintf("./users/%d", userID), os.ModeDir)
	}

	//send keyboard
	keyboard, err := usr.generateKeyboard()

	var msg *tg.SendMessageParams
	switch err {
	case ErrorNothingToDisplay:
		msg = tu.Message(
			tu.ID(userID),
			"Upload audio and cover image",
		)
	case ErrorUnknownState:
		msg = tu.Message(
			tu.ID(userID),
			"Something wrong happened",
		)
	case nil:
		msg = tu.Message(
			tu.ID(userID),
			"Choose an option below",
		).WithReplyMarkup(keyboard)
	}

	bot.SendMessage(msg)
}

func handleDeleteImage(bot *tg.Bot, update tg.Update) {

}

func handleDeleteAudio(bot *tg.Bot, update tg.Update) {

}

func handleUpload(bot *tg.Bot, update tg.Update) {
	userID := update.Message.From.ID
	user, ok := users[userID]
	if !ok {
		msg := tu.Message(
			update.Message.Chat.ChatID(),
			"Send /start command to begin",
		)

		bot.SendMessage(msg)

		return
	}

	if update.Message.Photo != nil {
		photosCount := len(update.Message.Photo)
		if photosCount != 1 {
			bot.SendMessage(tu.Message(update.Message.Chat.ChatID(), "Please, send a single image"))
			return
		}

		file, err := bot.GetFile(&tg.GetFileParams{FileID: update.Message.Photo[0].FileID})
		if err != nil {
			log.Panic("can't retreive image file url")
			return
		}
		user.ImageURL = bot.FileDownloadURL(file.FilePath)
	}

	if update.Message.Audio != nil {
		file, err := bot.GetFile(&tg.GetFileParams{FileID: update.Message.Audio.FileID})
		if err != nil {
			log.Panic("can't retreive audio file url")
			return
		}
		user.AudioURL = bot.FileDownloadURL(file.FilePath)
	}
}
