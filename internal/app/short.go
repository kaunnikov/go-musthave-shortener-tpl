package app

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (m *app) ShortHandler(w http.ResponseWriter, r *http.Request) {
	d := chi.URLParam(r, "id")

	URLMapSync.Lock()
	full, ok := URLMap[d]
	URLMapSync.Unlock()

	if ok {
		w.Header().Add("Location", full)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	http.Error(w, "Url not found!", http.StatusBadRequest)
}
