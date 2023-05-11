package app_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"kaunnikov/go-musthave-shortener-tpl/internal/app"
	"kaunnikov/go-musthave-shortener-tpl/internal/config"
	"kaunnikov/go-musthave-shortener-tpl/internal/logging"
	"kaunnikov/go-musthave-shortener-tpl/internal/storage/fs"
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

	cfg := config.LoadConfig()
	cfg.FileStoragePath = "/tmp/short-url-bd.json"

	if err := logging.Init(); err != nil {
		log.Fatalf("logger don't Run!: %s", err)
	}

	fs.Init(cfg)
	newApp := app.NewApp(cfg)

	ts := httptest.NewServer(newApp)
	defer ts.Close()

	statusCode, body := testRequest(t, ts, "POST", "/")
	assert.Equal(t, http.StatusCreated, statusCode)

	pat := regexp.MustCompile(`:\d{2,}/(\w+)`)
	if len(pat.FindSubmatch([]byte(body))) == 2 {
		shortURL = "/" + string(pat.FindSubmatch([]byte(body))[1])
	}

	statusCode, _ = testRequest(t, ts, "GET", shortURL)

	assert.Equal(t, http.StatusOK, statusCode)

}
