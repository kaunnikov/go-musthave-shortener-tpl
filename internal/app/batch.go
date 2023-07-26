package app

import (
	"encoding/json"
	"fmt"
	"io"
	"kaunnikov/go-musthave-shortener-tpl/internal/auth"
	"kaunnikov/go-musthave-shortener-tpl/internal/logging"
	"kaunnikov/go-musthave-shortener-tpl/internal/storage"
	"net/http"
)

func (m *app) BatchHandler(w http.ResponseWriter, r *http.Request) {
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

	var t []batchRequestMessage
	err = json.Unmarshal(body, &t)
	if err != nil {
		logging.Errorf("cannot decode request body to `JSON`: %s", err)
		http.Error(w, fmt.Sprintf("cannot decode request body to `JSON`: %s", err), http.StatusBadRequest)
		return
	}

	var result []batchResponseMessage
	token, _ := auth.GetUserToken(w, r)
	for _, item := range t {
		short, err := storage.SaveURLInStorage(token, item.URL)
		if err != nil {
			logging.Errorf("error write data: %s", err)
			http.Error(w, "Error in server!", http.StatusBadRequest)
			return
		}

		resItem := batchResponseMessage{item.CorrelationID, m.cfg.ResultURL + "/" + short}
		result = append(result, resItem)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp, err := json.Marshal(result)
	if err != nil {
		logging.Errorf("cannot encode response: %s", err)
		http.Error(w, fmt.Sprintf("cannot encode response: %s", err), http.StatusBadRequest)
	}

	_, err = w.Write(resp)
	if err != nil {
		logging.Errorf("cannot write response to the client: %s", err)
		http.Error(w, fmt.Sprintf("cannot write response to the client: %s", err), http.StatusBadRequest)

	}
}
