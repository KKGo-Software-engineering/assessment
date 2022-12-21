package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	go func() {
		if err := e.Start(":" + os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
