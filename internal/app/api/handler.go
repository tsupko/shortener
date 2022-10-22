package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/tsupko/shortener/internal/app/util"

	"github.com/tsupko/shortener/internal/app/service"
)

type RequestHandler struct {
	service service.ShorteningService
	baseURL string
}

func NewRequestHandler(service *service.ShorteningService, baseURL string) *RequestHandler {
	return &RequestHandler{
		service: *service,
		baseURL: baseURL,
	}
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
	id := h.service.Put(originalURL)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, err := w.Write([]byte(h.makeShortURL(id)))
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
	originalURL := h.service.Get(id)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *RequestHandler) handleJSONPost(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Printf("Error while closing request body: %v\n", err)
		}
	}()
	resBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error while reading request body: %v\n", err)
		return
	}

	value := request{}
	if err := json.Unmarshal(resBody, &value); err != nil {
		http.Error(w, "Could not unmarshal request: "+err.Error(), http.StatusBadRequest)
		return
	}
	hash := h.service.Put(value.URL)
	response := response{h.makeShortURL(hash)}
	responseString, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Could not marshal response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(responseString)
	if err != nil {
		log.Printf("Error while writing response: %v\n", err)
	}
}

func (h *RequestHandler) makeShortURL(id string) string {
	return h.baseURL + "/" + id
}

type request struct {
	URL string `json:"url"`
}

type response struct {
	Result string `json:"result"`
}
