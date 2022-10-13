package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/caarlos0/env/v6"

	"github.com/tsupko/shortener/internal/app/api"
	"github.com/tsupko/shortener/internal/app/service"
	"github.com/tsupko/shortener/internal/app/storage"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func main() {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Printf("error while parsing environment variables: %s\n", err)
	}

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "-a serverAddress")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "-b baseUrl")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "-f fileStoragePath")
	flag.Parse()

	var store storage.Storage
	if cfg.FileStoragePath != "" {
		log.Printf("environment variable `FILE_STORAGE_PATH` is found: %s\n", cfg.FileStoragePath)
		store = storage.NewFileStorage(cfg.FileStoragePath)
	} else {
		store = storage.NewMemoryStorage()
	}
	shorteningService := service.NewShorteningService(store)
	handler := api.NewRequestHandler(shorteningService, cfg.BaseURL)
	router := api.NewRouter(handler)
	err = http.ListenAndServe(cfg.ServerAddress, router)
	if err != nil {
		log.Printf("server returned error: %s\n", err)
	}
}
