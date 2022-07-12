package main

import (
	"log"
	"net/http"

	"github.com/tsupko/shortener/internal/app/api"
	"github.com/tsupko/shortener/internal/app/service"
	"github.com/tsupko/shortener/internal/app/storage"
)

func main() {
	memoryStorage := storage.NewMemoryStorage()
	shorteningService := service.NewShorteningService(memoryStorage)
	handler := api.NewRequestHandler(shorteningService)
	router := api.NewRouter(handler)
	log.Fatalln(http.ListenAndServe(":8080", router))
}
