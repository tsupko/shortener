package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tsupko/shortener/internal/app/exceptions"
	"github.com/tsupko/shortener/internal/app/storage"
)

func TestShortServiceSaveGet(t *testing.T) {
	s := NewShorteningService(storage.NewMemoryStorage())

	hash, err := s.Save("https://ya.ru", "")
	assert.Nil(t, err)
	assert.Len(t, hash, 8)
	url, err2 := s.Get(hash)
	assert.Equal(t, storage.User{URL: "https://ya.ru"}, url)
	assert.Nil(t, err2)

	url2, err3 := s.Get("hashDoesNotExist")
	assert.Equal(t, storage.User{URL: ""}, url2)
	assert.NotNil(t, err3)
	assert.Equal(t, exceptions.ErrURLNotFound, err3)
}
