package tracemoe

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type tgbot interface {
	GetFile(map[string]string) string
	SendMessage(map[string]string) error
	SendDocument(map[string]string, string, []byte) string
	SendAnimation(paras map[string]string, filename string, data []byte) (fileID string)

	GetConfig(name string) map[string]interface{}
	Log(interface{}, int)
}
type message interface {
	GetChatIDStr() string
	GetPhotoFileID() string
	GetMsgIDStr() string
	GetReplyMsgIDStr() string
	GetReplyToPhotoFileID() string
}

type traceresp struct {
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
	client = &http.Client{
		Timeout: 10 * time.Second,
	}
	json   = jsoniter.ConfigCompatibleWithStandardLibrary
	prefix = "https://trace.moe/api/search?token="
	token  = ""
	//mediaprefix = "https://media.trace.moe/video/${anilist_id}/${encodeURIComponent(filename)}?t=${at}&token=${tokenthumb}&mute"

	ErrNoPhoto = errors.New("tracemoe: No photo in this message.")
)

func New(b interface{}) func(interface{}, interface{}) error {
	conf := b.(tgbot).GetConfig("tracemoe")
	prefix += conf["token"].(string) + "&url="
	token = conf["token"].(string)
	b.(tgbot).Log("tracemoe: set prefix -> "+prefix, 0)
	return Handle
}

func Handle(m interface{}, b interface{}) error {
	return handle(m.(message), b.(tgbot))
}

func handle(msg message, bot tgbot) error {
	ID := msg.GetPhotoFileID()
	paras := map[string]string{
		"reply_to_message_id": msg.GetReplyMsgIDStr(),
		"chat_id":             msg.GetChatIDStr(),
		"parse_mode":          "Markdown",
	}

	if ID == "" {
		ID = msg.GetReplyToPhotoFileID()
		if ID == "" {
			paras["text"] = "未找到图片，请对图片内容回复该命令。"
			bot.SendMessage(paras)
			return ErrNoPhoto
		}
	}
	u := bot.GetFile(map[string]string{"file_id": ID})
	bot.Log("tracemoe: pic url -> "+u, 0)

	fileresp, _ := client.Get(u)
	body, _ := ioutil.ReadAll(fileresp.Body)
	fileresp.Body.Close()

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	w.WriteField("token", token)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="image"; filename="%s"`, getFilename(u)))
	p, _ := w.CreatePart(h)
	_, _ = p.Write(body)
	w.Close()
	req, _ := http.NewRequest("POST", "https://trace.moe/api/search", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		bot.Log("tracemoe: Making request failed.", 0)
		return err
	}

	var tresp traceresp
	err = json.NewDecoder(resp.Body).Decode(&tresp)
	resp.Body.Close()
	if err != nil {
		bot.Log("tracemoe: Json decode failed.", 0)
		resp, _ = client.Get(prefix + url.QueryEscape(u))
		bot.Log("tracemoe: Dump resp.Body ->", 0)
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		bot.Log(string(body), 0)

		return err
	}

	if len(tresp.Docs) != 0 {
		doc0 := tresp.Docs[0]
		mediaUrl := "https://media.trace.moe/video/" + strconv.Itoa(doc0.AnilistID) +
			"/" + url.QueryEscape(doc0.Filename) + "?t=" +
			strconv.FormatFloat(doc0.At, 'f', 2, 64) + "&token=" + doc0.TokenThumb + "&mute"

		bot.Log(mediaUrl, 0)

		paras["animation"] = mediaUrl
		paras["caption"] = "*Similarity:* `" + strconv.FormatFloat(doc0.Similarity*100, 'f', 2, 64) +
			"`\n*タイトル:* `" + doc0.Title + "`\n*Title:* `" + doc0.TitleRomaji +
			"`\n*File:* `" + doc0.Filename +
			"`\n*From* `" + strconv.FormatFloat(doc0.From, 'f', 2, 64) + "s` *to* `" +
			strconv.FormatFloat(doc0.To, 'f', 2, 64) + "s`"

		_ = bot.SendAnimation(paras, "", nil)
		return nil
	} else {
		paras["text"] = fmt.Sprintf("tracemoe: no Docs! Limit left: %d, Quota left %d", tresp.Limit, tresp.Quota)
		return bot.SendMessage(paras)
	}
}

func getFilename(s string) string {
	i := len(s) - 1
	for i > 0 && rune(s[i]) != '/' {
		i--
	}
	return s[i+1:]
}
