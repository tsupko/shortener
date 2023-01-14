package service

import "github.com/tsupko/shortener/internal/app/storage"

type ShorteningService interface {
	Get(hash string) (storage.User, error)
	GetAll(userID string) (map[string]string, error)
	Save(url, userID string) (string, error)
	SaveBatch(hashes, urls, userIds []string) ([]string, error)
}
