package config

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

type Postgres struct {
	Db *sql.DB
}

func New() (*Postgres, error) {
	var DATABASE_URL = os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", DATABASE_URL)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Postgres{Db: db}, nil
}
