package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type MyMap struct {
	idToOriginalURL map[string]string
}

// HandlePostRequest handles POST requests without path parameters, i.e. `POST /`,
// and does not support other HTTP methods
func (myMap *MyMap) HandlePostRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Requests other than POST are not allowed to `/`", http.StatusMethodNotAllowed)
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()
	originalURL, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error while reading request body: %v\n", err)
	}
	id := GenerateUniqueID(myMap)
	myMap.idToOriginalURL[string(id)] = string(originalURL)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(append([]byte("http://localhost:8080/"), id...))
	if err != nil {
		log.Printf("Error while writing response body: %v\n", err)
	}
}

// HandleGetRequest handles GET requests with a path parameter `id`, i.e. `GET /{id}`,
// and does not support other HTTP methods
func (myMap *MyMap) HandleGetRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Requests other than GET are not allowed to `/{id}`", http.StatusMethodNotAllowed)
		return
	}
	id := strings.TrimLeft(r.URL.Path, "/")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if originalURL, ok := myMap.idToOriginalURL[id]; ok {
		w.Header().Set("Location", originalURL)
	}
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	registerEndpoints()
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

func registerEndpoints() {
	router := mux.NewRouter()
	dict := &MyMap{make(map[string]string)}
	router.HandleFunc("/", dict.HandlePostRequest)
	router.HandleFunc("/{id}", dict.HandleGetRequest)
	http.Handle("/", router)
}

const (
	alphabet       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	alphabetLength = len(alphabet)
	shortURLLength = 8
)

func GenerateUniqueID(myMap *MyMap) []byte {
	id := make([]byte, 0, shortURLLength)
generateRandomSequence:
	for i := 0; i < shortURLLength; i++ {
		randomSymbol := alphabet[rand.Intn(alphabetLength)]
		id = append(id, randomSymbol)
	}
	if _, ok := myMap.idToOriginalURL[string(id)]; ok {
		goto generateRandomSequence
	}
	return id
}
