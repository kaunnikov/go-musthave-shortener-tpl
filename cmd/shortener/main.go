package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"kaunnikov/go-musthave-shortener-tpl/config"
	"kaunnikov/go-musthave-shortener-tpl/internal/app"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	cfg := &config.AppConfig{}

	loadFromArgs(cfg)
	loadFromENV(cfg)

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("logger don't Run! %s", err)
	}

	app.Sugar = logger.Sugar()

	newApp := app.NewApp(cfg)

	r := chi.NewRouter()
	r.Use(app.CustomMiddlewareLogger)
	r.Use(app.CustomCompression)

	r.Post("/", newApp.CreateHandler)
	r.Get("/{id}", newApp.ShortHandler)
	r.Post("/api/shorten", newApp.JSONHandler)

	log.Println("Running server on", cfg.Host)
	log.Fatal(http.ListenAndServe(cfg.Host, r))
}

func loadFromArgs(cfg *config.AppConfig) {
	flag.StringVar(&cfg.Host, "a", "localhost:8080", "Default Host:port")
	flag.StringVar(&cfg.ResultURL, "b", "http://localhost:8080", "Default result URL")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/short-url-db.json", "Default File Storage Path")
	flag.Parse()
}

func loadFromENV(cfg *config.AppConfig) {
	envRunAddr := strings.TrimSpace(os.Getenv("SERVER_ADDRESS"))
	if envRunAddr != "" {
		cfg.Host = envRunAddr
	}

	envBaseURL := strings.TrimSpace(os.Getenv("BASE_URL"))
	if envBaseURL != "" {
		cfg.ResultURL = envBaseURL
	}

	fileStorageFile := strings.TrimSpace(os.Getenv("FILE_STORAGE_PATH"))
	if fileStorageFile != "" {
		cfg.FileStoragePath = fileStorageFile
	}
}
