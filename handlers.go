package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Nekhaevalex/vinilify/types"
	"github.com/Nekhaevalex/vinilify/utils"
	"github.com/Nekhaevalex/vinilify/utils/converters"
	tg "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

const (
	//start command
	MessageSendStart      = "Send \"/start\" command to begin"
	MessageInstruction    = "Upload audio and cover image"
	MessageUnknownCommand = "Unknown command"

	//upload notificators
	MessageUploadedImage   = "Image uploaded"
	MessageUplodadedAudio  = "Audio uploaded"
	MessageReadyToGenerate = "Now send \"/generate\" to generate the vinyl"

	//error notificators
	MessageNoAudio = "You have not uploaded audio file"
	MessageNoImage = "You have not uploaded image file"

	//block notificators
	MessageGenerating = "Your video is being processed right now, wait for it to complete"
	MessageCooldown   = "You are in cooldown, wait..."

	//generator notificators
	MessageAudioDownloadFailed = "Could not download audio"
	MessageImageDownloadFailed = "Could not download image"
	MessageDownloadComplete    = "Files downloaded, generating video..."
	MessageDownloadStarted     = "Downloading files..."
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
		SendMessage(bot, update, MessageSendStart)
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

	if user.HasAudioURL() && user.HasImageURL() {
		SendMessage(bot, update, MessageReadyToGenerate)
	}

	fmt.Printf("\n\n\n%+v\n\n\n", user)

}

func handleGenerateVideo(bot *tg.Bot, update tg.Update) {
	userId := update.Message.From.ID
	user, ok := users[userId]
	if !ok {
		SendMessage(bot, update, MessageSendStart)
		return
	}

	//0. check if user is already generating video
	if user.Generating {
		SendMessage(bot, update, MessageGenerating)
		return
	}

	//0.1 check if user is in cooldown
	if time.Now().Compare(user.Cooldown) <= 0 {
		SendMessage(bot, update, MessageCooldown)
		return
	}

	//1. check if user has both audio and video file links
	if !user.HasAudioURL() {
		SendMessage(bot, update, MessageNoAudio)
		return
	}

	if !user.HasImageURL() {
		SendMessage(bot, update, MessageNoImage)
		return
	}

	//2. Go generate video note
	user.Generating = true

	go GenerateVideo(bot, update, user)
}

func GenerateVideo(bot *tg.Bot, update tg.Update, user *types.User) {

	defer func() {
		user.Generating = false
		user.Cooldown = time.Now().Add(time.Second * 100)
	}()

	assetsPath := utils.GetAssets()
	userPath := utils.GetUserPath(user.Id)

	//1. Download audio and image to the folder
	audioPath := user.GetAudioPath()
	imagePath := user.GetImagePath()
	SendMessage(bot, update, MessageDownloadStarted)

	err := utils.DownloadAttachment(audioPath, user.AudioURL)
	if err != nil {
		SendMessage(bot, update, MessageAudioDownloadFailed)
		return
	}
	err = utils.DownloadAttachment(imagePath, user.ImageURL)
	if err != nil {
		SendMessage(bot, update, MessageImageDownloadFailed)
		return
	}

	SendMessage(bot, update, MessageDownloadComplete)
	SendMessage(bot, update, "Preparing for generation...")

	//2. Mix audio with effect

	effect := assetsPath + "/sounds/vinyl.mp3"
	music := userPath + "/audio.mp3"
	mix := userPath + "/mix.mp3" //mixed audio is stored in users/.../mix.mp3

	err = converters.Mix(effect, music, mix)
	if err != nil {
		SendMessage(bot, update, "Error mixing audio "+err.Error())
		return
	}

	//3. Generate images

	image := filepath.Join(utils.GetRoot(), "users", fmt.Sprintf("%d", user.Id), "image.jpg")
	imageOut := filepath.Join(utils.GetRoot(), "users", fmt.Sprintf("%d", user.Id))
	err = converters.AssembleImages(image, imageOut) //video frames are stored in users/.../01...32.png
	if err != nil {
		SendMessage(bot, update, "Error generating images "+err.Error())
		return
	}

	//4. Generate 1 second video
	patternPath := userPath + "/%02d.png"
	secondVideoPath := userPath + "/secondvideo.mp4"
	err = converters.SecondVideo(patternPath, secondVideoPath) //video is stored in users/.../secondvideo.mp4
	if err != nil {
		SendMessage(bot, update, "Error generating a second-long video "+err.Error())
		return
	}

	//5. Generate minute long video
	minuteVideoPath := userPath + "/minutevideo.mp4"
	err = converters.LoopVideo(secondVideoPath, minuteVideoPath)
	if err != nil {
		SendMessage(bot, update, "Error generating a minute-long video "+err.Error())
		return
	}

	//6. Mix audio and video together
	videoPath := userPath + "/output.mp4"
	err = converters.AddAudio(mix, minuteVideoPath, videoPath)
	if err != nil {
		SendMessage(bot, update, "Error generating the final video "+err.Error())
		return
	}

	//7. Send the video note to the user
	SendMessage(bot, update, "Video has been generated, sending...")

	videoFile, err := os.Open(videoPath)
	if err != nil {
		SendMessage(bot, update, "Can't open the generated video file "+err.Error())
		return
	}

	bot.SendVideoNote(tu.VideoNote(update.Message.Chat.ChatID(), tg.InputFile{
		File: videoFile,
	}))

}

func SendMessage(bot *tg.Bot, update tg.Update, message string) {
	msg := tu.Message(
		update.Message.Chat.ChatID(),
		message,
	)
	bot.SendMessage(msg)
}
