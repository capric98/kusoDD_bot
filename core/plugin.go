package core

import "net/http"

type Plugin interface {
	Handle(*Message, http.ResponseWriter) (bool, error)
	Name() string
}
