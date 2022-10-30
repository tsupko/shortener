package db

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Source struct {
	db *sql.DB
}

func NewDB(DatabaseDsn string) (*Source, error) {
	db, err := sql.Open("pgx", DatabaseDsn)
	if err != nil {
		return nil, err
	}
	return &Source{db: db}, nil
}

func (dbSource *Source) Close() error {
	err := dbSource.db.Close()
	if err != nil {
		log.Println("error closing connection to DB:", err)
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
	_, err := dbSource.db.ExecContext(ctx, "create table urls (hash varchar(30) not null constraint urls_pk primary key, url varchar(500))")
	if err != nil {
		log.Println("init tables are NOT created - ", err)
		return
	}
	log.Println("init tables are created")
}

func (dbSource *Source) Save(hash string, url string) {
	log.Println("try to save; hash=", hash, "url=", url)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	row, err := dbSource.db.ExecContext(ctx, "insert into urls (hash, url) values ($1, $2)", hash, url)
	if err != nil {
		log.Println("error while Save:", err)
	}
	log.Println("db.Saved ", row)
}

func (dbSource *Source) Get(hash string) (string, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var url string
	row := dbSource.db.QueryRowContext(ctx, "select url from urls where hash = $1", hash)
	err := row.Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("Get from DB return nothing:", err)
			return "", false
		}
		log.Println("error while Get from DB:", err)
	}
	return url, true
}

func (dbSource *Source) GetAll() map[string]string {
	var hash string
	var url string
	var data = make(map[string]string)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := dbSource.db.QueryContext(ctx, "select hash, url from urls limit 20")
	if err != nil {
		log.Println(err)
		return data
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&hash, &url)
		if err != nil {
			log.Println(err)
			return data
		}
		data[hash] = url
	}
	err = rows.Err()
	if err != nil {
		return data
	}
	return data
}
