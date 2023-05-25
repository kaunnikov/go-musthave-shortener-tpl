package app

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"kaunnikov/go-musthave-shortener-tpl/internal/logging"
	"kaunnikov/go-musthave-shortener-tpl/internal/storage/db"
	"net/http"
)

func (m *app) PingHandler(w http.ResponseWriter, r *http.Request) {
	err := db.Ping()

	if err != nil {
		logging.Errorf("DB don't open: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
