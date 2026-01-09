package main

import (
	"fmt"
	"net/http"
	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()

	router.Get("/", basicHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	err:= server.ListenAndServe()
	if err != nil {
		fmt.Println("Server error:", err)
	}
}

func basicHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}
