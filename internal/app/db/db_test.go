package db

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tsupko/shortener/internal/app/exceptions"
	"github.com/tsupko/shortener/internal/app/util"
)

// launch docker before run tests otherwise most tests will be ignored
// docker run --name postgresql -e POSTGRES_USER=shortener -e POSTGRES_PASSWORD=pass -p 5432:5432 -d postgres

func TestPingWrongDBPort(t *testing.T) {
	db, err := NewDB("postgres://shortener:pass@localhost:5433/shortener")
	require.NoError(t, err)
	err = db.Ping()
	require.Error(t, err)
}

func TestPingOk(t *testing.T) {
	db := initDBStorage(t)

	err := db.Ping()
	require.NoError(t, err)
}

func TestSaveGet(t *testing.T) {
	db := initDBStorage(t)

	hash := util.GenerateUniqueID()
	url := "https://" + util.GenerateUniqueID()

	assert.Nil(t, db.Save(hash, url))

	urlFromDB, err := db.Get(hash)
	assert.Nil(t, err)
	assert.Equal(t, url, urlFromDB)
}

func TestGetEmpty(t *testing.T) {
	db := initDBStorage(t)

	hash := util.GenerateUniqueID()

	urlFromDB, err := db.Get(hash)
	assert.NotNil(t, err)
	assert.Equal(t, exceptions.ErrURLNotFound, err)
	assert.Equal(t, "", urlFromDB)
}

func TestSaveSameHash(t *testing.T) {
	db := initDBStorage(t)

	hash := util.GenerateUniqueID()
	url1 := "https://" + util.GenerateUniqueID()
	url2 := "https://" + util.GenerateUniqueID()

	err1 := db.Save(hash, url1)
	assert.Nil(t, err1)

	err2 := db.Save(hash, url2)
	assert.NotNil(t, err2)
	assert.Contains(t, err2.Error(), "urls_pk")
}

func TestSaveSameUrl(t *testing.T) {
	db := initDBStorage(t)

	hash := util.GenerateUniqueID()
	url1 := "https://" + util.GenerateUniqueID()

	err1 := db.Save(hash, url1)
	assert.Nil(t, err1)

	hash2 := util.GenerateUniqueID()

	err2 := db.Save(hash2, url1)
	assert.NotNil(t, err2)
	assert.Contains(t, err2.Error(), "urls_url_uindex")
}

func TestSave_GetHashByURL(t *testing.T) {
	db := initDBStorage(t)

	hash1 := util.GenerateUniqueID()
	url := "https://" + util.GenerateUniqueID()

	err1 := db.Save(hash1, url)
	assert.Nil(t, err1)

	hash2, err2 := db.GetHashByURL(url)
	assert.Nil(t, err2)
	assert.Equal(t, hash1, hash2)
}

func TestSave_GetHashByURL_Empty(t *testing.T) {
	db := initDBStorage(t)

	url := "https://" + util.GenerateUniqueID()

	hash, err := db.GetHashByURL(url)
	assert.NotNil(t, err)
	assert.Equal(t, "", hash)
}

func TestGetAll(t *testing.T) {
	db := initDBStorage(t)

	hash1 := util.GenerateUniqueID()
	hash2 := util.GenerateUniqueID()
	url1 := "https://" + hash1
	url2 := "https://" + hash2

	assert.Nil(t, db.Save(hash1, url1))
	assert.Nil(t, db.Save(hash2, url2))

	data, err := db.GetAll()
	assert.Nil(t, err)
	assert.Equal(t, data[hash1], url1)
	assert.Equal(t, data[hash2], url2)
}

func initDBStorage(t *testing.T) *Source {
	db, err := NewDB("postgres://shortener:pass@localhost:5432/shortener")
	if err != nil {
		t.Skip("no db connection")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = db.db.PingContext(ctx)

	if err != nil {
		log.Println("exceptions ping DB:", err)
		t.Skip("no db connection")
	}
	return db
}
