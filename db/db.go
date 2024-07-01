package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func ConnectDB(dbLocation string) {
	var err error
	DB, err = sql.Open("sqlite3", dbLocation)
	if err != nil {
		log.Fatal(err.Error())
	}
}
