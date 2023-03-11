package storage

import (
	"strings"

	"github.com/tsupko/shortener/internal/app/db"
	"github.com/tsupko/shortener/internal/app/exceptions"
)

type DBStorage struct {
	dbSource *db.Source
}

var _ Storage = &DBStorage{}

func NewDBStorage(db *db.Source) *DBStorage {
	db.InitTables()
	return &DBStorage{dbSource: db}
}

func (s *DBStorage) Save(hash, url, _ string) (string, error) {
	err := s.dbSource.Save(hash, url)
	if err != nil {
		if isHashUniqueViolation(err) {
			return hash, exceptions.ErrHashAlreadyExist
		}
		if isURLUniqueViolation(err) {
			oldHash, err2 := s.dbSource.GetHashByURL(url)
			if err2 != nil {
				return "", err2 // somebody deleted URL form db?
			}
			return oldHash, exceptions.ErrURLAlreadyExist
		}
		return "", err // unexpected error
	}
	return hash, nil
}

func isHashUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "urls_pk")
}

func isURLUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "urls_url_uindex")
}

func (s *DBStorage) SaveBatch(hashes, urls, _ []string) ([]string, error) {
	// TODO there is no collision hashed check
	err := s.dbSource.SaveBatch(hashes, urls)
	if err != nil {
		return nil, err
	}
	return hashes, nil
}

func (s *DBStorage) Get(hash string) (User, error) {
	url, err := s.dbSource.Get(hash)
	return User{url, ""}, err
}

func (s *DBStorage) GetAll(string) (map[string]string, error) {
	return s.dbSource.GetAll()
}
