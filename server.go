package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

func main() {
	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
