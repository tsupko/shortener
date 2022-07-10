package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type URL struct {
	Value string `json:"url"`
}

const LettersToKeep = 5 // length of `://` plus the number of URL letters to keep in a shortened URL

var idToOriginalURL = map[string]URL{"1": {"https://yandex.ru"}} // map for storing pairs of `id` -> `original URL`

func handleURLWithoutPathParameter(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body := r.Body
		defer func(body io.ReadCloser) {
			err := body.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(body)
		content, _ := ioutil.ReadAll(body)
		uri := &URL{}
		err := json.Unmarshal(content, uri)
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated) // 201
		_, err = w.Write(shorten(uri))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		http.Error(w, "Requests other than POST are not allowed to `/`", http.StatusBadRequest)
	}
}

func shorten(uri *URL) []byte {
	trimUpToIndex := strings.Index(uri.Value, "://") + LettersToKeep
	dotIndex := strings.LastIndex(uri.Value, ".")
	u := []byte(uri.Value)
	return append(u[:trimUpToIndex], u[dotIndex:]...)
}

func handleURLWithPathParameter(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		path := r.URL.Path
		slashIndex := strings.Index(path, "/")
		id := path[slashIndex+1:]
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Location", idToOriginalURL[id].Value)
		w.WriteHeader(http.StatusTemporaryRedirect) // 307
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
	router.HandleFunc("/", handleURLWithoutPathParameter)
	router.HandleFunc("/{id}", handleURLWithPathParameter)
	http.Handle("/", router)
}
