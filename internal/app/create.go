package app

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func (m *app) CreateHandler(w http.ResponseWriter, r *http.Request) {

	responseData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot read request body: %s", err), http.StatusBadRequest)
		return
	}
	if string(responseData) == "" {
		http.Error(w, "Empty POST request body!", http.StatusBadRequest)
		return
	}
	url := string(responseData)

	short := randSeq(5)
	s := StorageItem{URL: url, ShortURL: short}
	URLStorageSync.Lock()
	short, err = m.SaveURLInStorage(&s)
	if err != nil {
		Sugar.Errorf("error write data: %s", err)
		http.Error(w, "Error in server!", http.StatusBadRequest)
		return
	}
	URLStorageSync.Unlock()

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(m.cfg.ResultURL + "/" + short))
	if err != nil {
		log.Fatalf("cannot write response to the client: %s", err)
	}
}
