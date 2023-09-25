package storage

import (
	"kaunnikov/go-musthave-shortener-tpl/internal/config"
	"kaunnikov/go-musthave-shortener-tpl/internal/logging"
	"kaunnikov/go-musthave-shortener-tpl/internal/storage/db"
	"kaunnikov/go-musthave-shortener-tpl/internal/storage/fs"
	"kaunnikov/go-musthave-shortener-tpl/internal/storage/mem"
)

type Storage interface {
	Save(token string, full string) (string, error)
	Get(short string) (string, error)
	GetUrlsByUser(token string) ([]db.UrlsByUserResponseMessage, error)
	DeleteURLs(string, token string) error
	Ping() error
}

var defaultStorage Storage

func Init(cfg *config.AppConfig) {
	var err error
	if cfg.DatabaseDSN != "" {
		// "host=localhost port=5433 user=postgres password=password dbname=postgres sslmode=disable"
		defaultStorage, err = db.Init(cfg)
		if err != nil {
			logging.Fatalf("DB don't init: %s", err)
		}
	} else if cfg.FileStoragePath != "" {
		defaultStorage, err = fs.Init(cfg)
		if err != nil {
			logging.Fatalf("file storage don't init: %s", err)
		}
	} else {
		defaultStorage, err = mem.Init(cfg)
		if err != nil {
			logging.Fatalf("memory storage don't init: %s", err)
		}
	}
}

func SaveURLInStorage(token string, full string) (string, error) {
	return defaultStorage.Save(token, full)
}
func GetFullURL(short string) (string, error) {
	return defaultStorage.Get(short)
}

func GetURLsByUser(token string) ([]db.UrlsByUserResponseMessage, error) {
	return defaultStorage.GetUrlsByUser(token)
}

func DeleteURLs(URL string, token string) error {
	return defaultStorage.DeleteURLs(URL, token)
}

func Ping() error {
	return defaultStorage.Ping()
}
