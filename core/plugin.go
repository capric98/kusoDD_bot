package core

type Plugin interface {
	Handle(*Message) (bool, error)
}
