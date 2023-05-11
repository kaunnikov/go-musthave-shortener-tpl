package app

import (
	"github.com/go-chi/chi/v5"
	"kaunnikov/go-musthave-shortener-tpl/internal/storage/fs"
	"net/http"
)

func (m *app) ShortHandler(w http.ResponseWriter, r *http.Request) {
	d := chi.URLParam(r, "id")

	full := fs.GetFullURLFromStorage(d)

	if full != "" {
		w.Header().Add("Location", full)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	http.Error(w, "Url not found!", http.StatusBadRequest)
}
