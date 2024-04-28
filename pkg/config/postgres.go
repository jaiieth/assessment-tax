package config

import (
	"database/sql"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type Postgres struct {
	Db *sql.DB
}

func New() (*Postgres, error) {
	var env map[string]string
	env, err := godotenv.Read()

	if err != nil {
		return nil, err
	}

	var DATABASE_URL = env["DATABASE_URL"]

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
