package db

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/tsupko/shortener/internal/app/exceptions"
)

type Source struct {
	db *sql.DB
}

func NewDB(DatabaseDsn string) (*Source, error) {
	db, err := sql.Open("pgx", DatabaseDsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxIdleTime(time.Second * 30)
	db.SetConnMaxLifetime(time.Minute * 2)

	return &Source{db: db}, nil
}

func (dbSource *Source) Close() error {
	err := dbSource.db.Close()
	if err != nil {
		log.Println("exceptions closing connection to DB:", err)
	}
	return err
}

func (dbSource *Source) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := dbSource.db.PingContext(ctx); err != nil {
		log.Println("error while ping DB:", err)
		return err
	}
	return nil
}

func (dbSource *Source) InitTables() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := dbSource.db.ExecContext(ctx, `
        create table urls (hash varchar(30) not null constraint urls_pk primary key, url varchar(500));
		create unique index urls_url_uindex on urls (url);
    `)
	if err != nil {
		log.Println("init tables are NOT created - ", err)
		return
	}
	log.Println("init tables are created")
}

func (dbSource *Source) Save(hash string, url string) error {
	log.Println("try to save; hash=", hash, "url=", url)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	row, err := dbSource.db.ExecContext(ctx, "insert into urls (hash, url) values ($1, $2)", hash, url)
	if err != nil {
		log.Println("error while Save:", err)
		return err
	}
	log.Println("db.Saved ", row)
	return nil
}

func (dbSource *Source) SaveBatch(hashes []string, urls []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := dbSource.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := dbSource.db.PrepareContext(ctx, "insert into urls (hash, url) values ($1, $2)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := range hashes {
		if _, err = stmt.ExecContext(ctx, hashes[i], urls[i]); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (dbSource *Source) Get(hash string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var url string
	row := dbSource.db.QueryRowContext(ctx, "select url from urls where hash = $1", hash)
	err := row.Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("Get from DB return nothing:", err)
			return "", exceptions.ErrURLNotFound
		}
		log.Println("error while Get from DB:", err)
	}
	return url, nil
}

func (dbSource *Source) GetAll() (map[string]string, error) {
	Limit := 2000

	var hash string
	var url string
	data := make(map[string]string)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := dbSource.db.QueryContext(ctx, "select hash, url from urls limit $1", Limit)
	if err != nil {
		log.Println(err)
		return data, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&hash, &url)
		if err != nil {
			log.Println(err)
			return data, err
		}
		data[hash] = url
	}
	err = rows.Err()
	if err != nil {
		return data, err
	}
	return data, nil
}

func (dbSource *Source) GetHashByURL(url string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var hash string
	row := dbSource.db.QueryRowContext(ctx, "select hash from urls where url = $1", url)
	err := row.Scan(&hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("Get hash from DB return nothing:", err)
			return "", err
		}
		log.Println("exceptions while Get hash from DB:", err)
	}
	return hash, nil
}
