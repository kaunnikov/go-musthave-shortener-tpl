package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func mainHandle(w http.ResponseWriter, r *http.Request) {
	d := strings.TrimPrefix(r.URL.Path, "/")

	if d == "" {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
			return
		}
		if r.Header.Get("Content-Type") != "text/plain" {
			http.Error(w, "Invalid Content-Type!", http.StatusBadRequest)
			return
		}

		responseData, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Print("Получено значение: ", string(responseData))

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://localhost:8080/EwHXdJfB"))
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Add("Location", "https://practicum.yandex.ru/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainHandle)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
