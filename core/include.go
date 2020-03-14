package core

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) NewMessage(chatid int64, text string) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatid, text)
}
