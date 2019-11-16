package plugins

type Message interface {
	GetCommands() (int, []string)
	GetFromUserName() string
	GetChatIDStr() string
}
type Tgbot interface {
	SetWebHook() error
	CancelWebHook() error
	SendChatAction([]string, []string) error
	SendText([]string, []string) error
	Log(interface{}, int)
}
