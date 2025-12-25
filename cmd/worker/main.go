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

	"github.com/kawai-network/veridium/internal/image"
	"github.com/kawai-network/veridium/internal/whisper"
	"github.com/kawai-network/veridium/pkg/gateway"
)

func main() {
	// CLI flags
	host := flag.String("host", "0.0.0.0", "Host to bind the server")
	port := flag.Int("port", 8080, "Port to bind the server")
	modelPath := flag.String("model", "", "Path to model file (optional, auto-select if empty)")
	maxTokens := flag.Int("max-tokens", 2048, "Max tokens to generate")
	flag.Parse()

	// Initialize LlamaExecutor
	ctx := context.Background()
	log.Printf("Initializing llama.cpp executor (model: %s)...", *modelPath)

	executor, err := gateway.NewLlamaExecutor(ctx, *modelPath, int32(*maxTokens))
	if err != nil {
		log.Fatalf("Failed to initialize executor: %v", err)
	}

	// Initialize WhisperExecutor
	log.Printf("Initializing whisper.cpp service...")
	whisperService, err := whisper.NewService()
	if err != nil {
		log.Printf("Warning: Failed to initialize whisper service: %v", err)
		// Continue without whisper
	}
	whisperExecutor := gateway.NewWhisperExecutor(whisperService)

	// Initialize ImageExecutor
	log.Printf("Initializing stable-diffusion.cpp engine...")
	sdEngine := image.NewEngine()
	imageExecutor := gateway.NewSDLocalExecutor(sdEngine)

	// Create gateway server
	cfg := gateway.ServerConfig{
		Host:      *host,
		Port:      *port,
		ModelName: "llama.cpp", // TODO: Get actual model name from executor
	}
	server := gateway.NewServer(cfg, executor, whisperExecutor, imageExecutor)

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
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Stop(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	fmt.Println("Server stopped")
}
