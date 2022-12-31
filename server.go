package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/eskimoburger/assessment/handler"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("can't create table", err)
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
	h := handler.NewApplication(db)
	e := echo.New()

	e.POST("/expenses", h.CreateExpenseHandler)

	log.Fatal(e.Start(os.Getenv("PORT")))

}
