package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/caarlos0/env/v6"

	"github.com/tsupko/shortener/internal/app/api"
	"github.com/tsupko/shortener/internal/app/db"
	"github.com/tsupko/shortener/internal/app/service"
	"github.com/tsupko/shortener/internal/app/storage"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDsn     string `env:"DATABASE_DSN"`
}

func main() {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Println("error reading config file:", err)
	}

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "-a serverAddress")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "-b baseUrl")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "-f fileStoragePath")
	flag.StringVar(&cfg.DatabaseDsn, "d", cfg.DatabaseDsn, "-d DatabaseDsn")
	flag.Parse()

	log.Printf("config: %+v", cfg)

	var (
		store    storage.Storage
		dbSource *db.Source
	)
	switch {
	case cfg.DatabaseDsn != "":
		log.Println("init store as database based")
		dbSource, err = db.NewDB(cfg.DatabaseDsn)
		if err != nil {
			log.Fatal("failed to init dbSource: " + err.Error())
		}
		defer dbSource.Close()
		store = storage.NewDBStorage(dbSource)
	case cfg.FileStoragePath != "":
		log.Println("environment var FILE_STORAGE_PATH is found: " + cfg.FileStoragePath)
		log.Println("init store as file store based")
		store = storage.NewFileStorage(cfg.FileStoragePath)
	default:
		log.Println("init store as memory store based")
		store = storage.NewMemoryStorage()
	}

	shortService := service.NewShorteningService(store)

	handler := api.NewRequestHandler(shortService, cfg.BaseURL, dbSource)
	router := api.NewRouter(handler)

	server := &http.Server{Addr: cfg.ServerAddress, Handler: router}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := server.ListenAndServe(); err != nil {
			log.Println("listen and serve failed: " + err.Error())
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		ch := make(chan os.Signal, 1)
		signal.Notify(ch,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
		)

		s := <-ch

		log.Println("received signal: " + s.String())
		if err := server.Close(); err != nil {
			log.Println("close failed: " + err.Error())
		}
	}()

	wg.Wait()
}
