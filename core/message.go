package core

import "strconv"

type MessageIf interface {
	GetCommands() (int, []string)
	GetFromUserName() string
	GetChatIDStr() string
}

type Message struct {
	UpdateID int64   `json:"update_id"`
	Message  MsgType `json:"message"`
}

type MsgType struct {
	MessageID int64  `json:"message_id"`
	Text      string `json:"text"`
	Caption   string `json:"Caption"`
	From      struct {
		ID        int64  `json:"id"`
		IsBot     bool   `json:"is_bot"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		UserName  string `json:"username"`
		LangCode  string `json:"language_code"`
	} `json:"from"`
	Chat struct {
		ID        int64  `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		UserName  string `json:"username"`
		Type      string `json:"type"`
	} `json:"chat"`
	ForwardFrom struct {
		ID        int64  `json:"id"`
		IsBot     bool   `json:"is_bot"`
		FirstName string `json:"first_name"`
		UserName  string `json:"username"`
	} `json:"forward_from"`
	ForwardFromChat struct {
		ID       int64  `json:"id"`
		Title    string `json:"title"`
		UserName string `json:"username"`
		Type     string `json:"channel"`
	} `json:"forward_from_chat"`

	Date          int64  `json:"date"`
	ForwardSender string `json:"forward_sender_name"`
	ForwardMsgID  int64  `json:"forward_from_message_id"`
	ForwardDate   int64  `json:"forward_date"`

	Photo           []MsgPhoto  `json:"photo"`
	MediaGroupID    string      `json:"media_group_id"`
	Audio           MsgAudio    `json:"audio"`
	Video           MsgVideo    `json:"video"`
	Sticker         MsgSticker  `json:"sticker"`
	Animation       MsgAnime    `json:"animation"`
	Entities        []MsgEntity `json:"entities"`
	CaptionEntities []MsgEntity `json:"caption_entities"`
}

type MsgPhoto struct {
	FileID   string `json:"file_id"`
	FileSize int64  `json:"file_size"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

type MsgFile struct {
	FileName string   `json:"file_name"`
	MimeType string   `json:"mime_type"`
	Thumb    MsgPhoto `json:"thumb"`
	FileID   string   `json:"file_id"`
	FileSize int64    `json:"file_size"`
}

type MsgEntity struct {
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	Type   string `json:"type"`
}

type MsgAudio struct {
	Duration  int    `json:"duration"`
	Title     string `json:"title"`
	Performer string `json:"performer"`
	MsgFile
}

type MsgVideo struct {
	Duration int `json:"duration"`
	Width    int `json:"width"`
	Height   int `json:"height"`
	MsgFile
}

type MsgAnime struct {
	Duration int `json:"duration"`
	Width    int `json:"width"`
	Height   int `json:"height"`
	MsgFile
}

type MsgSticker struct {
	MsgFile
	Emoji      string `json:"emoji"`
	SetName    string `json:"set_name"`
	IsAnimated bool   `json:"is_animated"`
}

func (msg *Message) GetStrMsgID() string {
	return ""
}

// func (msg *Message) GetMsgLog() (result string) {
// 	defer func() { recover() }()

// 	result = msg.Message.From.UserName
// 	switch {
// 	case msg.Message.ForwardDate != 0:
// 		result += " forwards: "
// 		if msg.Message.Text != "" {
// 			result += msg.Message.Text
// 		} else {
// 			result += "something rather than plain text."
// 		}
// 	case msg.Message.Text != "":
// 		result += " says: " + msg.Message.Text
// 	case msg.Message.Sticker.FileID != "":
// 		result += " sends a sticker: " + msg.Message.Sticker.FileID
// 	}

// 	return
// }

func (msg *Message) GetCommands() (int, []string) {
	if l := len(msg.Message.Entities) + len(msg.Message.CaptionEntities); l != 0 {
		commands := []string{}
		for _, v := range msg.Message.Entities {
			if v.Type == "bot_command" {
				commands = append(commands, msg.Message.Text[v.Offset:v.Offset+v.Length])
			}
		}
		for _, v := range msg.Message.CaptionEntities {
			if v.Type == "bot_command" {
				commands = append(commands, msg.Message.Caption[v.Offset:v.Offset+v.Length])
			}
		}
		return len(commands), commands
	} else {
		return 0, nil
	}
}

func (msg *Message) GetChatIDStr() string {
	return strconv.FormatInt(msg.Message.Chat.ID, 10)
}
