package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

//func Test_mainHandle(t *testing.T) {
//	type want struct {
//		statusCode int
//	}
//	tests := []struct {
//		name    string
//		method  string
//		request string
//		body    string
//		want    want
//	}{
//		{
//			name:    "GET запрос к методу, который принимает только POST",
//			request: "/",
//			method:  http.MethodGet,
//			want:    want{http.StatusBadRequest},
//		},
//		{
//			name:    "POST запрос без тела",
//			request: "/",
//			method:  http.MethodPost,
//			want:    want{http.StatusBadRequest},
//		},
//		{
//			name:    "Хороший POST запрос",
//			request: "/",
//			method:  http.MethodPost,
//			body:    "https://yandex.ru",
//			want:    want{http.StatusCreated},
//		},
//		{
//			name:    "POST запрос к методу, который принимает только GET",
//			request: "/fdaw43",
//			method:  http.MethodPost,
//			want:    want{http.StatusBadRequest},
//		},
//		{
//			name:    "GET запрос к методу с несуществующим id",
//			request: "/not-allowed-id",
//			method:  http.MethodGet,
//			want:    want{http.StatusBadRequest},
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//
//			body := strings.NewReader(tt.body)
//			request := httptest.NewRequest(tt.method, tt.request, body)
//			w := httptest.NewRecorder()
//			mainHandle(w, request)
//			res := w.Result()
//			assert.Equal(t, res.StatusCode, tt.want.statusCode)
//			defer res.Body.Close()
//		})
//	}
//}

var testkeymap = map[string]string{}
var fullUrl string
var shortUrl string

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (int, string) {
	body := strings.NewReader(fullUrl)
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	if method == http.MethodPost {
		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		defer resp.Body.Close()
		return resp.StatusCode, string(respBody)
	} else {
		t.Log(resp)
		respBody := resp.Header.Get("Location")
		require.NoError(t, err)
		defer resp.Body.Close()
		return resp.StatusCode, respBody
	}
}

func TestRouter(t *testing.T) {
	fullUrl = "https://yandex.ru"
	r := Router()
	ts := httptest.NewServer(r)
	defer ts.Close()
	statusCode, body := testRequest(t, ts, "POST", "/")

	pat := regexp.MustCompile(`:\d{2,}/(\w+)`)
	if len(pat.FindSubmatch([]byte(body))) == 2 {
		shortUrl = "/" + string(pat.FindSubmatch([]byte(body))[1])
	}

	assert.Equal(t, http.StatusCreated, statusCode)
	statusCode, _ = testRequest(t, ts, "GET", shortUrl)

	assert.Equal(t, http.StatusOK, statusCode)

}
