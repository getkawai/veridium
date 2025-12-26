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

	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/internal/image"
	"github.com/kawai-network/veridium/internal/services"
	"github.com/kawai-network/veridium/internal/whisper"
	"github.com/kawai-network/veridium/pkg/fantasy/llamalib"
	"github.com/kawai-network/veridium/pkg/gateway"
	"github.com/kawai-network/veridium/pkg/hardware"
	"github.com/kawai-network/veridium/pkg/store"
	"github.com/kawai-network/veridium/pkg/tunnelkit"
)

func main() {
	// CLI flags

	walletPassword := flag.String("password", "", "Wallet password")
	walletAddress := flag.String("wallet", "", "Wallet address to use")
	importMnemonic := flag.String("import-mnemonic", "", "Import wallet from mnemonic")
	flag.Parse()

	// ============================================
	// Step 1: Setup Wallet
	// ============================================
	wallet := services.NewWalletService("")

	if !wallet.HasWallet() && *walletPassword == "" {
		log.Fatal("No wallet found. Use --password to create one.")
	}

	if wallet.HasWallet() && *walletPassword == "" {
		for _, w := range wallet.GetWallets() {
			log.Printf("  - %s (%s)", w.Address, w.Description)
		}
		log.Fatal("Use --password to unlock.")
	}

	var address string
	var err error

	if *importMnemonic != "" {
		address, err = wallet.CreateWallet(*walletPassword, *importMnemonic, "Kawai Contributor")
	} else if !wallet.HasWallet() {
		mnemonic, _ := wallet.GenerateMnemonic()
		fmt.Println("\n⚠️  SAVE YOUR MNEMONIC:")
		fmt.Printf("\n  %s\n\n", mnemonic)
		address, err = wallet.CreateWallet(*walletPassword, mnemonic, "Kawai Contributor")
	} else if *walletAddress != "" {
		address, err = wallet.SwitchWallet(*walletAddress, *walletPassword)
	} else {
		address, err = wallet.UnlockWallet(*walletPassword)
	}

	if err != nil {
		log.Fatalf("Wallet error: %v", err)
	}
	log.Printf("✓ Wallet: %s", address)

	// ============================================
	// Step 2: Detect Hardware
	// ============================================
	log.Println("Detecting hardware...")
	hwSpecs := hardware.DetectHardwareSpecs()
	hardwareInfo := fmt.Sprintf("%s, %d cores, %dGB RAM, GPU: %s (%dGB VRAM)",
		hwSpecs.CPU, hwSpecs.CPUCores, hwSpecs.TotalRAM, hwSpecs.GPUModel, hwSpecs.GPUMemory)
	log.Printf("✓ Hardware: %s", hardwareInfo)

	// ============================================
	// Step 3: Initialize KV Store & Register
	// ============================================
	ctx := context.Background()

	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		log.Fatalf("Failed to connect to KV: %v", err)
	}
	log.Println("✓ Connected to Cloudflare KV")

	// ============================================
	// Step 4: Start Tunnel (get public URL first)
	// ============================================
	tunnelCtx, tunnelCancel := context.WithCancel(context.Background())
	defer tunnelCancel()

	tunnelURL := startTunnel(tunnelCtx)
	if tunnelURL != "" {
		log.Printf("✓ Tunnel: %s", tunnelURL)
	}

	endpointURL := tunnelURL

	// ============================================
	// Step 5: Register Contributor to KV
	// ============================================
	contributor, err := kv.RegisterContributor(ctx, address, endpointURL, hardwareInfo)
	if err != nil {
		log.Fatalf("Failed to register: %v", err)
	}
	log.Printf("✓ Registered: %s (since %s)", contributor.WalletAddress, contributor.RegisteredAt.Format("2006-01-02"))

	// ============================================
	// Step 6: Start Heartbeat (direct to KV)
	// ============================================
	heartbeatCtx, heartbeatCancel := context.WithCancel(context.Background())
	defer heartbeatCancel()

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-heartbeatCtx.Done():
				return
			case <-ticker.C:
				if err := kv.UpdateHeartbeat(ctx, address); err != nil {
					log.Printf("⚠️ Heartbeat failed: %v", err)
				}
			}
		}
	}()
	log.Println("✓ Heartbeat started (30s)")

	// ============================================
	// Step 7: Initialize AI Services
	// ============================================
	llm := llamalib.NewService()

	// Wait for LLM initialization
	initCtx, initCancel := context.WithTimeout(ctx, 2*time.Minute)
	defer initCancel()
	if err := llm.WaitForInitialization(initCtx); err != nil {
		log.Fatalf("Failed to initialize LLM: %v", err)
	}

	// Load chat model (auto-select best)
	if err := llm.LoadChatModel(""); err != nil {
		log.Fatalf("Failed to load model: %v", err)
	}
	log.Println("✓ LLM ready")

	whisperService, _ := whisper.NewService()
	whisperExecutor := gateway.NewWhisperExecutor(whisperService)

	sdEngine := image.NewEngine()
	imageExecutor := gateway.NewSDLocalExecutor(sdEngine)

	// ============================================
	// Step 8: Start Server
	// ============================================
	server := gateway.NewServer(gateway.ServerConfig{
		Host:      "0.0.0.0",
		Port:      constant.LocalContributorPort,
		StaticDir: sdEngine.GetOutputsPath(),
	}, llm, whisperExecutor, imageExecutor, kv)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	fmt.Printf("\n  Wallet: %s\n  Local:  %s\n", address, constant.LocalContributorURL)
	if tunnelURL != "" {
		fmt.Printf("  Public: %s\n", tunnelURL)
	}
	fmt.Println()

	<-quit

	// Cleanup: mark offline
	if err := kv.UpdateHeartbeat(ctx, address); err == nil {
		// Could set status to offline here if needed
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	server.Stop(shutdownCtx)
}

func startTunnel(ctx context.Context) string {
	tunnels := tunnelkit.GetTunnels()
	for _, tunnel := range tunnels {
		if ok, _ := tunnelkit.HasActiveConnections(tunnel.TunnelID); !ok {
			go tunnelkit.RunTunnel(ctx, tunnel.TunnelToken)
			return tunnel.PublicURL
		}
	}
	return ""
}
