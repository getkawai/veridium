// Package kronk is the model server.
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
	"strings"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/getsentry/sentry-go"
	"github.com/kawai-network/veridium/cmd/server/api/services/kronk/build"
	"github.com/kawai-network/veridium/cmd/server/app/sdk/cache"
	"github.com/kawai-network/veridium/cmd/server/app/sdk/mux"
	"github.com/kawai-network/veridium/cmd/server/foundation/logger"
	"github.com/kawai-network/veridium/cmd/server/foundation/web"
	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/internal/image"
	"github.com/kawai-network/veridium/internal/services"
	"github.com/kawai-network/veridium/internal/whisper"
	"github.com/kawai-network/veridium/pkg/blockchain"
	"github.com/kawai-network/veridium/pkg/hardware"
	"github.com/kawai-network/veridium/pkg/kronk"
	pkglogger "github.com/kawai-network/veridium/pkg/logger"
	"github.com/kawai-network/veridium/pkg/store"
	"github.com/kawai-network/veridium/pkg/tools/catalog"
	"github.com/kawai-network/veridium/pkg/tools/defaults"
	"github.com/kawai-network/veridium/pkg/tools/libs"
	"github.com/kawai-network/veridium/pkg/tools/models"
	"github.com/kawai-network/veridium/pkg/tools/templates"
	"github.com/kawai-network/veridium/pkg/tunnelkit"
)

//go:embed static
var static embed.FS

var tag = "develop"

const SentryDSN = "https://6d138acbdde2516e32e24f016b472031@o4510620614983680.ingest.us.sentry.io/4510620618850304"

func Run(showHelp bool) error {
	var log *logger.Logger

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
		baseHandler := slog.NewJSONHandler(os.Stdout, nil)
		sentryHandler = pkglogger.NewSentryHandler(baseHandler)
	}

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* SEND ALERT *******")
		},
	}

	if sentryHandler != nil {
		log = logger.NewWithSentry(os.Stdout, logger.LevelInfo, "KRONK", web.GetTraceID, sentryHandler)
	} else {
		log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "KRONK", web.GetTraceID, events)
	}

	log.Info(context.Background(), "sentry", "enabled", err == nil)

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

	log.Info(ctx, "startup", "status", "downloading libraries")

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

	libs, err := libs.New(
		libs.WithBasePath(cfg.LibPath),
		libs.WithArch(arch),
		libs.WithOS(opSys),
		libs.WithProcessor(processor),
		libs.WithAllowUpgrade(cfg.AllowUpgrade),
		libs.WithVersion(defaults.LibVersion(cfg.LibVersion)),
	)
	if err != nil {
		return fmt.Errorf("unable to create libs api: %w", err)
	}

	log.Info(ctx, "startup", "status", "installing/updating libraries", "libPath", libs.LibsPath(), "arch", libs.Arch(), "os", libs.OS(), "processor", libs.Processor(), "update", true)

	downloadCtx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	if _, err := libs.Download(downloadCtx, log.Info); err != nil {
		return fmt.Errorf("unable to install llama.cpp: %w", err)
	}

	// -------------------------------------------------------------------------
	// Model System

	models, err := models.NewWithPaths(cfg.BasePath)
	if err != nil {
		return fmt.Errorf("unable to create catalog system: %w", err)
	}

	log.Info(ctx, "startup", "status", "model integrity checks, may take a few seconds")

	models.BuildIndex(log.Info)

	// -------------------------------------------------------------------------
	// Catalog System

	log.Info(ctx, "startup", "status", "downloading catalog")

	ctlg, err := catalog.New(
		catalog.WithBasePath(cfg.BasePath),
		catalog.WithGithubRepo(cfg.Catalog.GithubRepo))
	if err != nil {
		return fmt.Errorf("unable to create catalog system: %w", err)
	}

	if err := ctlg.Download(ctx, catalog.WithLogger(log.Info)); err != nil {
		return fmt.Errorf("unable to download catalog: %w", err)
	}

	// -------------------------------------------------------------------------
	// Template System

	log.Info(ctx, "startup", "status", "downloading templates")

	tmplts, err := templates.New(
		templates.WithBasePath(cfg.BasePath),
		templates.WithGithubRepo(cfg.Templates.GithubRepo),
		templates.WithCatalog(ctlg))
	if err != nil {
		return fmt.Errorf("unable to create template system: %w", err)
	}

	if err := tmplts.Download(ctx, templates.WithLogger(log.Info)); err != nil {
		return fmt.Errorf("unable to download templates: %w", err)
	}

	// -------------------------------------------------------------------------
	// Init Kronk

	log.Info(ctx, "startup", "status", "initializing kronk")

	if err := kronk.Init(); err != nil {
		return fmt.Errorf("installation invalid: %w", err)
	}

	cache, err := cache.New(cache.Config{
		Log:                  log.Info,
		BasePath:             cfg.BasePath,
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

	// Setup Wallet (Interactive only)
	wallet := services.NewWalletService("", kv)

	// Interactive wallet setup
	if !wallet.HasWallet() {
		// No wallet exists - create new one
		printInfo("No wallet found. Let's create one!")

		choice, err := promptChoice("Choose setup method:", []string{
			"Generate new mnemonic (recommended)",
			"Import existing mnemonic",
		})
		if err != nil {
			return fmt.Errorf("failed to get user choice: %w", err)
		}

		// Get password
		password, err := promptPassword("Enter password (min 8 characters): ")
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		if err := validatePassword(password); err != nil {
			return fmt.Errorf("invalid password: %w", err)
		}

		confirmPassword, err := promptPassword("Confirm password: ")
		if err != nil {
			return fmt.Errorf("failed to read password confirmation: %w", err)
		}
		if password != confirmPassword {
			return fmt.Errorf("passwords do not match")
		}

		var mnemonic string
		if choice == 0 {
			// Generate new mnemonic
			mnemonic, err = wallet.GenerateMnemonic()
			if err != nil {
				return fmt.Errorf("failed to generate mnemonic: %w", err)
			}
			printMnemonic(mnemonic)

			if !promptYesNo("Have you written down your mnemonic?") {
				return fmt.Errorf("please write down your mnemonic before continuing")
			}
		} else {
			// Import existing mnemonic
			printInfo("Enter your 12 or 24 word mnemonic phrase")
			mnemonic, err = promptPassword("Mnemonic (hidden): ")
			if err != nil {
				return fmt.Errorf("failed to read mnemonic: %w", err)
			}
			// CRITICAL: Trim whitespace to ensure consistent wallet generation
			// Copy-paste can introduce trailing spaces/newlines that would generate
			// a different wallet address than the standard mnemonic
			mnemonic = strings.Join(strings.Fields(mnemonic), " ")
			if err := validateMnemonic(mnemonic); err != nil {
				return fmt.Errorf("invalid mnemonic: %w", err)
			}
		}

		description, err := promptInput("Wallet name (e.g. My Contributor Wallet): ")
		if err != nil {
			return fmt.Errorf("failed to read wallet name: %w", err)
		}
		if description == "" {
			description = "Kronk Contributor"
		}

		walletAddress, err = wallet.CreateWallet(password, mnemonic, description)
		if err != nil {
			return fmt.Errorf("failed to create wallet: %w", err)
		}
		printSuccess(fmt.Sprintf("Wallet created: %s", walletAddress))
	} else {
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
	}

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

	// Start Heartbeat (fixed 30s interval)
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
				if err := kv.UpdateHeartbeat(ctx, walletAddress); err != nil {
					log.Info(ctx, "heartbeat", "status", "failed", "error", err)
				}
			}
		}
	}()
	log.Info(ctx, "startup", "status", "heartbeat started", "interval", "30s")

	// Initialize Whisper Service
	_, err = whisper.NewService()
	if err != nil {
		log.Info(ctx, "whisper", "status", "initialization failed", "error", err)
	} else {
		log.Info(ctx, "startup", "status", "whisper service ready")
	}

	// Initialize Stable Diffusion
	imageEngine := image.NewEngine()
	log.Info(ctx, "startup", "status", "stable diffusion ready")

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
██╗  ██╗██████╗  ██████╗ ███╗   ██╗██╗  ██╗    ███╗   ███╗███████╗
██║ ██╔╝██╔══██╗██╔═══██╗████╗  ██║██║ ██╔╝    ████╗ ████║██╔════╝
█████╔╝ ██████╔╝██║   ██║██╔██╗ ██║█████╔╝     ██╔████╔██║███████╗
██╔═██╗ ██╔══██╗██║   ██║██║╚██╗██║██╔═██╗     ██║╚██╔╝██║╚════██║
██║  ██╗██║  ██║╚██████╔╝██║ ╚████║██║  ██╗    ██║ ╚═╝ ██║███████║
╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═══╝╚═╝  ╚═╝    ╚═╝     ╚═╝╚══════╝                                                                                         
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
