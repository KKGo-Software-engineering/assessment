package expenses

import (
	"database/sql"
	"log"
	"os"
)

var Db *sql.DB

func InitDB() {
	Db = connectDB()
	// Create table
	createTb := `CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);`
	_, err := Db.Exec(createTb)

	if err != nil {
		log.Fatal("can't create table", err)
	}

}

func connectDB() (db *sql.DB) {
	var err error
	db_url := os.Getenv("DATABASE_URL")
	db_url += "?ssl.mode=disable"
	db, err = sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	return db
}

func GetDB() *sql.DB {
	return Db
}
