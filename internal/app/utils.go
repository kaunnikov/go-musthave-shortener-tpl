package app

import (
	"math/rand"
	"time"
)

func randSeq(letters int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	alphabet := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

	b := make([]rune, letters)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}
