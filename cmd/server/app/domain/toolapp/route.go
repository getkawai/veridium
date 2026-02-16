package toolapp

import (
	"net/http"

	"github.com/kawai-network/veridium/cmd/server/app/sdk/cache"
	"github.com/kawai-network/veridium/cmd/server/app/sdk/mid"
	"github.com/kawai-network/veridium/cmd/server/foundation/logger"
	"github.com/kawai-network/veridium/cmd/server/foundation/web"
	"github.com/kawai-network/veridium/pkg/tools/catalog"
	"github.com/kawai-network/veridium/pkg/tools/libs"
	"github.com/kawai-network/veridium/pkg/tools/models"
	"github.com/kawai-network/veridium/pkg/tools/templates"
)

// Config contains all the mandatory systems required by handlers.
// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log       *logger.Logger
	Cache     *cache.Cache
	Libs      *libs.Libs
	Models    *models.Models
	Catalog   *catalog.Catalog
	Templates *templates.Templates
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	api := newApp(cfg)

	// Use wallet-based authentication for read-only operations
	walletAuth := mid.WalletAuthenticate()

	app.HandlerFunc(http.MethodGet, version, "/libs", api.listLibs, walletAuth)
	app.HandlerFunc(http.MethodPost, version, "/libs/pull", api.pullLibs, walletAuth)

	app.HandlerFunc(http.MethodGet, version, "/models", api.listModels, walletAuth)
	app.HandlerFunc(http.MethodGet, version, "/models/", api.missingModel, walletAuth)
	app.HandlerFunc(http.MethodGet, version, "/models/{model}", api.showModel, walletAuth)
	app.HandlerFunc(http.MethodGet, version, "/models/ps", api.modelPS, walletAuth)

	app.HandlerFunc(http.MethodGet, version, "/catalog", api.listCatalog, walletAuth)
	app.HandlerFunc(http.MethodGet, version, "/catalog/filter/{filter}", api.listCatalog, walletAuth)
	app.HandlerFunc(http.MethodGet, version, "/catalog/{model}", api.showCatalogModel, walletAuth)

	app.HandlerFunc(http.MethodPost, version, "/vram/calculate", api.calculateVRAM, walletAuth)
}
