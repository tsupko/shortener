package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tsupko/shortener/internal/app/util"
)

func TestPingWrongPort(t *testing.T) {
	db, err := NewDB("postgres://shortener:pass@localhost:5433/shortener")
	require.NoError(t, err)
	err = db.Ping()
	require.Error(t, err)
}

//before run this
//docker run --name postgresql -e POSTGRES_USER=shortener -e POSTGRES_PASSWORD=pass -p 5432:5432 -d postgres

func _TestPingOk(t *testing.T) {
	db, err := NewDB("postgres://shortener:pass@localhost:5432/shortener")
	require.NoError(t, err)
	err = db.Ping()
	require.NoError(t, err)
}

func _TestSaveGet(t *testing.T) {
	db, err := NewDB("postgres://shortener:pass@localhost:5432/shortener")
	require.NoError(t, err)

	hash := util.GenerateUniqueID()

	db.Save(hash, "someUrl")

	urlFromDB, _ := db.Get(hash)
	assert.Equal(t, "someUrl", urlFromDB)
}

func _TestGetEmpty(t *testing.T) {
	db, err := NewDB("postgres://shortener:pass@localhost:5432/shortener")
	require.NoError(t, err)

	hash := util.GenerateUniqueID()

	urlFromDB, _ := db.Get(hash)
	assert.Equal(t, "someUrl", urlFromDB)
}
