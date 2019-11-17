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

type tgresp struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
	Result      struct {
		MsgType
		FileID   string `json:"file_id"`
		FileSize int64  `json:"file_size"`
		FilePath string `json:"file_path"`
	} `json:"result"`
}

var (
	prefix   = "https://api.telegram.org/bot"
	apiUrl   = make(map[string]string)
	bytePool = sync.Pool{
		New: func() interface{} { return new(bytes.Buffer) },
	}
	ErrKVnotFit = errors.New("telegram: k and v have different length.")
)

func (b *tgbot) Init() {
	prefix += b.token + "/"
	apiUrl["SetWebHook"] = prefix + "setWebhook"
	apiUrl["CancelWebHook"] = prefix + "deleteWebhook"
	apiUrl["SendChatAction"] = prefix + "sendChatAction"
	apiUrl["SendMessage"] = prefix + "sendmessage"
	apiUrl["SendDocument"] = prefix + "sendDocument"
	apiUrl["GetFile"] = prefix + "getFile"
}

func (b *tgbot) SetWebHook() error {
	return b.simpleCall("SetWebHook", []string{"url"}, []string{b.hookSuffix + b.hookPath})
}

func (b *tgbot) CancelWebHook() (e error) {
	//return b.simpleCall("CancelWebHook", nil, nil)
	r, e := b.client.Get(apiUrl["CancelWebHook"])
	if e == nil {
		result := check(r)
		if !result.Ok {
			e = errors.New("telegram: " + result.Description)
		}
	}
	return
}

func (b *tgbot) SendChatAction(k []string, v []string) error {
	return b.simpleCall("SendChatAction", k, v)
}

func (b *tgbot) SendMessage(k []string, v []string) error {
	return b.simpleCall("SendMessage", k, v)
}

func (b *tgbot) SendDocument(k []string, v []string, filename string, data []byte) (fileID string) {
	r, e := b.call("SendDocument", k, v, "document", filename, data)
	if e != nil {
		b.Log(e, 1)
	} else {
		if !r.Ok {
			b.Log(r.Description, 1)
		} else {
			fileID = r.Result.Document.FileID
		}

	}
	return
}

func (b *tgbot) SendAudio(k []string, v []string, filename string, data []byte) (fileID string) {
	r, e := b.call("SendAudio", k, v, "audio", filename, data)
	if e != nil {
		b.Log(e, 1)
	} else {
		if !r.Ok {
			b.Log(r.Description, 1)
		} else {
			fileID = r.Result.Audio.FileID
		}

	}
	return
}

func (b *tgbot) SendPhoto(k []string, v []string, filename string, data []byte) (fileID []string) {
	r, e := b.call("SendPhoto", k, v, "photo", filename, data)
	if e != nil {
		b.Log(e, 1)
	} else {
		if !r.Ok {
			b.Log(r.Description, 1)
		} else {
			for _, v := range r.Result.Photo {
				fileID = append(fileID, v.FileID)
			}
		}
	}
	return
}

func (b *tgbot) SendVideo(k []string, v []string, filename string, data []byte) (fileID string) {
	r, e := b.call("SendVideo", k, v, "video", filename, data)
	if e != nil {
		b.Log(e, 1)
	} else {
		if !r.Ok {
			b.Log(r.Description, 1)
		} else {
			fileID = r.Result.Video.FileID
		}

	}
	return
}

func (b *tgbot) SendAnimation(k []string, v []string, filename string, data []byte) (fileID string) {
	r, e := b.call("SendAnimation", k, v, "animation", filename, data)
	if e != nil {
		b.Log(e, 1)
	} else {
		if !r.Ok {
			b.Log(r.Description, 1)
		} else {
			fileID = r.Result.Animation.FileID
		}

	}
	return
}

func (b *tgbot) GetFile(k []string, v []string) (r string) {
	resp, e := b.call("GetFile", k, v, "", "", nil)
	if e == nil {
		r = "https://api.telegram.org/file/bot" + b.token + "/" + resp.Result.FilePath
	}
	//b.Log(resp, 1)
	return
}

func check(resp *http.Response) (result tgresp) {
	_ = json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()
	return result
}

func NewMultipart(api string, k []string, v []string, name string, filename string, data []byte) (req *http.Request, ack func()) {
	buf := (bytePool.Get()).(*bytes.Buffer)
	w := multipart.NewWriter(buf)

	ack = func() {
		buf.Truncate(0)
		bytePool.Put(buf)
	}

	for i := range k {
		_ = w.WriteField(k[i], v[i])
	}

	if filename != "" {
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition",
			fmt.Sprintf(`form-data; name="%s"; filename="%s"`, name, filename))
		p, _ := w.CreatePart(h)
		_, _ = p.Write(data)
	}
	w.Close()
	req, _ = http.NewRequest("POST", api, buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	return
}

func (b *tgbot) call(fname string, k []string, v []string, name string, filename string, data []byte) (*tgresp, error) {
	if len(k) != len(v) {
		return nil, ErrKVnotFit
	}
	b.Log("telegram: call "+fname, 0)
	req, ack := NewMultipart(apiUrl[fname], k, v, name, filename, data)
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

func extension(filename string) string {
	i := len(filename) - 1
	for ; i > 0 && rune(filename[i]) != '.'; i-- {
	}
	fmt.Println(filename[i:])
	return filename[i:]
}

func toStr(n int64) string {
	return strconv.FormatInt(n, 10)
}
