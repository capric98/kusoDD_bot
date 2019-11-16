package commands

func printInfo(msg message, bot tgbot) error {
	return bot.SendText(msg, "自我介绍", false)
}
