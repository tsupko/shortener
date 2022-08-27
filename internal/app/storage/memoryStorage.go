package storage

import "sync"

type MemoryStorage struct {
	concurrentMap sync.Map
}

var _ Storage = &MemoryStorage{}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (s *MemoryStorage) Put(id string, originalURL string) string {
	s.concurrentMap.Store(id, originalURL)
	return id
}

func (s *MemoryStorage) Get(id string) (string, bool) {
	value, ok := s.concurrentMap.Load(id)
	originalURL := ""
	if ok {
		originalURL = value.(string)
	}
	return originalURL, ok
}
