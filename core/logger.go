package core

import "log"

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
