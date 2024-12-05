package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ganyacc/ap.git/server"
	"github.com/labstack/echo/v4"
)

func main() {

	// Start server in a goroutine

	go func() {
		if err := server.StartServer(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Shutting down the server: %s", err)
		}
	}()

	log.Println("Server started on :8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Print("Shutdown signal received")

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	//Attempt to gracefully shutdown the server
	if err := echo.New().Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}

	log.Print("Server exited gracefully")

}
