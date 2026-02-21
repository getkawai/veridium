// Package healthcheck maintains the app layer api for contributor health status.
package healthcheck

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/kawai-network/veridium/cmd/server/app/sdk/cache"
	"github.com/kawai-network/veridium/cmd/server/foundation/logger"
	"github.com/kawai-network/veridium/cmd/server/foundation/web"
	"github.com/kawai-network/y/types"
)

type app struct {
	cache *cache.Cache
	log   *logger.Logger
}

func newApp(cfg Config) *app {
	return &app{
		cache: cfg.Cache,
		log:   cfg.Log,
	}
}

// HealthResponse represents the health status of the contributor
type HealthResponse struct {
	Status         types.ContributorStatus `json:"status"`          // "online" or "offline"
	ActiveRequests int64                         `json:"active_requests"` // Current number of active requests
	IsBusy         bool                          `json:"is_busy"`         // True if handling requests
}

// Encode implements the encoder interface.
func (h HealthResponse) Encode() ([]byte, string, error) {
	data, err := json.Marshal(h)
	return data, "application/json", err
}

func (a *app) health(ctx context.Context, r *http.Request) web.Encoder {
	// Get model status from cache
	modelStatus, err := a.cache.ModelStatus()
	if err != nil {
		a.log.Info(ctx, "health", "status", "error", "error", err)
		// Return offline status
		return HealthResponse{
			Status:         types.StatusOffline,
			ActiveRequests: 0,
			IsBusy:         false,
		}
	}

	// Calculate total active requests
	var totalActive int64
	
	for _, model := range modelStatus {
		totalActive += int64(model.ActiveStreams)
	}

	// Determine if contributor is busy
	isBusy := totalActive > 0

	// Contributor is online if we can successfully query the cache
	response := HealthResponse{
		Status:         types.StatusOnline,
		ActiveRequests: totalActive,
		IsBusy:         isBusy,
	}

	return response
}
