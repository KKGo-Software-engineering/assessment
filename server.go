package main

import (
	"log"

	"github.com/labstack/echo/v4"

	"github.com/labstack/echo/v4/middleware"
	"github.com/nattha-wim/assessment/expense"
)

func main() {

	expense.InitDB()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover()) // เผื่อ server เรา down
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "admin" && password == "45678" {
			return true, nil
		}
		return false, nil
	}))

	e.POST("/expenses", expense.CreateExpense)

	log.Println("server start at :2565")
	log.Fatal(e.Start(":2565"))
	log.Println("bye bye")

	// fmt.Println("Please use server.go for main file")
	// fmt.Println("start at port:", os.Getenv("PORT"))
}
