package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tsupko/shortener/internal/app/api/mocks"
)

func Test_handlePostRequest(t *testing.T) {
	type fields struct {
		myMap   *MyMap
		service *mocks.Shorten
	}
	type want struct {
		statusCode  int
		contentType string
		body        string
	}
	tests := []struct {
		name    string
		fields  fields
		args    string
		request string
		want    want
	}{
		{
			name: "test `POST /` returns 201, body <short_URL>",
			fields: fields{
				&MyMap{idToOriginalURL: map[string]string{"c2WD8F2q": "https://yandex.ru/really-long-link.htm#abc"}},
				&mocks.Shorten{},
			},
			args:    "https://yandex.ru/really-long-link.htm#abc",
			request: "/",
			want: want{
				statusCode:  201,
				contentType: "text/plain; charset=utf-8",
				body:        "http://localhost:8080/BpLnfgDs",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// todo mock correctly
			tt.fields.service.On("GenerateUniqueID", tt.fields.myMap).Return("BpLnfgDs")
			h := http.HandlerFunc(tt.fields.myMap.HandlePostRequest)
			requestBody := strings.NewReader(tt.args)
			request := httptest.NewRequest(http.MethodPost, tt.request, requestBody)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			body, err := io.ReadAll(result.Body)
			if err != nil {
				t.Errorf("Error while reading response body: %v\n", err)
			}
			err = result.Body.Close()
			if err != nil {
				t.Errorf("Error while closing response body: %v\n", err)
			}
			assert.Equal(t, tt.want.body, string(body))
		})
	}
}

func Test_handleGetRequest(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
		location    string
	}
	tests := []struct {
		name   string
		fields MyMap
		args   string
		want   want
	}{
		{
			name:   "test `GET /xIJw3Sw1` returns 307, header `Location: https://yandex.ru/really-long-link.htm#abc`",
			fields: MyMap{map[string]string{"xIJw3Sw1": "https://yandex.ru/really-long-link.htm#abc"}},
			args:   "/xIJw3Sw1",
			want: want{
				statusCode:  307,
				contentType: "text/plain; charset=utf-8",
				location:    "https://yandex.ru/really-long-link.htm#abc",
			},
		},
		{
			name:   "test `GET /w92OWjw2` returns 307, header `Location: https://google.com/search?item=abc#def`",
			fields: MyMap{map[string]string{"w92OWjw2": "https://google.com/search?item=abc#def`"}},
			args:   "/w92OWjw2",
			want: want{
				statusCode:  307,
				contentType: "text/plain; charset=utf-8",
				location:    "https://google.com/search?item=abc#def`",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.HandlerFunc(tt.fields.HandleGetRequest)
			w := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, tt.args, nil)
			h.ServeHTTP(w, request)
			result := w.Result()
			err := result.Body.Close()
			if err != nil {
				t.Errorf("Error while closing response body: %v\n", err)
			}
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))
		})
	}
}

func TestRouter(t *testing.T) {
	r := NewRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/id")
	defer func() {
		_ = resp.Body.Close()
	}()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "id: id", body)

	resp, body = testRequest(t, ts, "GET", "/other-id")
	defer func() {
		_ = resp.Body.Close()
	}()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "id: other-id", body)
}

func NewRouter() chi.Router {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			_, err := w.Write([]byte("id: " + id))
			if err != nil {
				log.Printf("Error while writing response body: %v\n", err)
			}
		})
	})
	return r
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer func() {
		_ = resp.Body.Close()
	}()

	return resp, string(respBody)
}
