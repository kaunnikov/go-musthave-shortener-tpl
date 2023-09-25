package app

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"kaunnikov/go-musthave-shortener-tpl/internal/errs"
	"kaunnikov/go-musthave-shortener-tpl/internal/logging"
	"kaunnikov/go-musthave-shortener-tpl/internal/storage"
	"net/http"
)

func (m *app) ShortHandler(w http.ResponseWriter, r *http.Request) {
	d := chi.URLParam(r, "id")

	full, err := storage.GetFullURL(d)

	// Если ссылка удалена - отдаём 410 статус
	var deletedErr *errs.URLIsDeletedError
	if errors.As(err, &deletedErr) {
		w.WriteHeader(http.StatusGone)
		return
	}

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
