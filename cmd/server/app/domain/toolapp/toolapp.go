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
	libTypeParam := r.URL.Query().Get("type")

	if libTypeParam != "" {
		libType := libs.ParseLibraryType(libTypeParam)
		lib, err := libs.New(
			libs.WithLibraryType(libType),
		)
		if err != nil {
			return errs.New(errs.Internal, err)
		}

		versionTag, err := lib.VersionInformation()
		if err != nil {
			return errs.New(errs.Internal, err)
		}

		return toAppVersionTag("retrieve", versionTag)
	}

	var allTags []libs.VersionTag
	for _, libType := range libs.AllLibraryTypes() {
		lib, err := libs.New(
			libs.WithLibraryType(libType),
		)
		if err != nil {
			continue
		}

		versionTag, err := lib.VersionInformation()
		if err != nil {
			versionTag = libs.VersionTag{
				Library: libType,
				Version: "",
				Latest:  "",
			}
		}
		allTags = append(allTags, versionTag)
	}

	return toLibsListResponse(allTags)
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

func (a *app) pullLibs(ctx context.Context, r *http.Request) web.Encoder {
	libTypeParam := r.URL.Query().Get("type")

	if libTypeParam == "" {
		return a.pullAllLibs(ctx, r)
	}

	libType := libs.ParseLibraryType(libTypeParam)
	return a.pullSingleLib(ctx, r, libType)
}

func (a *app) pullSingleLib(ctx context.Context, r *http.Request, libType libs.LibraryType) web.Encoder {
	w := web.GetWriter(ctx)
	f, ok := w.(http.Flusher)
	if !ok {
		return errs.New(errs.Internal, fmt.Errorf("streaming not supported"))
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	lib, err := libs.New(
		libs.WithLibraryType(libType),
		libs.WithAllowUpgrade(true),
	)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	log := func(ctx context.Context, msg string, args ...any) {
		vt := libs.VersionTag{Library: libType}
		status := fmt.Sprintf(msg, args...)
		resp := toAppVersion(status, vt)
		w.Write([]byte(resp))
		f.Flush()
	}

	vt, err := lib.Download(ctx, log)
	if err != nil {
		errResp := fmt.Sprintf("data: {\"status\":\"error\",\"error\":%q}\n\n", err.Error())
		w.Write([]byte(errResp))
		f.Flush()
		return web.NewNoResponse()
	}

	finalResp := fmt.Sprintf("data: {\"status\":\"complete\",\"library\":%q,\"version\":%q}\n\n", vt.Library.String(), vt.Version)
	w.Write([]byte(finalResp))
	w.Write([]byte("data: [DONE]\n\n"))
	f.Flush()

	return web.NewNoResponse()
}

func (a *app) pullAllLibs(ctx context.Context, r *http.Request) web.Encoder {
	w := web.GetWriter(ctx)
	f, ok := w.(http.Flusher)
	if !ok {
		return errs.New(errs.Internal, fmt.Errorf("streaming not supported"))
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for _, libType := range libs.AllLibraryTypes() {
		lib, err := libs.New(
			libs.WithLibraryType(libType),
			libs.WithAllowUpgrade(true),
		)
		if err != nil {
			errResp := fmt.Sprintf("data: {\"status\":\"error\",\"library\":%q,\"error\":%q}\n\n", libType.String(), err.Error())
			w.Write([]byte(errResp))
			f.Flush()
			continue
		}

		log := func(ctx context.Context, msg string, args ...any) {
			vt := libs.VersionTag{Library: libType}
			status := fmt.Sprintf(msg, args...)
			resp := toAppVersion(status, vt)
			w.Write([]byte(resp))
			f.Flush()
		}

		vt, err := lib.Download(ctx, log)
		if err != nil {
			errResp := fmt.Sprintf("data: {\"status\":\"error\",\"library\":%q,\"error\":%q}\n\n", libType.String(), err.Error())
			w.Write([]byte(errResp))
			f.Flush()
			continue
		}

		finalResp := fmt.Sprintf("data: {\"status\":\"complete\",\"library\":%q,\"version\":%q}\n\n", vt.Library.String(), vt.Version)
		w.Write([]byte(finalResp))
		f.Flush()
	}

	w.Write([]byte("data: [DONE]\n\n"))
	f.Flush()

	return web.NewNoResponse()
}
