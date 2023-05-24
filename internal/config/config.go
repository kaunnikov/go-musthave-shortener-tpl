package config

import (
	"flag"
	"os"
	"strings"
)

type AppConfig struct {
	Host            string
	ResultURL       string
	FileStoragePath string
	DatabaseDSN     string
}

func LoadConfig() *AppConfig {
	cfg := &AppConfig{}
	loadFromArgs(cfg)
	loadFromENV(cfg)
	return cfg
}

func loadFromArgs(cfg *AppConfig) {
	flag.StringVar(&cfg.Host, "a", "localhost:8080", "Default Host:port")
	flag.StringVar(&cfg.ResultURL, "b", "http://localhost:8080", "Default result URL")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/short-url-db-kaunnikov.json", "Default File Storage Path")
	flag.StringVar(&cfg.DatabaseDSN, "d", "host=localhost port=5432 user=postgres password=root dbname=postgres sslmode=disable", "Default Database DSN")
	flag.Parse()
}

func loadFromENV(cfg *AppConfig) {
	envRunAddr := strings.TrimSpace(os.Getenv("SERVER_ADDRESS"))
	if envRunAddr != "" {
		cfg.Host = envRunAddr
	}

	envBaseURL := strings.TrimSpace(os.Getenv("BASE_URL"))
	if envBaseURL != "" {
		cfg.ResultURL = envBaseURL
	}

	fileStorageFile := strings.TrimSpace(os.Getenv("FILE_STORAGE_PATH"))
	if fileStorageFile != "" {
		cfg.FileStoragePath = fileStorageFile
	}

	databaseDSN := strings.TrimSpace(os.Getenv("DATABASE_DSN"))
	if databaseDSN != "" {
		cfg.DatabaseDSN = databaseDSN
	}
}
