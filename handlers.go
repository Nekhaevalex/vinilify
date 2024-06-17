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

const (
	MessageGenerating       = "Your video is being processed right now, wait for it to complete"
	MessageCooldown         = "You are in cooldown, wait..."
	MessageUnknownCommand   = "Unknown command"
	MessageUploadedImage    = "Image uploaded"
	MessageUplodadedAudio   = "Audio uploaded"
	MessageReadyToGenerate  = "Now send \"/generate\" to generate the vinyl"
	MessageSendStartCommand = "Send \"/start\" command to begin"
	MessageInstruction      = "Upload audio and cover image"
)

func handleStart(bot *tg.Bot, update tg.Update) {

	userID := update.Message.From.ID

	//check if user is in the map, if not - add
	_, ok := users[userID]
	if !ok {
		users[userID] = &types.User{
			Id:       userID,
			State:    0,
			Cooldown: time.Now(),
			AudioURL: "",
			ImageURL: "",
		}
	}

	// Checking if "users" directory exists
	usersExists, _ := utils.DirExists("./users")
	if !usersExists {
		// if not -- create
		createdErr := os.Mkdir("./users", 0755)
		if createdErr != nil {
			log.Fatal(createdErr)
		}
	}

	// Checking if user exists
	userDirExists, _ := utils.DirExists(fmt.Sprintf("./users/%d", userID))
	if !userDirExists {
		createdErr := os.Mkdir(fmt.Sprintf("./users/%d", userID), 0755)
		if createdErr != nil {
			log.Fatal(createdErr)
		}
	}

	SendMessage(bot, update, MessageInstruction)
}

func handleUpload(bot *tg.Bot, update tg.Update) {
	userID := update.Message.From.ID
	user, ok := users[userID]
	if !ok {
		SendMessage(bot, update, MessageSendStartCommand)
		return
	}

	if update.Message.Audio == nil && update.Message.Photo == nil {
		SendMessage(bot, update, MessageUnknownCommand)
		return
	}

	if update.Message.Photo != nil {
		photosCount := len(update.Message.Photo)

		file, err := bot.GetFile(&tg.GetFileParams{FileID: update.Message.Photo[photosCount-1].FileID})
		if err != nil {
			log.Panic("can't retreive image file url")
			return
		}
		user.ImageURL = bot.FileDownloadURL(file.FilePath)
		SendMessage(bot, update, MessageUploadedImage)
	}

	if update.Message.Audio != nil {
		file, err := bot.GetFile(&tg.GetFileParams{FileID: update.Message.Audio.FileID})
		if err != nil {
			log.Panic("can't retreive audio file url")
			return
		}
		user.AudioURL = bot.FileDownloadURL(file.FilePath)
		SendMessage(bot, update, MessageUplodadedAudio)
	}

	fmt.Printf("\n\n\n%+v\n\n\n", user)

	if user.HasAudioURL() && user.HasImageURL() {
		SendMessage(bot, update, MessageReadyToGenerate)
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

func SendMessage(bot *tg.Bot, update tg.Update, message string) {
	msg := tu.Message(
		update.Message.Chat.ChatID(),
		message,
	)
	bot.SendMessage(msg)
}
