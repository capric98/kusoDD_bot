package core

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"sync"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Tgbot interface {
	SetWebHook() error
	CancelWebHook() error
	SendChatAction([]string, []string) error
	//action: typing, upload_photo, record_video, upload_video, record_audio, upload_audio,
	// upload_document, find_location, record_video_note, upload_video_note
	SendMessage([]string, []string) error
	//chat_id, text, parse_mode, disable_web_page_preview, disable_notification, reply_to_message_id, reply_markup
}

type tgresp struct {
	Ok          bool    `json:"ok"`
	Description string  `json:"description"`
	Result      MsgType `json:"result"`
}

var (
	apiUrl   = make(map[string]string)
	bytePool = sync.Pool{
		New: func() interface{} { return new(bytes.Buffer) },
	}
	ErrKVnotFit = errors.New("telegram: k and v have different length.")
)

func (b *tgbot) Init() {
	prefix := "https://api.telegram.org/bot" + b.token + "/"
	apiUrl["SetWebHook"] = prefix + "setWebhook"
	apiUrl["CancelWebHook"] = prefix + "deleteWebhook"
	apiUrl["SendChatAction"] = prefix + "sendChatAction"
	apiUrl["SendMessage"] = prefix + "sendmessage"
}

func NewMultipart(api string, k []string, v []string, ftype string, filename string, data []byte) (req *http.Request, ack func()) {
	buf := (bytePool.Get()).(*bytes.Buffer)
	w := multipart.NewWriter(buf)

	ack = func() {
		buf.Truncate(0)
		bytePool.Put(buf)
	}

	for i := range k {
		_ = w.WriteField(k[i], v[i])
	}

	if ftype != "" {
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition",
			fmt.Sprintf(`form-data; name="%s"; filename="%s"`, ftype, filename))
		p, _ := w.CreatePart(h)
		_, _ = p.Write(data)
	}
	w.Close()
	req, _ = http.NewRequest("POST", api, buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	return
}

func (b *tgbot) call(fname string, k []string, v []string, filetype string, filename string, data []byte) (*tgresp, error) {
	if len(k) != len(v) {
		return nil, ErrKVnotFit
	}
	req, ack := NewMultipart(apiUrl[fname], k, v, "", "", nil)
	defer ack()

	resp, e := b.client.Do(req)
	if e != nil {
		return nil, e
	}

	if result := check(resp); result.Ok {
		return &result, nil
	} else {
		return nil, errors.New("telegram: " + result.Description)
	}
}

func (b *tgbot) simpleCall(fname string, k []string, v []string) (e error) {
	_, e = b.call(fname, k, v, "", "", nil)
	return
}

func (b *tgbot) SetWebHook() error {
	return b.simpleCall("SetWebHook", []string{"url"}, []string{b.hookSuffix + b.hookPath})
}

func (b *tgbot) CancelWebHook() (e error) {
	return b.simpleCall("CancelWebHook", nil, nil)
}

func (b *tgbot) SendChatAction(k []string, v []string) error {
	return b.simpleCall("SendChatAction", k, v)
}

func (b *tgbot) SendMessage(k []string, v []string) error {
	return b.simpleCall("SendMessage", k, v)
}

func check(resp *http.Response) (result tgresp) {
	_ = json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()
	return result
}

func toStr(n int64) string {
	return strconv.FormatInt(n, 10)
}
