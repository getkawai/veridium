package mid

import (
	"context"
	"net/http"
	"runtime/debug"

	"github.com/kawai-network/veridium/cmd/server/app/sdk/errs"
	"github.com/kawai-network/veridium/cmd/server/foundation/web"
)

// Panics recovers from panics and converts the panic to an error so it is
// reported in Metrics and handled in Errors.
func Panics() web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) (resp web.Encoder) {

			// Defer a function to recover from a panic and set the err return
			// variable after the fact.
			defer func() {
				if rec := recover(); rec != nil {
					trace := debug.Stack()
					resp = errs.Errorf(errs.InternalOnlyLog, "PANIC [%v] TRACE[%s]", rec, string(trace))
				}
			}()

			return next(ctx, r)
		}

		return h
	}

	return m
}
