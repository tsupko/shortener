package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleUrlWithoutPathParameter(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
		body        string
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "test `POST /` returns 201, body <short_URL>",
			request: "/",
			want: want{
				statusCode:  201,
				contentType: "text/plain; charset=utf-8",
				body:        "id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody := strings.NewReader("https://yandex.ru/id")
			request := httptest.NewRequest(http.MethodPost, tt.request, requestBody)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(HandleURLWithoutPathParameter)
			h.ServeHTTP(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
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

func TestHandleUrlWithPathParameter(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
		location    string
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "test `GET /{id}` returns 307, header `Location:<full_URL>`",
			request: "/id",
			want: want{
				statusCode:  307,
				contentType: "text/plain; charset=utf-8",
				location:    "https://yandex.ru/id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("", tt.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(HandleURLWithPathParameter)
			h.ServeHTTP(w, request)
			result := w.Result()
			err := result.Body.Close()
			if err != nil {
				t.Errorf("Error while closing response body: %v", err)
			}
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))
		})
	}
}
