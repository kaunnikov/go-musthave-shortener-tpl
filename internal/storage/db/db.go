package db

import (
	"context"
	"database/sql"
	"fmt"
	"kaunnikov/go-musthave-shortener-tpl/internal/config"
	"kaunnikov/go-musthave-shortener-tpl/internal/errs"
	"kaunnikov/go-musthave-shortener-tpl/internal/logging"
	"kaunnikov/go-musthave-shortener-tpl/internal/utils"
)

var (
	tableName = "url_storage_iter_15"
	storage   DataBaseStorage
)

type DataBaseStorage struct {
	connect   *sql.DB
	resultURL string
}

type UrlsByUserResponseMessage struct {
	ShortURL string `json:"short_url"`
	FullURL  string `json:"original_url"`
}

func Init(cfg *config.AppConfig) (*DataBaseStorage, error) {
	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		logging.Errorf("DB don't open: %s", err)
		return nil, fmt.Errorf("DB don't open: %w", err)
	}
	storage = DataBaseStorage{
		connect:   db,
		resultURL: cfg.ResultURL,
	}

	err = checkTables()
	if err != nil {
		logging.Errorf("Don't check tables: %s", err)
		return nil, err
	}

	return &storage, nil
}

func (db *DataBaseStorage) Save(token string, full string) (string, error) {
	// Сначала попробуем найти старую запись в БД
	shortFromDB, err := getShortByFullURL(token, full)
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
	short, err = insert(token, short, full)
	if err != nil {
		logging.Infof("Err in insert: %s", err)
		return "", err
	}
	return short, nil
}

func (db *DataBaseStorage) Get(short string) (string, error) {
	var fullURL string
	var isDeletet bool
	res := storage.connect.QueryRowContext(context.Background(), "SELECT full_url, is_deleted FROM "+tableName+" WHERE short_url = $1;", short)
	err := res.Scan(&fullURL, &isDeletet)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	if isDeletet {
		return "", &errs.URLIsDeletedError{
			Err: fmt.Errorf("URL %s is deleted", short),
		}
	}

	return fullURL, nil
}

func (db *DataBaseStorage) GetUrlsByUser(token string) ([]UrlsByUserResponseMessage, error) {
	items := make([]UrlsByUserResponseMessage, 0)
	rows, err := storage.connect.QueryContext(context.Background(), "SELECT short_url, full_url FROM "+tableName+" WHERE user_token = $1;", token)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var u UrlsByUserResponseMessage
		err = rows.Scan(&u.ShortURL, &u.FullURL)
		if err != nil {
			return nil, err
		}
		u.ShortURL = db.resultURL + "/" + u.ShortURL

		items = append(items, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (db *DataBaseStorage) DeleteURLs(URL string, token string) error {
	query := "UPDATE " + tableName + " SET is_deleted = true WHERE short_url = $1 AND user_token = $2"

	_, err := storage.connect.ExecContext(context.Background(), query, URL, token)
	if err != nil {
		return err
	}
	return nil
}

func (db *DataBaseStorage) Ping() error {
	err := storage.connect.Ping()
	if err != nil {
		return err
	}

	return nil
}

func checkTables() error {
	query := "CREATE TABLE IF NOT EXISTS " + tableName + " " +
		"(user_token varchar(36), short_url varchar(16) not null, full_url varchar(128) not null, is_deleted boolean DEFAULT false)"

	_, err := storage.connect.ExecContext(context.Background(), query)
	if err != nil {
		return fmt.Errorf("table "+tableName+" don't created: %w", err)
	}
	return nil
}

func getShortByFullURL(token string, fullURL string) (string, error) {
	var shortURL string
	row := storage.connect.QueryRowContext(context.Background(), "SELECT short_url FROM "+tableName+" WHERE full_url = $1 AND user_token = $2;", fullURL, token)
	err := row.Scan(&shortURL)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return shortURL, nil
}

func insert(token string, shortURL string, fullURL string) (string, error) {
	_, err := storage.connect.ExecContext(context.Background(), "INSERT INTO "+tableName+" (user_token, short_url, full_url) VALUES ($1, $2, $3);", token, shortURL, fullURL)
	if err != nil {
		return "", err
	}
	return shortURL, nil
}
