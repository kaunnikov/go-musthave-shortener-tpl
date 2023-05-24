package mem

import (
	"sync"
)

var (
	URLMap         = make(map[string]string, 1000)
	URLStorageSync = sync.Mutex{}
)

func Append(fullURL string, shortURL string) {
	URLStorageSync.Lock()
	defer URLStorageSync.Unlock()

	URLMap[shortURL] = fullURL
}

func GetByShort(shortURL string) string {
	URLStorageSync.Lock()
	defer URLStorageSync.Unlock()

	return URLMap[shortURL]
}

func GetByFull(fullURL string) string {
	URLStorageSync.Lock()
	defer URLStorageSync.Unlock()

	for s, f := range URLMap {
		if f == fullURL {
			return s
		}
	}
	return ""
}
