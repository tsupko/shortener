package main

import (
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	slash         = "/"
	doubleSlashes = "//"
)

var idToOriginalURL = map[string]string{"id": "https://yandex.ru/id"}

// HandleURLWithoutPathParameter handles POST requests without path parameters, i.e. `POST /`,
// and does not support other HTTP methods
func HandleURLWithoutPathParameter(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(r.Body)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		var originalURL = body
		shortenedURL := shorten(originalURL)
		idToOriginalURL[string(shortenedURL)] = string(originalURL)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(shortenedURL)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		http.Error(w, "Requests other than POST are not allowed to `/`", http.StatusBadRequest)
	}
}

// HandleURLWithPathParameter handles GET requests with a path parameter `id`, i.e. `GET /{id}`,
// and does not support other HTTP methods
func HandleURLWithPathParameter(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		path := r.URL.Path
		slashIndex := strings.Index(path, slash)
		id := path[slashIndex+len(slash):]
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if originalURL, ok := idToOriginalURL[id]; ok {
			w.Header().Set("Location", originalURL)
		}
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		http.Error(w, "Requests other than GET are not allowed to `/{id}`", http.StatusBadRequest)
	}
}

func main() {
	registerEndpoints()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func registerEndpoints() {
	router := mux.NewRouter()
	router.HandleFunc("/", HandleURLWithoutPathParameter)
	router.HandleFunc("/{id}", HandleURLWithPathParameter)
	http.Handle("/", router)
}

func shorten(url []byte) []byte {
	stringURL := string(url)
	doubleSlashesIndex := strings.Index(stringURL, doubleSlashes)
	sliceURLWithoutProtocol := url[doubleSlashesIndex+len(doubleSlashes):]
	slashIndex := strings.Index(string(sliceURLWithoutProtocol), slash)
	sliceURLWithoutHostPort := sliceURLWithoutProtocol[slashIndex+len(slash):]
	return sliceURLWithoutHostPort
}
