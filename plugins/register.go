package plugins

import (
	"time"

	"github.com/capric98/kusoDD_bot/core"
	"github.com/capric98/kusoDD_bot/plugins/getsticker"
	"github.com/capric98/kusoDD_bot/plugins/helper"
)

func Register(name string, settings map[string]interface{}) (ok bool, f func(core.Message) <-chan bool, timeout time.Duration, opts []core.Opt) {
	ok = true
	switch name {
	case "helper":
		f, timeout, opts = helper.New(settings)
	case "getsticker":
		f, timeout, opts = getsticker.New(settings)
	// case "tracemoe":
	// case "saucenao":
	default:
		ok = false
	}
	return
}
