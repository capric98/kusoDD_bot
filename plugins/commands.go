package plugins

import (
	"net/http"
	"time"

	"github.com/capric98/kusoDD_bot/plugins/getsticker"
	"github.com/capric98/kusoDD_bot/plugins/saucenao"
	"github.com/capric98/kusoDD_bot/plugins/tracemoe"
)

type complug struct {
	commands map[string](func(interface{}, interface{}) error)
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

func Register(b interface{}) *complug {
	c := &complug{
		commands: make(map[string](func(interface{}, interface{}) error)),
	}
	c.commands["/info"] = printInfo
	c.commands["/help"] = printHelp
	c.commands["/getsticker"] = getsticker.Handle
	c.commands["/whatpic"] = saucenao.New(b)
	c.commands["/whatanime"] = tracemoe.New(b)

	// Add /command@kusoDD_bot
	comlist := []string{}
	for k, _ := range c.commands {
		comlist = append(comlist, k)
	}
	for _, v := range comlist {
		c.commands[v+"@"+b.(Tgbot).GetBotName()] = c.commands[v]
	}
	return c
}
