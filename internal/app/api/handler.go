package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/tsupko/shortener/internal/app/util"

	"github.com/tsupko/shortener/internal/app/service"
)

type RequestHandler struct {
	storage service.ShorteningService
}

func NewRequestHandler(storage *service.ShorteningService) *RequestHandler {
	return &RequestHandler{*storage}
}

// handlePostRequest handles POST requests without path parameters, i.e. `POST /`,
// and does not support other HTTP methods
func (h *RequestHandler) handlePostRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Requests other than POST are not allowed to `/`", http.StatusMethodNotAllowed)
		return
	}

	originalURL, _ := util.ReadRequestBody(r)
	if len(originalURL) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := h.storage.Put(originalURL)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	shortenedURL := util.ServerAddress + id
	_, err := w.Write([]byte(shortenedURL))
	if err != nil {
		log.Printf("Error while writing body body: %v\n", err)
	}
}

// handleGetRequest handles GET requests with a path parameter `id`, i.e. `GET /{id}`,
// and does not support other HTTP methods
func (h *RequestHandler) handleGetRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Requests other than GET are not allowed to `/{id}`", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimLeft(r.URL.Path, "/")
	originalURL := h.storage.Get(id)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
