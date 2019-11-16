package commands

import (
	. "github.com/capric98/kusoDD_bot/plugins"
)

func printInfo(msg Message, bot Tgbot) error {
	return bot.SendText(msg, "这是一个看到新功能就想往上加的臭DD bot\n查看帮助 /help", false)
}
