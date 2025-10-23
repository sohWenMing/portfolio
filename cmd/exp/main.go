package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	srv := http.Server{
		Addr:    ":8000",
		Handler: http.HandlerFunc(TestHandler),
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Println("listening on port 8000")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("error occured in server: %v", err)
		}
	}()
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("error on server shutdown: %v", err)
	} else {
		fmt.Println("server shutdown gracefully.")
	}

}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Able to receive on port 8000"))
}
