package app

import (
	"errors"
	"fmt"
	"io"
	"kaunnikov/go-musthave-shortener-tpl/internal/auth"
	"kaunnikov/go-musthave-shortener-tpl/internal/errs"
	"kaunnikov/go-musthave-shortener-tpl/internal/logging"
	"kaunnikov/go-musthave-shortener-tpl/internal/storage"
	"net/http"
)

func (m *app) CreateHandler(w http.ResponseWriter, r *http.Request) {
	responseData, err := io.ReadAll(r.Body)
	if err != nil {
		logging.Errorf("cannot read request body: %s", err)
		http.Error(w, fmt.Sprintf("cannot read request body: %s", err), http.StatusBadRequest)
		return
	}

	if string(responseData) == "" {
		logging.Errorf("Empty POST request body! %s", r.URL)
		http.Error(w, "Empty POST request body!", http.StatusBadRequest)
		return
	}

	token, err := auth.GetUserToken(w, r)
	if err != nil {
		logging.Errorf("cannot get user token: %s", err)
		http.Error(w, fmt.Sprintf("cannot get user token: %s", token), http.StatusBadRequest)
		return
	}

	short, err := storage.SaveURLInStorage(token, string(responseData))
	// Если нашли запись в БД, то отдадим с нужным статусом
	var doubleErr *errs.DoubleError
	if errors.As(err, &doubleErr) {
		w.WriteHeader(http.StatusConflict)
		_, err = w.Write([]byte(m.cfg.ResultURL + "/" + doubleErr.ShortURL))
		if err != nil {
			logging.Errorf("cannot write response to the client: %s", err)
			http.Error(w, fmt.Sprintf("cannot write response to the client: %s", err), http.StatusBadRequest)
		}
		return
	}
	if err != nil {
		logging.Errorf("error write data: %s", err)
		http.Error(w, "Error in server!", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(m.cfg.ResultURL + "/" + short))
	if err != nil {
		logging.Errorf("cannot write response to the client: %s", err)
		http.Error(w, fmt.Sprintf("cannot write response to the client: %s", err), http.StatusBadRequest)
	}
}
