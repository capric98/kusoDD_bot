package commands

import (
	"bytes"
	"errors"
	"image/png"

	. "github.com/capric98/kusoDD_bot/plugins"
	"golang.org/x/image/webp"
)

var (
	ErrNoSticker     = errors.New("commands: No sticker in the message.")
	ErrEmptyResponse = errors.New("commands: getFile call returns an empty response.")
)

func sendStickerFile(msg Message, bot Tgbot) error {
	ID := msg.GetReplyToStickerFileID()
	paras := map[string]string{
		"chat_id":             msg.GetChatIDStr(),
		"reply_to_message_id": msg.GetReplyMsgIDStr(),
	}

	if ID == "" {
		bot.Log(ErrNoSticker, 0)
		paras["text"] = "使用方式：对一个sticker回复该指令"
		bot.SendMessage(paras)
	} else {
		u := bot.GetFile(map[string]string{"file_id": ID})
		if u == "" {
			return ErrEmptyResponse
		}
		filename := msg.GetReplyToStickerSetName() + "-" + ID + ".png"
		resp, e := client.Get(u)
		if e != nil {
			return e
		}
		image, e := webp.Decode(resp.Body)
		resp.Body.Close()
		if e != nil {
			return e
		}

		var buf bytes.Buffer
		png.Encode(&buf, image)
		_ = bot.SendDocument(paras, filename, buf.Bytes())
	}
	return nil
}
