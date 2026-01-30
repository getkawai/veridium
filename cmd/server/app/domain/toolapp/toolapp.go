// Package toolapp provides endpoints to handle tool management.
package toolapp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kawai-network/veridium/cmd/server/app/sdk/cache"
	"github.com/kawai-network/veridium/cmd/server/app/sdk/errs"
	"github.com/kawai-network/veridium/cmd/server/foundation/logger"
	"github.com/kawai-network/veridium/cmd/server/foundation/web"
	"github.com/kawai-network/veridium/pkg/tools/catalog"
	"github.com/kawai-network/veridium/pkg/tools/libs"
	"github.com/kawai-network/veridium/pkg/tools/models"
	"github.com/kawai-network/veridium/pkg/tools/templates"
)

type app struct {
	log       *logger.Logger
	cache     *cache.Cache
	libs      *libs.Libs
	models    *models.Models
	catalog   *catalog.Catalog
	templates *templates.Templates
}

func newApp(cfg Config) *app {
	return &app{
		log:       cfg.Log,
		cache:     cfg.Cache,
		libs:      cfg.Libs,
		models:    cfg.Models,
		catalog:   cfg.Catalog,
		templates: cfg.Templates,
	}
}

func (a *app) listLibs(ctx context.Context, r *http.Request) web.Encoder {
	versionTag, err := a.libs.VersionInformation()
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return toAppVersionTag("retrieve", versionTag)
}

func (a *app) listModels(ctx context.Context, r *http.Request) web.Encoder {
	models, err := a.models.RetrieveFiles()
	if err != nil {
		return errs.Errorf(errs.Internal, "unable to retrieve model list: %s", err)
	}

	return toListModelsInfo(models)
}

func (a *app) missingModel(ctx context.Context, r *http.Request) web.Encoder {
	return errs.New(errs.InvalidArgument, fmt.Errorf("model parameter is required"))
}

func (a *app) showModel(ctx context.Context, r *http.Request) web.Encoder {
	modelID := web.Param(r, "model")

	mi, err := a.models.RetrieveInfo(modelID)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	krn, err := a.cache.AquireModel(ctx, mi.ID)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return toModelInfo(mi, krn.ModelInfo())
}

func (a *app) modelPS(ctx context.Context, r *http.Request) web.Encoder {
	models, err := a.cache.ModelStatus()
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	a.log.Info(ctx, "models", "len", len(models))

	return toModelDetails(models)
}

func (a *app) listCatalog(ctx context.Context, r *http.Request) web.Encoder {
	filterCategory := web.Param(r, "filter")

	list, err := a.catalog.CatalogModelList(filterCategory)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return toCatalogModelsResponse(list)
}

func (a *app) showCatalogModel(ctx context.Context, r *http.Request) web.Encoder {
	modelID := web.Param(r, "model")

	model, err := a.catalog.RetrieveModelDetails(modelID)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return toCatalogModelResponse(model)
}

func (a *app) calculateVRAM(ctx context.Context, r *http.Request) web.Encoder {
	var req VRAMRequest
	if err := web.Decode(r, &req); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	cfg := models.VRAMConfig{
		ContextWindow:   req.ContextWindow,
		BytesPerElement: req.BytesPerElement,
		Slots:           req.Slots,
		CacheSequences:  req.CacheSequences,
	}

	vram, err := a.models.CalculateVRAM(req.ModelID, cfg)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return toVRAMResponse(req.ModelID, vram)
}
