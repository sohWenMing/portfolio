package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/sohWenMing/portfolio/internal/controllers/handlers"
	loadenv "github.com/sohWenMing/portfolio/internal/env"
	csrfmiddleware "github.com/sohWenMing/portfolio/internal/middleware/csrf_middleware"
	csrf_protect "github.com/sohWenMing/portfolio/internal/security/csrf_protect"
)

func main() {
	envGetter, err := loadenv.LoadEnv(".env")
	if err != nil {
		panic(err)
	}
	csrfKey := csrf_protect.GetCSRFKey(envGetter)
	csrf := csrf.Protect(csrfKey)

	r := chi.NewRouter()
	r.Use(csrfmiddleware.CSRFMWGetToken)
	r.Get("/", handlers.TestHandler)
	r.Post("/", handlers.TestReceiveFormHandler)
	fmt.Println("listening on port 3000")
	http.ListenAndServe(":3000", csrf(r))
}
