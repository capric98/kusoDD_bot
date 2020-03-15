package core

import (
	"bytes"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) Println(v ...interface{}) {
	log.Println(v...)
}

func (b *Bot) Printf(format string, v ...interface{}) {
	if lv, ok := v[0].(string); ok {
		if LV[lv] < b.level {
			return
		}
	}
	log.Printf(format, v...)
}

func (b *Bot) printUpdate(u tgbotapi.Update) {
	var buf bytes.Buffer
	if b.level == -1 {
		write(&buf, " debug - \"", u.Message.From.String())
		write(&buf, "\" SAYS \"", u.Message.Text, "\"")
		if u.Message.Chat != nil {
			write(&buf, " IN chatID=", strconv.FormatInt(u.Message.Chat.ID, 10))
		}
		log.Println(buf.String())
	}
}

func write(b *bytes.Buffer, s ...string) {
	for k := range s {
		_, _ = b.Write([]byte(s[k]))
	}
}
