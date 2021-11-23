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
	if err != nil {
		log.Fatalf("get migrate driver with instance error: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "sqlite3", driver)
	if err != nil {
		log.Fatalf("new instance error: %v", err)
	}
	if err = m.Up(); err != nil {
		log.Fatalf("up error: %v", err)
	}
	curV, dirty, err := driver.Version()
	if err != nil {
		log.Fatalf("get version error: %v", err)
	}
	log.Printf("current version is %d, dirty is %v", curV, dirty)
}
