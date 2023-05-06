package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io"
	"kaunnikov/go-musthave-shortener-tpl/config"
	"kaunnikov/go-musthave-shortener-tpl/internal/app"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

var fullURL string
var shortURL string

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (int, string) {
	body := strings.NewReader(fullURL)
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	if err != nil {
		log.Println(err)
	}
	if method == http.MethodPost {
		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		defer resp.Body.Close()
		return resp.StatusCode, string(respBody)
	} else {
		respBody := resp.Header.Get("Location")
		require.NoError(t, err)
		defer resp.Body.Close()
		return resp.StatusCode, respBody
	}
}

func TestRouter(t *testing.T) {
	fullURL = "https://yandex.ru"

	cfg := &config.AppConfig{Host: ":8080", ResultURL: ":8080", FileStoragePath: "/tmp/short-url-db.json"}

	loadFromENV(cfg)

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("logger don't Run! %s", err)
	}
	app.Sugar = logger.Sugar()

	newApp := app.NewApp(cfg)

	r := chi.NewRouter()
	r.Post("/", newApp.CreateHandler)
	r.Get("/{id}", newApp.ShortHandler)
	r.Post("/api/shorten", newApp.JSONHandler)

	ts := httptest.NewServer(r)
	defer ts.Close()
	statusCode, body := testRequest(t, ts, "POST", "/")

	pat := regexp.MustCompile(`:\d{2,}/(\w+)`)
	if len(pat.FindSubmatch([]byte(body))) == 2 {
		shortURL = "/" + string(pat.FindSubmatch([]byte(body))[1])
	}

	assert.Equal(t, http.StatusCreated, statusCode)
	statusCode, _ = testRequest(t, ts, "GET", shortURL)

	assert.Equal(t, http.StatusOK, statusCode)

}
