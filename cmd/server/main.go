package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sohWenMing/portfolio/internal/controllers/handlers"
	loadenv "github.com/sohWenMing/portfolio/internal/env"
	headerinspectionmiddleware "github.com/sohWenMing/portfolio/internal/middleware/header_inspection_middleware"
	csrf_protect "github.com/sohWenMing/portfolio/internal/security/csrf_protect"
	templating "github.com/sohWenMing/portfolio/internal/views/templating"
)

func main() {
	tplExecutor, err := templating.InitTemplateExecutor()
	if err != nil {
		//TODO: replace with logging function
		fmt.Println("error occured in InitTemplateExecutor: ", err)
		panic(err)
	}

	envGetter, err := loadenv.LoadEnv(".env")
	if err != nil {
		panic(err)
	}
	csrfMW := csrf_protect.LoadCSRFMW(".env", envGetter)

	r := chi.NewRouter()

	//here we are casting the mux retunred by chi.NewRouter() - it has the ServeHTTP method which fulfils
	//the interface, so we can pass the whole mux as a handler into further middlewares
	handler := http.Handler(r)

	handler = headerinspectionmiddleware.InspectHeaders(handler)
	handler = csrfMW(handler)
	// r.Use(csrfmiddleware.CSRFMWGetToken)
	r.Get("/", handlers.TestHandler(tplExecutor))
	r.Post("/", handlers.TestReceiveFormHandler)
	fmt.Println("listening on port 3000")
	http.ListenAndServe(":3000", handler)
}
