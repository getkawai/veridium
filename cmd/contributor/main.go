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
	"github.com/kawai-network/veridium/pkg/tunnelkit"
)

func main() {
	// CLI flags
	host := flag.String("host", "0.0.0.0", "Host to bind the server")
	port := flag.Int("port", 8080, "Port to bind the server")
	modelPath := flag.String("model", "", "Path to model file (optional, auto-select if empty)")
	maxTokens := flag.Int("max-tokens", 2048, "Max tokens to generate")
	enableTunnel := flag.Bool("tunnel", true, "Enable Cloudflare Tunnel (default: true)")
	tunnelIndex := flag.Int("tunnel-index", 0, "Tunnel index to use from generated tunnels (default: 0)")
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
		StaticDir: sdEngine.GetOutputsPath(),
	}
	server := gateway.NewServer(cfg, executor, whisperExecutor, imageExecutor)

	// Start Cloudflare Tunnel if enabled
	tunnelCtx, tunnelCancel := context.WithCancel(context.Background())
	defer tunnelCancel()

	if *enableTunnel {
		go startTunnel(tunnelCtx, *tunnelIndex)
	}

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

// startTunnel starts a Cloudflare Tunnel for the contributor service
func startTunnel(ctx context.Context, preferredIndex int) {
	log.Println("Starting Cloudflare Tunnel...")

	// Get available tunnels
	tunnels := tunnelkit.GetTunnels()
	if len(tunnels) == 0 {
		log.Println("⚠️  No tunnels configured. Run: go run cmd/bulk-tunnels/main.go")
		log.Println("⚠️  Contributor will only be accessible locally")
		return
	}

	// Try preferred tunnel first
	tunnel, tunnelIndex := selectAvailableTunnel(tunnels, preferredIndex)
	if tunnel == nil {
		log.Println("❌ All tunnels are in use. Please generate more tunnels or wait.")
		return
	}

	log.Printf("✓ Using tunnel [%d]: %s", tunnelIndex, tunnel.Hostname)
	log.Printf("✓ Public URL: %s", tunnel.PublicURL)

	// Run tunnel with retry logic
	maxRetries := len(tunnels)
	for attempt := 0; attempt < maxRetries; attempt++ {
		err := tunnelkit.RunTunnel(ctx, tunnel.TunnelToken)

		if ctx.Err() != nil {
			// Context cancelled, normal shutdown
			log.Println("Tunnel stopped")
			return
		}

		if err == nil {
			// Tunnel stopped normally
			return
		}

		// Check if error is due to tunnel already in use
		errMsg := err.Error()
		if containsAny(errMsg, []string{"already in use", "duplicate connection", "connection exists"}) {
			log.Printf("⚠️  Tunnel [%d] is in use by another contributor", tunnelIndex)

			// Try next available tunnel
			remainingTunnels := make([]*tunnelkit.TunnelInfo, 0)
			for i, t := range tunnels {
				if i != tunnelIndex {
					remainingTunnels = append(remainingTunnels, t)
				}
			}

			if len(remainingTunnels) == 0 {
				log.Println("❌ No more tunnels available")
				return
			}

			// Select next tunnel
			tunnel, tunnelIndex = selectAvailableTunnel(tunnels, (tunnelIndex+1)%len(tunnels))
			if tunnel == nil {
				log.Println("❌ All tunnels are in use")
				return
			}

			log.Printf("🔄 Retrying with tunnel [%d]: %s", tunnelIndex, tunnel.Hostname)
			continue
		}

		// Other error, log and exit
		log.Printf("❌ Tunnel error: %v", err)
		return
	}

	log.Println("❌ Failed to connect to any tunnel after all retries")
}

// containsAny checks if string contains any of the substrings
func containsAny(s string, substrs []string) bool {
	for _, substr := range substrs {
		if len(s) >= len(substr) {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
		}
	}
	return false
}

// selectAvailableTunnel tries to find an available tunnel
// Returns the tunnel and its index, or nil if all are in use
func selectAvailableTunnel(tunnels []*tunnelkit.TunnelInfo, preferredIndex int) (*tunnelkit.TunnelInfo, int) {
	// Validate preferred index
	if preferredIndex >= len(tunnels) {
		log.Printf("⚠️  Tunnel index %d out of range (available: %d tunnels)", preferredIndex, len(tunnels))
		preferredIndex = 0
	}

	// Try preferred tunnel first
	if isTunnelAvailable(tunnels[preferredIndex]) {
		log.Printf("✓ Preferred tunnel [%d] is available", preferredIndex)
		return tunnels[preferredIndex], preferredIndex
	}
	log.Printf("⚠️  Preferred tunnel [%d] is in use, searching for alternatives...", preferredIndex)

	// Try other tunnels
	for i, tunnel := range tunnels {
		if i == preferredIndex {
			continue // Already tried
		}
		if isTunnelAvailable(tunnel) {
			log.Printf("✓ Found available tunnel [%d]", i)
			return tunnel, i
		}
	}

	// All tunnels are in use
	return nil, -1
}

// isTunnelAvailable checks if a tunnel is currently available (not in use)
func isTunnelAvailable(tunnel *tunnelkit.TunnelInfo) bool {
	// Check via Cloudflare API if tunnel has active connections
	hasActiveConnections, err := tunnelkit.HasActiveConnections(tunnel.TunnelID)
	if err != nil {
		// If error checking, assume available (optimistic)
		// Will fail fast if actually in use
		log.Printf("⚠️  Could not check tunnel status: %v", err)
		return true
	}

	return !hasActiveConnections
}
