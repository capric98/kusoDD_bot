package commands

type message interface {
	GetCommands() (int, []string)
}
type tgbot interface {
	SendText(interface{}, string, bool) error
	Log(interface{}, int)
}

type complug struct {
	commands map[string](func(message, tgbot) error)
}

func (c *complug) Handle(m interface{}, b interface{}) {
	msg := m.(message)
	bot := b.(tgbot)
	clen, clist := msg.GetCommands()
	bot.Log(clist, 0)
	if clen != 0 {
		for _, v := range clist {
			if c.commands[v] != nil {
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

func NewPlugin() *complug {
	c := &complug{
		commands: make(map[string](func(message, tgbot) error)),
	}
	c.commands["/info"] = printInfo
	return c
}
