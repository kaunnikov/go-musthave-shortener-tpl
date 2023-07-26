package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"kaunnikov/go-musthave-shortener-tpl/internal/auth"
	"kaunnikov/go-musthave-shortener-tpl/internal/errs"
	"kaunnikov/go-musthave-shortener-tpl/internal/logging"
	"kaunnikov/go-musthave-shortener-tpl/internal/storage"
	"net/http"
)

func (m *app) UserURLsHandler(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetUserToken(w, r)
	var authError *errs.TokenNotFoundInCookie

	// Если кука не содержит ID пользователя - возвращаем 401
	if errors.As(err, &authError) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	result, err := storage.GetURLsByUser(token)

	// При отсутствии сокращённых пользователем URL  - возвращаем 204
	if len(result) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp, err := json.Marshal(result)
	fmt.Println(resp)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		logging.Errorf("cannot write response to the client: %s", err)
		http.Error(w, fmt.Sprintf("cannot write response to the client: %s", err), http.StatusBadRequest)
	}

}
