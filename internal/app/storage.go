package app

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"
)

type StorageItem struct {
	URL      string `json:"full"`
	ShortURL string `json:"short"`
}

var URLStorageSync = sync.Mutex{}

func (m *app) SaveURLInStorage(item *StorageItem) (string, error) {
	// Проверим, есть ли уже такая ссылка
	if shortURL := m.getShortURLFromStorage(item.URL); shortURL != "" {
		return shortURL, nil
	}
	file, err := os.OpenFile(m.cfg.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Sugar.Errorf("storage don't open to write! Error: %s. Path: %s", err, m.cfg.FileStoragePath)
	}

	data, err := json.Marshal(item)
	if err != nil {
		return "", err
	}

	data = append(data, '\n')
	_, err = file.Write(data)
	return item.ShortURL, err
}

func (m *app) GetFullURLFromStorage(shortURL string) string {
	file, err := os.OpenFile(m.cfg.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		Sugar.Errorf("storage don't open to read! Error: %s. Path: %s", err, m.cfg.FileStoragePath)
	}

	r := bufio.NewReader(file)
	s, e := readLine(r)
	var item StorageItem
	for e == nil {
		err = json.Unmarshal([]byte(s), &item)
		if err != nil {
			Sugar.Errorf("storage don't open to read! Error: %s. Path: %s", err, m.cfg.FileStoragePath)
		}

		if item.ShortURL == shortURL {
			return item.URL
		}
		s, e = readLine(r)
	}
	return ""
}
func (m *app) getShortURLFromStorage(fullURL string) string {
	file, err := os.OpenFile(m.cfg.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		Sugar.Errorf("storage don't open to read! Error: %s. Path: %s", err, m.cfg.FileStoragePath)
	}

	r := bufio.NewReader(file)
	s, e := readLine(r)
	var item StorageItem
	for e == nil {
		err = json.Unmarshal([]byte(s), &item)
		if err != nil {
			Sugar.Errorf("storage don't open to read! Error: %s. Path: %s", err, m.cfg.FileStoragePath)
		}

		if item.URL == fullURL {
			return item.ShortURL
		}
		s, e = readLine(r)
	}
	return ""
}

func readLine(r *bufio.Reader) (string, error) {
	var (
		isPrefix       = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}
