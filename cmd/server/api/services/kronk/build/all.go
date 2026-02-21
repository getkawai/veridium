// Package build binds all the routes into the specified app.
package build

import (
	"github.com/kawai-network/veridium/cmd/server/app/domain/chatapp"
	"github.com/kawai-network/veridium/cmd/server/app/domain/embedapp"
	"github.com/kawai-network/veridium/cmd/server/app/domain/healthcheck"
	"github.com/kawai-network/veridium/cmd/server/app/domain/imageapp"
	"github.com/kawai-network/veridium/cmd/server/app/domain/msgsapp"
	"github.com/kawai-network/veridium/cmd/server/app/domain/rerankapp"
	"github.com/kawai-network/veridium/cmd/server/app/domain/respapp"
	"github.com/kawai-network/veridium/cmd/server/app/domain/toolapp"
	"github.com/kawai-network/veridium/cmd/server/app/domain/ttsapp"
	"github.com/kawai-network/veridium/cmd/server/app/domain/whisperapp"
	"github.com/kawai-network/veridium/cmd/server/app/sdk/mux"
	"github.com/kawai-network/veridium/cmd/server/foundation/web"
)

// Routes constructs the all value which provides the implementation of
// of RouteAdder for specifying what routes to bind to this instance.
func Routes() all {
	return all{}
}

type all struct{}

// Add implements the RouterAdder interface.
func (all) Add(app *web.App, cfg mux.Config) {
	healthcheck.Routes(app, healthcheck.Config{
		Log:   cfg.Log,
		Cache: cfg.Cache,
	})

	toolapp.Routes(app, toolapp.Config{
		Log: cfg.Log,

		Cache:     cfg.Cache,
		Libs:      cfg.Libs,
		Models:    cfg.Models,
		Catalog:   cfg.Catalog,
		Templates: cfg.Templates,
	})

	chatapp.Routes(app, chatapp.Config{
		Log: cfg.Log,

		Cache: cfg.Cache,
	})

	embedapp.Routes(app, embedapp.Config{
		Log: cfg.Log,

		Cache: cfg.Cache,
	})

	rerankapp.Routes(app, rerankapp.Config{
		Log: cfg.Log,

		Cache: cfg.Cache,
	})

	respapp.Routes(app, respapp.Config{
		Log: cfg.Log,

		Cache: cfg.Cache,
	})

	msgsapp.Routes(app, msgsapp.Config{
		Log: cfg.Log,

		Cache: cfg.Cache,
	})

	imageapp.Routes(app, imageapp.Config{
		Log: cfg.Log,

		Engine: cfg.ImageEngine,
	})

	whisperapp.Routes(app, whisperapp.RouteConfig{
		Log: cfg.Log,
	})

	ttsapp.Routes(app, ttsapp.RouteConfig{
		Log: cfg.Log,
	})
}
