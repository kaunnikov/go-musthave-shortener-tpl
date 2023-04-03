package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var urlList = make(map[string]string, 1000)
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func mainHandle(w http.ResponseWriter, r *http.Request) {
	d := strings.TrimPrefix(r.URL.Path, "/")

	if d == "" {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
			return
		}
		if r.Header.Get("Content-Type") != "text/plain" {
			http.Error(w, "Invalid Content-Type!", http.StatusBadRequest)
			return
		}

		responseData, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		url := string(responseData)

		shortUrl := randSeq(10)
		urlList[shortUrl] = url

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://" + r.Host + "/" + shortUrl))
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	fullUrl, ok := urlList[d]
	if !ok {
		http.Error(w, "Url not found!", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Add("Location", fullUrl)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainHandle)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
