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

func (s *ShorteningService) Save(originalURL string) (string, error) {
	shorteningIdentifier := s.generateShorteningIdentifier()
	log.Printf("storage: put original URL %s identified by its shortening ID %s\n", originalURL, shorteningIdentifier)
	return s.storage.Save(shorteningIdentifier, originalURL)
}

func (s *ShorteningService) Get(shorteningIdentifier string) (string, error) {
	originalURL, err := s.storage.Get(shorteningIdentifier)
	if err != nil {
		log.Printf("storage: got original URL %s identified by its shortening ID %s\n", originalURL, shorteningIdentifier)
	}
	return originalURL, err
}

func (s *ShorteningService) GetAll() (interface{}, interface{}) {
	return s.storage.GetAll()
}

func (s *ShorteningService) generateShorteningIdentifier() string {
	hash := util.GenerateUniqueID()
	if _, err := s.storage.Get(hash); err != nil {
		return hash
	}
	log.Printf("hash %s already exists, generating a new one", hash)
	return s.generateShorteningIdentifier()
}
