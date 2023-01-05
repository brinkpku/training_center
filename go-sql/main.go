package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type User struct {
	Id                 int
	First, Last, Email string
}

func main() {
	db, err := sqlx.Connect("mysql", "name:password@tcp(ip:port)/ras_edge?charset=utf8mb4")
	if err != nil {
		log.Fatalln(err)
	}

	user := []User{}

	db.Select(&user, "select * from users")

	log.Println("users...")
	log.Println(user)

}
