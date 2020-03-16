package getsticker

import (
	"bytes"
	"image/png"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/capric98/kusoDD_bot/core"
	"golang.org/x/image/webp"
)

var (
	ack    = make(chan bool, 1)
	client *http.Client

	python, script string
)

func New(settings map[string]interface{}) (func(core.Message) <-chan bool, time.Duration, []core.Opt) {
	if settings["python"] != nil {
		python = settings["python"].(string)
	}
	if settings["script"] != nil {
		script = settings["script"].(string)
	}

	checkTmp()

	client = &http.Client{
		Timeout: 5 * time.Second,
	}
	return func(msg core.Message) <-chan bool {
		return handle(msg)
	}, 5 * time.Second, nil
}

func handle(msg core.Message) <-chan bool {
	select {
	case <-ack:
	default:
	}
	if msg.Message == nil {
		ack <- false
		return ack
	}

	go func() {
		if msg.Message.IsCommand() && msg.Message.Command() == "getsticker" {
			ack <- true
			sticker := msg.Message.Sticker
			if sticker == nil && msg.Message.ReplyToMessage != nil {
				sticker = msg.Message.ReplyToMessage.Sticker
			}
			if sticker == nil {
				resp := core.NewMessage(msg.Message.Chat.ID, "未找到sticker,请对一个sticker回复该指令")
				resp.ReplyToMessageID = msg.Message.MessageID
				if _, e := msg.Bot.Send(resp); e != nil {
					msg.Bot.Printf("%6s - getsticker failed to send response: \"%v\".\n", "warn", e)
				}
			} else {
				go func() { _, _ = msg.Bot.Send(core.NewChatAction(msg.Message.Chat.ID, "UPLOAD_DOCUMENT")) }()

				filename := sticker.SetName + "-" + sticker.FileID + ".png"

				slink, e := msg.Bot.GetFileDirectURL(sticker.FileID)
				if e != nil {
					msg.Bot.Printf("%6s - getsticker failed to get file link: \"%v\".\n", "warn", e)
					return
				}

				sresp, e := client.Get(slink)
				if e != nil {
					msg.Bot.Printf("%6s - getsticker failed to download file: \"%v\".\n", "warn", e)
					return
				}
				ibody, e := ioutil.ReadAll(sresp.Body)
				sresp.Body.Close()
				if e != nil {
					msg.Bot.Printf("%6s - getsticker failed to download file: \"%v\".\n", "warn", e)
					return
				}

				image, e := webp.Decode(bytes.NewReader(ibody))
				if e != nil {
					//msg.Bot.Printf("%6s - getsticker failed to decode sticker file: \"%v\".\n", "warn", e)
					decodeTGS(ibody, sticker.SetName+"-"+sticker.FileID, msg)
					return
				}
				var buf bytes.Buffer
				e = png.Encode(&buf, image)
				if e != nil {
					msg.Bot.Printf("%6s - getsticker failed to encode sticker file: \"%v\".\n", "warn", e)
					return
				}
				resp := core.NewDocumentUpload(
					msg.Message.Chat.ID,
					core.NewFileBytes(filename, &buf, int64(buf.Len())),
				)

				if _, e := msg.Bot.Send(resp); e != nil {
					msg.Bot.Printf("%6s - getsticker failed to send response: \"%v\".\n", "warn", e)
				}
			}
		} else {
			ack <- false
		}
	}()
	return ack
}
