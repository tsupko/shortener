package api

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tsupko/shortener/internal/app/service"
)

func TestHandlers(t *testing.T) {
	h := NewRequestHandler(
		service.NewMockShorteningService(),
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
				body: "redirect to ",
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
				body: "redirect to https://ya.ru",
			},
		},
		{
			name: "POST /",
			request: request{
				method: http.MethodPost,
				path:   "/",
				body:   "https://ya.ru",
			},
			want: want{
				statusCode: 201,
				headers: map[string][]string{
					"Content-Type": {""},
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
		{
			name: "POST existing URL",
			request: request{
				method: http.MethodPost,
				path:   "/api/shorten",
				body:   `{"url":"https://already.exist"}`,
			},
			want: want{
				statusCode: 409,
				headers: map[string][]string{
					"Content-Type": {"application/json"},
					"Location":     {""},
				},
				body: `{"result":"http://localhost:8080/urlAlreadyExistHash"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.request.method, tt.request.path, strings.NewReader(tt.request.body))
			request = request.WithContext(context.WithValue(context.Background(), UserIDContextKey, "123456789"))
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
				t.Errorf("error reading response body: %v", err)
			}
			err = result.Body.Close()
			if err != nil {
				t.Errorf("error closing response body: %v", err)
			}
			assert.Equal(t, tt.want.body, string(body))
		})
	}
}
