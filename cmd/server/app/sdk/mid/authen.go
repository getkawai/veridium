package mid

import (
	"context"
	"net/http"

	"github.com/kawai-network/veridium/cmd/server/app/sdk/authclient"
	"github.com/kawai-network/veridium/cmd/server/app/sdk/errs"
	"github.com/kawai-network/veridium/cmd/server/foundation/web"
)

// Authenticate calls out to the auth service to authenticate the call.
func Authenticate(client *authclient.Client, admin bool, endpoint string) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			ar, err := client.Authenticate(ctx, r.Header.Get("authorization"), admin, endpoint)
			if err != nil {
				return errs.New(errs.Unauthenticated, err)
			}

			ctx = setSubject(ctx, ar.Subject)

			return next(ctx, r)
		}

		return h
	}

	return m
}
