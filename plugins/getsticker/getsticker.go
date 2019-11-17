package getsticker

import (
	"bytes"
	"errors"
	"image/png"
	"net/http"
	"time"

	"golang.org/x/image/webp"
)

type tgbot interface {
	GetFile(map[string]string) string
	SendMessage(map[string]string) error
	SendDocument(map[string]string, string, []byte) string
	Log(interface{}, int)
}
type message interface {
	GetChatIDStr() string
	GetReplyToStickerFileID() string
	GetReplyToStickerSetName() string
	GetReplyMsgIDStr() string
}

var (
	client = &http.Client{
		Timeout: 10 * time.Second,
	}
	ErrNoSticker     = errors.New("commands: No sticker in the message.")
	ErrEmptyResponse = errors.New("commands: getFile call returns an empty response.")
)

func Handle(m interface{}, b interface{}) error {
	msg := m.(message)
	bot := b.(tgbot)

	ID := msg.GetReplyToStickerFileID()
	paras := map[string]string{
		"chat_id":             msg.GetChatIDStr(),
		"reply_to_message_id": msg.GetReplyMsgIDStr(),
	}

	if ID == "" {
		bot.Log(ErrNoSticker, 0)
		bot.Log(msg, 0)
		paras["text"] = "请对一个sticker回复该指令"
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
