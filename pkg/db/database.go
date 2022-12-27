package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Doittikorn/assessment/pkg/config"
)

var db *sql.DB

func ConnectDB() {
	var err error
	c := config.NewConfig()
	db, err = sql.Open("postgres", c.Database())

	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	} else {
		fmt.Println("DB Connected...")
	}

	createTb := `
	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);
	`
	_, err = db.Exec(createTb)

	if err != nil {
		log.Fatal("can't create table", err)
	}
}

func GetDB() *sql.DB {
	return db
}
