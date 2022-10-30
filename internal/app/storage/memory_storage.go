package storage

import (
	"sync"

	"github.com/tsupko/shortener/internal/app/exceptions"
)

type MemoryStorage struct {
	data map[string]string
	mtx  sync.RWMutex
}

var _ Storage = &MemoryStorage{}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{data: make(map[string]string)}
}

func (s *MemoryStorage) Save(hash string, url string) (string, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	s.data[hash] = url
	return hash, nil
}

func (s *MemoryStorage) SaveBatch(hashes []string, urls []string) ([]string, error) {
	values := make([]string, 0, len(hashes))
	for i := range hashes {
		hash, _ := s.Save(hashes[i], urls[i])
		values = append(values, hash)
	}
	return values, nil
}

func (s *MemoryStorage) Get(hash string) (string, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	value, ok := s.data[hash]
	if ok {
		return value, nil
	}
	return value, exceptions.ErrURLNotFound
}

func (s *MemoryStorage) GetAll() (map[string]string, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.data, nil
}
