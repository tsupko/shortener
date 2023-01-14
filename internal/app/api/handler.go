package api

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/tsupko/shortener/internal/app/db"
	"github.com/tsupko/shortener/internal/app/exceptions"
	"github.com/tsupko/shortener/internal/app/storage"
	"github.com/tsupko/shortener/internal/app/util"

	"github.com/tsupko/shortener/internal/app/service"
)

type RequestHandler struct {
	service  service.ShorteningService
	baseURL  string
	dbSource *db.Source
}

func NewRequestHandler(service service.ShorteningService, baseURL string, db *db.Source) *RequestHandler {
	return &RequestHandler{
		service:  service,
		baseURL:  baseURL,
		dbSource: db,
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
	userID := getUserID(r.Context())
	hash, err := h.service.Save(originalURL, userID)
	if errors.Is(err, exceptions.ErrURLAlreadyExist) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	shortURL := h.makeShortURL(hash)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err = w.Write([]byte(shortURL))
	if err != nil {
		log.Printf("error writing response body: %v\n", err)
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
	originalURL, err := h.service.Get(id)
	if err != nil {
		originalURL = storage.User{} // TODO if not found what to do?
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Location", originalURL.URL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	_, err = w.Write([]byte("redirect to " + originalURL.URL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *RequestHandler) handleJSONPost(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Printf("error closing request body: %v\n", err)
		}
	}()
	resBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading request body: %v\n", err)
		return
	}

	value := request{}
	if err := json.Unmarshal(resBody, &value); err != nil {
		http.Error(w, "Could not unmarshal request: "+err.Error(), http.StatusBadRequest)
		return
	}

	userID := getUserID(r.Context())

	w.Header().Set("Content-Type", "application/json")

	hash, err := h.service.Save(value.URL, userID)
	if errors.Is(err, exceptions.ErrURLAlreadyExist) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	if err != nil {
		log.Printf("error saving original URL: %v\n", err)
	}

	response := response{h.makeShortURL(hash)}
	responseString, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Could not marshal response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(responseString)
	if err != nil {
		log.Printf("error writing response: %v\n", err)
	}
}

func (h *RequestHandler) handleBatch(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Printf("error closing request body: %v\n", err)
		}
	}()
	resBody, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var batchRequests []BatchRequest
	if err := json.Unmarshal(resBody, &batchRequests); err != nil {
		log.Println("can not unmarshal body:[", string(resBody), "] ", err)
	}

	batchResponses := make([]BatchResponse, 0, len(batchRequests))

	userID := getUserID(r.Context())
	// TODO make it through batch SaveBatch
	for i := range batchRequests {
		hash, err := h.service.Save(batchRequests[i].OriginalURL, userID)
		if err != nil {
			log.Println("unexpected exceptions", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		batchResponse := BatchResponse{batchRequests[i].CorrelationID, h.makeShortURL(hash)}
		batchResponses = append(batchResponses, batchResponse)
	}

	responseString, err := json.Marshal(batchResponses)
	if err != nil {
		log.Println("can not marshal batchResponses:[", string(resBody), "] ", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(responseString)
	if err != nil {
		log.Printf("error writing response: %v\n", err)
	}
}

func (h *RequestHandler) handlePing(w http.ResponseWriter) {
	if h.dbSource == nil {
		log.Println("db ping error, db is not initialized")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := h.dbSource.Ping()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("db ping error:", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *RequestHandler) makeShortURL(id string) string {
	return h.baseURL + "/" + id
}

func (h *RequestHandler) getUserUrls(w http.ResponseWriter, r *http.Request) {
	userIDCookie, err := r.Cookie(util.UserIDCookieName)
	if err != nil {
		log.Println("error getting userID cookie:", err)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	userID, err2 := decodeID(userIDCookie.Value)
	if err2 != nil {
		log.Println("error decoding userID cookie:", err2)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	log.Println(strconv.Itoa(int(userID)) + " userID found")

	pairs, err3 := h.service.GetAll(strconv.Itoa(int(userID)))

	if err3 != nil {
		log.Println("error getting data from storage:", err3)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var list []listResponse

	for hash, value := range pairs {
		listResponse := listResponse{h.makeShortURL(hash), value}
		list = append(list, listResponse)
	}
	responseString, err := json.Marshal(list)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(responseString)
	if err != nil {
		panic(err)
	}
}

func getUserID(ctx context.Context) string {
	value := ctx.Value(UserIDContextKey)
	return value.(string)
}

type request struct {
	URL string `json:"url"`
}

type response struct {
	Result string `json:"result"`
}

type listResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type BatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
