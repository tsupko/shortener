package service

import (
	"log"

	"github.com/tsupko/shortener/internal/app/storage"
	"github.com/tsupko/shortener/internal/app/util"
)

type ShorteningService struct {
	storage storage.Storage
}

func NewShorteningService(storage storage.Storage) *ShorteningService {
	return &ShorteningService{storage: storage}
}

func (s *ShorteningService) Put(originalURL string) string {
	shorteningIdentifier := s.generateShorteningIdentifier()
	log.Printf("storage: put original URL %s identified by its shortening ID %s\n", originalURL, shorteningIdentifier)
	return s.storage.Put(shorteningIdentifier, originalURL)
}

func (s *ShorteningService) Get(shorteningIdentifier string) string {
	originalURL, _ := s.storage.Get(shorteningIdentifier)
	log.Printf("storage: got original URL %s identified by its shortening ID %s\n", originalURL, shorteningIdentifier)
	return originalURL
}

func (s *ShorteningService) generateShorteningIdentifier() string {
	id := util.GenerateUniqueID()
	if _, ok := s.storage.Get(id); !ok {
		return id
	}
	return s.generateShorteningIdentifier()
}
