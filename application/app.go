package application

import (
	"context"
	"fmt"
	"net/http"
)

type App struct {
	// Application fields go here
	router http.Handler
}

func NewApp() *App {
	return &App{
		router: loadRoutes(),
	}
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}

	err := server.ListenAndServe()

	if err != nil {
		return fmt.Errorf("Failed to start server: %w", err)
	}
	
	return nil
}