package helper

import (
	"time"

	"github.com/capric98/kusoDD_bot/core"
)

var (
	ack      = make(chan bool, 1)
	info     = "这是一个看到新功能就想往上加的臭DD bot\n查看帮助 /help"
	helpText = `*kusoDD_bot v0.1.0*
by [capric98](https://github.com/capric98)

/help 输出当前内容
/getsticker 获取一张sticker
/whatpic 使用sauceNAO搜索图片
/whatanime 使用tracemoe搜索番剧`
)

func New(settings map[string]interface{}) (func(core.Message) <-chan bool, time.Duration, []core.Opt) {
	return func(msg core.Message) <-chan bool {
		// Make sure ack is clear.
		select {
		case <-ack:
		default:
			break
		}
		if msg.Message.IsCommand() {
			switch msg.Message.Command() {
			case "help":
				ack <- true
				go func() {
					resp := core.NewMessage(msg.Message.Chat.ID, helpText)
					resp.ReplyToMessageID = msg.Message.MessageID
					resp.ParseMode = "Markdown"
					resp.DisableWebPagePreview = true
					_, _ = msg.Bot.Send(resp)
				}()
			case "info":
				ack <- true
				go func() {
					resp := core.NewMessage(msg.Message.Chat.ID, info)
					resp.ReplyToMessageID = msg.Message.MessageID
					_, _ = msg.Bot.Send(resp)
				}()
			default:
				ack <- false
			}
		} else {
			ack <- false
		}
		return ack
	}, 10 * time.Millisecond, nil
}
