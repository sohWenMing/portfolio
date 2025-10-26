package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	businesslogic "github.com/sohWenMing/portfolio/internal/business_logic"
	"github.com/sohWenMing/portfolio/internal/controllers/handlers"
	dbinfra "github.com/sohWenMing/portfolio/internal/db_infra"
	dbinterface "github.com/sohWenMing/portfolio/internal/db_interface"
	loadenv "github.com/sohWenMing/portfolio/internal/env"
	headerinspectionmiddleware "github.com/sohWenMing/portfolio/internal/middleware/header_inspection_middleware"
	csrf_protect "github.com/sohWenMing/portfolio/internal/security/csrf_protect"
	"github.com/sohWenMing/portfolio/internal/security/passwordhashing"
	templating "github.com/sohWenMing/portfolio/internal/views/templating"
)

type services struct {
	userservice *businesslogic.UserService
}

const portAddr string = ":3000"

func main() {
	appDb, handler, _, err := Run(false, ".env")
	if err != nil {
		panic(err)
	}
	err = appDb.Migrate()
	if err != nil {
		fmt.Println("error occured during migration: ", err)
	} else {
		fmt.Println("migrations successfully ran")
	}

	srv := &http.Server{
		Addr:    portAddr,
		Handler: handler,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)

	go func() {
		fmt.Println("listening on port: ", portAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("server error: %v", err)
		}
	}()

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("error during server shutdown: ", err)
	}
	if appDb != nil && appDb.DB != nil {
		appDb.DB.Close()
		fmt.Println("db closed successfully.")
	}

	fmt.Println("server shutdown gracefully.")
}

func Run(isTestingDBLocally bool, envPath string) (
	appDb *dbinfra.AppDB,
	handler http.Handler,
	returnedServices *services,
	err error) {
	tplExecutor, err := loadTemplateExecutor()
	if err != nil {
		return nil, nilHandler{}, nil, err
	}
	envGetter, err := loadEnvGetter(envPath)
	if err != nil {
		return nil, nilHandler{}, nil, err
	}
	csrfMW := csrf_protect.LoadCSRFMW(".env", envGetter)
	appDb, err = loadDb(envGetter, isTestingDBLocally)
	if err != nil {
		return nil, nilHandler{}, nil, err
	}
	// services will have acccess to the CRUD operations, as defined by the db interface
	// note that because the term service is used as a whole to define both business logic
	// this would mean that a business logic level service has access to a CRUD level
	// service
	returnedServices = &services{
		userservice: businesslogic.InitUserService(
			10*time.Second,
			dbinterface.InitUserServiceWithPostgres(appDb.DB),
			&passwordhashing.BcryptInterfacer{},
		),
	}
	r := chi.NewRouter()

	//here we are casting the mux retunred by chi.NewRouter() - it has the ServeHTTP method which fulfils
	//the interface, so we can pass the whole mux as a handler into further middlewares
	handler = http.Handler(r)

	handler = headerinspectionmiddleware.InspectHeaders(handler)
	handler = csrfMW(handler)
	// r.Use(csrfmiddleware.CSRFMWGetToken)
	r.Get("/", handlers.TestHandler(tplExecutor))
	r.Post("/", handlers.TestReceiveFormHandler)
	return appDb, handler, returnedServices, nil
}

func loadDb(envGetter *loadenv.EnvGetter, isTestingDBLocally bool) (*dbinfra.AppDB, error) {
	appDB, err := dbinfra.InitDB(envGetter.GetDBConfig(isTestingDBLocally))
	if err != nil {
		return nil, err
	}
	return appDB, nil
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
