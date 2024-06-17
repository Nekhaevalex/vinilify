package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Nekhaevalex/vinilify/types"
	"github.com/Nekhaevalex/vinilify/utils"
	tg "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func handleStart(bot *tg.Bot, update tg.Update) {

	userID := update.Message.From.ID

	//check if user is in the map, if not - add
	_, ok := users[userID]
	if !ok {
		users[userID] = types.User{
			Id:       userID,
			State:    0,
			Cooldown: time.Now(),
			AudioURL: "",
			ImageURL: "",
		}
	}
	usr := users[userID]

	// Checking if "users" directory exists
	usersExists, _ := utils.DirExists("./users")
	if !usersExists {
		// if not -- create
		os.Mkdir("./users", os.ModeDir)
	}

	// Checking if user directory exists
	userDirExists, _ := utils.DirExists(fmt.Sprintf("./users/%d", userID))
	if !userDirExists {
		os.Mkdir(fmt.Sprintf("./users/%d", userID), os.ModeDir)
	}

	//send keyboard
	keyboard, err := usr.GenerateKeyboard()

	var msg *tg.SendMessageParams
	switch err {
	case types.ErrorNothingToDisplay:
		msg = tu.Message(
			tu.ID(userID),
			"Upload audio and cover image",
		)
	case types.ErrorUnknownState:
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

func handleGenerateVideo(bot *tg.Bot, update tg.Update) {
	userId := update.Message.From.ID
	user, ok := users[userId]
	if !ok {
		msg := tu.Message(
			update.Message.Chat.ChatID(),
			"Send /start command to begin",
		)

		bot.SendMessage(msg)

		return
	}

	//0. check if user is already generating video
	if user.Generating {
		msg := tu.Message(
			update.Message.Chat.ChatID(),
			MessageGenerating,
		)
		bot.SendMessage(msg)
		return
	}

	//0.1 check if user is in cooldown
	if time.Now().Compare(user.Cooldown) <= 0 {
		msg := tu.Message(
			update.Message.Chat.ChatID(),
			MessageCooldown,
		)
		bot.SendMessage(msg)
		return
	}

	//1. check if user has both audio and video file links
	if !user.HasAudioURL() {
		msg := tu.Message(
			update.Message.Chat.ChatID(),
			"You have not uploaded audio file",
		)
		bot.SendMessage(msg)
		return
	}

	if !user.HasImageURL() {
		msg := tu.Message(
			update.Message.Chat.ChatID(),
			"You have not uploaded image file",
		)
		bot.SendMessage(msg)
		return
	}

	//2. Run the thread for generation of the video
	// go func(bot *tg.Bot, update tg.Update, u User) {
	// 	video := user.GenerateVideo()
	//  msg := tu.VideoNote(video)
	//  bot.SendMessage(
	//		...
	//	)
	// }()
}

const MessageGenerating = "Your video is being processed right now, wait for it to complete"
const MessageCooldown = "You are in cooldown, wait..."
