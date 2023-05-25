package storage

import (
	"kaunnikov/go-musthave-shortener-tpl/internal/storage/db"
	"kaunnikov/go-musthave-shortener-tpl/internal/storage/fs"
)

func SaveURLInStorage(full string, short string) (string, error) {
	// Если подключена БД - Сохраняем в БД
	if err := db.Ping(); err == nil {
		shortURL, err := db.GetOrSave(full, short)
		if err != nil {
			return "", err
		}
		return shortURL, nil
	}

	return fs.SaveURLInFileStorage(full)
}

func GetFullURL(short string) (string, error) {
	// Если подключена БД - Сохраняем в БД
	if err := db.Ping(); err == nil {
		fullURL, err := db.GetFullByShortURL(short)
		if err != nil {
			return "", err
		}
		return fullURL, nil
	}

	return fs.GetFullURLFromStorage(short), nil
}
