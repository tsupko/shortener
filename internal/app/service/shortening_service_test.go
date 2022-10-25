package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tsupko/shortener/internal/app/storage"
	"github.com/tsupko/shortener/internal/app/storage/mocks"
)

func TestShorteningServicePutGet(t *testing.T) {
	s := NewShorteningService(storage.NewMemoryStorage())

	id := s.Put("https://ya.ru")
	assert.Len(t, id, 8)
	url := s.Get(id)
	assert.Equal(t, "https://ya.ru", url)
	assert.Equal(t, "", s.Get("idDoesNotExist"))
}

func TestShorteningServiceDuplicateID(t *testing.T) {
	s := NewShorteningService(&mocks.MockStorage{})
	id := s.Put("https://ya.ru")
	assert.Len(t, id, 8)
}
