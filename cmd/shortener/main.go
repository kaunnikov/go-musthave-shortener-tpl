package main

import (
	"kaunnikov/go-musthave-shortener-tpl/internal/app"
	"kaunnikov/go-musthave-shortener-tpl/internal/config"
	"kaunnikov/go-musthave-shortener-tpl/internal/logging"
	"kaunnikov/go-musthave-shortener-tpl/internal/storage/db"
	"kaunnikov/go-musthave-shortener-tpl/internal/storage/fs"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()

	if err := logging.Init(); err != nil {
		log.Fatalf("logger don't Run!: %s", err)
	}

	fs.Init(cfg)
	if cfg.DatabaseDSN != "" {
		// "host=localhost port=5433 user=postgres password=password dbname=postgres sslmode=disable"
		db.Init(cfg)
	}
	newApp := app.NewApp(cfg)

	logging.Infof("Running server on %s", cfg.Host)
	logging.Fatalf("cannot listen and serve: %s", http.ListenAndServe(cfg.Host, newApp))
}
