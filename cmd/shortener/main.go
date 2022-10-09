package main

import (
	"log"
	"net/http"

	"github.com/caarlos0/env/v6"

	"github.com/tsupko/shortener/internal/app/api"
	"github.com/tsupko/shortener/internal/app/service"
	"github.com/tsupko/shortener/internal/app/storage"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

func main() {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	memoryStorage := storage.NewMemoryStorage()
	shorteningService := service.NewShorteningService(memoryStorage)
	handler := api.NewRequestHandler(shorteningService, cfg.BaseURL)
	router := api.NewRouter(handler)
	log.Fatalln(http.ListenAndServe(cfg.ServerAddress, router))
}
