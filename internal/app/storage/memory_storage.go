package storage

import (
	"fmt"
	"sync"
)

type MemoryStorage struct {
	concurrentMap sync.Map
}

func (s *MemoryStorage) GetAll() (interface{}, interface{}) {
	regularMap := map[string]interface{}{}
	s.concurrentMap.Range(func(key, value interface{}) bool {
		regularMap[fmt.Sprint(key)] = value
		return true
	})
	return regularMap, nil
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
