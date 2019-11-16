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
	SendText([]string, []string) error
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
	apiUrl["SetWebHook"] = prefix + "setWebhook?url=" + b.hookSuffix + b.hookPath
	apiUrl["CancelWebHook"] = prefix + "deleteWebhook"
	apiUrl["SendChatAction"] = prefix + "sendChatAction"
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
		w.Close()
		req, _ = http.NewRequest("POST", api, buf)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	return
}

func check(resp *http.Response) (result tgresp) {
	_ = json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()
	return result
}

func (b *tgbot) SetWebHook() error {
	resp, err := b.client.Get(apiUrl["SetWebHook"])
	b.Log("Set webhook.", 0)
	if err != nil {
		return err
	}

	if result := check(resp); result.Ok {
		return nil
	} else {
		return errors.New("telegram: " + result.Description)
	}
}

func (b *tgbot) CancelWebHook() error {
	resp, err := b.client.Get(apiUrl["CancelWebHook"])
	b.Log("Delete webhook.", 0)
	if err != nil {
		return err
	}

	if result := check(resp); result.Ok {
		return nil
	} else {
		return errors.New("telegram: " + result.Description)
	}
}

func (b *tgbot) SendChatAction(k []string, v []string) error {
	if len(k) != len(v) {
		return ErrKVnotFit
	}
	req, ack := NewMultipart(apiUrl["SendChatAction"], k, v, "", "", nil)
	defer ack()

	resp, e := b.client.Do(req)
	if e != nil {
		return e
	}

	if result := check(resp); result.Ok {
		return nil
	} else {
		return errors.New("telegram: " + result.Description)
	}
}

func (b *tgbot) SendText(k []string, v []string) error {
	if len(k) != len(v) {
		return ErrKVnotFit
	}
	req, ack := NewMultipart(apiUrl["SendText"], k, v, "", "", nil)
	defer ack()

	resp, e := b.client.Do(req)
	if e != nil {
		return e
	}

	if result := check(resp); result.Ok {
		return nil
	} else {
		return errors.New("telegram: " + result.Description)
	}
}

func toStr(n int64) string {
	return strconv.FormatInt(n, 10)
}
