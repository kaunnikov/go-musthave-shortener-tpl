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
	flag.Parse()
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
