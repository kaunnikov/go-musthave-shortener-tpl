package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"math/rand"
	"net/http"
)

var urlList = make(map[string]string, 1000)
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var prefix string

type AppConfig struct {
	Host   string
	Prefix string
}

type Message struct {
	Name string
}

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

	urlPrefix := prefix
	if len(urlPrefix) > 1 {
		lastChar := string(urlPrefix[len(urlPrefix)-1])
		if lastChar != "/" {
			urlPrefix += "/"
		}
	}
	if urlPrefix == "" {
		urlPrefix = "/"
	}

	_, errWrite := w.Write([]byte("http://" + r.Host + urlPrefix + short))
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

func parseFlags() *AppConfig {
	appConfig := AppConfig{Prefix: ""}

	flag.StringVar(&appConfig.Host, "a", ":8080", "Default Host:port")
	flag.Func("b", "App prefix", func(s string) error {

		if string(s[0]) == "/" {
			s = s[1:]
		}

		if s != "" {
			appConfig.Prefix = "/" + s
		}

		return nil
	})

	flag.Parse()
	return &appConfig
}

func main() {
	m := Message{"Hello"}
	b, _ := json.Marshal(m)
	fmt.Println(b)

	appConfig := parseFlags()
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
