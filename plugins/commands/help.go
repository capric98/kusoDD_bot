package commands

import (
	. "github.com/capric98/kusoDD_bot/plugins"
)

var helpText = `**kusoDD_bot v0.1.0**
by **[capric98](https://github.com/capric98)**

/info bot信息
/help 输出当前内容
/getsticker 获取一张sticker`

func printHelp(msg Message, bot Tgbot) error {
	k := []string{"chat_id", "text"}
	v := []string{msg.GetChatIDStr(), helpText}
	return bot.SendMessage(k, v)
}
