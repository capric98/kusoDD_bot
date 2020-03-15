package saucenao

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/capric98/kusoDD_bot/core"
)

type Sauce struct {
	Header struct {
		UserID     string `json:"user_id"`
		ShortLimit string `json:"short_limit"`
		LongLimit  string `json:"long_limit"`
	} `json:"header"`
	Results []struct {
		Header struct {
			Similarity string `json:"similarity"`
			Thumbnail  string `json:"thumbnail"`
		} `json:"header"`
		Data struct {
			ExtUrls []string `json:"ext_urls"`
			Title   string   `json:"title"`
			Author  string   `json:"author_name"`
			PixivID int      `json:"pixiv_id"`
			MemName string   `json:"member_name"`
		} `json:"data"`
	} `json:"results"`
}

var (
	client *http.Client
	ack    = make(chan bool, 1)
	prefix = "https://saucenao.com/search.php?db=999&output_type=2&numres=1"
)

func New(settings map[string]interface{}) (func(core.Message) <-chan bool, time.Duration, []core.Opt) {
	client = &http.Client{Timeout: 5 * time.Second}
	prefix += "&api_key=" + settings["key"].(string) + "&url="
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
	if (msg.Message.IsCommand() && msg.Message.Command() == "whatpic") || msg.CaptionCommand() == "whatpic" {
		ack <- true
		var maxsize int
		var fid string

		resp := core.NewMessage(msg.Message.Chat.ID, "")
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
			resp.Text = "消息内未发现图片，请尝试对一张图片回复/whatpic或者带上该命令直接发送一张图片。"
		} else {
			var e error
			u, e := msg.Bot.GetFileDirectURL(fid)
			if e != nil {
				msg.Bot.Printf("%6s - saucenao failed to get direct url: \"%v\".\n", "info", e)
			}
			go func() { _, _ = msg.Bot.Send(core.NewChatAction(msg.Message.Chat.ID, "TYPING")) }()
			resp.Text, e = search(u, msg)
			resp.ParseMode = "Markdown"
			if e != nil {
				msg.Bot.Printf("%6s - saucenao failed to search pic: \"%v\".\n", "info", e)
				resp.Text = fmt.Sprintf("查询失败：%v", e)
			}
		}
		if _, e := msg.Bot.Send(resp); e != nil {
			msg.Bot.Printf("%6s - saucenao failed to send response: \"%v\".\n", "info", e)
		}
	} else {
		ack <- false
	}
}

func search(u string, msg core.Message) (result string, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = err.(error)
		}
	}()
	resp, err := client.Get(prefix + url.PathEscape(u))
	if err != nil {
		return "", err
	}
	var sresp Sauce
	err = msg.Bot.Json.NewDecoder(resp.Body).Decode(&sresp)
	resp.Body.Close()
	if err != nil {
		return "", err
	}

	if len(sresp.Results) == 0 {
		return "查询无结果，这可能是由于图片中带有黑边、遮挡物，或者SauceNAO未收录。", nil
	}
	if len(sresp.Results[0].Data.ExtUrls) == 0 {
		return "无可靠结果，这可能是由于图片中带有黑边、遮挡物，或者SauceNAO未收录。", nil
	}
	if sresp.Results[0].Data.PixivID != 0 {
		result = "\n*Similarity :* " + sresp.Results[0].Header.Similarity
		result += "%\n*Illustrator:* " + sresp.Results[0].Data.MemName +
			"\n*Pixiv ID     :* [" + strconv.Itoa(sresp.Results[0].Data.PixivID) +
			"](" + result + sresp.Results[0].Data.ExtUrls[0] + ")"
	} else {
		result = sresp.Results[0].Data.ExtUrls[0]
		result += "\n*Similarity:* " + sresp.Results[0].Header.Similarity
	}
	return result, nil
}
