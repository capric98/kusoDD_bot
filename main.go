package main

import (
	"flag"

	"github.com/capric98/kusoDD_bot/core"
)

var (
	config   = flag.String("config", "config.json", "config file path")
	loglevel = flag.String("log", "normal", "debug/normal/none")
)

func main() {
	flag.Parse()

	//go func() {
	//	log.Println(http.ListenAndServe("localhost:6060", nil))
	//}()

	bot := core.Newbot(config, loglevel)
	if bot != nil {
		bot.Run()
	}
}
