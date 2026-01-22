package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/getsentry/sentry-go"
	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/internal/image"
	"github.com/kawai-network/veridium/internal/services"
	"github.com/kawai-network/veridium/internal/whisper"
	"github.com/kawai-network/veridium/pkg/blockchain"
	"github.com/kawai-network/veridium/pkg/fantasy/llamalib"
	"github.com/kawai-network/veridium/pkg/gateway"
	"github.com/kawai-network/veridium/pkg/hardware"
	"github.com/kawai-network/veridium/pkg/logger"
	"github.com/kawai-network/veridium/pkg/store"
	"github.com/kawai-network/veridium/pkg/tunnelkit"
)

const (
	SentryDSN = "https://6d138acbdde2516e32e24f016b472031@o4510620614983680.ingest.us.sentry.io/4510620618850304"
)

func main() {
	// CLI flags

	walletPassword := flag.String("password", "", "Wallet password")
	walletAddress := flag.String("wallet", "", "Wallet address to use")
	importMnemonic := flag.String("import-mnemonic", "", "Import wallet from mnemonic")
	flag.Parse()

	// ============================================
	// Step 0: Initialize Sentry
	// ============================================
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              SentryDSN,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		EnableLogs:       true, // Enable Logger API as requested
		BeforeSendLog: func(log *sentry.Log) *sentry.Log {
			// filter all logs below warning
			if log.Severity <= sentry.LogSeverityWarning {
				return nil
			}
			return log
		},
	})
	if err != nil {
		slog.Error("Sentry initialization failed", "error", err)
	}

	// Always defer flush, even if init failed (it's a no-op then) or if we error out later
	defer sentry.Flush(2 * time.Second)

	if err == nil {
		// properties of the logger
		handler := slog.NewTextHandler(os.Stderr, nil)
		// wrap the handler with SentryHandler
		sentryHandler := logger.NewSentryHandler(handler)
		// create a new logger with the SentryHandler
		logger := slog.New(sentryHandler)
		// set the default logger to the new logger
		slog.SetDefault(logger)
	}

	// ============================================
	// Step 1: Initialize KV Store (Required for Wallet)
	// ============================================
	ctx := context.Background()

	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		slog.Error("Failed to connect to KV", "error", err)
		os.Exit(1)
	}
	slog.Info("✓ Connected to Cloudflare KV")

	// ============================================
	// Step 1.5: Initialize Blockchain Client (For Halving Logic)
	// ============================================
	blockchainBC, err := blockchain.NewClient(blockchain.Config{
		RPCUrl:           constant.MonadRpcUrl,
		TokenAddress:     constant.KawaiTokenAddress,
		OTCMarketAddress: constant.OTCMarketAddress,
		USDTAddress:      constant.StablecoinAddress,
	})
	if err != nil {
		slog.Warn("⚠️ Failed to initialize blockchain client, halving logic will use default rates", "error", err)
	} else {
		kv.SetSupplyQuerier(blockchainBC)
		slog.Info("✓ Blockchain client initialized and injected for halving logic", "rpc", constant.MonadRpcUrl)
	}

	// ============================================
	// Step 2: Setup Wallet
	// ============================================
	wallet := services.NewWalletService("", kv)

	if !wallet.HasWallet() && *walletPassword == "" {
		slog.Error("No wallet found. Use --password to create one.")
		os.Exit(1)
	}

	if wallet.HasWallet() && *walletPassword == "" {
		for _, w := range wallet.GetWallets() {
			slog.Info("Available wallet", "address", w.Address, "description", w.Description)
		}
		slog.Error("Use --password to unlock.")
		os.Exit(1)
	}

	var address string

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
		slog.Error("Wallet error", "error", err)
		os.Exit(1)
	}
	slog.Info("✓ Wallet", "address", address)

	// ============================================
	// Step 2.5: Register Holder (KAWAI Token Holder Registry)
	// ============================================
	holderRegistry := blockchain.NewHolderRegistry(kv)
	if err := holderRegistry.RegisterHolder(ctx, common.HexToAddress(address), "cli"); err != nil {
		slog.Warn("⚠️ Failed to register holder", "error", err)
		// Non-fatal - continue anyway
	} else {
		slog.Info("✓ Holder registered")
	}

	// ============================================
	// Step 3: Detect Hardware
	// ============================================
	slog.Info("Detecting hardware...")
	hwSpecs := hardware.DetectHardwareSpecs()
	hardwareInfo := fmt.Sprintf("%s, %d cores, %dGB RAM, GPU: %s (%dGB VRAM)",
		hwSpecs.CPU, hwSpecs.CPUCores, hwSpecs.TotalRAM, hwSpecs.GPUModel, hwSpecs.GPUMemory)
	slog.Info("✓ Hardware", "info", hardwareInfo)

	// ============================================
	// Step 4: Register Contributor (KV already init)
	// ============================================
	// ctx is already created above

	// ============================================
	// Step 5: Start Tunnel (get public URL first)
	// ============================================
	tunnelCtx, tunnelCancel := context.WithCancel(context.Background())
	defer tunnelCancel()

	tunnelURL := startTunnel(tunnelCtx)
	if tunnelURL != "" {
		slog.Info("✓ Tunnel", "url", tunnelURL)
	} else {
		// Fatal error: contributor needs public endpoint
		slog.Error("No public tunnel available. Cannot start contributor.")
		os.Exit(1)
	}

	endpointURL := tunnelURL

	// ============================================
	// Step 6: Register Contributor to KV
	// ============================================
	contributor, err := kv.RegisterContributor(ctx, address, endpointURL, hardwareInfo)
	if err != nil {
		slog.Error("Failed to register", "error", err)
		os.Exit(1)
	}
	slog.Info("✓ Registered", "wallet", contributor.WalletAddress, "since", contributor.RegisteredAt.Format("2006-01-02"))

	// ============================================
	// Step 7: Start Heartbeat (direct to KV)
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
					slog.Warn("⚠️ Heartbeat failed", "error", err)
				}
			}
		}
	}()
	slog.Info("✓ Heartbeat started (30s)")

	// ============================================
	// Step 8: Initialize AI Services
	// ============================================
	llm := llamalib.NewService()

	// Wait for LLM initialization
	initCtx, initCancel := context.WithTimeout(ctx, 2*time.Minute)
	defer initCancel()
	if err := llm.WaitForInitialization(initCtx); err != nil {
		slog.Error("Failed to initialize LLM", "error", err)
		os.Exit(1)
	}

	// Load chat model (auto-select best)
	if err := llm.LoadChatModel(""); err != nil {
		slog.Error("Failed to load model", "error", err)
		os.Exit(1)
	}
	slog.Info("✓ LLM ready")

	whisperService, err := whisper.NewService()
	if err != nil {
		slog.Error("Failed to initialize Whisper service", "error", err)
		os.Exit(1)
	}
	whisperExecutor := gateway.NewWhisperExecutor(whisperService)

	sdEngine := image.NewEngine()
	imageExecutor := gateway.NewSDLocalExecutor(sdEngine)

	// ============================================
	// Step 9: Start Server
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
			slog.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	fmt.Printf("\n  Wallet: %s\n  Local:  %s\n", address, constant.LocalContributorURL)
	if tunnelURL != "" {
		fmt.Printf("  Public: %s\n", tunnelURL)
	}
	fmt.Println()

	<-quit

	// Cleanup: mark offline
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := kv.MarkContributorOffline(shutdownCtx, address); err != nil {
		slog.Warn("Failed to mark contributor offline", "error", err)
	} else {
		slog.Info("✓ Contributor marked offline")
	}

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
