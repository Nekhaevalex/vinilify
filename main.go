package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	tg "github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func getToken() (string, error) {
	dat, err := os.ReadFile("token")
	return string(dat), err
}

// Global hashmap for storing user requests
var users map[int64]user

func main() {

	// Initializing the map of the users
	users = make(map[int64]user)

	// Telegram API token
	botToken, err := getToken()
	if err != nil {
		log.Fatal("Failed to find token file:", err)
	}

	// Initializing the bot
	bot, err := tg.NewBot(botToken, tg.WithDefaultDebugLogger())
	if err != nil {
		log.Fatal(err)
	}

	// Launching updates and handler
	updates, err := bot.UpdatesViaLongPolling(&tg.GetUpdatesParams{Timeout: 10})
	if err != nil {
		log.Panic(err)
	}

	// Launching bot handler
	bh, err := th.NewBotHandler(bot, updates)
	if err != nil {
		log.Panic(err)
	}

	defer func() {
		bh.Stop()
		bot.StopLongPolling()
	}()

	bh.Handle(
		handleStart,
		th.CommandEqual("start"),
	)

	bh.Handle(
		HandleDeleteImage,
		th.CommandEqual("remove_image"),
	)

	bh.Handle(
		HandleDeleteAudio,
		th.CommandEqual("remove_audio"),
	)

	bh.Handle(
		HandleUpload,
		th.Any(),
	)

	bh.Start()

}

func dirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return false, err
}

func handleStart(bot *tg.Bot, update tg.Update) {

	userID := update.Message.From.ID

	//check if user is in the map, if not - add
	_, ok := users[userID]
	if !ok {
		users[userID] = user{
			State:     0,
			Audiofile: nil,
			Image:     nil,
			Cooldown:  time.Now(),
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
	keyboard, err := usr.GenerateKeyboard()

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

func HandleDeleteImage(bot *tg.Bot, update tg.Update) {

}

func HandleDeleteAudio(bot *tg.Bot, update tg.Update) {

}

func HandleUpload(bot *tg.Bot, update tg.Update) {
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

		num_photos := len(update.Message.Photo)
		if num_photos != 1 {
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

// func EnsureStart

func (u user) DownloadAttachment(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		log.Panic(err)
		return err
	}

	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		log.Panic(err)
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Panic(err)
	}

	return err
}
