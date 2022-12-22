package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err := createExpenseTable(db); err != nil {
		log.Fatal(err)
	}
}

func createExpenseTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS expenses (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			amount INT NOT NULL,
			note VARCHAR(255),
			tags TEXT[]
		)
	`)
	return err
}
