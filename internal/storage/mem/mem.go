package mem

import (
	"fmt"
	"kaunnikov/go-musthave-shortener-tpl/internal/config"
	"kaunnikov/go-musthave-shortener-tpl/internal/storage/db"
	"kaunnikov/go-musthave-shortener-tpl/internal/utils"
	"sync"
)

var (
	URLMap         = make(map[string][]map[string]string, 1000)
	URLStorageSync = sync.Mutex{}
	storage        MemoryStorage
)

type MemoryStorage struct {
	resultURL string
}

func Init(cfg *config.AppConfig) (*MemoryStorage, error) {
	storage = MemoryStorage{
		resultURL: cfg.ResultURL,
	}
	return &storage, nil
}

func (mem *MemoryStorage) Save(token string, full string) (string, error) {
	URLStorageSync.Lock()
	defer URLStorageSync.Unlock()

	short := utils.RandSeq(5)

	item := make(map[string]string)
	item[short] = full
	URLMap[token] = append(URLMap[token], item)

	fmt.Println(URLMap)
	return short, nil
}

func (mem *MemoryStorage) Get(short string) (string, error) {
	URLStorageSync.Lock()
	defer URLStorageSync.Unlock()

	for _, m := range URLMap {
		for _, u := range m {
			full, isFind := u[short]
			if isFind {
				return full, nil
			}
		}
	}

	return "", nil
}

func (mem *MemoryStorage) GetUrlsByUser(token string) ([]db.UrlsByUserResponseMessage, error) {
	URLStorageSync.Lock()
	defer URLStorageSync.Unlock()
	URLs, ok := URLMap[token]

	if !ok {
		return nil, nil
	}

	var userURLs = make([]db.UrlsByUserResponseMessage, 0)
	for _, item := range URLs {
		for short, full := range item {
			userURLs = append(userURLs, db.UrlsByUserResponseMessage{
				ShortURL: mem.resultURL + "/" + short,
				FullURL:  full,
			})

		}
	}
	return userURLs, nil

}

func (mem *MemoryStorage) Ping() error {
	return nil
}
