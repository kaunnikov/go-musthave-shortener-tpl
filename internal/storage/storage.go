package storage

import (
	"kaunnikov/go-musthave-shortener-tpl/internal/config"
	"kaunnikov/go-musthave-shortener-tpl/internal/logging"
	"kaunnikov/go-musthave-shortener-tpl/internal/storage/db"
	"kaunnikov/go-musthave-shortener-tpl/internal/storage/fs"
	"kaunnikov/go-musthave-shortener-tpl/internal/storage/mem"
)

type Storage interface {
	Save(full string) (string, error)
	Get(short string) (string, error)
	Ping() error
}

var defaultStorage Storage

func Init(cfg *config.AppConfig) {
	var err error
	if cfg.DatabaseDSN != "" {
		// "host=localhost port=5433 user=postgres password=password dbname=postgres sslmode=disable"
		defaultStorage, err = db.Init(cfg)
		if err != nil {
			logging.Fatalf("DB don't init: %w", err)
		}
	} else if cfg.FileStoragePath != "" {
		defaultStorage, err = fs.Init(cfg)
		if err != nil {
			logging.Fatalf("file storage don't init: %w", err)
		}
	} else {
		defaultStorage, err = mem.Init()
		if err != nil {
			logging.Fatalf("memory storage don't init: %w", err)
		}
	}
}

func SaveURLInStorage(full string) (string, error) {
	short, err := defaultStorage.Save(full)
	if err != nil {
		logging.Errorf("Don't save full URL: %w", err)
		return "", err
	}
	return short, nil
}
func GetFullURL(short string) (string, error) {
	full, err := defaultStorage.Get(short)
	if err != nil {
		return "", err
	}
	return full, nil
}
func Ping() error {
	return defaultStorage.Ping()
}
