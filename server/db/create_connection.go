package db

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func CreateConnection() *sql.DB {
	DB, err := sql.Open("postgres", os.Getenv("POSTGRES_CONNECTION"))
	if err != nil {
		panic(err)
	}

	return DB
}
