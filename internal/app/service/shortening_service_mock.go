package service

import (
	"github.com/tsupko/shortener/internal/app/exceptions"
)

var _ ShorteningService = &MockShorteningService{}

type MockShorteningService struct{}

func NewMockShorteningService() *MockShorteningService {
	return &MockShorteningService{}
}

func (s *MockShorteningService) Get(hash string) (string, error) {
	if hash == "12345" {
		return "https://ya.ru", nil
	}
	return "", exceptions.ErrURLNotFound
}

func (s *MockShorteningService) GetAll() (map[string]string, error) {
	data := make(map[string]string)
	data["12345"] = "https://ya.ru"
	return data, nil
}

func (s *MockShorteningService) Save(url string) (string, error) {
	if url == "https://ya.ru" {
		return "12345", nil
	}
	if url == "https://already.exist" {
		return "urlAlreadyExistHash", exceptions.ErrURLAlreadyExist
	}
	return "67890", nil
}

func (s *MockShorteningService) SaveBatch(hashes []string, _ []string) ([]string, error) {
	values := make([]string, 0, len(hashes))
	for i := range hashes {
		values = append(values, hashes[i])
	}
	return values, nil
}
