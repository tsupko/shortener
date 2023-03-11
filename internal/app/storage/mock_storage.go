package storage

import (
	"log"

	"github.com/tsupko/shortener/internal/app/exceptions"
)

type MockStorage struct {
	requestCount int
}

var _ Storage = &MockStorage{}

func (m *MockStorage) Save(hash, _, _ string) (string, error) {
	return hash, nil
}

func (m *MockStorage) SaveBatch(hashes, urls, userIds []string) ([]string, error) {
	values := make([]string, 0, len(hashes))
	for i := range hashes {
		hash, _ := m.Save(hashes[i], urls[i], userIds[i])
		values = append(values, hash)
	}
	return values, nil
}

func (m *MockStorage) Get(hash string) (User, error) {
	log.Default().Println("mock storage get with hash: ", hash)
	if m.requestCount > 0 {
		return User{"", ""}, exceptions.ErrURLNotFound
	}
	m.requestCount++
	return User{"hashExists", ""}, nil
}

func (m *MockStorage) GetAll(string) (map[string]string, error) {
	return make(map[string]string), nil
}
