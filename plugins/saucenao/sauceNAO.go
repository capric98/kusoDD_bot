package saucenao

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"net/url"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type tgbot interface {
	GetFile(map[string]string) string
	SendMessage(map[string]string) error
	SendDocument(map[string]string, string, []byte) string

	GetConfig(name string) map[string]interface{}
	Log(interface{}, int)
}
type message interface {
	GetChatIDStr() string
	GetPhotoFileID() string
	GetReplyToPhotoFileID() string
	GetReplyMsgIDStr() string
}

type saucenaoresp struct {
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
		} `json:"data"`
	} `json:"results"`
} //simple

var (
	client = &http.Client{
		Timeout: 10 * time.Second,
	}
	prefix     = "https://saucenao.com/search.php?db=999&output_type=2&numres=1"
	ErrNoPhoto = errors.New("saucenao: No photo in this message.")
)

func New(b interface{}) func(interface{}, interface{}) error {
	conf := b.(tgbot).GetConfig("sauceNAO")
	prefix += "&api_key=" + conf["key"].(string) + "&url="
	return Handle
}

func Handle(m interface{}, b interface{}) error {
	return handle(m.(message), b.(tgbot))
}

func handle(msg message, bot tgbot) error {
	ID := msg.GetPhotoFileID()
	if ID == "" {
		ID = msg.GetReplyToPhotoFileID()
		if ID == "" {
			return ErrNoPhoto
		}
	}
	u := bot.GetFile(map[string]string{"file_id": ID})
	bot.Log("sauceNAO: pic url -> "+u, 0)
	resp, err := client.Get(prefix + url.QueryEscape(u))
	if err != nil {
		return err
	}

	var sresp saucenaoresp
	err = json.NewDecoder(resp.Body).Decode(&sresp)
	resp.Body.Close()
	if err != nil {
		return err
	}

	var paras map[string]string
	paras["reply_to_message_id"] = msg.GetReplyMsgIDStr()
	paras["chat_id"] = msg.GetChatIDStr()
	paras["text"] = "author: " + sresp.Results[0].Data.Author + "\nSimilarity:" + sresp.Results[0].Header.Similarity

	return bot.SendMessage(paras)
}
