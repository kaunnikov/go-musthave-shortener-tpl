package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (m *app) JSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid Content Type!", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot read request body: %s", err), http.StatusBadRequest)
		return
	}

	var t jsonStruct
	err = json.Unmarshal(body, &t)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot decode request body to `JSON`: %s", err), http.StatusBadRequest)
		return
	}

	short := randSeq(5)
	s := StorageItem{URL: t.URL, ShortURL: short}
	URLStorageSync.Lock()
	short, err = m.SaveURLInStorage(&s)
	if err != nil {
		Sugar.Errorf("error write data: %s", err)
		http.Error(w, "Error in server!", http.StatusBadRequest)
		return
	}
	URLStorageSync.Unlock()

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	shortRes := shortenResponse{
		Result: m.cfg.ResultURL + "/" + short,
	}

	resp, err := json.Marshal(shortRes)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot encode responce: %s", err), http.StatusBadRequest)
	}

	_, err = w.Write(resp)
	if err != nil {
		log.Fatalf("cannot write response to the client: %s", err)
	}
}
