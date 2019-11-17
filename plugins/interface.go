package plugins

type Message interface {
	GetCommands() (int, []string)
	GetFromUserName() string
	GetChatIDStr() string
	GetPhotoFileID() string
	GetStickerFileID() string
}
type Tgbot interface {
	SetWebHook() error
	CancelWebHook() error
	SendChatAction([]string, []string) error
	//action: typing, upload_photo, record_video, upload_video, record_audio, upload_audio,
	// upload_document, find_location, record_video_note, upload_video_note
	SendMessage([]string, []string) error
	//chat_id, text, parse_mode, disable_web_page_preview, disable_notification, reply_to_message_id, reply_markup
	GetFile(k []string, v []string) error
	//file_id
	Log(interface{}, int)
}
