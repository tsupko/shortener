package storage

import (
	"github.com/tsupko/shortener/internal/app/db"
)

type DBStorage struct {
	dbSource *db.Source
}

func NewDBStorage(db *db.Source) *DBStorage {
	db.InitTables()
	return &DBStorage{dbSource: db}
}

func (s *DBStorage) Put(hash string, url string) string {
	_, ok := s.dbSource.Get(hash)
	if ok {
		return hash
	}
	s.dbSource.Save(hash, url)
	return hash
}

func (s *DBStorage) Get(hash string) (string, bool) {
	return s.dbSource.Get(hash)

}

func (s *DBStorage) GetAll() (interface{}, interface{}) {
	urls := s.dbSource.GetAll()
	return urls, nil
}
