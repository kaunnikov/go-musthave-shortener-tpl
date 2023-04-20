package app

import (
	"kaunnikov/go-musthave-shortener-tpl/config"
	"sync"
)

type app struct {
	cfg *config.AppConfig
}

type jsonStruct struct {
	URL string `json:"URL"`
}

type shortenResponse struct {
	Result string `json:"result"`
}

var URLMap = make(map[string]string, 1000)
var URLMapSync = sync.Mutex{}

func NewApp(cfg *config.AppConfig) *app {
	return &app{cfg: cfg}
}
