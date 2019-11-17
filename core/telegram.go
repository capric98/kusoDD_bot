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
	apiUrl["SendAudio"] = prefix + "sendAudio"
	apiUrl["SendVideo"] = prefix + "sendVideo"
	apiUrl["SendPhoto "] = prefix + "sendPhoto"
	apiUrl["GetFile"] = prefix + "getFile"
}

func (b *tgbot) SetWebHook() error {
	return b.simpleCall("SetWebHook", map[string]string{
		"url": b.hookSuffix + b.hookPath,
	})
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

func (b *tgbot) SendChatAction(paras map[string]string) error {
	return b.simpleCall("SendChatAction", paras)
}

func (b *tgbot) SendMessage(paras map[string]string) error {
	return b.simpleCall("SendMessage", paras)
}

func (b *tgbot) SendDocument(paras map[string]string, filename string, data []byte) (fileID string) {
	_ = b.SendChatAction(map[string]string{"action": "upload_document", "chat_id": paras["chat_id"]})

	r, e := b.call("SendDocument", paras, "document", filename, data)
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

func (b *tgbot) SendAudio(paras map[string]string, filename string, data []byte) (fileID string) {
	_ = b.SendChatAction(map[string]string{"action": "upload_audio", "chat_id": paras["chat_id"]})

	r, e := b.call("SendAudio", paras, "audio", filename, data)
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

func (b *tgbot) SendPhoto(paras map[string]string, filename string, data []byte) (fileID []string) {
	_ = b.SendChatAction(map[string]string{"action": "upload_photo", "chat_id": paras["chat_id"]})

	r, e := b.call("SendPhoto", paras, "photo", filename, data)
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

func (b *tgbot) SendVideo(paras map[string]string, filename string, data []byte) (fileID string) {
	_ = b.SendChatAction(map[string]string{"action": "upload_video", "chat_id": paras["chat_id"]})

	r, e := b.call("SendVideo", paras, "video", filename, data)
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

func (b *tgbot) SendAnimation(paras map[string]string, filename string, data []byte) (fileID string) {
	_ = b.SendChatAction(map[string]string{"action": "upload_video", "chat_id": paras["chat_id"]})

	r, e := b.call("SendAnimation", paras, "animation", filename, data)
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

func (b *tgbot) GetFile(paras map[string]string) (r string) {
	resp, e := b.call("GetFile", paras, "", "", nil)
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

func NewMultipart(api string, paras map[string]string, name string, filename string, data []byte) (req *http.Request, ack func()) {
	buf := (bytePool.Get()).(*bytes.Buffer)
	w := multipart.NewWriter(buf)

	ack = func() {
		buf.Truncate(0)
		bytePool.Put(buf)
	}

	for k, v := range paras {
		_ = w.WriteField(k, v)
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

func (b *tgbot) call(fname string, paras map[string]string, name string, filename string, data []byte) (*tgresp, error) {
	b.Log("telegram: call "+fname, 0)
	req, ack := NewMultipart(apiUrl[fname], paras, name, filename, data)
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

func (b *tgbot) simpleCall(fname string, paras map[string]string) (e error) {
	_, e = b.call(fname, paras, "", "", nil)
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

func (b *tgbot) GetBotName() string {
	return b.username
}

func (b *tgbot) GetConfig(name string) map[string]interface{} {
	return b.pluginconf[name]
}
