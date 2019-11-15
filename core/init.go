package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type tgbot struct {
	client     *http.Client
	apiUrl     string
	hookSuffix string
	hookPath   string
	port       int
	loglevel   int // -1:debug 0:normal 100:none
}

type settings struct {
	Token      string `json:"token"`
	HookSuffix string `json:"hookSuffix"`
	HookPath   string `json:"hookPath"`
	Port       int    `json:"port"`
}

func Newbot(conf *string, loglevel *string) *tgbot {
	nb := &tgbot{
		loglevel: 0,
	}
	switch *loglevel {
	case "debug":
		nb.loglevel = -1
	case "none":
		nb.loglevel = 100
	}

	if _, err := os.Stat(*conf); os.IsNotExist(err) {
		nb.Log("core: "+*conf+" file does not exist!", 1)
		return nil
	}
	var config settings
	if f, err := os.Open(*conf); err == nil {
		j := json.NewDecoder(f)
		if e := j.Decode(&config); e != nil {
			nb.Log(e, 1)
			return nil
		}

		// Validation
		resp, e := http.Get("https://api.telegram.org/bot" + config.Token + "/getMe")
		if e != nil {
			nb.Log(e, 1)
			return nil
		}
		var obj interface{}
		_ = json.NewDecoder(resp.Body).Decode(&obj)
		resp.Body.Close()

		if omap := obj.(map[string]interface{}); omap["ok"].(bool) {
			if result := omap["result"].(map[string]interface{}); result["is_bot"].(bool) {
				nb.Log("Validated: "+result["first_name"].(string), 1)
			} else {
				nb.Log("It is not a bot!", 1)
				return nil
			}
		} else {
			nb.Log(omap, 1)
			return nil
		}

		// Validate success, make bot.
		nb.client = &http.Client{}
		nb.apiUrl = "https://api.telegram.org/bot" + config.Token + "/"
		nb.port = config.Port
		nb.hookSuffix = config.HookSuffix
		nb.hookPath = config.HookPath
		if nb.port < 1000 {
			nb.Log("Port low than 1000.", 1)
			return nil
		}
		return nb
	} else {
		nb.Log(err, 1)
	}

	return nil
}

func (bot *tgbot) Run() {
	bot.CancelWebHook()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != bot.hookPath {
			http.Error(w, "Bad request.", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case "POST":
			body, _ := ioutil.ReadAll(r.Body)
			fmt.Println("")
			fmt.Println(string(body))
		default:
			http.Error(w, "Only support POST method.", http.StatusBadRequest)
		}
	})

	srv := &http.Server{
		Addr:    "127.0.0.1:" + strconv.Itoa(bot.port),
		Handler: mux,
		//TLSConfig:    cfg,
		//TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	//log.Fatal(srv.ListenAndServeTLS("tls.crt", "tls.key"))

	if e := bot.SetWebHook(); e != nil {
		bot.Log(e, 1)
		return
	}

	srv.ListenAndServe()

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	bot.CancelWebHook()
}

func (bot *tgbot) Log(body interface{}, level int) {
	if level > bot.loglevel {
		log.Println(body)
	}
}
