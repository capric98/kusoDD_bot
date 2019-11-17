package plugins

var (
	info     = "这是一个看到新功能就想往上加的臭DD bot\n查看帮助 /help"
	helpText = `*kusoDD_bot v0.1.0*
by [capric98](https://github.com/capric98)

/help 输出当前内容
/getsticker 获取一张sticker
/whatpic 使用sauceNAO搜索图片`
)

func printHelp(m interface{}, b interface{}) error {
	msg := m.(Message)
	bot := b.(Tgbot)
	return bot.SendMessage(map[string]string{
		"chat_id":    msg.GetChatIDStr(),
		"parse_mode": "Markdown",
		"text":       helpText,

		"disable_web_page_preview": "true",
	})
}

func printInfo(m interface{}, b interface{}) error {
	msg := m.(Message)
	bot := b.(Tgbot)
	return bot.SendMessage(map[string]string{
		"chat_id": msg.GetChatIDStr(),
		"text":    info,
	})
}
