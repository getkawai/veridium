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
	"path/filepath"
	"runtime"
	"strings"
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
	sdmodels "github.com/kawai-network/veridium/pkg/stablediffusion/models"
	"github.com/kawai-network/x/store"
	"github.com/kawai-network/veridium/pkg/tools/catalog"
	"github.com/kawai-network/veridium/pkg/tools/defaults"
	"github.com/kawai-network/veridium/pkg/tools/libs"
	"github.com/kawai-network/veridium/pkg/tools/models"
	"github.com/kawai-network/veridium/pkg/tools/templates"
	"github.com/kawai-network/x/constant"
	"github.com/kawai-network/x/tunnelkit"
	"github.com/kawai-network/contracts"
)

//go:embed static
var static embed.FS

var tag = "develop"

const (
	SentryDSN             = "https://6d138acbdde2516e32e24f016b472031@o4510620614983680.ingest.us.sentry.io/4510620618850304"
	SentryFlushTimeout    = 2 * time.Second
	ServerShutdownTimeout = 1 * time.Minute
	ServerReadTimeout     = 30 * time.Second
	ServerWriteTimeout    = 15 * time.Minute
	ServerIdleTimeout     = 1 * time.Minute
	CacheTTL              = 20 * time.Minute
	CacheModelsInCache    = 3
)

const (
	CatalogGithubRepo   = "https://api.github.com/repos/ardanlabs/kronk_catalogs/contents/catalogs"
	TemplatesGithubRepo = "https://api.github.com/repos/ardanlabs/kronk_catalogs/contents/templates"
	SetupRequiredMsg    = "Please run 'kawai-contributor setup' first to %s"
)

// StartCommand runs the server (equivalent to the old Run function)
func StartCommand(args []string) error {
	var showHelp bool

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

	logWriter := logger.NewWriter(paths.ContributorLog())

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

	defer sentry.Flush(SentryFlushTimeout)

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

	ctx := context.Background()

	if err := run(ctx, log, showHelp); err != nil {
		return err
	}

	return nil
}

func run(ctx context.Context, log *logger.Logger, showHelp bool) error {
	if !showHelp {
		log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))
	}

	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout     time.Duration `conf:"default:30s"`
			WriteTimeout    time.Duration `conf:"default:15m"`
			IdleTimeout     time.Duration `conf:"default:1m"`
			ShutdownTimeout time.Duration `conf:"default:1m"`
		}

		Catalog struct {
			GithubRepo string `conf:"default:https://api.github.com/repos/ardanlabs/kronk_catalogs/contents/catalogs"`
		}
		Templates struct {
			GithubRepo string `conf:"default:https://api.github.com/repos/ardanlabs/kronk_catalogs/contents/templates"`
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
			return fmt.Errorf("%s library not found. "+SetupRequiredMsg, libType.DisplayName(), "install required libraries")
		}

		log.Info(ctx, "startup", "status", "library verified", "library", libType.DisplayName())
	}

	llamaLibs, err := libs.New(
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

	models, err := models.NewWithPaths(paths.Base())
	if err != nil {
		return fmt.Errorf("unable to create catalog system: %w", err)
	}

	log.Info(ctx, "startup", "status", "model integrity checks, may take a few seconds")
	models.BuildIndex(log.Info)

	catalogSvc, err := catalog.New(
		catalog.WithBasePath(paths.Base()),
		catalog.WithGithubRepo(cfg.Catalog.GithubRepo))
	if err != nil {
		return fmt.Errorf("unable to create catalog system: %w", err)
	}

	catalogs, err := catalogSvc.RetrieveCatalogs()
	if err != nil {
		return fmt.Errorf("catalog not found. "+SetupRequiredMsg, "download catalog")
	}

	if len(catalogs) == 0 {
		return fmt.Errorf("catalog is empty. "+SetupRequiredMsg, "download catalog")
	}

	log.Info(ctx, "startup", "status", "catalog verified", "catalogs", len(catalogs))

	templatesSvc, err := templates.New(
		templates.WithBasePath(paths.Base()),
		templates.WithGithubRepo(cfg.Templates.GithubRepo),
		templates.WithCatalog(catalogSvc))
	if err != nil {
		return fmt.Errorf("unable to create template system: %w", err)
	}

	templatesPath := templatesSvc.TemplatesPath()
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
		return fmt.Errorf("no templates found. "+SetupRequiredMsg, "download templates")
	}

	log.Info(ctx, "startup", "status", "templates verified", "count", templateCount)

	log.Info(ctx, "startup", "status", "initializing kronk")

	if err := kronk.Init(); err != nil {
		return fmt.Errorf("installation invalid: %w", err)
	}

	cacheSvc, err := cache.New(cache.Config{
		Log:                  log.Info,
		BasePath:             paths.Base(),
		Templates:            templatesSvc,
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

		if err := cacheSvc.Shutdown(ctx); err != nil {
			log.Error(ctx, "kronk manager", "ERROR", err)
		}
	}()

	printBanner()

	log.Info(ctx, "startup", "status", "initializing contributor features")

	kv, walletAddress, hwSpecs, err := initContributorFeatures(ctx, log, cfg.Web.ShutdownTimeout)
	if err != nil {
		return err
	}

	tunnelURL := startTunnel(ctx, log)
	if tunnelURL != "" {
		log.Info(ctx, "startup", "status", "tunnel started", "url", tunnelURL)
	} else {
		log.Info(ctx, "tunnel", "status", "no tunnel available")
	}

	endpointURL := tunnelURL
	if endpointURL == "" {
		endpointURL = constant.LocalContributorURL
	}

	hardwareInfo := fmt.Sprintf("%s, %d cores, %dGB RAM, GPU: %s (%dGB VRAM)",
		hwSpecs.CPU, hwSpecs.CPUCores, hwSpecs.TotalRAM, hwSpecs.GPUModel, hwSpecs.GPUMemory)
	contributor, err := kv.RegisterContributor(ctx, walletAddress, endpointURL, hardwareInfo)
	if err != nil {
		return fmt.Errorf("failed to register contributor: %w", err)
	}
	log.Info(ctx, "startup", "status", "contributor registered", "wallet", contributor.WalletAddress, "since", contributor.RegisteredAt.Format("2006-01-02"))

	// No heartbeat - client will check /v1/health directly for real-time status

	whisperModelsDir := paths.WhisperModels()
	if err := os.MkdirAll(whisperModelsDir, 0755); err != nil {
		log.Info(ctx, "whisper", "status", "failed to create models directory", "error", err)
	} else {
		models, err := whisperapp.ListDownloadedModels()
		if err != nil {
			log.Info(ctx, "whisper", "status", "failed to list models", "error", err)
		}
		if len(models) == 0 {
			return fmt.Errorf("whisper models not found. "+SetupRequiredMsg, "download whisper models")
		}
		log.Info(ctx, "startup", "status", "whisper service ready", "models", len(models))
	}

	imageEngine, imageEditEngine, err := initStableDiffusion(ctx, log, hwSpecs)
	if err != nil {
		return err
	}

	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := kv.MarkContributorOffline(shutdownCtx, walletAddress); err != nil {
			log.Info(ctx, "shutdown", "status", "failed to mark contributor offline", "error", err)
		} else {
			log.Info(ctx, "shutdown", "status", "contributor marked offline")
		}
	}()

	log.Info(ctx, "startup", "status", "initializing V1 API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	cfgMux := mux.Config{
		Build: tag,
		Log:   log,

		Tracer:          nil,
		Cache:           cacheSvc,
		Libs:            llamaLibs,
		Models:          models,
		Catalog:         catalogSvc,
		Templates:       templatesSvc,
		ImageEngine:     imageEngine,
		ImageEditEngine: imageEditEngine,
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

// initContributorFeatures initializes wallet, KV store, blockchain, and hardware detection
func initContributorFeatures(ctx context.Context, log *logger.Logger, shutdownTimeout time.Duration) (*store.KVStore, string, *hardware.HardwareSpecs, error) {
	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		log.Error(ctx, "kv store", "ERROR", err)
		return nil, "", nil, fmt.Errorf("failed to connect to KV: %w", err)
	}
	log.Info(ctx, "startup", "status", "connected to Cloudflare KV")

	blockchainClient, err := blockchain.NewClient(blockchain.Config{
		RPCUrl:           contracts.MonadRpcUrl,
		TokenAddress:     contracts.KawaiTokenAddress,
		OTCMarketAddress: contracts.OTCMarketAddress,
		USDTAddress:      contracts.StablecoinAddress,
	})
	if err != nil {
		log.Warn(ctx, "blockchain", "status", "failed to initialize, using default rates", "error", err)
	} else {
		kv.SetSupplyQuerier(blockchainClient)
		log.Info(ctx, "startup", "status", "blockchain client initialized", "rpc", contracts.MonadRpcUrl)
	}

	wallet := services.NewWalletService("", kv)

	if !wallet.HasWallet() {
		return nil, "", nil, fmt.Errorf("no wallet found. "+SetupRequiredMsg, "configure your wallet")
	}

	wallets := wallet.GetWallets()
	printInfo("Wallet found!")

	var walletAddress string
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
			return nil, "", nil, fmt.Errorf("failed to select wallet: %w", err)
		}

		selectedWallet := wallets[choice].Address
		password, err := promptPassword("Enter password: ")
		if err != nil {
			return nil, "", nil, fmt.Errorf("failed to read password: %w", err)
		}

		walletAddress, err = wallet.SwitchWallet(selectedWallet, password)
		if err != nil {
			return nil, "", nil, fmt.Errorf("failed to switch wallet: %w", err)
		}
	} else {
		password, err := promptPassword("Enter password to unlock: ")
		if err != nil {
			return nil, "", nil, fmt.Errorf("failed to read password: %w", err)
		}

		walletAddress, err = wallet.UnlockWallet(password)
		if err != nil {
			return nil, "", nil, fmt.Errorf("invalid password: %w", err)
		}
	}
	printSuccess(fmt.Sprintf("Wallet unlocked: %s", walletAddress))

	log.Info(ctx, "startup", "status", "wallet ready", "address", walletAddress)

	holderRegistry := blockchain.NewHolderRegistry(kv)
	if err := holderRegistry.RegisterHolder(ctx, common.HexToAddress(walletAddress), "kronk"); err != nil {
		log.Info(ctx, "holder", "status", "registration failed", "error", err)
	} else {
		log.Info(ctx, "startup", "status", "holder registered")
	}

	hwSpecs := hardware.DetectHardwareSpecs()
	log.Info(ctx, "startup", "status", "hardware detected", "cpu", hwSpecs.CPU, "cores", hwSpecs.CPUCores, "ram_gb", hwSpecs.TotalRAM, "gpu", hwSpecs.GPUModel, "vram_gb", hwSpecs.GPUMemory)

	return kv, walletAddress, hwSpecs, nil
}

type stableDiffusionModelBundle struct {
	selectedModel sdmodels.ModelSpec
	diffusionPath string
	llmPath       string
	vaePath       string
	editModelPath string
}

func validateStableDiffusionModelBundle(bundle *stableDiffusionModelBundle) error {
	if bundle == nil {
		return errors.New("stable diffusion model bundle is nil")
	}

	required := []struct {
		name string
		path string
	}{
		{name: "diffusion_model", path: bundle.diffusionPath},
		{name: "llm_model", path: bundle.llmPath},
		{name: "vae_model", path: bundle.vaePath},
	}

	for _, model := range required {
		if err := validateStableDiffusionModelFile(model.name, model.path); err != nil {
			return err
		}
	}

	// Edit model is optional at startup, but if selected it must be valid.
	if bundle.editModelPath != "" {
		if err := validateStableDiffusionModelFile("edit_model", bundle.editModelPath); err != nil {
			return err
		}
	}

	return nil
}

func validateStableDiffusionModelFile(name, modelPath string) error {
	if modelPath == "" {
		return fmt.Errorf("stable diffusion required model path is empty (%s)", name)
	}

	info, err := os.Stat(modelPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("stable diffusion required model file not found (%s): %s", name, modelPath)
		}
		return fmt.Errorf("failed to stat stable diffusion model file (%s): %w", name, err)
	}

	if info.IsDir() {
		return fmt.Errorf("stable diffusion required model path points to a directory (%s): %s", name, modelPath)
	}
	if info.Size() == 0 {
		return fmt.Errorf("stable diffusion required model file is empty (%s): %s", name, modelPath)
	}

	return nil
}

// initStableDiffusion initializes generation and edit Stable Diffusion engines.
func initStableDiffusion(ctx context.Context, log *logger.Logger, hwSpecs *hardware.HardwareSpecs) (*sd.StableDiffusion, *sd.StableDiffusion, error) {
	log.Info(ctx, "startup", "status", "initializing stable diffusion")

	libPath := sd.GetLibraryPath()
	if !sd.IsLibraryInstalled() {
		return nil, nil, fmt.Errorf("stable diffusion backend library not found at %s. "+SetupRequiredMsg, libPath, "install SD library")
	}

	modelsPath := paths.Models()
	bundle, err := resolveStableDiffusionModelBundle(modelsPath, hwSpecs)
	if err != nil {
		return nil, nil, err
	}
	if err := validateStableDiffusionModelBundle(bundle); err != nil {
		return nil, nil, fmt.Errorf("stable diffusion model validation failed: %w", err)
	}

	log.Info(ctx, "startup", "status", "selected SD model based on hardware", "model", bundle.selectedModel.Name, "path", bundle.diffusionPath)

	ctxParams := &sd.ContextParams{
		DiffusionModelPath: bundle.diffusionPath,
		LLMPath:            bundle.llmPath,
		VAEPath:            bundle.vaePath,
		DiffusionFlashAttn: true,
		OffloadParamsToCPU: true,
		FlowShift:          3.0,
	}

	if err := sd.InitLibrary(libPath); err != nil {
		return nil, nil, fmt.Errorf("failed to initialize stable diffusion backend library at %s: %w", libPath, err)
	}

	generationEngine, err := sd.NewStableDiffusion(ctxParams)
	if err != nil {
		log.Warn(ctx, "startup", "status", "failed to init SD engine", "error", err)
		return nil, nil, fmt.Errorf("failed to initialize stable diffusion generation engine: %w", err)
	}
	if !generationEngine.IsReady() {
		return nil, nil, errors.New("stable diffusion generation engine not ready after initialization")
	}

	var editEngine *sd.StableDiffusion
	if bundle.editModelPath != "" {
		editCtxParams := &sd.ContextParams{
			DiffusionModelPath: bundle.editModelPath,
			LLMPath:            bundle.llmPath,
			VAEPath:            bundle.vaePath,
			DiffusionFlashAttn: true,
			OffloadParamsToCPU: true,
			FlowShift:          3.0,
		}

		editEngine, err = sd.NewStableDiffusion(editCtxParams)
		if err != nil {
			log.Warn(ctx, "startup", "status", "failed to init SD edit engine", "error", err)
			editEngine = nil
		} else if !editEngine.IsReady() {
			log.Warn(ctx, "startup", "status", "stable diffusion edit engine created but not ready")
			editEngine = nil
		}
	}

	log.Info(ctx, "startup", "status", "stable diffusion ready", "edit_model_available", editEngine != nil)
	return generationEngine, editEngine, nil
}

func resolveStableDiffusionModelBundle(modelsPath string, hwSpecs *hardware.HardwareSpecs) (*stableDiffusionModelBundle, error) {
	downloader := modeldownloader.New(modelsPath)

	models := sdmodels.GetAvailableModels()
	if len(models) == 0 {
		return nil, errors.New("stable diffusion catalog is empty")
	}

	selectedModel := models[0]
	if hwSpecs != nil {
		selectedModel = sdmodels.SelectOptimalModel(&sdmodels.HardwareSpecs{
			TotalRAM:     hwSpecs.TotalRAM,
			AvailableRAM: hwSpecs.AvailableRAM,
			CPU:          hwSpecs.CPU,
			CPUCores:     hwSpecs.CPUCores,
			GPUMemory:    hwSpecs.GPUMemory,
			GPUModel:     hwSpecs.GPUModel,
		})
	}

	diffusionPath, err := findModelByFilename(modelsPath, selectedModel.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed searching diffusion model: %w", err)
	}
	if diffusionPath == "" {
		diffusionPath, err = downloader.DiscoverModel()
		if err != nil {
			return nil, fmt.Errorf("stable diffusion model not found: %w", err)
		}
	}
	if diffusionPath == "" {
		return nil, fmt.Errorf("stable diffusion model not found. "+SetupRequiredMsg, "download SD model")
	}

	llmPath, err := findModelByFilename(modelsPath, selectedModel.LLMFilename)
	if err != nil {
		return nil, fmt.Errorf("failed searching LLM model: %w", err)
	}
	if llmPath == "" && selectedModel.LLMFilename != "" {
		return nil, fmt.Errorf("required LLM model not found: %s. %s", selectedModel.LLMFilename, fmt.Sprintf(SetupRequiredMsg, "download SD model bundle"))
	}

	vaePath, err := findModelByFilename(modelsPath, selectedModel.VAEFilename)
	if err != nil {
		return nil, fmt.Errorf("failed searching VAE model: %w", err)
	}
	if vaePath == "" && selectedModel.VAEFilename != "" {
		return nil, fmt.Errorf("required VAE model not found: %s. %s", selectedModel.VAEFilename, fmt.Sprintf(SetupRequiredMsg, "download SD model bundle"))
	}

	editModelPath, err := findModelByFilename(modelsPath, selectedModel.EditModelFile)
	if err != nil {
		return nil, fmt.Errorf("failed searching edit model: %w", err)
	}
	if editModelPath == "" && selectedModel.EditFallbackFile != "" {
		editModelPath, err = findModelByFilename(modelsPath, selectedModel.EditFallbackFile)
		if err != nil {
			return nil, fmt.Errorf("failed searching fallback edit model: %w", err)
		}
	}
	if editModelPath == "" && selectedModel.EditModelFile != "" {
		return nil, fmt.Errorf("required edit model not found: %s. %s", selectedModel.EditModelFile, fmt.Sprintf(SetupRequiredMsg, "download SD model bundle"))
	}

	return &stableDiffusionModelBundle{
		selectedModel: selectedModel,
		diffusionPath: diffusionPath,
		llmPath:       llmPath,
		vaePath:       vaePath,
		editModelPath: editModelPath,
	}, nil
}

func findModelByFilename(root string, filename string) (string, error) {
	if filename == "" {
		return "", nil
	}

	var match string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.EqualFold(info.Name(), filename) {
			match = path
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return match, nil
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
		// Check if tunnel has active connections
		hasActive, err := tunnelkit.HasActiveConnections(tunnel.TunnelID)
		if err != nil {
			log.Warn(ctx, "tunnel", "status", "failed to check active connections", "tunnel", tunnel.Hostname, "error", err)
			continue
		}
		if hasActive {
			continue
		}

		// Start tunnel in goroutine with timeout
		tunnelURL := make(chan string, 1)
		errChan := make(chan error, 1)

		go func() {
			if err := tunnelkit.RunTunnel(ctx, tunnel.TunnelToken); err != nil {
				errChan <- err
			} else {
				tunnelURL <- tunnel.PublicURL
			}
		}()

		select {
		case url := <-tunnelURL:
			log.Info(ctx, "tunnel", "status", "started successfully", "url", tunnel.PublicURL)
			return url
		case err := <-errChan:
			log.Warn(ctx, "tunnel", "status", "failed to start", "tunnel", tunnel.Hostname, "error", err)
			// Continue to next tunnel
		case <-time.After(10 * time.Second):
			log.Warn(ctx, "tunnel", "status", "timeout waiting for tunnel", "tunnel", tunnel.Hostname)
			// Continue to next tunnel
		case <-ctx.Done():
			log.Info(ctx, "tunnel", "status", "context cancelled")
			return ""
		}
	}
	return ""
}
