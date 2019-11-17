package commands

import (
	"errors"
	. "github.com/capric98/kusoDD_bot/plugins"
)

var (
	ErrNoSticker = errors.New("commands: No sticker in the message.")
)

func sendStickerFile(msg Message, bot Tgbot) error {
	ID := msg.GetPhotoFileID()
	if ID == "" {
		bot.Log(ErrNoSticker, 0)
	} else {
		return bot.GetFile([]string{"file_id"}, []string{ID})
	}
	return nil
}
