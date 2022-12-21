package main

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func ConnectAndCreateTable() *sql.DB {
	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS expenses (
			id SERIAL PRIMARY KEY,
			title TEXT,
			amount FLOAT,
			note TEXT,
			tags TEXT[]
		);
	`)
	if err != nil {
		panic(err)
	}
	return db
}
