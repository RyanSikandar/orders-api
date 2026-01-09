package main

import (
	"context"
	"fmt"

	"github.com/RyanSikandar/orders-api/application"
)

func main() {
	app := application.NewApp()

	err := app.Start(context.TODO())

	if err != nil {
		fmt.Println("Error starting the app:", err)
	}
}