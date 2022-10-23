package mocks

import (
	"log"

	"github.com/tsupko/shortener/internal/app/storage"
)

type MockStorage struct {
	requestCount int
}

func (m *MockStorage) GetAll() (interface{}, interface{}) {
	//TODO implement me
	panic("implement me")
}

var _ storage.Storage = &MockStorage{}

func (m *MockStorage) Put(id string, _ string) string {
	return id
}

func (m *MockStorage) Get(id string) (string, bool) {
	log.Default().Println("mock storage: got with id:", id)
	if m.requestCount > 0 {
		return "", false
	}
	m.requestCount++
	return "idExists", true
}
