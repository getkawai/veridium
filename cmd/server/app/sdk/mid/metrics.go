package mid

import (
	"context"
	"net/http"

	"github.com/kawai-network/veridium/cmd/server/foundation/web"
)

// Metrics updates program counters.
func Metrics() web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			resp := next(ctx, r)

			// Metrics collection removed

			return resp
		}

		return h
	}

	return m
}
