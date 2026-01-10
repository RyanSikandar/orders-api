package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/RyanSikandar/orders-api/application"
)

func main() {
	app := application.NewApp(application.LoadConfig())

	// Create a context that will be cancelled when SIGINT or SIGKILL is received
	// This enables graceful shutdown when user presses Ctrl+C or system sends kill signal
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	// Ensure signal listener cleanup happens even if app.Start returns early
	defer cancel()

	err := app.Start(ctx)

	if err != nil {
		fmt.Println("Error starting the app:", err)
	}
}
