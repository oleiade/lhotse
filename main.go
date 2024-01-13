package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
)

//nolint:forbidigo
func main() {
	// Create a logger instance
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create a new Echo instance
	e := echo.New()

	// Configure the Echo instance
	e.HideBanner = true

	// Register middleware
	e.Use(slogecho.New(logger))
	e.Use(middleware.Recover())

	// Register route handlers
	RegisterHandlers(e, &ServerImpl{})

	// Setup signal handling for graceful shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	go func() {
		<-signalCh
		// Initiate graceful shutdown
		if shutdownErr := e.Shutdown(context.Background()); shutdownErr != nil {
			slog.Error("Failed to shutdown server", "error_message", shutdownErr.Error())
			return
		}
	}()

	// Start the server
	if err := e.Start(":3434"); err != nil {
		slog.Error("Failed to start server", "error_message", err.Error())
		return
	}
}
