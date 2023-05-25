package app

import (
	"context"
	_ "github.com/jackc/pgx/v5/stdlib"
	"kaunnikov/go-musthave-shortener-tpl/internal/logging"
	"kaunnikov/go-musthave-shortener-tpl/internal/storage/db"
	"net/http"
	"time"
)

func (m *app) PingHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	err := db.Ping(ctx)

	if err != nil {
		logging.Errorf("DB don't open: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
