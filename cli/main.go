package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/capric98/kusoDD_bot/core"
	"github.com/capric98/kusoDD_bot/plugins"
)

var (
	config   = flag.String("conf", "config.json", "config file")
	loglevel = flag.String("log", "normal", "debug/normal/none")
	verbose  = flag.Bool("verbose", false, "enable tgbotapi debug")
)

func main() {
	flag.Parse()

	conf, e := core.ResolvConf(*config)
	if e != nil {
		log.Fatal(e)
	}
	bot, e := core.NewBot(conf, *loglevel, *verbose)
	if e != nil {
		log.Fatal(e)
	}

	for n, s := range conf.Plugins {
		ok, f, t, opts := plugins.Register(n, s)
		if ok {
			log.Println("Found plugin:", n)
			bot.NewPlugin(n, f, t, opts)
		}
	}

	if e = bot.Run(); e != nil {
		log.Fatal(e)
	}

	c := make(chan os.Signal, 10)
	signal.Notify(c, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)
	<-c
	bot.Stop()
	log.Println("Bye~")
}
