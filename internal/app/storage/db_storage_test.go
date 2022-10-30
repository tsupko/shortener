package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tsupko/shortener/internal/app/db"
	"github.com/tsupko/shortener/internal/app/exceptions"
	"github.com/tsupko/shortener/internal/app/util"
)

func TestDBStorage(t *testing.T) {
	storage := initDBStorage(t)

	generatedHash1 := util.GenerateUniqueID()
	url1 := "https://" + generatedHash1
	savedHash1, err1 := storage.Save(generatedHash1, url1)
	assert.Nil(t, err1)
	assert.Equal(t, generatedHash1, savedHash1)

	// save new url with same hash
	url2 := "https://" + util.GenerateUniqueID()
	savedHash2, err2 := storage.Save(generatedHash1, url2)
	assert.NotNil(t, err2)
	assert.Error(t, exceptions.ErrHashAlreadyExist, err2)
	assert.Equal(t, generatedHash1, savedHash2)

	// save same url with new hash
	generatedHash3 := util.GenerateUniqueID()
	savedHash3, err3 := storage.Save(generatedHash3, url1)
	assert.NotNil(t, err3)
	assert.Error(t, exceptions.ErrURLAlreadyExist, err3)
	assert.NotEqual(t, generatedHash1, generatedHash3)
	assert.Equal(t, generatedHash1, savedHash3)

	savedURL1, err := storage.Get(savedHash1)
	assert.Nil(t, err)
	assert.Equal(t, url1, savedURL1) // old value

	data, err4 := storage.GetAll()
	assert.Nil(t, err4)
	assert.GreaterOrEqual(t, len(data), 1)

	urlFromMap := data[generatedHash1]
	assert.Equal(t, url1, urlFromMap)
}

func TestDBStorageSaveBatch(t *testing.T) {
	storage := initDBStorage(t)

	hash1 := util.GenerateUniqueID()
	hash2 := util.GenerateUniqueID()
	url1 := "https://" + hash1
	url2 := "https://" + hash2
	hashes := []string{hash1, hash2}
	urls := []string{url1, url2}

	savedHashes, err := storage.SaveBatch(hashes, urls)
	assert.Nil(t, err)
	assert.Equal(t, hashes, savedHashes)

	savedURL1, err1 := storage.Get(hash1)
	assert.Nil(t, err1)
	assert.Equal(t, url1, savedURL1)

	savedURL2, err2 := storage.Get(hash2)
	assert.Nil(t, err2)
	assert.Equal(t, url2, savedURL2)
}

func initDBStorage(t *testing.T) *DBStorage {
	dbSource, err := db.NewDB("postgres://shortener:pass@localhost:5432/shortener")
	if err != nil {
		t.Error("exceptions when NewDB", err)
	}

	err = dbSource.Ping()
	if err != nil {
		t.Skip("no db connection")
	}
	storage := NewDBStorage(dbSource)
	return storage
}
