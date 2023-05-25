package db

import (
	"context"
	"database/sql"
	"kaunnikov/go-musthave-shortener-tpl/internal/config"
	"kaunnikov/go-musthave-shortener-tpl/internal/logging"
)

var conf *config.AppConfig

func Init(cfg *config.AppConfig) {
	conf = cfg
	checkTables()
}

func checkTables() {
	db, err := sql.Open("pgx", conf.DatabaseDSN)
	if err != nil {
		logging.Fatalf("DB don't open: %s", err)
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logging.Fatalf("DB don't Close: %s", err)
		}
	}(db)

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS url_storage (short_url varchar(16) not null, full_url varchar(128) not null)")
	if err != nil {
		logging.Fatalf("Table url_storage don't created: %s", err)
	}
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
