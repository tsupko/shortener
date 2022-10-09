package storage

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type FileStorage struct {
	data            map[string]string
	fileStoragePath string
	producer        *Producer
	mtx             sync.RWMutex
}

var _ Storage = &FileStorage{}

func NewFileStorage(fileStoragePath string) *FileStorage {
	checkDirExistOrCreate(fileStoragePath)
	mapStore := readFromFileIntoMap(fileStoragePath)

	fileProducer, err := NewProducer(fileStoragePath)
	if err != nil {
		log.Fatal(err)
	}
	return &FileStorage{data: mapStore, fileStoragePath: fileStoragePath, producer: fileProducer}
}

func (s *FileStorage) Put(hash string, url string) string {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.data[hash] = url
	s.writeToFile(hash, url)
	return hash
}

func (s *FileStorage) Get(hash string) (string, bool) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	value, ok := s.data[hash]
	return value, ok
}

func checkDirExistOrCreate(fileStoragePath string) {
	dir, _ := filepath.Split(fileStoragePath)
	if dir == "" {
		return
	}
	if _, err := os.Stat(fileStoragePath); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0700)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func readFromFileIntoMap(fileStoragePath string) map[string]string {
	consumer, err := NewConsumer(fileStoragePath)
	if err != nil {
		log.Fatal(err)
	}
	defer func(consumer *Consumer) {
		err := consumer.Close()
		if err != nil {
			fmt.Printf("Error while closing consumer: %v\n", err)
		}
	}(consumer)

	mapStore := make(map[string]string)
	for i := 0; ; i++ {
		record, err := consumer.ReadRecord()
		if err != nil {
			break
		}
		mapStore[record.Hash] = record.URL
	}
	return mapStore
}

func (s *FileStorage) writeToFile(hash string, url string) {
	record := Record{hash, url}
	err := s.producer.WriteRecord(&record)
	if err != nil {
		log.Fatal(err)
	}
}
