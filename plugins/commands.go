package plugins

import (
	"net/http"
	"time"

	"github.com/capric98/plugins/getsticker"
)

type complug struct {
	commands map[string](func(Message, Tgbot) error)
}

var (
	client = &http.Client{
		Timeout: 60 * time.Second,
	}
)

func (c *complug) Handle(m interface{}, b interface{}) {
	msg := m.(Message)
	bot := b.(Tgbot)
	clen, clist := msg.GetCommands()
	if clen != 0 {
		for _, v := range clist {
			if c.commands[v] != nil {
				bot.Log("commands: Handled "+msg.GetFromUserName()+"'s command \""+v+"\"", 0)
				if e := c.commands[v](msg, bot); e != nil {
					bot.Log(e, 1)
				}
			}
		}
	}
	return
}

func (c *complug) Name() string {
	return "commands"
}

func RegisterPlugins() *complug {
	c := &complug{
		commands: make(map[string](func(Message, Tgbot) error)),
	}
	c.commands["/info"] = printInfo
	c.commands["/help"] = printHelp
	c.commands["/getsticker"] = getsticker.Handle()
	return c
}
