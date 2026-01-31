package imageapp

import (
	"net/http"

	"github.com/kawai-network/veridium/cmd/server/app/sdk/mid"
	"github.com/kawai-network/veridium/cmd/server/foundation/logger"
	"github.com/kawai-network/veridium/cmd/server/foundation/web"
	sd "github.com/kawai-network/veridium/pkg/stablediffusion"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log *logger.Logger

	Engine *sd.StableDiffusion
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	api := newApp(cfg)

	// Use wallet-based authentication
	walletAuth := mid.WalletAuthenticate()

	app.HandlerFunc(http.MethodPost, version, "/images/generations", api.generations, walletAuth)

	// Add file server route for generated images
	AddFileServerRoute(app, cfg)
}
