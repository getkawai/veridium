package apitest

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/kawai-network/veridium/cmd/server/api/services/kronk/build"
	"github.com/kawai-network/veridium/cmd/server/app/sdk/cache"
	"github.com/kawai-network/veridium/cmd/server/app/sdk/mux"
	"github.com/kawai-network/veridium/cmd/server/foundation/logger"
	"github.com/kawai-network/veridium/cmd/server/foundation/web"
	"github.com/kawai-network/veridium/pkg/kronk"
	"github.com/kawai-network/veridium/pkg/tools/catalog"
	"github.com/kawai-network/veridium/pkg/tools/defaults"
	"github.com/kawai-network/veridium/pkg/tools/libs"
	"github.com/kawai-network/veridium/pkg/tools/models"
	"github.com/kawai-network/veridium/pkg/tools/templates"
)

// New initialized the system to run a test.
func New(t *testing.T, testName string) *Test {
	ctx := context.Background()

	// -------------------------------------------------------------------------

	var buf bytes.Buffer
	log := logger.New(&buf, logger.LevelInfo, "TEST", web.GetTraceID)

	// -------------------------------------------------------------------------

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	libs, err := libs.New(
		libs.WithVersion(defaults.LibVersion("")),
		libs.WithLibraryType(libs.LibraryLlama),
	)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := libs.Download(ctx, log.Info); err != nil {
		t.Fatal(err)
	}

	// -------------------------------------------------------------------------

	models, err := models.New()
	if err != nil {
		t.Fatal(err)
	}

	// -------------------------------------------------------------------------
	// Catalog System

	ctlg, err := catalog.New()
	if err != nil {
		t.Fatal(err)
	}

	if err := ctlg.Download(ctx, catalog.WithLogger(log.Info)); err != nil {
		t.Fatal(err)
	}

	// -------------------------------------------------------------------------
	// Template System

	tmplts, err := templates.New(templates.WithCatalog(ctlg))
	if err != nil {
		t.Fatal(err)
	}

	if err := tmplts.Download(ctx, templates.WithLogger(log.Info)); err != nil {
		t.Fatal(err)
	}

	// -------------------------------------------------------------------------
	// Init Kronk

	if err := kronk.Init(); err != nil {
		t.Fatal(err)
	}

	cache, err := cache.New(cache.Config{
		Log:             log.Info,
		Templates:       tmplts,
		ModelsInCache:   1,
		CacheTTL:        5 * time.Minute,
		ModelConfigFile: "../../../../../../zarf/kms/model_config.yaml",
	})

	if err != nil {
		t.Fatal(err)
	}

	// -------------------------------------------------------------------------

	t.Cleanup(func() {
		t.Helper()

		ctx := context.Background()

		if err := cache.Shutdown(ctx); err != nil {
			t.Fatal(err)
		}

		t.Logf("******************** LOGS (%s) ********************\n\n", testName)
		t.Log(buf.String())
		t.Logf("******************** LOGS (%s) ********************\n", testName)
	})

	// -------------------------------------------------------------------------

	cfgMux := mux.Config{
		Build:     "test",
		Log:       log,
		Tracer:    nil,
		Cache:     cache,
		Libs:      libs,
		Models:    models,
		Catalog:   ctlg,
		Templates: tmplts,
	}

	mux := mux.WebAPI(cfgMux,
		build.Routes(),
	)

	return &Test{
		mux: mux,
	}
}
