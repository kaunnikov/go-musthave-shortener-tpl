package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	"kaunnikov/go-musthave-shortener-tpl/config"
	"kaunnikov/go-musthave-shortener-tpl/internal/app"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {
	cfg := &config.AppConfig{}

	flag.StringVar(&cfg.Host, "a", "localhost:8080", "Default Host:port")
	flag.StringVar(&cfg.ResultURL, "b", "http://localhost:8080", "Default result URL")
	flag.Parse()

	loadFromENV(cfg)

	newApp := app.NewApp(cfg)

	re := regexp.MustCompile(`:\d{2,}/(\w+)`)
	patternResultURL := "/"
	if len(re.FindSubmatch([]byte(cfg.ResultURL))) == 2 {
		patternResultURL = "/" + string(re.FindSubmatch([]byte(cfg.ResultURL))[1])
	}

	r := chi.NewRouter()
	r.Post("/", newApp.CreateHandler)
	r.Get(patternResultURL+"{id}", newApp.ShortHandler)
	r.Post("/api/shorten", newApp.JSONHandler)

	log.Println("Running server on", cfg.Host)
	log.Fatal(http.ListenAndServe(cfg.Host, r))
}

func loadFromENV(cfg *config.AppConfig) {
	envRunAddr := os.Getenv("SERVER_ADDRESS")
	envRunAddr = strings.TrimSpace(envRunAddr)
	if envRunAddr != "" {
		cfg.Host = envRunAddr
	}

	envBaseURL := os.Getenv("BASE_URL")
	envBaseURL = strings.TrimSpace(envBaseURL)
	if envBaseURL != "" {
		cfg.ResultURL = envBaseURL
	}

}
