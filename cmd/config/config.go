package config

import "flag"

type AppConfig struct {
	Host   string
	Prefix string
}

func ParseFlags() *AppConfig {
	appConfig := AppConfig{Prefix: ""}

	flag.StringVar(&appConfig.Host, "a", ":8080", "Default Host:port")
	flag.Func("b", "App prefix", func(s string) error {

		if s != "" {
			appConfig.Prefix = "/" + s
		}

		return nil
	})

	flag.Parse()
	return &appConfig
}
