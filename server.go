package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	zlog, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer zlog.Sync()

	// Replace the global logger with the one we just created
	// so that we can use the zap.L() function to get the logger
	// anywhere in the code.
	zap.ReplaceGlobals(zlog)

	e := echo.New()

	// Start the server
	errChan := make(chan error, 1)
	go func() {
		errChan <- e.Start(":" + GetEnv("PORT", "2565"))
	}()

	// Wait for an interrupt or kill signal
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	// Shutdown the server gracefully
	select {
	case err := <-errChan:
		if err != nil && err != http.ErrServerClosed {
			zlog.Fatal("Server failed to start", zap.Error(err))
		}
		zlog.Info("Server gracefully stopped")

	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		zlog.Info("Shutting down the server")
		if err := e.Shutdown(ctx); err != nil {
			zlog.Fatal("Server failed to shutdown", zap.Error(err))
		}
		zlog.Info("Server gracefully stopped")
	}
}

// GetEnv gets the value of the environment variable named by the key.
// It returns the value, which will be empty if the variable is not present.
// If no fallback is given, it will return an empty string.
func GetEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
