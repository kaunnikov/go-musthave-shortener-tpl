package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kaunnikov/go-musthave-shortener-tpl/internal/auth"
	"kaunnikov/go-musthave-shortener-tpl/internal/errs"
	"kaunnikov/go-musthave-shortener-tpl/internal/logging"
	"kaunnikov/go-musthave-shortener-tpl/internal/storage"
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

	token, err := auth.GetUserToken(w, r)
	if err != nil {
		logging.Errorf("cannot get user token: %s", err)
		http.Error(w, fmt.Sprintf("cannot get user token: %s", err), http.StatusBadRequest)
		return
	}

	short, err := storage.SaveURLInStorage(token, t.URL)
	if err != nil {
		logging.Errorf("cannot save URL in storage: %s", err)
		http.Error(w, fmt.Sprintf("cannot save URL in storage: %s", err), http.StatusBadRequest)
	}

	// Если нашли запись в БД, то отдадим с нужным статусом
	var doubleErr *errs.DoubleError
	if errors.As(err, &doubleErr) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)

		shortRes := shortenResponse{
			Result: m.cfg.ResultURL + "/" + doubleErr.ShortURL,
		}
		resp, err := json.Marshal(shortRes)
		if err != nil {
			logging.Errorf("cannot encode response: %s", err)
			http.Error(w, fmt.Sprintf("cannot encode response: %s", err), http.StatusBadRequest)
			return
		}

		_, err = w.Write(resp)
		if err != nil {
			logging.Errorf("cannot write response to the client: %s", err)
			http.Error(w, "Error in server!", http.StatusBadRequest)
		}
		return
	}

	resp, err := json.Marshal(shortenResponse{
		Result: m.cfg.ResultURL + "/" + short,
	})
	if err != nil {
		logging.Errorf("cannot encode response: %s", err)
		http.Error(w, fmt.Sprintf("cannot encode response: %s", err), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(resp)
	if err != nil {
		logging.Errorf("cannot write response to the client: %s", err)
		http.Error(w, fmt.Sprintf("cannot write response to the client: %s", err), http.StatusBadRequest)
	}
}
