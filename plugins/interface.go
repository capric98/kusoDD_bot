package plugins

type Message interface {
	GetCommands() (int, []string)
	GetFromUserName() string
}
type Tgbot interface {
	SendText(interface{}, string, bool) error
	Log(interface{}, int)
}
