package service

import (
	"log"

	"github.com/tsupko/shortener/internal/app/storage"
	"github.com/tsupko/shortener/internal/app/util"
)

var _ ShorteningService = &ShorteningServiceImpl{}

type ShorteningServiceImpl struct {
	storage storage.Storage
}

func NewShorteningService(storage storage.Storage) *ShorteningServiceImpl {
	return &ShorteningServiceImpl{storage: storage}
}

func (s *ShorteningServiceImpl) Save(originalURL, userID string) (string, error) {
	shorteningIdentifier := s.generateShorteningIdentifier()
	log.Printf("storage: put original URL %s identified by its shortening ID %s\n", originalURL, shorteningIdentifier)
	return s.storage.Save(shorteningIdentifier, originalURL, userID)
}

func (s *ShorteningServiceImpl) SaveBatch(hashes, urls, userIds []string) ([]string, error) {
	return s.storage.SaveBatch(hashes, urls, userIds)
}

func (s *ShorteningServiceImpl) Get(shorteningIdentifier string) (storage.User, error) {
	originalURL, err := s.storage.Get(shorteningIdentifier)
	if err != nil {
		log.Printf("storage: got original URL %s identified by its shortening ID %s\n", originalURL, shorteningIdentifier)
	}
	return originalURL, err
}

func (s *ShorteningServiceImpl) GetAll(userID string) (map[string]string, error) {
	return s.storage.GetAll(userID)
}

func (s *ShorteningServiceImpl) generateShorteningIdentifier() string {
	hash := util.GenerateUniqueID()
	if _, err := s.storage.Get(hash); err != nil {
		return hash
	}
	log.Printf("hash %s already exists, generating a new one", hash)
	return s.generateShorteningIdentifier()
}
