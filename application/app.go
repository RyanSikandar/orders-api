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

	err = server.ListenAndServe()

	if err != nil {
		return fmt.Errorf("Failed to start server: %w", err)
	}

	return nil
}
