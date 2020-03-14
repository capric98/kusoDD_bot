package getsticker

import (
	"bytes"
	"image/png"
	"net/http"
	"time"

	"github.com/capric98/kusoDD_bot/core"
	"golang.org/x/image/webp"
)

// type tgbot interface {
// 	GetFile(map[string]string) string
// 	SendMessage(map[string]string) error
// 	SendDocument(map[string]string, string, []byte) string
// 	Log(interface{}, int)
// }
// type message interface {
// 	GetChatIDStr() string
// 	GetReplyToStickerFileID() string
// 	GetReplyToStickerSetName() string
// 	GetReplyMsgIDStr() string
// }

// var (
// 	client = &http.Client{
// 		Timeout: 10 * time.Second,
// 	}
// 	ErrNoSticker     = errors.New("commands: No sticker in the message.")
// 	ErrEmptyResponse = errors.New("commands: getFile call returns an empty response.")
// )

// func Handle(m interface{}, b interface{}) error {
// 	msg := m.(message)
// 	bot := b.(tgbot)

// 	ID := msg.GetReplyToStickerFileID()
// 	paras := map[string]string{
// 		"chat_id":             msg.GetChatIDStr(),
// 		"reply_to_message_id": msg.GetReplyMsgIDStr(),
// 	}

// 	if ID == "" {
// 		bot.Log(ErrNoSticker, 0)
// 		//bot.Log(msg, 0)
// 		paras["text"] = "请对一个sticker回复该指令"
// 		bot.SendMessage(paras)
// 	} else {
// 		u := bot.GetFile(map[string]string{"file_id": ID})
// 		if u == "" {
// 			return ErrEmptyResponse
// 		}
// 		filename := msg.GetReplyToStickerSetName() + "-" + ID + ".png"
// 		resp, e := client.Get(u)
// 		if e != nil {
// 			return e
// 		}
// 		image, e := webp.Decode(resp.Body)
// 		resp.Body.Close()
// 		if e != nil {
// 			return e
// 		}

// 		var buf bytes.Buffer
// 		png.Encode(&buf, image)
// 		_ = bot.SendDocument(paras, filename, buf.Bytes())
// 	}
// 	return nil
// }

var (
	ack    = make(chan bool, 1)
	client *http.Client
)

func New(settings map[string]interface{}) (func(core.Message) <-chan bool, time.Duration, []core.Opt) {
	client = &http.Client{
		Timeout: 5 * time.Second,
	}
	return func(msg core.Message) <-chan bool {
		return handle(msg)
	}, 5 * time.Second, nil
}

func handle(msg core.Message) (c <-chan bool) {
	select {
	case <-ack:
	default:
	}
	c = ack
	if msg.Message.IsCommand() && msg.Message.Command() == "getsticker" {
		ack <- true
		sticker := msg.Message.Sticker
		if sticker == nil {
			sticker = msg.Message.ReplyToMessage.Sticker
		}
		if sticker == nil {
			resp := core.NewMessage(msg.Message.Chat.ID, "未找到sticker,请对一个sticker回复该指令")
			resp.ReplyToMessageID = msg.Message.MessageID
			if _, e := msg.Bot.Send(resp); e != nil {
				msg.Bot.Printf("%6s - getsticker failed to send response: \"%v\".\n", "info", e)
			}
		} else {
			filename := sticker.SetName + "-" + sticker.FileID + ".png"

			slink, e := msg.Bot.GetFileDirectURL(sticker.FileID)
			if e != nil {
				msg.Bot.Printf("%6s - getsticker failed to get file link: \"%v\".\n", "info", e)
				return
			}

			sresp, e := client.Get(slink)
			if e != nil {
				msg.Bot.Printf("%6s - getsticker failed to download file: \"%v\".\n", "info", e)
				return
			}
			defer sresp.Body.Close()

			// msg.Bot.Println(sresp.Header["Content-Type"])

			image, e := webp.Decode(sresp.Body)
			if e != nil {
				msg.Bot.Printf("%6s - getsticker failed to decode sticker file: \"%v\".\n", "info", e)
				resp := core.NewMessage(msg.Message.Chat.ID, "图像解码失败，这可能是由于该sticker是使用了TGS格式的动态图像，getsticker暂时还不支持这种格式。")
				resp.ReplyToMessageID = msg.Message.MessageID
				_, _ = msg.Bot.Send(resp)
				return
			}
			var buf bytes.Buffer
			e = png.Encode(&buf, image)
			if e != nil {
				msg.Bot.Printf("%6s - getsticker failed to encode sticker file: \"%v\".\n", "info", e)
				return
			}
			resp := core.NewDocumentUpload(
				msg.Message.Chat.ID,
				core.NewFileBytes(filename, &buf, int64(buf.Len())),
			)

			if _, e := msg.Bot.Send(resp); e != nil {
				msg.Bot.Printf("%6s - getsticker failed to send response: \"%v\".\n", "info", e)
			}
		}
	} else {
		ack <- false
	}
	return
}
