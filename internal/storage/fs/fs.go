package fs

import (
	"bufio"
	"encoding/json"
	"fmt"
	"kaunnikov/go-musthave-shortener-tpl/internal/config"
	"kaunnikov/go-musthave-shortener-tpl/internal/logging"
	"kaunnikov/go-musthave-shortener-tpl/internal/storage/mem"
	"kaunnikov/go-musthave-shortener-tpl/internal/utils"
	"os"
)

type StorageItem struct {
	URL      string `json:"full"`
	ShortURL string `json:"short"`
}

var conf *config.AppConfig

func Init(cfg *config.AppConfig) {
	conf = cfg
}

func SaveURLInStorage(full string) (string, error) {
	// Проверим, есть ли уже такая ссылка
	if shortURL := getShortURLFromStorage(full); shortURL != "" {
		return shortURL, nil
	}
	file, err := os.OpenFile(conf.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return "", fmt.Errorf("storage don't open to write! Error: %s. Path: %s", err, conf.FileStoragePath)
	}

	item := StorageItem{URL: full, ShortURL: utils.RandSeq(5)}
	data, err := json.Marshal(item)
	if err != nil {
		return "", fmt.Errorf("cannot encode storage item %s", err)
	}

	data = append(data, '\n')

	_, err = file.Write(data)

	// Запишем в кеш и отдадим результат
	mem.Append(item.URL, item.ShortURL)
	return item.ShortURL, err
}

func GetFullURLFromStorage(shortURL string) string {
	//Проверяем запись в кеше, если есть - отдаём
	if fullURL := mem.GetByShort(shortURL); fullURL != "" {
		return fullURL
	}

	file, err := os.OpenFile(conf.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logging.Errorf("storage don't open to read! Error: %s. Path: %s", err, conf.FileStoragePath)
	}

	r := bufio.NewReader(file)
	s, e := readLine(r)
	var item StorageItem
	for e == nil {
		err = json.Unmarshal([]byte(s), &item)
		if err != nil {
			logging.Errorf("storage don't open to read! Error: %s. Path: %s", err, conf.FileStoragePath)
		}

		if item.ShortURL == shortURL {
			return item.URL
		}
		s, e = readLine(r)
	}
	return ""
}
func getShortURLFromStorage(fullURL string) string {
	// Проверяем запись в кеше, если есть - отдаём
	if short := mem.GetByFull(fullURL); short != "" {
		return short
	}

	file, err := os.OpenFile(conf.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logging.Errorf("storage don't open to read! Error: %s. Path: %s", err, conf.FileStoragePath)
	}

	r := bufio.NewReader(file)
	s, e := readLine(r)
	var item StorageItem
	for e == nil {
		err = json.Unmarshal([]byte(s), &item)
		if err != nil {
			logging.Errorf("Cannot decode urls: %s", err)
		}

		if item.URL == fullURL {
			mem.Append(item.URL, item.ShortURL)
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
