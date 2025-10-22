package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sohWenMing/portfolio/internal/controllers/handlers"
	dbinfra "github.com/sohWenMing/portfolio/internal/db_infra"
	loadenv "github.com/sohWenMing/portfolio/internal/env"
	headerinspectionmiddleware "github.com/sohWenMing/portfolio/internal/middleware/header_inspection_middleware"
	csrf_protect "github.com/sohWenMing/portfolio/internal/security/csrf_protect"
	templating "github.com/sohWenMing/portfolio/internal/views/templating"
)

func main() {
	db, handler, err := Run(false, ".env")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("listening on port 3000")
	http.ListenAndServe(":3000", handler)
}

func Run(isTestingDBLocally bool, envPath string) (db *sql.DB, handler http.Handler, err error) {
	tplExecutor, err := loadTemplateExecutor()
	if err != nil {
		return nil, nilHandler{}, err
	}
	envGetter, err := loadEnvGetter(envPath)
	if err != nil {
		return nil, nilHandler{}, err
	}
	csrfMW := csrf_protect.LoadCSRFMW(".env", envGetter)
	db, err = loadDb(envGetter, isTestingDBLocally)
	if err != nil {
		return nil, nilHandler{}, err
	}
	defer db.Close()
	r := chi.NewRouter()

	//here we are casting the mux retunred by chi.NewRouter() - it has the ServeHTTP method which fulfils
	//the interface, so we can pass the whole mux as a handler into further middlewares
	handler = http.Handler(r)

	handler = headerinspectionmiddleware.InspectHeaders(handler)
	handler = csrfMW(handler)
	// r.Use(csrfmiddleware.CSRFMWGetToken)
	r.Get("/", handlers.TestHandler(tplExecutor))
	r.Post("/", handlers.TestReceiveFormHandler)
	return db, handler, nil
}

func loadDb(envGetter *loadenv.EnvGetter, isTestingDBLocally bool) (*sql.DB, error) {
	db, err := dbinfra.InitDB(envGetter.GetDBConfig(isTestingDBLocally))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func loadEnvGetter(envPath string) (*loadenv.EnvGetter, error) {
	envGetter, err := loadenv.LoadEnv(envPath)
	if err != nil {
		return nil, err
	}
	return envGetter, err
}

func loadTemplateExecutor() (*templating.TemplateExecutor, error) {
	tplExecutor, err := templating.InitTemplateExecutor()
	if err != nil {
		return nil, err
	}
	return tplExecutor, nil
}

type nilHandler struct{}

func (n nilHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
