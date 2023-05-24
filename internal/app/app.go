package app

import (
	"github.com/go-chi/chi/v5"
	"kaunnikov/go-musthave-shortener-tpl/internal/compression"
	"kaunnikov/go-musthave-shortener-tpl/internal/config"
	"kaunnikov/go-musthave-shortener-tpl/internal/logging"
)

type app struct {
	*chi.Mux
	cfg *config.AppConfig
}

type requestMessage struct {
	URL string `json:"URL"`
}

type shortenResponse struct {
	Result string `json:"result"`
}

func NewApp(cfg *config.AppConfig) *app {
	a := &app{
		chi.NewRouter(),
		cfg,
	}
	a.registerRouetes()
	return a
}

func (m *app) registerRouetes() {
	m.Use(logging.CustomMiddlewareLogger)
	m.Use(compression.CustomCompression)

	m.Post("/", m.CreateHandler)
	m.Get("/{id}", m.ShortHandler)
	m.Post("/api/shorten", m.JSONHandler)
	m.Get("/ping", m.PingHandler)
}
