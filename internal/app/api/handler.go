package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/tsupko/shortener/internal/app/db"
	"github.com/tsupko/shortener/internal/app/util"

	"github.com/tsupko/shortener/internal/app/service"
)

type RequestHandler struct {
	service   service.ShorteningService
	baseURL   string
	dbService *db.DB
}

func NewRequestHandler(service *service.ShorteningService, baseURL string, db *db.DB) *RequestHandler {
	return &RequestHandler{
		service:   *service,
		baseURL:   baseURL,
		dbService: db,
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

func (h *RequestHandler) handlePing(w http.ResponseWriter, r *http.Request) {
	log.Println("db ping, DatabaseDsn:" + h.dbService.DatabaseDsn)
	err := h.dbService.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("db ping error:", err)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *RequestHandler) makeShortURL(id string) string {
	return h.baseURL + "/" + id
}

func (h *RequestHandler) getUserUrls(w http.ResponseWriter, r *http.Request) {
	userIDCookie, err := r.Cookie(util.UserIDCookieName)
	if err != nil {
		log.Println("Error while getting userID cookie:", err)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	userID, err2 := decodeID(userIDCookie.Value)
	if err2 != nil {
		log.Println("Error while decoding userID cookie:", err2)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	log.Println(strconv.Itoa(int(userID)) + " userID found")

	pairs, err3 := h.service.GetAll()

	if err3 != nil {
		log.Println("Error while getting data from storage:", err3)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var list []listResponse

	for hash, value := range pairs.(map[string]string) {
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
