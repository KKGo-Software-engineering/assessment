package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Ic3Sandy/assessment/expenses"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))

	db = ConnectAndCreateTable()
	expenses.SetDB(db)
	defer db.Close()

	e := echo.New()
	e.GET("/expenses", expenses.GetExpenses)
	e.GET("/expenses/:id", expenses.GetExpensesById)
	e.POST("/expenses", expenses.CreateExpense)
	e.PUT("/expenses/:id", expenses.UpdateExpense)
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
