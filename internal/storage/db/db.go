package db

import (
	"database/sql"
	"fmt"
	"kaunnikov/go-musthave-shortener-tpl/internal/config"
	"kaunnikov/go-musthave-shortener-tpl/internal/logging"
)

var tableName = "url_storage"
var conf *config.AppConfig

func Init(cfg *config.AppConfig) {
	conf = cfg
	checkTables()
}

func connect() (*sql.DB, error) {
	db, err := sql.Open("pgx", conf.DatabaseDSN)
	if err != nil {
		logging.Fatalf("DB don't open: %s", err)
		return nil, err
	}

	return db, nil
}

func Ping() error {
	if conf == nil {
		return fmt.Errorf("DB dot't init")
	}
	db, err := connect()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logging.Fatalf("DB don't Close: %s", err)
		}
	}(db)

	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func GetOrSave(fullURL string, short string) (string, error) {
	db, err := connect()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logging.Fatalf("DB don't Close: %s", err)
		}
	}(db)
	if err != nil {
		return "", err
	}

	// Сначала попробуем найти старую запись в БД
	shortFromDB, err := getShortByFullURL(db, fullURL)
	if err != nil {
		logging.Infof("Err in getByFullURL: %s", err)
		return "", err
	}

	// Если нашли - отдаём
	if shortFromDB != "" {
		return shortFromDB, nil
	}

	// Если не нашли - создаём
	short, err = insert(db, short, fullURL)
	if err != nil {
		logging.Infof("Err in insert: %s", err)
		return "", err
	}
	return short, nil
}

func checkTables() {
	db, err := connect()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logging.Fatalf("DB don't Close: %s", err)
		}
	}(db)
	if err != nil {
		logging.Fatalf("DB don't work: %s", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS " + tableName + " (short_url varchar(16) not null, full_url varchar(128) not null)")
	if err != nil {
		logging.Fatalf("Table "+tableName+" don't created: %s", err)
	}
}

func GetFullByShortURL(shortURL string) (string, error) {
	db, err := connect()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logging.Fatalf("DB don't Close: %s", err)
		}
	}(db)
	if err != nil {
		logging.Fatalf("DB don't work: %s", err)
	}

	var fullURL string
	res := db.QueryRow("SELECT full_url FROM "+tableName+" WHERE short_url = $1;", shortURL)
	err = res.Scan(&fullURL)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return fullURL, nil

}

func getShortByFullURL(db *sql.DB, fullURL string) (string, error) {
	var shortURL string
	row := db.QueryRow("SELECT short_url FROM "+tableName+" WHERE full_url = $1;", fullURL)
	err := row.Scan(&shortURL)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return shortURL, nil
}

func insert(db *sql.DB, shortURL string, fullURL string) (string, error) {
	_, err := db.Exec("INSERT INTO "+tableName+" (short_url, full_url) VALUES ($1, $2);", shortURL, fullURL)
	if err != nil {
		return "", err
	}
	return shortURL, nil
}
