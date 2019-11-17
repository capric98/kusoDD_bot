package plugins

type Message interface {
	GetCommands() (int, []string)
	GetFromUserName() string
	GetChatIDStr() string
	GetPhotoFileID() string
	GetStickerFileID() string
	GetReplyToStickerFileID() string
	GetReplyToStickerFileName() string
	GetReplyMsgIDStr() string
	GetReplyToStickerSetName() string
	GetReplyToStickerEmoji() string
	GetReplyToPhotoFileID() string
}
type Tgbot interface {
	SetWebHook() error
	CancelWebHook() error
	SendChatAction(paras map[string]string) error
	//action: typing, upload_photo, record_video, upload_video, record_audio, upload_audio,
	// upload_document, find_location, record_video_note, upload_video_note
	SendMessage(paras map[string]string) error
	//chat_id, text, parse_mode, disable_web_page_preview, disable_notification, reply_to_message_id, reply_markup
	GetFile(paras map[string]string) string
	//file_id
	SendDocument(paras map[string]string, filename string, data []byte) (fileID string)
	SendAudio(paras map[string]string, filename string, data []byte) (fileID string)
	SendVideo(paras map[string]string, filename string, data []byte) (fileID string)
	SendAnimation(paras map[string]string, filename string, data []byte) (fileID string)
	SendPhoto(paras map[string]string, filename string, data []byte) (fileID []string)

	GetConfig(name string) map[string]interface{}
	Log(interface{}, int)
}
