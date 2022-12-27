package main

import (
	"fmt"

	_ "github.com/lib/pq"

	"github.com/Doittikorn/assessment/config"
)

func main() {
	config := config.NewConfig()

	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", config.Port)
}
