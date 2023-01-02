package main

import (
	"fmt"
	"os"

	"github.com/cs4begas/assessment/expenses"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Please use server.go for main file")
	expenses.InitDB()
	os_port := os.Getenv("PORT")
	fmt.Println("start at port:", os_port)
}
