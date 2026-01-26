package apitest

import (
	"bytes"
	"context"
	"net"
	"testing"
	"time"

	"github.com/kawai-network/veridium/cmd/server/api/services/kronk/build"
	"github.com/kawai-network/veridium/cmd/server/app/domain/authapp"
	"github.com/kawai-network/veridium/cmd/server/app/sdk/authclient"
	"github.com/kawai-network/veridium/cmd/server/app/sdk/cache"
	"github.com/kawai-network/veridium/cmd/server/app/sdk/mux"
	"github.com/kawai-network/veridium/cmd/server/app/sdk/security"
	"github.com/kawai-network/veridium/cmd/server/app/sdk/security/auth"
	"github.com/kawai-network/veridium/cmd/server/foundation/logger"
	"github.com/kawai-network/veridium/cmd/server/foundation/web"
	"github.com/kawai-network/veridium/pkg/kronk"
	"github.com/kawai-network/veridium/pkg/tools/catalog"
	"github.com/kawai-network/veridium/pkg/tools/defaults"
	"github.com/kawai-network/veridium/pkg/tools/libs"
	"github.com/kawai-network/veridium/pkg/tools/models"
	"github.com/kawai-network/veridium/pkg/tools/templates"
	"google.golang.org/grpc/test/bufconn"
)

// New initialized the system to run a test.
func New(t *testing.T, testName string) *Test {
	ctx := context.Background()

	// -------------------------------------------------------------------------

	var buf bytes.Buffer
	log := logger.New(&buf, logger.LevelInfo, "TEST", web.GetTraceID)

	// -------------------------------------------------------------------------

	auth, err := auth.New(auth.Config{
		KeyLookup: &keyStore{},
		Issuer:    "kronk project",
	})

	if err != nil {
		t.Fatal(err)
	}

	// -------------------------------------------------------------------------

	var authClientOpts []func(*authclient.Client)

	// If no host is provided for the auth service, we will start it ourselves
	// with a bufconn listener.
	sec, err := security.New(security.Config{
		Issuer: auth.Issuer(),
	})

	if err != nil {
		t.Fatal(err)
	}

	log.Info(ctx, "startup", "status", "starting auth server")

	lis := bufconn.Listen(1024 * 1024)

	authApp := authapp.Start(ctx, authapp.Config{
		Log:      log,
		Security: sec,
		Listener: lis,
		Tracer:   nil,
		Enabled:  true,
	})

	authClientOpts = append(authClientOpts, authclient.WithDialer(func(ctx context.Context, _ string) (net.Conn, error) {
		return lis.Dial()
	}))

	// -------------------------------------------------------------------------

	authHost := ""
	if len(authClientOpts) > 0 {
		authHost = "passthrough:///bufnet"
	}

	authClient, err := authclient.New(log, authHost, authClientOpts...)
	if err != nil {
		t.Fatal(err)
	}

	// -------------------------------------------------------------------------

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	libs, err := libs.New(
		libs.WithVersion(defaults.LibVersion("")),
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

		authClient.Close()
		authApp.Shutdown(ctx)
		sec.Close()

		t.Logf("******************** LOGS (%s) ********************\n\n", testName)
		t.Log(buf.String())
		t.Logf("******************** LOGS (%s) ********************\n", testName)
	})

	// -------------------------------------------------------------------------

	cfgMux := mux.Config{
		Build:      "test",
		Log:        log,
		AuthClient: authClient,
		Tracer:     nil,
		Cache:      cache,
		Libs:       libs,
		Models:     models,
		Catalog:    ctlg,
		Templates:  tmplts,
	}

	mux := mux.WebAPI(cfgMux,
		build.Routes(),
	)

	return &Test{
		Sec: sec,
		mux: mux,
	}
}
