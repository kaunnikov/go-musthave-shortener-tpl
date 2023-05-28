package mem

import (
	"kaunnikov/go-musthave-shortener-tpl/internal/utils"
	"sync"
)

var (
	URLMap         = make(map[string]string, 1000)
	URLStorageSync = sync.Mutex{}
	storage        MemStorage
)

type MemStorage struct {
}

func Init() (*MemStorage, error) {
	storage = MemStorage{}
	return &storage, nil
}

func (mem *MemStorage) Save(full string) (string, error) {
	URLStorageSync.Lock()
	defer URLStorageSync.Unlock()

	short := utils.RandSeq(5)
	URLMap[short] = full
	return short, nil
}

func (mem *MemStorage) Get(short string) (string, error) {
	URLStorageSync.Lock()
	defer URLStorageSync.Unlock()

	return URLMap[short], nil
}

func (mem *MemStorage) Ping() error {
	return nil
}
