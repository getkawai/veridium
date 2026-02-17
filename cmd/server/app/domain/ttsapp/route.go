package ttsapp

import (
	"net/http"

	"github.com/kawai-network/veridium/cmd/server/app/sdk/mid"
	"github.com/kawai-network/veridium/cmd/server/foundation/logger"
	"github.com/kawai-network/veridium/cmd/server/foundation/web"
)

// RouteConfig contains all the mandatory systems required by handlers.
type RouteConfig struct {
	Log *logger.Logger
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg RouteConfig) {
	const version = "v1"

	api := newApp(Config{
		Log: cfg.Log,
	})

	// Use wallet-based authentication
	walletAuth := mid.WalletAuthenticate()

	// OpenAI-compatible text-to-speech endpoint
	app.HandlerFunc(http.MethodPost, version, "/audio/speech", api.generations, walletAuth)
}
