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
	if ID == "" {
		bot.Log(ErrNoSticker, 0)
	} else {
		u := bot.GetFile([]string{"file_id"}, []string{ID})
		if u == "" {
			return ErrEmptyResponse
		}
		filename := msg.GetReplyToStickerFileName() + ".png"
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
		k := []string{"chat_id", "reply_to_message_id"}
		v := []string{msg.GetChatIDStr(), msg.GetReplyMsgIDStr()}
		_ = bot.SendDocument(k, v, filename, buf.Bytes())
	}
	return nil
}
