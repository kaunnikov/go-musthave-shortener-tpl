package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"kaunnikov/go-musthave-shortener-tpl/cmd/config"
	"log"
	"math/rand"
	"net/http"
)

var urlList = make(map[string]string, 1000)
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var prefix string

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func mainHandle(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}

	responseData, err := io.ReadAll(r.Body)
	if err != nil || string(responseData) == "" {
		http.Error(w, "Invalid POST body!", http.StatusBadRequest)
		return
	}
	url := string(responseData)

	short := randSeq(10)
	urlList[short] = url

	w.WriteHeader(http.StatusCreated)

	lastChar := string(prefix[len(prefix)-1])
	if lastChar != "/" {
		prefix += "/"
	}

	_, errWrite := w.Write([]byte("http://" + r.Host + prefix + short))
	if errWrite != nil {
		panic(errWrite)
	}
}

func shortHandle(w http.ResponseWriter, r *http.Request) {
	d := chi.URLParam(r, "id")

	if full, ok := urlList[d]; ok {
		w.Header().Add("Location", full)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	http.Error(w, "Url not found!", http.StatusBadRequest)
}

func main() {
	appConfig := config.ParseFlags()
	prefix = appConfig.Prefix
	defaultRoute := "/"
	if prefix != "" {
		defaultRoute = prefix
	}

	r := chi.NewRouter()
	r.Route(defaultRoute, func(r chi.Router) {
		r.Post("/", mainHandle)
		r.Get("/{id}", shortHandle)
	})
	fmt.Println("Running server on", appConfig.Host)
	log.Fatal(http.ListenAndServe(appConfig.Host, r))
}
