package app

import (
	"encoding/json"
	"fmt"
	"io"
	"kaunnikov/go-musthave-shortener-tpl/internal/logging"
	"kaunnikov/go-musthave-shortener-tpl/internal/storage"
	"kaunnikov/go-musthave-shortener-tpl/internal/utils"
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
	for _, item := range t {
		short, err := storage.SaveURLInStorage(item.URL, utils.RandSeq(5))
		if err != nil {
			logging.Errorf("error write data: %s", err)
			http.Error(w, "Error in server!", http.StatusBadRequest)
			return
		}

		resItem := batchResponseMessage{item.CorrelationId, m.cfg.ResultURL + "/" + short}
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
		logging.Fatalf("cannot write response to the client: %s", err)
	}
}
