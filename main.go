package main

import (
	"log"
	"os"

	"github.com/Nekhaevalex/vinilify/types"
	tg "github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func getToken() (string, error) {
	dat, err := os.ReadFile("token")
	return string(dat), err
}

// Global hashmap for storing user requests
var users map[int64]types.User

func main() {

	// Initializing the map of the users
	users = make(map[int64]types.User)

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
		handleDeleteImage,
		th.CommandEqual("remove_image"),
	)

	bh.Handle(
		handleDeleteAudio,
		th.CommandEqual("remove_audio"),
	)

	bh.Handle(
		handleGenerateVideo,
		th.CommandEqual("generate"),
	)

	bh.Handle(
		handleUpload,
		th.Any(),
	)

	bh.Start()

}
