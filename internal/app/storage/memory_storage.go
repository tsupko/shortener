package storage

import (
	"sync"

	"github.com/tsupko/shortener/internal/app/exceptions"
)

type User struct {
	UserID string
	URL    string
}

type MemoryStorage struct {
	data map[string]User
	mtx  sync.RWMutex
}

var _ Storage = &MemoryStorage{}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{data: make(map[string]User)}
}

func (s *MemoryStorage) Save(hash, url, userID string) (string, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	s.data[hash] = User{
		UserID: userID,
		URL:    url,
	}
	return hash, nil
}

func (s *MemoryStorage) SaveBatch(hashes, urls, userIds []string) ([]string, error) {
	values := make([]string, 0, len(hashes))
	for i := range hashes {
		hash, _ := s.Save(hashes[i], urls[i], userIds[i])
		values = append(values, hash)
	}
	return values, nil
}

func (s *MemoryStorage) Get(hash string) (User, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	value, ok := s.data[hash]
	if ok {
		return value, nil
	}
	return value, exceptions.ErrURLNotFound
}

func (s *MemoryStorage) GetAll(userID string) (map[string]string, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	m := make(map[string]string)
	for hash, item := range s.data {
		if item.UserID == userID {
			m[hash] = item.URL
		}
	}
	return m, nil
}
