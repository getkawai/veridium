package whisperapp

import (
	"net/http"

	"github.com/kawai-network/veridium/cmd/server/app/sdk/mid"
	"github.com/kawai-network/veridium/cmd/server/foundation/logger"
	"github.com/kawai-network/veridium/cmd/server/foundation/web"
)

// RouteConfig contains all the mandatory systems required by handlers.
type RouteConfig struct {
	Log              *logger.Logger
	WhisperModelsDir string
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg RouteConfig) {
	const version = "v1"

	api := newApp(Config{
		Log: cfg.Log,
	})

	// Use wallet-based authentication
	walletAuth := mid.WalletAuthenticate()

	// OpenAI-compatible audio transcription endpoint
	app.HandlerFunc(http.MethodPost, version, "/audio/transcriptions", api.transcriptions, walletAuth)
}
