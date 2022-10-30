package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tsupko/shortener/internal/app/service"
	"github.com/tsupko/shortener/internal/app/storage"
)

func TestHandlers(t *testing.T) {
	h := NewRequestHandler(
		service.NewShorteningService(storage.NewTestStorage()),
		"http://localhost:8080",
		nil,
	)
	type request struct {
		method string
		path   string
		body   string
	}
	type want struct {
		statusCode int
		headers    http.Header
		body       string
	}
	type test struct {
		name    string
		request request
		want    want
	}
	tests := []test{
		{
			name: "GET non-existing value",
			request: request{
				method: http.MethodGet,
				path:   "/",
				body:   "",
			},
			want: want{
				statusCode: 307,
				headers: map[string][]string{
					"Content-Type": {"text/plain; charset=utf-8"},
					"Location":     {""},
				},
				body: "",
			},
		},
		{
			name: "GET existing value",
			request: request{
				method: http.MethodGet,
				path:   "/12345",
				body:   "",
			},
			want: want{
				statusCode: 307,
				headers: map[string][]string{
					"Content-Type": {"text/plain; charset=utf-8"},
					"Location":     {"https://ya.ru"},
				},
				body: "",
			},
		},
		{
			name: "POST /",
			request: request{
				method: http.MethodPost,
				path:   "/",
				body:   "https://github.com/tsupko/shortener/runs/7826862296?check_suite_focus=true",
			},
			want: want{
				statusCode: 201,
				headers: map[string][]string{
					"Content-Type": {"text/plain; charset=utf-8"},
					"Location":     {""},
				},
				body: "http://localhost:8080/12345",
			},
		},
		{
			name: "POST url",
			request: request{
				method: http.MethodPost,
				path:   "/api/shorten",
				body:   `{"url":"https://ya.ru"}`,
			},
			want: want{
				statusCode: 201,
				headers: map[string][]string{
					"Content-Type": {"application/json"},
					"Location":     {""},
				},
				body: `{"result":"http://localhost:8080/12345"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.request.method, tt.request.path, strings.NewReader(tt.request.body))

			w := httptest.NewRecorder()
			switch tt.request.method {
			case http.MethodGet:
				h.handleGetRequest(w, request)
			case http.MethodPost:
				if tt.request.path == "/" {
					h.handlePostRequest(w, request)
				} else if tt.request.path == "/api/shorten" {
					h.handleJSONPost(w, request)
				} else {
					t.Fatalf("unexpected path: %s", tt.request.path)
				}
			}

			result := w.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.headers.Get("Content-Type"), result.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.headers.Get("Location"), result.Header.Get("Location"))
			body, err := io.ReadAll(result.Body)
			if err != nil {
				t.Errorf("Error while reading response body: %v", err)
			}
			err = result.Body.Close()
			if err != nil {
				t.Errorf("Error while closing response body: %v", err)
			}
			assert.Equal(t, tt.want.body, string(body))
		})
	}
}
