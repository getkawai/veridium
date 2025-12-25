package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kawai-network/veridium/pkg/gateway"
)

func main() {
	// CLI flags
	host := flag.String("host", "0.0.0.0", "Host to bind the server")
	port := flag.Int("port", 8080, "Port to bind the server")
	model := flag.String("model", "llama-3.2-8b", "Model name to advertise")
	flag.Parse()

	// Create executor (using MockExecutor for now, replace with llama.cpp integration)
	executor := &gateway.MockExecutor{}

	// Create gateway server
	cfg := gateway.ServerConfig{
		Host:      *host,
		Port:      *port,
		ModelName: *model,
	}
	server := gateway.NewServer(cfg, executor)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := server.Start(); err != nil {
			errChan <- err
		}
	}()

	// Wait for signal or error
	select {
	case <-quit:
		fmt.Println("\nShutting down server...")
	case err := <-errChan:
		log.Fatalf("Server error: %v", err)
	}

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	fmt.Println("Server stopped")
}
