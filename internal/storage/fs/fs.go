package fs

import (
	"bufio"
	"encoding/json"
	"fmt"
	"kaunnikov/go-musthave-shortener-tpl/internal/config"
	"kaunnikov/go-musthave-shortener-tpl/internal/logging"
	"kaunnikov/go-musthave-shortener-tpl/internal/storage/db"
	"kaunnikov/go-musthave-shortener-tpl/internal/utils"
	"os"
)

var storage FileStorage

type StorageItem struct {
	Token    string `json:"token"`
	URL      string `json:"full"`
	ShortURL string `json:"short"`
}

type FileStorage struct {
	path      string
	resultURL string
}

func Init(cfg *config.AppConfig) (*FileStorage, error) {
	storage = FileStorage{
		path:      cfg.FileStoragePath,
		resultURL: cfg.ResultURL,
	}
	return &storage, nil
}

func (fs *FileStorage) Save(token string, full string) (string, error) {
	// Проверим, есть ли уже такая ссылка
	if shortURL := getShortURLFromStorage(token, full); shortURL != "" {
		return shortURL, nil
	}

	// Если нет - создаём запись в файле
	file, err := os.OpenFile(storage.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		logging.Errorf("storage don't open to write! Error: %s. Path: %s", err, storage.path)
		return "", fmt.Errorf("storage don't open to write! Error: %w. Path: %s", err, storage.path)
	}

	item := StorageItem{Token: token, URL: full, ShortURL: utils.RandSeq(5)}
	data, err := json.Marshal(item)
	if err != nil {
		return "", fmt.Errorf("cannot encode storage item %s", err)
	}

	data = append(data, '\n')

	_, err = file.Write(data)

	return item.ShortURL, err
}

func (fs *FileStorage) Get(token string, short string) (string, error) {
	file, err := os.OpenFile(storage.path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logging.Errorf("storage don't open to read! Error: %s", err)
		return "", err
	}

	r := bufio.NewReader(file)
	s, e := readLine(r)
	var item StorageItem
	for e == nil {
		err = json.Unmarshal([]byte(s), &item)
		if err != nil {
			logging.Errorf("storage don't open to read! Error: %s. Path: %s", err, storage.path)
			return "", err
		}

		if item.ShortURL == short && item.Token == token {
			return item.URL, nil
		}
		s, e = readLine(r)
	}
	return "", nil
}

func (fs *FileStorage) GetUrlsByUser(token string) ([]db.UrlsByUserResponseMessage, error) {
	file, err := os.OpenFile(storage.path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logging.Errorf("storage don't open to read! Error: %s", err)
		return nil, err
	}

	r := bufio.NewReader(file)
	s, e := readLine(r)
	userURLs := make([]db.UrlsByUserResponseMessage, 0)
	var item StorageItem

	for e == nil {
		err = json.Unmarshal([]byte(s), &item)
		if err != nil {
			logging.Errorf("storage don't open to read! Error: %s. Path: %s", err, storage.path)
			return nil, err
		}

		if item.Token == token {
			userURLs = append(userURLs, db.UrlsByUserResponseMessage{
				ShortURL: fs.resultURL + "/" + item.ShortURL,
				FullURL:  item.URL,
			})

		}
		s, e = readLine(r)
	}
	return userURLs, nil
}

func (fs *FileStorage) Ping() error {
	return nil
}

func getShortURLFromStorage(token string, fullURL string) string {
	file, err := os.OpenFile(storage.path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logging.Errorf("storage don't open to read! Error: %s. Path: %s", err, storage.path)
		return ""
	}

	r := bufio.NewReader(file)
	s, e := readLine(r)
	var item StorageItem
	for e == nil {
		err = json.Unmarshal([]byte(s), &item)
		if err != nil {
			logging.Errorf("Cannot decode urls: %s", err)
			return ""
		}

		if item.URL == fullURL && item.Token == token {
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
