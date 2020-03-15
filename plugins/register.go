package plugins

import (
	"time"

	"github.com/capric98/kusoDD_bot/core"
	"github.com/capric98/kusoDD_bot/plugins/getsticker"
	"github.com/capric98/kusoDD_bot/plugins/helper"
	"github.com/capric98/kusoDD_bot/plugins/saucenao"
	"github.com/capric98/kusoDD_bot/plugins/tracemoe"
)

func Register(name string, settings map[string]interface{}) (ok bool, f func(core.Message) <-chan bool, timeout time.Duration, opts []core.Opt) {
	ok = true
	switch name {
	case "helper":
		f, timeout, opts = helper.New(settings)
	case "getsticker":
		f, timeout, opts = getsticker.New(settings)
	case "tracemoe":
		f, timeout, opts = tracemoe.New(settings)
	case "saucenao":
		f, timeout, opts = saucenao.New(settings)
	default:
		ok = false
	}
	return
}
