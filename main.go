package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/RyanSikandar/orders-api/application"
)

func main() {
	app := application.NewApp()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel() // this tells go to stop listening for signals anymore

	err := app.Start(ctx)

	if err != nil {
		fmt.Println("Error starting the app:", err)
	}
}