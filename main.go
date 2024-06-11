package main

import (
	"fmt"
	"os"
	"time"

	tg "github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type User struct {
	State      string
	Audiofile  string
	Image      string
	Generating bool
	Cooldown   time.Time
}

func (u *User) ChangeState(newState string) {
	u.State = newState
}

var m map[int]User

func main() {

	//Initializing the map of the users
	m = make(map[int]User)

	botToken := "7152618794:AAHS6f-lNvaW_jBKrBNl0jJv-jc_Y5GC7pM"

	bot, err := tg.NewBot(botToken, tg.WithDefaultDebugLogger())

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	updates, err := bot.UpdatesViaLongPolling(&tg.GetUpdatesParams{Timeout: 10})
	if err != nil {
		fmt.Println(err)
	}

	bh, err := th.NewBotHandler(bot, updates)
	if err != nil {
		fmt.Println(err)
	}

	defer bh.Stop()
	defer bot.StopLongPolling()

	bh.Handle(
		func(bot *tg.Bot, update tg.Update) {
			_, err := bot.SendMessage(
				tu.Message(
					tu.ID(update.Message.Chat.ID),
					"Response"))

			if err != nil {
				fmt.Println(err)
			}
		},
		th.CommandEqual("start"),
	)

	bh.Start()

}

// func HandleStart(bot *tg.Bot, update tg.Update) {
// 	userID :=
// }
