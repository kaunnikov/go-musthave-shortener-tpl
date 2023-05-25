package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kaunnikov/go-musthave-shortener-tpl/internal/logging"
	"kaunnikov/go-musthave-shortener-tpl/internal/storage"
	"kaunnikov/go-musthave-shortener-tpl/internal/storage/db"
	"kaunnikov/go-musthave-shortener-tpl/internal/utils"
	"net/http"
)

func (m *app) JSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		logging.Errorf("Invalid Content Type: %s", r.Header.Get("Content-Type"))
		http.Error(w, "Invalid Content Type!", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logging.Errorf("cannot read request body: %s", err)
		http.Error(w, fmt.Sprintf("cannot read request body: %s", err), http.StatusBadRequest)
		return
	}

	var t requestMessage
	err = json.Unmarshal(body, &t)
	if err != nil {
		logging.Errorf("cannot decode request body to `JSON`: %s", err)
		http.Error(w, fmt.Sprintf("cannot decode request body to `JSON`: %s", err), http.StatusBadRequest)
		return
	}

	short, err := storage.SaveURLInStorage(t.URL, utils.RandSeq(5))
	// Если нашли запись в БД, то отдадим с нужным статусом
	var doubleErr *db.DoubleError
	if errors.As(err, &doubleErr) {
		w.WriteHeader(http.StatusConflict)
		_, err = w.Write([]byte(m.cfg.ResultURL + "/" + doubleErr.ShortURL))
		if err != nil {
			logging.Fatalf("cannot write response to the client: %s", err)
		}
		return
	}

	if err != nil {
		logging.Errorf("error write data: %s", err)
		http.Error(w, "Error in server!", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	shortRes := shortenResponse{
		Result: m.cfg.ResultURL + "/" + short,
	}

	resp, err := json.Marshal(shortRes)
	if err != nil {
		logging.Errorf("cannot encode response: %s", err)
		http.Error(w, fmt.Sprintf("cannot encode response: %s", err), http.StatusBadRequest)
	}

	_, err = w.Write(resp)
	if err != nil {
		logging.Fatalf("cannot write response to the client: %s", err)
	}
}
