package core

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	jsoniter "github.com/json-iterator/go"
)

type Bot struct {
	Conf *Config
	Json jsoniter.API

	*tgbotapi.BotAPI
	msg     chan tgbotapi.Update
	plugins []Plugin

	srv    *http.Server
	level  int
	ctx    context.Context
	cancel func()
}

type Plugin struct {
	name    string
	timeout time.Duration
	handle  func(Message) <-chan bool
}

type Message struct {
	tgbotapi.Update
	Bot *Bot
}

var (
	LV = make(map[string]int)
)

func NewBot(conf *Config, mode string, debug bool) (b *Bot, e error) {
	LV["debug"] = -1
	LV["normal"] = 50
	LV["warn"] = 60
	LV["none"] = 100
	b = &Bot{
		Conf:    conf,
		Json:    jsoniter.ConfigCompatibleWithStandardLibrary,
		plugins: make([]Plugin, 0, 1),
		level:   LV[mode],
	}
	b.BotAPI, e = tgbotapi.NewBotAPI(conf.Token)
	if e != nil {
		return
	}
	b.ctx, b.cancel = context.WithCancel(context.Background())
	b.msg = make(chan tgbotapi.Update, 10)
	b.Debug = debug

	_, _ = b.Request(tgbotapi.RemoveWebhookConfig{})

	http.HandleFunc(b.Conf.Server.Path, func(w http.ResponseWriter, r *http.Request) {
		var update tgbotapi.Update
		_ = b.Json.NewDecoder(r.Body).Decode(&update)
		defer r.Body.Close()

		b.msg <- update
	})

	return
}

func (b *Bot) Run() (e error) {
	if b.Conf.Server.TLS == nil {
		b.srv = &http.Server{
			Addr: "127.0.0.1:" + strconv.Itoa(int(b.Conf.Server.Port)),
		}
		go func() { _ = b.srv.ListenAndServe() }()
	} else {
		// https://blog.bracebin.com/achieving-perfect-ssl-labs-score-with-go
		cfg := &tls.Config{
			MinVersion:       tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			// PreferServerCipherSuites: true,
			// CipherSuites: []uint16{
			// 	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			// 	tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			// 	tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			// 	tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			// },
		}
		b.srv = &http.Server{
			Addr:         ":" + strconv.Itoa(int(b.Conf.Server.Port)),
			TLSConfig:    cfg,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		}
		go func() { _ = b.srv.ListenAndServeTLS(b.Conf.Server.TLS.Cert, b.Conf.Server.TLS.Key) }()
	}
	_, e = b.Request(tgbotapi.NewWebhook(b.Conf.Server.Host + b.Conf.Server.Path))
	go b.handleUpdate()

	return
}

func (b *Bot) handleUpdate() {
	for u := range b.msg {
		b.printUpdate(u)
		go func(update tgbotapi.Update) {
			defer func() {
				if e := recover(); e != nil {
					log.Println(" core - Unexpected panic:", e)
				}
			}()
			var msg Message
			msg.Bot = b
			msg.Update = update
			for _, p := range b.plugins {
				select {
				case <-time.After(p.timeout):
					b.Printf("%6s - Plugin \"%s\" handled update timeout.\n", "info", p.name)
				case do := <-p.handle(msg):
					if do {
						b.Printf("%6s - Plugin \"%s\" handled update %d.\n", "debug", p.name, msg.UpdateID)
					}
				}
			}
		}(u)
	}
}

func (b *Bot) WithCancel() (context.Context, func()) {
	return context.WithCancel(b.ctx)
}

func (b *Bot) WithTimeout(timeout time.Duration) (context.Context, func()) {
	return context.WithTimeout(b.ctx, timeout)
}

func (b *Bot) Done() <-chan struct{} {
	return b.ctx.Done()
}

func (b *Bot) Stop() {
	_, _ = b.Request(tgbotapi.RemoveWebhookConfig{})
	_ = b.srv.Shutdown(context.Background())
	close(b.msg)
	b.cancel()
}
