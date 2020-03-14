package core

import (
	"errors"
	"os"

	jsoniter "github.com/json-iterator/go"
)

type Config struct {
	Token   string                            `json:"token"`
	Server  Server                            `json:"server"`
	Plugins map[string]map[string]interface{} `json:"plugins"`
}

type Server struct {
	Host string `json:"host"`
	Path string `json:"path"`
	Port int16  `json:"port"`
	TLS  *TLS   `json:"tls"`
}

type TLS struct {
	Cert string `json:"cert"`
	Key  string `json:"Key"`
}

func ResolvConf(conf string) (c *Config, e error) {
	if _, e = os.Stat(conf); os.IsNotExist(e) {
		return nil, errors.New("Config file does not exist!")
	}
	freader, e := os.Open(conf)
	if e != nil {
		return
	}
	defer freader.Close()
	decoder := jsoniter.ConfigCompatibleWithStandardLibrary.NewDecoder(freader)

	c = new(Config)
	e = decoder.Decode(c)
	return
}
