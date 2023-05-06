package app

import (
	"kaunnikov/go-musthave-shortener-tpl/config"
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

func NewApp(cfg *config.AppConfig) *app {
	return &app{cfg: cfg}
}
