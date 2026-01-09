package application

import (
	"fmt"
	"net/http"

	"github.com/RyanSikandar/orders-api/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World")
	})

	router.Route("/orders", loadOrderRoutes)

	return router
}

func loadOrderRoutes(r chi.Router) {
	orderHandler := &handler.Order{}

	r.Get("/", orderHandler.List)
	r.Post("/", orderHandler.Create)
	r.Get("/{id}", orderHandler.GetByID)
	r.Put("/{id}", orderHandler.UpdateByID)
	r.Delete("/{id}", orderHandler.DeleteByID)
}
