package app

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"kaunnikov/go-musthave-shortener-tpl/internal/logging"
	"net/http"
	"time"
)

func (m *app) PingHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	db, err := sql.Open("pgx", m.cfg.DatabaseDSN)
	if err != nil {
		logging.Errorf("DB don't open: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer db.Close()

	_, err = db.QueryContext(ctx, "SELECT * FROM information_schema.tables")
	if err != nil {
		logging.Errorf("DB don't work: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
