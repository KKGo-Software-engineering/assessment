package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"

	"github.com/Doittikorn/assessment/pkg/config"
	"github.com/Doittikorn/assessment/pkg/db"
)

func main() {
	c := config.NewConfig()
	db.ConnectDB()

	e := echo.New()

	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, `Hello`)
	})

	// graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		e.Logger.Fatal(e.Start(":" + c.Port()))
	}()
	fmt.Println("App started.")

	killSignal := <-signals
	switch killSignal {
	case os.Interrupt:
		fmt.Println("Got SIGINT...")
	case syscall.SIGTERM:
		fmt.Println("got SIGTERM...")
	}
	fmt.Println("App is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

}
