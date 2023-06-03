package db

import (
	"database/sql"
	"fmt"
	"kaunnikov/go-musthave-shortener-tpl/internal/config"
	"kaunnikov/go-musthave-shortener-tpl/internal/errs"
	"kaunnikov/go-musthave-shortener-tpl/internal/logging"
	"kaunnikov/go-musthave-shortener-tpl/internal/utils"
)

var (
	tableName = "url_storage"
	storage   DataBaseStorage
)

type DataBaseStorage struct {
	connect *sql.DB
}

func Init(cfg *config.AppConfig) (*DataBaseStorage, error) {
	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		logging.Errorf("DB don't open: %s", err)
		return nil, fmt.Errorf("DB don't open: %w", err)
	}
	storage = DataBaseStorage{connect: db}

	err = checkTables()
	if err != nil {
		logging.Errorf("Don't check tables: %s", err)
		return nil, err
	}

	return &storage, nil
}

func (db *DataBaseStorage) Save(full string) (string, error) {
	// Сначала попробуем найти старую запись в БД
	shortFromDB, err := getShortByFullURL(full)
	if err != nil {
		logging.Infof("Err in getByFullURL: %s", err)
		return "", err
	}

	// Если нашли - отдаём
	if shortFromDB != "" {
		return "", &errs.DoubleError{
			ShortURL: shortFromDB,
			Err:      fmt.Errorf("double for %s", full),
		}
	}

	// Если не нашли - создаём
	short := utils.RandSeq(5)
	short, err = insert(short, full)
	if err != nil {
		logging.Infof("Err in insert: %s", err)
		return "", err
	}
	return short, nil
}

func (db *DataBaseStorage) Get(short string) (string, error) {
	var fullURL string
	res := storage.connect.QueryRow("SELECT full_url FROM "+tableName+" WHERE short_url = $1;", short)
	err := res.Scan(&fullURL)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return fullURL, nil
}

func (db *DataBaseStorage) Ping() error {
	err := storage.connect.Ping()
	if err != nil {
		return err
	}

	return nil
}

func checkTables() error {
	_, err := storage.connect.Exec("CREATE TABLE IF NOT EXISTS " + tableName + " (short_url varchar(16) not null, full_url varchar(128) not null)")
	if err != nil {
		return fmt.Errorf("table "+tableName+" don't created: %w", err)
	}
	return nil
}

func getShortByFullURL(fullURL string) (string, error) {
	var shortURL string
	row := storage.connect.QueryRow("SELECT short_url FROM "+tableName+" WHERE full_url = $1;", fullURL)
	err := row.Scan(&shortURL)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return shortURL, nil
}

func insert(shortURL string, fullURL string) (string, error) {
	_, err := storage.connect.Exec("INSERT INTO "+tableName+" (short_url, full_url) VALUES ($1, $2);", shortURL, fullURL)
	if err != nil {
		return "", err
	}
	return shortURL, nil
}
