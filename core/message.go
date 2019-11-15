package core

type Message struct {
	UpdateID int64 `json:"update_id"`
	Message  struct {
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

		Photo        []MessagePhoto `json:"photo"`
		MediaGroupID string         `json:"media_group_id"`
	} `json:"message"`
}

type MessagePhoto struct {
	FileID   string `json:"file_id"`
	FileSize int64  `json:"file_size"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

type MessageFile struct {
	FileName string       `json:"file_name"`
	MimeType string       `json:"mime_type"`
	Thumb    MessagePhoto `json:"thumb"`
	FileID   string       `json:"file_id"`
	FileSize int64        `json:"file_size"`
}

func (msg *Message) GetStrMsgID() string {
	return ""
}

func (msg *Message) GetMsgLog() string {
	return msg.Message.Text
}
