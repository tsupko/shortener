package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tsupko/shortener/internal/app/db"
	"github.com/tsupko/shortener/internal/app/util"
)

func TestDBStorage(t *testing.T) {
	hash := util.GenerateUniqueID()

	dbSource, err := db.NewDB("postgres://shortener:pass@localhost:5432/shortener")

	if err != nil {
		t.Error("error when NewDB", err)
	}

	err = dbSource.Ping()
	if err != nil {
		t.Skip("no db connection")
	}
	storage := NewDBStorage(dbSource)

	storage.Put(hash, "http://url")
	storage.Put(hash, "http://url2")

	url, _ := storage.Get(hash)
	assert.Equal(t, "http://url", url) //old value

	data, _ := storage.GetAll()
	assert.GreaterOrEqual(t, len(data.(map[string]string)), 1)

	urlFromMap := data.(map[string]string)[hash]
	assert.Equal(t, "http://url", urlFromMap)
}
