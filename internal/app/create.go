package app

import (
	"fmt"
	"io"
	"kaunnikov/go-musthave-shortener-tpl/internal/logging"
	"kaunnikov/go-musthave-shortener-tpl/internal/storage/fs"
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
	url := string(responseData)

	short, err := fs.SaveURLInStorage(url)
	if err != nil {
		logging.Errorf("error write data: %s", err)
		http.Error(w, "Error in server!", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(m.cfg.ResultURL + "/" + short))
	if err != nil {
		logging.Fatalf("cannot write response to the client: %s", err)
	}
	return
}
