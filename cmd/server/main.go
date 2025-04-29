package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dekib/rePacks/internal/server"
)

func main() {
	// Channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	// Capture Ctrl+C and kill commands
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Create server with timeouts
	srv := server.CreateServer()

	// Configure port
	port := ":8086"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = ":" + envPort
	}

	// HTTP server with timeouts
	httpServer := &http.Server{
		Addr:              port,
		Handler:           srv,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	// Run server in goroutine to allow graceful shutdown
	go func() {
		log.Printf("Server running on %s", port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Block until we receive interrupt signal
	<-quit
	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited properly")
}
