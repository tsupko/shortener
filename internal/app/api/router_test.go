package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tsupko/shortener/internal/app/service"
	"github.com/tsupko/shortener/internal/app/storage"
	"github.com/tsupko/shortener/internal/app/util"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(body))
	require.NoError(t, err)

	client := http.DefaultClient
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(resp.Body)

	return resp, string(respBody)
}

func TestRouter(t *testing.T) {
	r := NewRouter(NewRequestHandler(service.NewShorteningService(storage.NewTestStorage()), util.ServerAddress))
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, _ := testRequest(t, ts, "GET", "/12345", "")
	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
	assert.Equal(t, "https://ya.ru", resp.Header.Get("Location"))
	closeBody(t, resp)

	resp, _ = testRequest(t, ts, "GET", "/98765", "")
	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
	assert.Equal(t, "", resp.Header.Get("Location"))
	closeBody(t, resp)

	resp, body := testRequest(t, ts, "POST", "/", "https://ya.ru")
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "http://localhost:8080/12345", body)
	closeBody(t, resp)

	resp, body = testRequest(t, ts, "POST", "/", "")
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "", body)
	closeBody(t, resp)

	resp, body = testRequest(t, ts, "POST", "/1/2", "https://ya.ru")
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, "404 page not found\n", body)
	closeBody(t, resp)
}

func TestPostApi(t *testing.T) {
	ts := getServer()
	defer ts.Close()

	resp, body := testRequest(t, ts, "POST", "/api/shorten", `{"url":"https://ya.ru"}`)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, `{"result":"http://localhost:8080/12345"}`, body)
	closeBody(t, resp)
}

func getServer() *httptest.Server {
	r := NewRouter(NewRequestHandler(service.NewShorteningService(storage.NewTestStorage()), util.ServerAddress))
	ts := httptest.NewServer(r)
	return ts
}

func closeBody(t *testing.T, resp *http.Response) {
	err := resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
}
