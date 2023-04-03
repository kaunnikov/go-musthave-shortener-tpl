package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_mainHandle(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name    string
		method  string
		request string
		body    string
		want    want
	}{
		{
			name:    "GET запрос к методу, который принимает только POST",
			request: "/",
			method:  http.MethodGet,
			want:    want{http.StatusBadRequest},
		},
		{
			name:    "POST запрос без тела",
			request: "/",
			method:  http.MethodPost,
			want:    want{http.StatusBadRequest},
		},
		{
			name:    "Хороший POST запрос",
			request: "/",
			method:  http.MethodPost,
			body:    "https://yandex.ru",
			want:    want{http.StatusCreated},
		},
		{
			name:    "POST запрос к методу, который принимает только GET",
			request: "/fdaw43",
			method:  http.MethodPost,
			want:    want{http.StatusBadRequest},
		},
		{
			name:    "GET запрос к методу с несуществующим id",
			request: "/not-allowed-id",
			method:  http.MethodGet,
			want:    want{http.StatusBadRequest},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body := strings.NewReader(tt.body)
			request := httptest.NewRequest(tt.method, tt.request, body)
			w := httptest.NewRecorder()
			mainHandle(w, request)
			res := w.Result()
			assert.Equal(t, res.StatusCode, tt.want.statusCode)

		})
	}
}
