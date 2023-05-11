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
	URLMap[fullURL] = shortURL
	URLStorageSync.Unlock()
}

func GetByFull(fullURL string) string {
	URLStorageSync.Lock()
	res := URLMap[fullURL]
	URLStorageSync.Unlock()
	return res
}

func GetByShort(shortURL string) string {
	URLStorageSync.Lock()
	var res string
	for f, s := range URLMap {
		if s == shortURL {
			res = f
			break
		}
	}
	URLStorageSync.Unlock()
	return res
}
