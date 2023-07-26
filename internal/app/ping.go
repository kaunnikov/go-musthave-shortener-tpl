package app

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"kaunnikov/go-musthave-shortener-tpl/internal/logging"
	"kaunnikov/go-musthave-shortener-tpl/internal/storage"
	"net/http"
)

func (m *app) PingHandler(w http.ResponseWriter, _ *http.Request) {
	err := storage.Ping()

	if err != nil {
		logging.Errorf("Ping fail: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
