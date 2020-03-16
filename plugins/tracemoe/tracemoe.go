package tracemoe

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strconv"
	"time"

	"github.com/capric98/kusoDD_bot/core"
)

type Trace struct {
	Docs  []doc `json:"docs"`
	Limit int   `json:"limit"`
	Quota int   `json:"quota"`
}

type doc struct {
	From        float64 `json:"from"`
	To          float64 `json:"to"`
	At          float64 `json:"at"`
	Similarity  float64 `json:"similarity"`
	AnilistID   int     `json:"anilist_id"`
	TokenThumb  string  `json:"tokenthumb"`
	Filename    string  `json:"filename"`
	Title       string  `json:"title"`
	TitleRomaji string  `json:"title_romaji"`
} //simple

var (
	client *http.Client
	ack    = make(chan bool, 1)
	prefix = "https://trace.moe/api/search?token="
	token  = ""
)

func New(settings map[string]interface{}) (func(core.Message) <-chan bool, time.Duration, []core.Opt) {
	client = &http.Client{Timeout: 10 * time.Second}

	token = settings["token"].(string)
	prefix += token + "&url="

	return func(msg core.Message) <-chan bool {
		select {
		case <-ack:
		default:
		}
		go handle(msg)
		return ack
	}, 5 * time.Second, nil
}

func handle(msg core.Message) {
	if msg.Message == nil {
		ack <- false
		return
	}
	if (msg.Message.IsCommand() && msg.Message.Command() == "whatanime") || msg.CaptionCommand() == "whatanime" {
		ack <- true
	} else {
		ack <- false
		return
	}

	resp := core.NewMessage(msg.Message.Chat.ID, "")
	defer func() {
		if e := recover(); e != nil {
			resp.Text = "查询失败，请稍后重试。"
			_, _ = msg.Bot.Send(resp)
			msg.Bot.Printf("%6s - tracemoe failed: \"%v\".\n", "warn", e)
		}
	}()

	var maxsize int
	var fid string

	if msg.Message.ReplyToMessage != nil {
		for _, p := range msg.Message.ReplyToMessage.Photo {
			size := p.Width * p.Height
			if size > maxsize {
				maxsize = size
				fid = p.FileID
			}
		}
		resp.ReplyToMessageID = msg.Message.ReplyToMessage.MessageID
	} else {
		for _, p := range msg.Message.Photo {
			size := p.Width * p.Height
			if size > maxsize {
				maxsize = size
				fid = p.FileID
			}
		}
	}

	if fid == "" {
		resp.Text = "消息内未发现图片，请尝试对一张图片回复/whatanime或者带上该命令直接发送一张图片。"
	} else {
		go func() { _, _ = msg.Bot.Send(core.NewChatAction(msg.Message.Chat.ID, "TYPING")) }()

		u, e := msg.Bot.GetFileDirectURL(fid)
		if e != nil {
			msg.Bot.Printf("%6s - tracemoe failed to get direct url: \"%v\".\n", "warn", e)
			return
		}

		fileresp, _ := client.Get(u)
		body, _ := ioutil.ReadAll(fileresp.Body)
		fileresp.Body.Close()

		buf := new(bytes.Buffer)
		w := multipart.NewWriter(buf)

		_ = w.WriteField("token", token)

		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition",
			`form-data; name="image"; filename="`+getFilename(u)+`"`)
		p, _ := w.CreatePart(h)
		_, _ = p.Write(body)
		w.Close()
		req, _ := http.NewRequest("POST", "https://trace.moe/api/search", buf)
		req.Header.Set("Content-Type", w.FormDataContentType())

		aresp, err := client.Do(req)
		for count := 0; count < 3 && err != nil; {
			count++
			aresp, err = client.Do(req)
		}
		if err != nil {
			panic(err)
		}

		var tresp Trace
		err = msg.Bot.Json.NewDecoder(aresp.Body).Decode(&tresp)
		aresp.Body.Close()
		if err != nil {
			msg.Bot.Println("tracemoe: Json decode failed.", 0)
			aresp, _ = client.Get(prefix + url.PathEscape(u))
			msg.Bot.Println("tracemoe: Dump resp.Body ->", 0)
			body, _ := ioutil.ReadAll(aresp.Body)
			aresp.Body.Close()
			msg.Bot.Println(string(body), 0)

			return
		}
		if len(tresp.Docs) != 0 {
			go func() { _, _ = msg.Bot.Send(core.NewChatAction(msg.Message.Chat.ID, "UPLOAD_VIDEO")) }()

			doc0 := tresp.Docs[0]
			mediaUrl := "https://media.trace.moe/video/" + strconv.Itoa(doc0.AnilistID) +
				"/" + url.PathEscape(doc0.Filename) + "?t=" +
				strconv.FormatFloat(doc0.At, 'f', -1, 64) + "&token=" + doc0.TokenThumb + "&mute"

			if mresp, err := client.Get(mediaUrl); err != nil {
				msg.Bot.Printf("%6s - tracemoe failed to download media: \"%v\".\n", "warn", e)
			} else {
				fr := core.NewFileBytes("capture.mp4", mresp.Body, mresp.ContentLength)
				ac := core.NewAnimationUpload(msg.Message.Chat.ID, fr)
				ac.ParseMode = "Markdown"
				ac.Caption = "*Similarity:* `" + strconv.FormatFloat(doc0.Similarity*100, 'f', 2, 64) +
					"%`\n*タイトル:* `" + doc0.Title + /*"`\n*Title:* `" + doc0.TitleRomaji +*/ //deprecated
					"`\n*File:* `" + doc0.Filename +
					"`\n*From* `" + strconv.FormatFloat(doc0.From, 'f', -1, 64) + "s` *to* `" +
					strconv.FormatFloat(doc0.To, 'f', -1, 64) + "s`"
				if _, e := msg.Bot.Send(ac); e != nil {
					msg.Bot.Printf("%6s - tracemoe failed to send response: \"%v\".\n", "warn", e)
				}
			}

			return
		}
		// If there are some results back, the function will return before these.
		msg.Bot.Printf("%6s - tracemoe: no Docs! Limit left: %d, Quota left %d\n", "warn", tresp.Limit, tresp.Quota)
		resp.Text = "无查询结果，这可能是由于图片中带有黑边、遮挡物，或者trace.moe未收录"
	}
	_, _ = msg.Bot.Send(resp)

}

func getFilename(s string) string {
	i := len(s) - 1
	for i > 0 && rune(s[i]) != '/' {
		i--
	}
	return s[i+1:]
}
