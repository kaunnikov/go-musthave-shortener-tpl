package config

import "flag"

type AppConfig struct {
	Host      string
	ResultURL string
}

func ParseFlags() *AppConfig {
	appConfig := new(AppConfig)

	flag.StringVar(&appConfig.Host, "a", "localhost:8080", "Default Host:port")
	flag.StringVar(&appConfig.ResultURL, "b", "http://localhost:8080", "Default result URL")
	flag.Parse()
	return appConfig
}
