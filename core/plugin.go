package core

import (
	"github.com/capric98/kusoDD_bot/plugins/commands"
)

type Plugin interface {
	Handle(interface{}, interface{})
	//     message      tgbot
	Name() string
}

func (b *tgbot) RegisterPlugins(conf settings) {
	//Default: Command Handle
	b.plugins = append(b.plugins, commands.RegisterPlugins())
}
