package mem

import (
	"kaunnikov/go-musthave-shortener-tpl/internal/utils"
	"sync"
)

var (
	URLMap         = make(map[string]string, 1000)
	URLStorageSync = sync.Mutex{}
	storage        MemoryStorage
)

type MemoryStorage struct {
}

func Init() (*MemoryStorage, error) {
	storage = MemoryStorage{}
	return &storage, nil
}

func (mem *MemoryStorage) Save(full string) (string, error) {
	URLStorageSync.Lock()
	defer URLStorageSync.Unlock()

	short := utils.RandSeq(5)
	URLMap[short] = full
	return short, nil
}

func (mem *MemoryStorage) Get(short string) (string, error) {
	URLStorageSync.Lock()
	defer URLStorageSync.Unlock()

	return URLMap[short], nil
}

func (mem *MemoryStorage) Ping() error {
	return nil
}
