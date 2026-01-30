package chatapp

import (
	"net/http"

	"github.com/kawai-network/veridium/cmd/server/app/sdk/cache"
	"github.com/kawai-network/veridium/cmd/server/app/sdk/mid"
	"github.com/kawai-network/veridium/cmd/server/foundation/logger"
	"github.com/kawai-network/veridium/cmd/server/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log *logger.Logger

	Cache *cache.Cache
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	api := newApp(cfg)

	// Use wallet-based authentication (API keys contain their own expiration)
	walletAuth := mid.WalletAuthenticate()

	app.HandlerFunc(http.MethodPost, version, "/chat/completions", api.chatCompletions, walletAuth)
}
