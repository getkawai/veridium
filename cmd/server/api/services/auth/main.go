package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/ardanlabs/conf/v3"
	"github.com/kawai-network/veridium/cmd/server/app/domain/authapp"
	"github.com/kawai-network/veridium/cmd/server/app/sdk/security"
	"github.com/kawai-network/veridium/cmd/server/foundation/logger"
	"github.com/kawai-network/veridium/cmd/server/foundation/web"
)

var tag = "develop"

func main() {
	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* SEND ALERT *******")
		},
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "AUTH", web.GetTraceID, events)

	// -------------------------------------------------------------------------

	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {

	// -------------------------------------------------------------------------
	// GOMAXPROCS

	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------
	// Configuration

	cfg := struct {
		conf.Version
		Auth struct {
			Host    string `conf:"default:localhost:6000"`
			Issuer  string `conf:"default:kronk project"`
			Enabled bool   `conf:"default:false"`
		}
	}{
		Version: conf.Version{
			Build: tag,
			Desc:  "Auth",
		},
	}

	const prefix = "AUTH"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
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
	// Initialize authentication support

	log.Info(ctx, "startup", "status", "initializing authentication support")

	sec, err := security.New(security.Config{
		Issuer: cfg.Auth.Issuer,
	})

	if err != nil {
		return fmt.Errorf("unable to initialize security system: %w", err)
	}

	defer sec.Close()

	// -------------------------------------------------------------------------
	// Start Auth Service

	log.Info(ctx, "startup", "status", "initializing auth server")

	lis, err := net.Listen("tcp", cfg.Auth.Host)
	if err != nil {
		return fmt.Errorf("failed to listen on host %s : %w", cfg.Auth.Host, err)
	}

	authApp := authapp.Start(ctx, authapp.Config{
		Log:      log,
		Security: sec,
		Listener: lis,
		Tracer:   nil,
		Enabled:  cfg.Auth.Enabled,
	})

	defer authApp.Shutdown(ctx)

	// -------------------------------------------------------------------------
	// Wait and Shutdown

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	sig := <-shutdown

	log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
	defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

	return nil
}

var logo = `
██╗  ██╗██████╗  ██████╗ ███╗   ██╗██╗  ██╗     █████╗ ██╗   ██╗████████╗██╗  ██╗
██║ ██╔╝██╔══██╗██╔═══██╗████╗  ██║██║ ██╔╝    ██╔══██╗██║   ██║╚══██╔══╝██║  ██║
█████╔╝ ██████╔╝██║   ██║██╔██╗ ██║█████╔╝     ███████║██║   ██║   ██║   ███████║
██╔═██╗ ██╔══██╗██║   ██║██║╚██╗██║██╔═██╗     ██╔══██║██║   ██║   ██║   ██╔══██║
██║  ██╗██║  ██║╚██████╔╝██║ ╚████║██║  ██╗    ██║  ██║╚██████╔╝   ██║   ██║  ██║
╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═══╝╚═╝  ╚═╝    ╚═╝  ╚═╝ ╚═════╝    ╚═╝   ╚═╝  ╚═╝                                                                                 
`
