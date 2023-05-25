package db

import (
	"context"
	"database/sql"
	"kaunnikov/go-musthave-shortener-tpl/internal/config"
)

var conf *config.AppConfig

func Init(cfg *config.AppConfig) {
	conf = cfg
}

func Ping(ctx context.Context) error {
	db, err := sql.Open("pgx", conf.DatabaseDSN)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.QueryContext(ctx, "SELECT * FROM information_schema.tables")
	if err != nil {
		return err
	}

	return nil

}
