package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// before run this
// docker run --name postgresql -e POSTGRES_USER=shortener -e POSTGRES_PASSWORD=pass -p 5432:5432 -d postgres
func _TestPingOk(t *testing.T) {
	db := NewDB("postgres://shortener:pass@localhost:5432/shortener")
	err := db.Ping()
	assert.Equal(t, nil, err)
}

func TestPingWrongPort(t *testing.T) {
	db := NewDB("postgres://shortener:pass@localhost:5433/shortener")
	err := db.Ping()
	assert.Contains(t, err.Error(), "failed to connect")
}
