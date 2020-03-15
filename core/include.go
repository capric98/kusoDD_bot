package core

import (
	"io"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewMessage(chatid int64, text string) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatid, text)
}

func NewFileConfig(fid string) tgbotapi.FileConfig {
	return tgbotapi.FileConfig{
		FileID: fid,
	}
}

func NewAnimationUpload(chatid int64, file interface{}) tgbotapi.AnimationConfig {
	return tgbotapi.NewAnimationUpload(chatid, file)
}

func NewDocumentUpload(chatid int64, file interface{}) tgbotapi.DocumentConfig {
	return tgbotapi.NewDocumentUpload(chatid, file)
}

func NewFileBytes(filename string, fr io.Reader, size int64) tgbotapi.FileReader {
	return tgbotapi.FileReader{
		Name:   filename,
		Reader: fr,
		Size:   size,
	}
}

func NewChatAction(chatid int64, action string) tgbotapi.ChatActionConfig {
	return tgbotapi.NewChatAction(chatid, action)
}
