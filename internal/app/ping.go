package app

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"kaunnikov/go-musthave-shortener-tpl/internal/logging"
	"net/http"
)

func (m *app) PingHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("pgx", m.cfg.DatabaseDSN)
	if err != nil {
		logging.Errorf("DB don't open: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer db.Close()

	_, err = db.Exec("SELECT * FROM information_schema.tables")
	if err != nil {
		logging.Errorf("DB don't work: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
