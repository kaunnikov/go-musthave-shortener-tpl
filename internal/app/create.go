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

	short := randSeq(10)
	URLMapSync.Lock()
	URLMap[short] = url
	URLMapSync.Unlock()

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(m.cfg.ResultURL + "/" + short))
	if err != nil {
		log.Fatalf("cannot write response to the client: %s", err)
	}
}
