package application

import (
	"context"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type App struct {
	// Application fields go here
	router http.Handler
	rdb    *redis.Client
}

func NewApp() *App {
	return &App{
		router: loadRoutes(),
		rdb:    redis.NewClient(&redis.Options{}),
	}
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}

	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("Failed to connect to Redis: %w", err)
	}

	fmt.Println("Connected to Redis successfully")

	fmt.Println("Starting server on :3000")

	go func() {
		<-ctx.Done()
		fmt.Println("Shutting down server...")

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

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("Server error: %w", err)
	}

	return nil
}
