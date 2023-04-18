package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"kaunnikov/go-musthave-shortener-tpl/config"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
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

type jsonStruct struct {
	URL string `json:"URL"`
}

type shortenResponse struct {
	Result string `json:"result"`
}

func jsonHandle(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid Content Type!", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid POST body!", http.StatusBadRequest)
		return
	}

	var t jsonStruct
	err = json.Unmarshal(body, &t)
	if err != nil {
		panic(err)
	}

	short := randSeq(10)
	urlList[short] = t.URL

	w.Header().Add("Content-Type", "application/json")
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

	shortRes := shortenResponse{
		Result: "http://" + r.Host + urlPrefix + short,
	}

	resp, err := json.Marshal(shortRes)
	if err != nil {
		panic(err)
	}

	_, errWrite := w.Write(resp)
	if errWrite != nil {
		panic(errWrite)
	}
}

func mainHandle(w http.ResponseWriter, r *http.Request) {

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

func main() {
	appConfig := config.ParseFlags()

	envRunAddr := os.Getenv("SERVER_ADDRESS")
	envRunAddr = strings.TrimSpace(envRunAddr)
	if envRunAddr != "" {
		appConfig.Host = envRunAddr
	}

	envBaseURL := os.Getenv("BASE_URL")
	envBaseURL = strings.TrimSpace(envBaseURL)
	if envBaseURL != "" {
		appConfig.Prefix = envBaseURL
	}

	prefix = appConfig.Prefix
	defaultRoute := "/"
	if prefix != "" {
		defaultRoute = prefix + "/"
	}
	log.Println("Prefix: " + defaultRoute)

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", mainHandle)
		r.Get(defaultRoute+"{id}", shortHandle)
		r.Post("/api/shorten", jsonHandle)
	})
	fmt.Println("Running server on", appConfig.Host)
	log.Fatal(http.ListenAndServe(appConfig.Host, r))
}
