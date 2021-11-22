package main

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var sqliteDB *sql.DB

	if sqliteDB != nil {
		sqliteDB.Close()
	}

	sqliteDB, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Fatalf("open data.db failed: %v", err)
	}

	if err = sqliteDB.Ping(); err != nil {
		return
	}
	driver, err := sqlite3.WithInstance(sqliteDB, &sqlite3.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"sqlite3", driver)
	m.Steps(2)
}
