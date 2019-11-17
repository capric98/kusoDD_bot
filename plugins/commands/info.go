package commands

import (
	. "github.com/capric98/kusoDD_bot/plugins"
)

var (
	info = "这是一个看到新功能就想往上加的臭DD bot\n查看帮助 /help"
)

func printInfo(msg Message, bot Tgbot) error {
	return bot.SendMessage(map[string]string{
		"chat_id": msg.GetChatIDStr(),
		"text":    info,
	})
}
