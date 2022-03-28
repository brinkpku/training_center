package main

import (
	"database/sql"
	"os"
	"runtime"
	"sync"
	"time"

	"log"

	// sqlite
	_ "github.com/mattn/go-sqlite3"
)

//go:generate go run main.go && python analyze.py && tail -n 1 test.log

type DBConfig struct {
	Dsl          string
	MaxOpenConns int
	MaxIdleConns int
}

var sqliteDB *sql.DB

// InitDB 初始化sqlite database
func InitDB(cfg *DBConfig) (err error) {
	if sqliteDB != nil {
		sqliteDB.Close()
	}

	sqliteDB, err = sql.Open("sqlite3", cfg.Dsl)
	if err != nil {
		return
	}
	if cfg.MaxOpenConns <= 0 {
		log.Println("get pu nums", runtime.NumCPU())
		cfg.MaxOpenConns = runtime.NumCPU() // error num reduced significantly after limiting conns num
		cfg.MaxIdleConns = runtime.NumCPU()
	}

	sqliteDB.SetMaxOpenConns(1) // no database is locked error
	// sqliteDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqliteDB.SetMaxIdleConns(cfg.MaxIdleConns)

	if err = sqliteDB.Ping(); err != nil {
		return
	}
	return nil
}

// GetSqlClient ...
func GetSqlClient() *sql.DB {
	return sqliteDB
}

var (
	writeErrNum int
	readErrNum  int
	mu          sync.Mutex
	threadNum   = int64(1e5)
)

func main() {
	start := time.Now().Unix()
	os.Remove("test.log")
	f, err := os.OpenFile("test.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetFlags(log.LstdFlags | log.Ltime | log.Lshortfile)
	log.SetOutput(f)
	if err := InitDB(&DBConfig{
		Dsl: "data.db?_journal_mode=WAL&_busy_timeout=60",
	}); err != nil {
		log.Fatal(err)
	}
	if _, err := sqliteDB.Exec("create table if not exists foo (`id` integer, `value` text);"); err != nil {
		log.Fatal(err)
	}
	wg := sync.WaitGroup{}
	for i := 0; int64(i) < threadNum; i++ {
		wg.Add(1)
		go func() {
			readRequest()
			wg.Done()
		}()
	}
	for i := 0; int64(i) < threadNum; i++ {
		wg.Add(1)
		go func() {
			writeRequest()
			wg.Done()
		}()
	}
	wg.Wait()
	sqliteDB.Exec("delete from foo;")
	sqliteDB.Close()
	log.Println("run", time.Now().Unix()-start, "s.", writeErrNum, "write error.", readErrNum, "read error.")
}

func writeRequest() {
	tx, err := sqliteDB.Begin()
	if err != nil {
		log.Println(err)
		return
	}
	if _, err := tx.Exec("insert into foo values (?, ?);", time.Now().Nanosecond(), "write a message"); err != nil {
		log.Println(err)
		// mu.Lock()
		// readErrNum += 1
		// mu.Unlock()
		if inErr := tx.Rollback(); inErr != nil {
			log.Println(inErr)
		}
	} else {
		log.Println("write a record")
		if err := tx.Commit(); err != nil {
			log.Println(err)
		}
	}
}

func readRequest() {
	tx, err := sqliteDB.Begin()
	if err != nil {
		log.Println(err)
		return
	}
	if _, err := tx.Exec("select * from foo limit 1;"); err != nil {
		log.Println(err)
		// mu.Lock()
		// readErrNum += 1
		// mu.Unlock()
		if inErr := tx.Rollback(); inErr != nil {
			log.Println(inErr)
		}
	} else {
		log.Println("read a record")
		if err := tx.Commit(); err != nil {
			log.Println(err)
		}
	}
}
