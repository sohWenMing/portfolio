package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sohWenMing/portfolio/internal/controllers/handlers"
)

func main() {
	r := chi.NewRouter()
	r.Get("/", handlers.TestHandler)
	r.Post("/", handlers.TestReceiveFormHandler)
	fmt.Println("listening on port 3000")
	http.ListenAndServe(":3000", r)
}
