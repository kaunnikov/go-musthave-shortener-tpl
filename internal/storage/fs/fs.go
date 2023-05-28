package fs

import (
	"bufio"
	"encoding/json"
	"fmt"
	"kaunnikov/go-musthave-shortener-tpl/internal/config"
	"kaunnikov/go-musthave-shortener-tpl/internal/logging"
	"kaunnikov/go-musthave-shortener-tpl/internal/utils"
	"os"
)

var storage FsStorage

type StorageItem struct {
	URL      string `json:"full"`
	ShortURL string `json:"short"`
}

type FsStorage struct {
	path string
}

func Init(cfg *config.AppConfig) (*FsStorage, error) {
	storage = FsStorage{
		path: cfg.FileStoragePath,
	}
	return &storage, nil
}

func (fs *FsStorage) Save(full string) (string, error) {
	// Проверим, есть ли уже такая ссылка
	if shortURL := getShortURLFromStorage(full); shortURL != "" {
		return shortURL, nil
	}

	// Если нет - создаём запись в файле
	file, err := os.OpenFile(storage.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return "", fmt.Errorf("storage don't open to write! Error: %s. Path: %s", err, storage.path)
	}

	item := StorageItem{URL: full, ShortURL: utils.RandSeq(5)}
	data, err := json.Marshal(item)
	if err != nil {
		return "", fmt.Errorf("cannot encode storage item %s", err)
	}

	data = append(data, '\n')

	_, err = file.Write(data)

	return item.ShortURL, err
}

func (fs *FsStorage) Get(short string) (string, error) {
	file, err := os.OpenFile(storage.path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logging.Errorf("storage don't open to read! Error: %w", err)
		return "", err
	}

	r := bufio.NewReader(file)
	s, e := readLine(r)
	var item StorageItem
	for e == nil {
		err = json.Unmarshal([]byte(s), &item)
		if err != nil {
			logging.Errorf("storage don't open to read! Error: %s. Path: %s", err, storage.path)
		}

		if item.ShortURL == short {
			return item.URL, nil
		}
		s, e = readLine(r)
	}
	return "", nil
}

func (fs *FsStorage) Ping() error {
	return nil
}

func getShortURLFromStorage(fullURL string) string {
	file, err := os.OpenFile(storage.path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logging.Errorf("storage don't open to read! Error: %s. Path: %s", err, storage.path)
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
