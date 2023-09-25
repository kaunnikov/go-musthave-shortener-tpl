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

func (m *app) UserDeleteURLHandler(w http.ResponseWriter, r *http.Request) {

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

	if len(body) == 0 {
		logging.Errorf("Empty request body!")
		http.Error(w, fmt.Sprintf("Empty request body: %s", body), http.StatusBadRequest)
		return
	}

	// Получаем список переданных URL
	var URLs []string
	err = json.Unmarshal(body, &URLs)
	if err != nil {
		logging.Errorf("cannot decode request body to `JSON`: %s", err)
		http.Error(w, fmt.Sprintf("cannot decode request body to `JSON`: %s", err), http.StatusBadRequest)
		return
	}

	if len(URLs) == 1 && URLs[0] == "" {
		logging.Errorf("Empty request body!")
		http.Error(w, fmt.Sprintf("Empty request body:%s", body), http.StatusBadRequest)
		return
	}

	token, err := auth.GetUserToken(w, r)
	if err != nil {
		logging.Errorf("cannot get user token: %s", err)
		http.Error(w, fmt.Sprintf("cannot get user token: %s", err), http.StatusBadRequest)
		return
	}

	// Наполняем канал входными данными
	ch := generator(URLs)
	deleteURL(ch, token)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

func generator(input []string) chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for _, u := range input {
			ch <- u
		}
	}()

	return ch
}

func deleteURL(ch <-chan string, token string) {
	for URL := range ch {
		err := storage.DeleteURLs(URL, token)
		if err != nil {
			logging.Errorf("cannot delete URL %s for user token: %s. Err: ", URL, token, err)
		}
	}
}
