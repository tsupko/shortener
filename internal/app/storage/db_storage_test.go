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
	userID1 := "user1"
	savedHash1, err1 := storage.Save(generatedHash1, url1, userID1)
	assert.Nil(t, err1)
	assert.Equal(t, generatedHash1, savedHash1)

	// save new URL with same hash
	url2 := "https://" + util.GenerateUniqueID()
	userID2 := "user2"
	savedHash2, err2 := storage.Save(generatedHash1, url2, userID2)
	assert.NotNil(t, err2)
	assert.Error(t, exceptions.ErrHashAlreadyExist, err2)
	assert.Equal(t, generatedHash1, savedHash2)

	// save same URL with new hash
	generatedHash3 := util.GenerateUniqueID()
	userID3 := "user3"
	savedHash3, err3 := storage.Save(generatedHash3, url1, userID3)
	assert.NotNil(t, err3)
	assert.Error(t, exceptions.ErrURLAlreadyExist, err3)
	assert.NotEqual(t, generatedHash1, generatedHash3)
	assert.Equal(t, generatedHash1, savedHash3)

	savedURL1, err := storage.Get(savedHash1)
	assert.Nil(t, err)
	assert.Equal(t, url1, savedURL1) // old value

	data, err4 := storage.GetAll(userID1)
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
	userID1 := "user1"
	userID2 := "user2"
	hashes := []string{hash1, hash2}
	urls := []string{url1, url2}
	userIds := []string{userID1, userID2}

	savedHashes, err := storage.SaveBatch(hashes, urls, userIds)
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
