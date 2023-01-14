package service

import (
	"github.com/tsupko/shortener/internal/app/exceptions"
	"github.com/tsupko/shortener/internal/app/storage"
)

var _ ShorteningService = &MockShorteningService{}

type MockShorteningService struct{}

func NewMockShorteningService() *MockShorteningService {
	return &MockShorteningService{}
}

func (s *MockShorteningService) Get(hash string) (storage.User, error) {
	if hash == "12345" {
		return storage.User{URL: "https://ya.ru"}, nil
	}
	return storage.User{}, exceptions.ErrURLNotFound
}

func (s *MockShorteningService) GetAll(_ string) (map[string]string, error) {
	data := make(map[string]string)
	data["12345"] = "https://ya.ru"
	return data, nil
}

func (s *MockShorteningService) Save(url, _ string) (string, error) {
	if url == "https://ya.ru" {
		return "12345", nil
	}
	if url == "https://already.exist" {
		return "urlAlreadyExistHash", exceptions.ErrURLAlreadyExist
	}
	return "67890", nil
}

func (s *MockShorteningService) SaveBatch(hashes, _, _ []string) ([]string, error) {
	values := make([]string, 0, len(hashes))
	for i := range hashes {
		values = append(values, hashes[i])
	}
	return values, nil
}
