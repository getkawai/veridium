package kronk

import (
	"context"
	"embed"
	"errors"
	"expvar"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/getsentry/sentry-go"
	"github.com/kawai-network/veridium/cmd/server/api/services/kronk/build"
	"github.com/kawai-network/veridium/cmd/server/app/domain/whisperapp"
	"github.com/kawai-network/veridium/cmd/server/app/sdk/cache"
	"github.com/kawai-network/veridium/cmd/server/app/sdk/mux"
	"github.com/kawai-network/veridium/cmd/server/foundation/logger"
	"github.com/kawai-network/veridium/cmd/server/foundation/web"
	"github.com/kawai-network/veridium/internal/paths"
	"github.com/kawai-network/veridium/internal/services"
	"github.com/kawai-network/veridium/pkg/blockchain"
	"github.com/kawai-network/veridium/pkg/hardware"
	"github.com/kawai-network/veridium/pkg/kronk"
	pkglogger "github.com/kawai-network/veridium/pkg/logger"
	sd "github.com/kawai-network/veridium/pkg/stablediffusion"
	"github.com/kawai-network/veridium/pkg/stablediffusion/modeldownloader"
	"github.com/kawai-network/veridium/pkg/store"
	"github.com/kawai-network/veridium/pkg/tools/catalog"
	"github.com/kawai-network/veridium/pkg/tools/defaults"
	"github.com/kawai-network/veridium/pkg/tools/libs"
	"github.com/kawai-network/veridium/pkg/tools/models"
	"github.com/kawai-network/veridium/pkg/tools/templates"
	"github.com/kawai-network/x/constant"
	"github.com/kawai-network/x/tunnelkit"
)

//go:embed static
var static embed.FS

var tag = "develop"

const SentryDSN = "https://6d138acbdde2516e32e24f016b472031@o4510620614983680.ingest.us.sentry.io/4510620618850304"

// StartCommand runs the server (equivalent to the old Run function)
func StartCommand(args []string) error {
	var showHelp bool

	// Parse flags
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--help", "-h":
			showHelp = true
		}
	}

	return Run(showHelp)
}

// Run runs the kronk server (used by both main and start command)
func Run(showHelp bool) error {
	var log *logger.Logger

	// Configure log writer (before Sentry init so we can log to file immediately)
	logWriter := logger.NewWriter(paths.ContributorLog())

	// Initialize Sentry
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              SentryDSN,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		EnableLogs:       true,
		BeforeSendLog: func(log *sentry.Log) *sentry.Log {
			if log.Severity < sentry.LogSeverityWarning {
				return nil
			}
			return log
		},
	})
	if err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	defer sentry.Flush(2 * time.Second)

	// Setup logger with Sentry
	var sentryHandler slog.Handler
	if err == nil {
		baseHandler := slog.NewJSONHandler(logWriter, nil)
		sentryHandler = pkglogger.NewSentryHandler(baseHandler)
	}

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* SEND ALERT *******")
		},
	}

	if sentryHandler != nil {
		log = logger.NewWithSentry(logWriter, logger.LevelInfo, "KRONK", web.GetTraceID, sentryHandler)
	} else {
		log = logger.NewWithEvents(logWriter, logger.LevelInfo, "KRONK", web.GetTraceID, events)
	}

	log.Info(context.Background(), "logging", "file", paths.ContributorLog())

	// -------------------------------------------------------------------------

	ctx := context.Background()

	if err := run(ctx, log, showHelp); err != nil {
		return err
	}

	return nil
}

func run(ctx context.Context, log *logger.Logger, showHelp bool) error {

	// -------------------------------------------------------------------------
	// GOMAXPROCS

	if !showHelp {
		log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))
	}

	// -------------------------------------------------------------------------
	// Configuration

	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout     time.Duration `conf:"default:30s"`
			WriteTimeout    time.Duration `conf:"default:15m"`
			IdleTimeout     time.Duration `conf:"default:1m"`
			ShutdownTimeout time.Duration `conf:"default:1m"`
		}

		Catalog struct {
			GithubRepo string `conf:"default:https://api.github.com/repos/kawai-network/veridium_catalogs/contents/catalogs"`
		}
		Templates struct {
			GithubRepo string `conf:"default:https://api.github.com/repos/kawai-network/veridium_catalogs/contents/templates"`
		}
		Cache struct {
			ModelsInCache        int           `conf:"default:3"`
			TTL                  time.Duration `conf:"default:20m"`
			IgnoreIntegrityCheck bool          `conf:"default:true"`
			ModelConfigFile      string
		}

		BasePath     string
		LibPath      string
		LibVersion   string
		Arch         string
		OS           string
		Processor    string
		HfToken      string `conf:"mask"`
		AllowUpgrade bool   `conf:"default:true"`
		LlamaLog     int    `conf:"default:1"`
	}{
		Version: conf.Version{
			Build: tag,
			Desc:  "Kronk",
		},
	}

	const prefix = "KRONK"
	if showHelp {
		help, err := conf.UsageInfo(prefix, &cfg)
		if err != nil {
			return fmt.Errorf("parsing config: %w", err)
		}
		return fmt.Errorf("%s", help)
	}

	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	// -------------------------------------------------------------------------
	// App Starting

	log.Info(ctx, "starting service", "version", cfg.Build)
	defer log.Info(ctx, "shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Info(ctx, "startup", "config", out)

	log.BuildInfo(ctx)

	expvar.NewString("build").Set(cfg.Build)

	fmt.Println(logo)

	// -------------------------------------------------------------------------

	// -------------------------------------------------------------------------
	// Library System

	log.Info(ctx, "startup", "status", "checking libraries")

	arch, err := defaults.Arch(cfg.Arch)
	if err != nil {
		return err
	}

	opSys, err := defaults.OS(cfg.OS)
	if err != nil {
		return err
	}

	processor, err := defaults.Processor(cfg.Processor)
	if err != nil {
		return err
	}

	// Check all required libraries (llama, whisper, stablediffusion)
	requiredLibs := []libs.LibraryType{libs.LibraryLlama, libs.LibraryWhisper, libs.LibraryStableDiffusion}
	for _, libType := range requiredLibs {
		lib, err := libs.New(
			libs.WithBasePath(paths.Base()),
			libs.WithArch(arch),
			libs.WithOS(opSys),
			libs.WithProcessor(processor),
			libs.WithLibraryType(libType),
		)
		if err != nil {
			return fmt.Errorf("unable to create libs api for %s: %w", libType.DisplayName(), err)
		}

		log.Info(ctx, "startup", "status", "checking library", "library", libType.DisplayName(), "libPath", lib.LibsPath(), "arch", lib.Arch(), "os", lib.OS(), "processor", lib.Processor())

		if _, err := lib.InstalledVersion(); err != nil {
			return fmt.Errorf("%s library not found. Please run 'kawai-contributor setup' first to install required libraries", libType.DisplayName())
		}

		log.Info(ctx, "startup", "status", "library verified", "library", libType.DisplayName())
	}

	// Keep llama libs for backward compatibility in mux.Config
	libs, err := libs.New(
		libs.WithBasePath(paths.Base()),
		libs.WithArch(arch),
		libs.WithOS(opSys),
		libs.WithProcessor(processor),
		libs.WithAllowUpgrade(cfg.AllowUpgrade),
		libs.WithVersion(defaults.LibVersion(cfg.LibVersion)),
		libs.WithLibraryType(libs.LibraryLlama),
	)
	if err != nil {
		return fmt.Errorf("unable to create libs api: %w", err)
	}

	log.Info(ctx, "startup", "status", "all libraries verified")

	// -------------------------------------------------------------------------
	// Model System

	// Use paths.Base() for consistent base path resolution
	models, err := models.NewWithPaths(paths.Base())
	if err != nil {
		return fmt.Errorf("unable to create catalog system: %w", err)
	}

	log.Info(ctx, "startup", "status", "model integrity checks, may take a few seconds")

	models.BuildIndex(log.Info)

	// -------------------------------------------------------------------------
	// Catalog System

	log.Info(ctx, "startup", "status", "checking catalog")

	ctlg, err := catalog.New(
		catalog.WithBasePath(paths.Base()),
		catalog.WithGithubRepo(cfg.Catalog.GithubRepo))
	if err != nil {
		return fmt.Errorf("unable to create catalog system: %w", err)
	}

	// Check if catalog exists, don't download
	catalogs, err := ctlg.RetrieveCatalogs()
	if err != nil {
		return fmt.Errorf("catalog not found. Please run 'kawai-contributor setup' first to download catalog: %w", err)
	}

	if len(catalogs) == 0 {
		return fmt.Errorf("catalog is empty. Please run 'kawai-contributor setup' first to download catalog")
	}

	log.Info(ctx, "startup", "status", "catalog verified", "catalogs", len(catalogs))

	// -------------------------------------------------------------------------
	// Template System

	log.Info(ctx, "startup", "status", "checking templates")

	tmplts, err := templates.New(
		templates.WithBasePath(paths.Base()),
		templates.WithGithubRepo(cfg.Templates.GithubRepo),
		templates.WithCatalog(ctlg))
	if err != nil {
		return fmt.Errorf("unable to create template system: %w", err)
	}

	// Check if templates exist, don't download
	templatesPath := tmplts.TemplatesPath()
	entries, err := os.ReadDir(templatesPath)
	if err != nil {
		return fmt.Errorf("unable to read templates directory: %w", err)
	}

	templateCount := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			templateCount++
		}
	}

	if templateCount == 0 {
		return fmt.Errorf("no templates found. Please run 'kawai-contributor setup' first to download templates")
	}

	log.Info(ctx, "startup", "status", "templates verified", "count", templateCount)

	// -------------------------------------------------------------------------
	// Init Kronk

	log.Info(ctx, "startup", "status", "initializing kronk")

	if err := kronk.Init(); err != nil {
		return fmt.Errorf("installation invalid: %w", err)
	}

	cache, err := cache.New(cache.Config{
		Log:                  log.Info,
		BasePath:             paths.Base(),
		Templates:            tmplts,
		ModelsInCache:        cfg.Cache.ModelsInCache,
		CacheTTL:             cfg.Cache.TTL,
		IgnoreIntegrityCheck: cfg.Cache.IgnoreIntegrityCheck,
		ModelConfigFile:      cfg.Cache.ModelConfigFile,
	})

	if err != nil {
		return fmt.Errorf("initializing kronk manager: %w", err)
	}

	defer func() {
		log.Info(ctx, "shutdown", "status", "shutting down kronk")

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := cache.Shutdown(ctx); err != nil {
			log.Error(ctx, "kronk manager", "ERROR", err)
		}
	}()

	// -------------------------------------------------------------------------
	// Contributor Features (Always Enabled)

	var walletAddress string

	// Print welcome banner
	printBanner()

	log.Info(ctx, "startup", "status", "initializing contributor features")

	// Initialize KV Store
	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		log.Error(ctx, "kv store", "ERROR", err)
		return fmt.Errorf("failed to connect to KV: %w", err)
	}
	log.Info(ctx, "startup", "status", "connected to Cloudflare KV")

	// Initialize Blockchain Client for halving logic
	blockchainClient, err := blockchain.NewClient(blockchain.Config{
		RPCUrl:           constant.MonadRpcUrl,
		TokenAddress:     constant.KawaiTokenAddress,
		OTCMarketAddress: constant.OTCMarketAddress,
		USDTAddress:      constant.StablecoinAddress,
	})
	if err != nil {
		log.Info(ctx, "blockchain", "status", "failed to initialize, using default rates", "error", err)
	} else {
		kv.SetSupplyQuerier(blockchainClient)
		log.Info(ctx, "startup", "status", "blockchain client initialized", "rpc", constant.MonadRpcUrl)
	}

	// Setup Wallet
	wallet := services.NewWalletService("", kv)

	// Check if wallet exists
	if !wallet.HasWallet() {
		return fmt.Errorf("no wallet found. Please run 'kawai-contributor setup' first to configure your wallet")
	}

	// Wallet exists - unlock it
	wallets := wallet.GetWallets()
	printInfo("Wallet found!")

	if len(wallets) > 1 {
		fmt.Println("\nAvailable wallets:")
		for i, w := range wallets {
			active := ""
			if w.IsActive {
				active = " (active)"
			}
			fmt.Printf("  %d. %s - %s%s\n", i+1, w.Description, w.Address[:10]+"...", active)
		}

		choice, err := promptChoice("\nSelect wallet:", func() []string {
			options := make([]string, len(wallets))
			for i, w := range wallets {
				options[i] = fmt.Sprintf("%s (%s...)", w.Description, w.Address[:10])
			}
			return options
		}())
		if err != nil {
			return fmt.Errorf("failed to select wallet: %w", err)
		}

		selectedWallet := wallets[choice].Address
		password, err := promptPassword("Enter password: ")
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}

		walletAddress, err = wallet.SwitchWallet(selectedWallet, password)
		if err != nil {
			return fmt.Errorf("failed to switch wallet: %w", err)
		}
	} else {
		password, err := promptPassword("Enter password to unlock: ")
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}

		walletAddress, err = wallet.UnlockWallet(password)
		if err != nil {
			return fmt.Errorf("invalid password: %w", err)
		}
	}
	printSuccess(fmt.Sprintf("Wallet unlocked: %s", walletAddress))

	log.Info(ctx, "startup", "status", "wallet ready", "address", walletAddress)

	// Register Holder
	holderRegistry := blockchain.NewHolderRegistry(kv)
	if err := holderRegistry.RegisterHolder(ctx, common.HexToAddress(walletAddress), "kronk"); err != nil {
		log.Info(ctx, "holder", "status", "registration failed", "error", err)
	} else {
		log.Info(ctx, "startup", "status", "holder registered")
	}

	// Detect Hardware
	hwSpecs := hardware.DetectHardwareSpecs()
	hardwareInfo := fmt.Sprintf("%s, %d cores, %dGB RAM, GPU: %s (%dGB VRAM)",
		hwSpecs.CPU, hwSpecs.CPUCores, hwSpecs.TotalRAM, hwSpecs.GPUModel, hwSpecs.GPUMemory)
	log.Info(ctx, "startup", "status", "hardware detected", "info", hardwareInfo)

	// Start Tunnel (always enabled for contributors)
	var tunnelURL string
	tunnelCtx, tunnelCancel := context.WithCancel(context.Background())
	defer tunnelCancel()

	tunnelURL = startTunnel(tunnelCtx, log)
	if tunnelURL != "" {
		log.Info(ctx, "startup", "status", "tunnel started", "url", tunnelURL)
	} else {
		log.Info(ctx, "tunnel", "status", "no tunnel available")
	}

	// Register Contributor
	endpointURL := tunnelURL
	if endpointURL == "" {
		endpointURL = constant.LocalContributorURL
	}

	contributor, err := kv.RegisterContributor(ctx, walletAddress, endpointURL, hardwareInfo)
	if err != nil {
		return fmt.Errorf("failed to register contributor: %w", err)
	}
	log.Info(ctx, "startup", "status", "contributor registered", "wallet", contributor.WalletAddress, "since", contributor.RegisteredAt.Format("2006-01-02"))

	// Start Heartbeat with Metrics (fixed 30s interval)
	heartbeatCtx, heartbeatCancel := context.WithCancel(context.Background())
	defer heartbeatCancel()

	// Track metrics for contributor discovery
	// TODO: These metrics need to be updated by the request handler
	// For now, we report basic availability metrics only
	// Future: Integrate with cache.AquireModel() to track actual request metrics

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-heartbeatCtx.Done():
				return
			case <-ticker.C:
				// Detect region (simple heuristic based on timezone)
				region := detectRegion()

				// Get available models from cache
				availableModels := getAvailableModels(cache)

				// Get current active streams from cache as proxy for active requests
				modelStatus, err := cache.ModelStatus()
				var activeRequests int64
				if err == nil {
					for _, model := range modelStatus {
						activeRequests += int64(model.ActiveStreams)
					}
				}

				// Update metrics
				// Note: TotalRequests, AvgResponseTime, and SuccessRate will be 0 until
				// request tracking is implemented in the API handlers
				metrics := &store.ContributorMetrics{
					Region:          region,
					AvailableModels: availableModels,
					ActiveRequests:  activeRequests,
					TotalRequests:   0,   // TODO: Track in request handler
					AvgResponseTime: 0,   // TODO: Track in request handler
					SuccessRate:     1.0, // Default to 1.0 until we have failure data
					CPUCores:        hwSpecs.CPUCores,
					TotalRAM:        hwSpecs.TotalRAM,
					AvailableRAM:    hwSpecs.AvailableRAM,
					GPUModel:        hwSpecs.GPUModel,
					GPUMemory:       hwSpecs.GPUMemory,
				}

				if err := kv.UpdateContributorMetrics(ctx, walletAddress, metrics); err != nil {
					log.Info(ctx, "heartbeat", "status", "failed", "error", err)
				}
			}
		}
	}()
	log.Info(ctx, "startup", "status", "heartbeat with metrics started", "interval", "30s")

	// Initialize Whisper Service (pkg/whisper)
	var whisperModelsDir string
	{
		log.Info(ctx, "startup", "status", "initializing whisper")

		whisperModelsDir = paths.WhisperModels()

		// Ensure models directory exists
		if err := os.MkdirAll(whisperModelsDir, 0755); err != nil {
			log.Info(ctx, "whisper", "status", "failed to create models directory", "error", err)
		} else {
			// Check if whisper models exist, don't download
			models, _ := whisperapp.ListDownloadedModels()
			if len(models) == 0 {
				return fmt.Errorf("whisper models not found. Please run 'kawai-contributor setup' first to download whisper models")
			}
			log.Info(ctx, "startup", "status", "whisper service ready", "models", len(models))
		}
	}

	// Initialize Stable Diffusion
	var imageEngine *sd.StableDiffusion
	{
		log.Info(ctx, "startup", "status", "initializing stable diffusion")

		// 1. Check Library
		if !sd.IsLibraryInstalled() {
			return fmt.Errorf("stable diffusion library not found. Please run 'kawai-contributor setup' first to install SD library")
		}

		// 2. Find Model
		// Models are organized by {author}/{repo}/ from HuggingFace URLs
		modelsPath := paths.Models()
		downloader := modeldownloader.New(modelsPath)

		modelFile, err := downloader.DiscoverModel()
		if err != nil {
			log.Info(ctx, "startup", "status", "error discovering models", "error", err)
		}

		if modelFile == "" {
			// No model found - don't download, return error
			return fmt.Errorf("stable diffusion model not found. Please run 'kawai-contributor setup' first to download SD model")
		}

		log.Info(ctx, "startup", "status", "found SD model", "path", modelFile)

		// 3. Initialize Engine
		ctxParams := &sd.ContextParams{
			DiffusionModelPath: modelFile,
			DiffusionFlashAttn: true,
			OffloadParamsToCPU: true,
		}

		eng, err := sd.NewStableDiffusion(ctxParams)
		if err != nil {
			log.Info(ctx, "startup", "status", "failed to init SD engine", "error", err)
		} else {
			imageEngine = eng
			log.Info(ctx, "startup", "status", "stable diffusion ready")
		}
	}

	// Cleanup on shutdown
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := kv.MarkContributorOffline(shutdownCtx, walletAddress); err != nil {
			log.Info(ctx, "shutdown", "status", "failed to mark contributor offline", "error", err)
		} else {
			log.Info(ctx, "shutdown", "status", "contributor marked offline")
		}
	}()

	// -------------------------------------------------------------------------
	// Start API Service

	log.Info(ctx, "startup", "status", "initializing V1 API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	cfgMux := mux.Config{
		Build: tag,
		Log:   log,

		Tracer:      nil,
		Cache:       cache,
		Libs:        libs,
		Models:      models,
		Catalog:     ctlg,
		Templates:   tmplts,
		ImageEngine: imageEngine,
	}

	webAPI := mux.WebAPI(cfgMux,
		build.Routes(),
		mux.WithCORS([]string{"*"}),
		mux.WithFileServer(true, static, "static", "/", []string{"v1"}),
	)

	api := http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", constant.LocalContributorPort),
		Handler:      webAPI,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", api.Addr)

		serverErrors <- api.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}

		log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)
	}

	return nil
}

var logo = `
██╗  ██╗ █████╗ ██╗    ██╗ █████╗ ██╗
██║ ██╔╝██╔══██╗██║    ██║██╔══██╗██║
█████╔╝ ███████║██║ █╗ ██║███████║██║
██╔═██╗ ██╔══██║██║███╗██║██╔══██║██║
██║  ██╗██║  ██║╚███╔███╔╝██║  ██║██║
╚═╝  ╚═╝╚═╝  ╚═╝ ╚══╝╚══╝ ╚═╝  ╚═╝╚═╝
`

// startTunnel attempts to start a tunnel and returns the public URL
func startTunnel(ctx context.Context, log *logger.Logger) string {
	tunnels := tunnelkit.GetTunnels()
	for _, tunnel := range tunnels {
		if ok, _ := tunnelkit.HasActiveConnections(tunnel.TunnelID); !ok {
			go tunnelkit.RunTunnel(ctx, tunnel.TunnelToken)
			return tunnel.PublicURL
		}
	}
	return ""
}

// detectRegion detects the contributor's geographic region
// Simple heuristic based on timezone offset
func detectRegion() string {
	_, offset := time.Now().Zone()
	offsetHours := offset / 3600

	// Simple region mapping based on UTC offset
	// Note: Boundaries are exclusive on upper end to avoid overlaps
	switch {
	case offsetHours >= -8 && offsetHours < -5:
		return "us-west"
	case offsetHours >= -5 && offsetHours < -3:
		return "us-east"
	case offsetHours >= 0 && offsetHours < 3:
		return "eu-west"
	case offsetHours >= 3 && offsetHours < 6:
		return "eu-east"
	case offsetHours >= 6 && offsetHours < 9:
		return "asia-west"
	case offsetHours >= 9 && offsetHours <= 12:
		return "asia-east"
	default:
		return "unknown"
	}
}

// getAvailableModels returns list of available model IDs from cache
func getAvailableModels(cache *cache.Cache) []string {
	if cache == nil {
		return []string{}
	}

	// Get model status from cache
	modelDetails, err := cache.ModelStatus()
	if err != nil {
		return []string{}
	}

	modelIDs := make([]string, 0, len(modelDetails))
	for _, model := range modelDetails {
		modelIDs = append(modelIDs, model.ID)
	}

	return modelIDs
}
