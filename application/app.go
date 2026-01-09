package application

import (
	"context"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type App struct {
	// router is stored as http.Handler interface for flexibility
	// allowing easy swapping of routing implementations if needed
	router http.Handler
	// rdb holds the Redis client for persistence layer operations
	rdb *redis.Client
}

func NewApp() *App {
	app := &App{
		rdb: redis.NewClient(&redis.Options{}),
	}
	app.loadRoutes()
	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}

	// Verify Redis connection is healthy before starting the HTTP server
	// This ensures we fail fast if Redis is unavailable rather than accepting requests
	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("Failed to connect to Redis: %w", err)
	}

	fmt.Println("Connected to Redis successfully")

	fmt.Println("Starting server on :3000")

	// Start a goroutine to handle graceful shutdown when context is cancelled
	// This allows the server to finish processing existing requests before stopping
	go func() {
		<-ctx.Done()
		fmt.Println("Shutting down server...")

		// Create a new context with 5 second timeout for shutdown operations
		// Using background context here since the parent context is already cancelled
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5_000_000_000)
		defer cancel()

		if shutdownErr := server.Shutdown(shutdownCtx); shutdownErr != nil {
			fmt.Println("Error during server shutdown:", shutdownErr)
		}

		if err := a.rdb.Close(); err != nil {
			fmt.Println("Error closing Redis client:", err)
		}

		fmt.Println("Server and Redis client gracefully stopped")
	}()

	// ListenAndServe blocks until server shuts down
	// http.ErrServerClosed is expected during graceful shutdown, so we ignore it
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("Server error: %w", err)
	}

	return nil
}
