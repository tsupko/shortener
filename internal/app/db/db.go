package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pkg/errors"
)

type DB struct {
	DatabaseDsn string
}

func NewDB(DatabaseDsn string) *DB {
	return &DB{DatabaseDsn: DatabaseDsn}
}

func (dbSource *DB) Ping() error {
	if len(dbSource.DatabaseDsn) < 2 {
		return errors.New("DatabaseDsn too short")
	}
	db, err := sql.Open("pgx", dbSource.DatabaseDsn)
	if err != nil {
		return err
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}
