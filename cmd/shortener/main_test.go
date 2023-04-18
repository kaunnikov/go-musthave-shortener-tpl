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
		t.Log(resp)
		respBody := resp.Header.Get("Location")
		require.NoError(t, err)
		defer resp.Body.Close()
		return resp.StatusCode, respBody
	}
}

func TestRouter(t *testing.T) {
	fullURL = "https://yandex.ru"
	r := NewRouter()
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
