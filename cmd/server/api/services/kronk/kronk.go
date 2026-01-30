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
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/getsentry/sentry-go"
	"github.com/kawai-network/veridium/cmd/server/api/services/kronk/build"
	"github.com/kawai-network/veridium/cmd/server/app/sdk/cache"
	"github.com/kawai-network/veridium/cmd/server/app/sdk/mux"

	"github.com/kawai-network/veridium/cmd/server/foundation/logger"
	"github.com/kawai-network/veridium/cmd/server/foundation/web"
	"github.com/kawai-network/veridium/pkg/kronk"
	pkglogger "github.com/kawai-network/veridium/pkg/logger"
	"github.com/kawai-network/veridium/pkg/tools/catalog"
	"github.com/kawai-network/veridium/pkg/tools/defaults"
	"github.com/kawai-network/veridium/pkg/tools/libs"
	"github.com/kawai-network/veridium/pkg/tools/models"
	"github.com/kawai-network/veridium/pkg/tools/templates"
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
			ReadTimeout        time.Duration `conf:"default:30s"`
			WriteTimeout       time.Duration `conf:"default:15m"`
			IdleTimeout        time.Duration `conf:"default:1m"`
			ShutdownTimeout    time.Duration `conf:"default:1m"`
			APIHost            string        `conf:"default:localhost:8080"`
			CORSAllowedOrigins []string      `conf:"default:*"`
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
	// Start API Service

	log.Info(ctx, "startup", "status", "initializing V1 API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	cfgMux := mux.Config{
		Build: tag,
		Log:   log,

		Tracer:    nil,
		Cache:     cache,
		Libs:      libs,
		Models:    models,
		Catalog:   ctlg,
		Templates: tmplts,
	}

	webAPI := mux.WebAPI(cfgMux,
		build.Routes(),
		mux.WithCORS(cfg.Web.CORSAllowedOrigins),
		mux.WithFileServer(true, static, "static", "/", []string{"v1"}),
	)

	api := http.Server{
		Addr:         cfg.Web.APIHost,
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
