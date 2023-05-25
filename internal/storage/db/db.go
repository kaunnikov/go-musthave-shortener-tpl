package db

import (
	"database/sql"
	"kaunnikov/go-musthave-shortener-tpl/internal/config"
)

var conf *config.AppConfig

func Init(cfg *config.AppConfig) {
	conf = cfg
}

func Ping() error {
	db, err := sql.Open("pgx", conf.DatabaseDSN)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}
