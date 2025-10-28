package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	integration "github.com/sohWenMing/portfolio/internal/integration"
)

const portAddr string = ":3000"

func main() {
	appDb, handler, _, err := integration.Run(false, ".env")
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
