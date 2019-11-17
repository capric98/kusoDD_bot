package saucenao

import "errors"

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
}

var (
	ErrNoPhoto = errors.New("saucenao: No photo in this message.")
)

func New(b interface{}) func(interface{}, interface{}) error {
	conf := b.(tgbot).GetConfig("sauceNAO")
	for k, v := range conf {
		b.(tgbot).Log("sauceNAO: "+k+v.(string), 0)
	}
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
	u = u + " "
	return nil
}
