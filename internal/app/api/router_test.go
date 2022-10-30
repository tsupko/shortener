package api

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tsupko/shortener/internal/app/service"
	"github.com/tsupko/shortener/internal/app/storage"
)

func TestGetEmpty(t *testing.T) {
	ts := getServer()
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/98765", "")

	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
	assert.Equal(t, "", resp.Header.Get("Location"))
	assert.Equal(t, "", body)
	closeBody(t, resp)
}

func TestGetPositive(t *testing.T) {
	ts := getServer()
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/12345", "")

	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
	assert.Equal(t, "https://ya.ru", resp.Header.Get("Location"))
	assert.Equal(t, "", body)
	closeBody(t, resp)
}

func TestPostPositive(t *testing.T) {
	ts := getServer()
	defer ts.Close()

	resp, body := testRequest(t, ts, "POST", "/", "https://ya.ru")

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "http://localhost:8080/12345", body)
	closeBody(t, resp)
}

func TestPostBadRequest(t *testing.T) {
	ts := getServer()
	defer ts.Close()

	resp, body := testRequest(t, ts, "POST", "/", "")

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "", body)
	closeBody(t, resp)
}

func TestRouterNotFound(t *testing.T) {
	ts := getServer()
	defer ts.Close()

	resp, body := testRequest(t, ts, "POST", "/1/2", "https://ya.ru")

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

func TestAcceptEncodingGzip(t *testing.T) {
	ts := getServer()
	defer ts.Close()

	resp, body := testRequest(t, ts, "POST", "/", "https://ya.ru", "Accept-Encoding", "gzip")

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "http://localhost:8080/12345", unzip(body))
	assert.Equal(t, "gzip", resp.Header.Get("Content-Encoding"))
	closeBody(t, resp)
}

func TestContentEncodingGzip(t *testing.T) {
	ts := getServer()
	defer ts.Close()

	requestBody := zip("https://ya.ru")
	resp, body := testRequest(t, ts, "POST", "/", requestBody, "Content-Encoding", "gzip")

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "http://localhost:8080/12345", body)
	closeBody(t, resp)
}

func TestPingDb(t *testing.T) {
	ts := getServer()
	defer ts.Close()

	resp, _ := testRequest(t, ts, "GET", "/ping", "")

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	closeBody(t, resp)
}

func getServer() *httptest.Server {
	r := NewRouter(NewRequestHandler(
		service.NewShorteningService(storage.NewTestStorage()),
		"http://localhost:8080",
		nil,
	))
	ts := httptest.NewServer(r)
	return ts
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body string, headers ...string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(body))
	if len(headers) == 2 {
		req.Header.Set(headers[0], headers[1])
	}
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

func closeBody(t *testing.T, resp *http.Response) {
	err := resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func zip(original string) string {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(original)); err != nil {
		log.Println(err)
	}
	if err := gz.Close(); err != nil {
		log.Println(err)
	}
	return b.String()
}

func unzip(original string) string {
	reader := bytes.NewReader([]byte(original))
	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		log.Println("error while unzip", err)
		return ""
	}
	output, err := io.ReadAll(gzReader)
	if err != nil {
		log.Println("error while unzip", err)
		return ""
	}
	return string(output)
}
