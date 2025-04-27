package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func CreateConnection() *sql.DB {
	DB, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		panic(err)
	}

	return DB
}
