package core

import (
	"time"
)

type Opt struct{}

func (b *Bot) NewPlugin(name string, f func(Message) <-chan bool, timeout time.Duration, opts []Opt) {
	b.plugins = append(b.plugins, Plugin{
		name:    name,
		handle:  f,
		timeout: timeout,
	})
}
