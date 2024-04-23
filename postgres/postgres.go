package postgres

import (
	"database/sql"
	"fmt"

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

	var (
		host     = env["DB_HOST"]
		port     = env["DB_PORT"]
		user     = env["DB_USER"]
		password = env["DB_PASSWORD"]
		dbname   = env["DB_NAME"]
	)

	databaseSource := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", databaseSource)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &Postgres{Db: db}, nil
}
