package core

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type tgbot struct {
	client     *http.Client
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
	} else {
		nb.Log(err, 1)
	}

	return nil
}

func (bot *tgbot) Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != bot.hookSuffix {
			http.Error(w, "Bad request.", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case "POST":
		default:
			http.Error(w, "Only support POST method.", http.StatusBadRequest)
		}
	})
	// stuck here
}

func (bot *tgbot) Log(body interface{}, level int) {
	if level > bot.loglevel {
		log.Println(body)
	}
}
