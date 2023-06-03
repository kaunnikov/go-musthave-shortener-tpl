package app

import (
	"github.com/go-chi/chi/v5"
	"kaunnikov/go-musthave-shortener-tpl/internal/logging"
	"kaunnikov/go-musthave-shortener-tpl/internal/storage"
	"net/http"
)

func (m *app) ShortHandler(w http.ResponseWriter, r *http.Request) {
	d := chi.URLParam(r, "id")

	full, err := storage.GetFullURL(d)
	if err != nil {
		logging.Errorf("Cannot find short url for full %q: %s", r.URL, err)
		http.Error(w, "Server Error!", http.StatusInternalServerError)
	}

	if full != "" {
		w.Header().Add("Location", full)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	logging.Errorf("Url not found: %s", r.URL)
	http.Error(w, "Url not found!", http.StatusBadRequest)
}
