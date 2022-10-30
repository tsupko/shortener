package storage

import (
	"log"
	"os"
	"path/filepath"
	"sync"
)

type FileStorage struct {
	data            map[string]string
	fileStoragePath string
	producer        *producer
	mtx             sync.RWMutex
}

func (s *FileStorage) GetAll() (interface{}, interface{}) {
	return s.data, nil
}

var _ Storage = &FileStorage{}

func NewFileStorage(fileStoragePath string) *FileStorage {
	checkDirExistOrCreate(fileStoragePath)
	mapStore := readFromFileIntoMap(fileStoragePath)

	fileProducer, err := NewProducer(fileStoragePath)
	if err != nil {
		log.Println("can not create NewFileStorage", err)
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
			log.Println("error accessing file system:", err)
		}
	}
}

func readFromFileIntoMap(fileStoragePath string) map[string]string {
	consumer, err := NewConsumer(fileStoragePath)
	if err != nil {
		log.Println("error reading file from disk:", err)
	}
	defer consumer.Close()

	mapStore := make(map[string]string)
	for {
		record, err := consumer.ReadRecord()
		if err != nil {
			break
		}
		mapStore[record.Hash] = record.URL
	}
	return mapStore
}

func (s *FileStorage) writeToFile(hash string, url string) {
	record := record{hash, url}
	err := s.producer.WriteRecord(&record)
	if err != nil {
		log.Println("error writing to file:", err)
	}
}
