package mid

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/kawai-network/veridium/cmd/server/app/sdk/errs"
	"github.com/kawai-network/veridium/cmd/server/foundation/web"
	"github.com/kawai-network/veridium/pkg/apikey"
)

// WalletAuthenticate validates wallet-based API keys and sets wallet address in context.
// This middleware supports stateless authentication using encoded wallet addresses.
//
// The API key itself contains the expiration time, so no configuration is needed here.
//
// Example usage:
//
//	walletAuth := mid.WalletAuthenticate()
//	app.HandlerFunc(http.MethodPost, "v1", "/chat/completions", api.chatCompletions, walletAuth)
func WalletAuthenticate() web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			// Extract Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				return errs.New(errs.Unauthenticated, errors.New("missing authorization header"))
			}

			// Remove "Bearer " prefix (OpenAI spec)
			apiKey := strings.TrimPrefix(authHeader, "Bearer ")
			if apiKey == authHeader {
				// No "Bearer " prefix found
				return errs.New(errs.Unauthenticated, errors.New("invalid authorization format"))
			}

			// Validate API key
			payload, err := apikey.Validate(apiKey)
			if err != nil {
				return errs.New(errs.Unauthenticated, err)
			}

			// Check expiration (enforced by the key itself)
			if payload.IsExpired() {
				return errs.New(errs.Unauthenticated, errors.New("API key expired"))
			}

			// Set wallet address and metadata in context
			ctx = setWalletAddress(ctx, payload.GetWalletAddress())
			ctx = setKeyIssuedAt(ctx, payload.GetIssuedAt())
			ctx = setKeyNonce(ctx, payload.GetNonce())

			// Also set subject for compatibility with existing code
			ctx = setSubject(ctx, payload.GetWalletAddress())

			return next(ctx, r)
		}

		return h
	}

	return m
}

// DualAuthenticate supports wallet-based authentication.
// It tries wallet-based auth (if key starts with "vk-").
//
// Currently this function is kept for structural compatibility but only
// supports wallet-based authentication.
//
// Example usage:
//
//	dualAuth := mid.DualAuthenticate(nil, false, "chat-completions")
func DualAuthenticate(client interface{}, admin bool, endpoint string) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				return errs.New(errs.Unauthenticated, errors.New("missing authorization header"))
			}

			apiKey := strings.TrimPrefix(authHeader, "Bearer ")
			if apiKey == authHeader {
				// No "Bearer " prefix found
				return errs.New(errs.Unauthenticated, errors.New("invalid authorization format"))
			}

			// Try wallet-based auth first (check for "vk-" prefix)
			if strings.HasPrefix(apiKey, "vk-") {
				payload, err := apikey.Validate(apiKey)
				if err == nil && !payload.IsExpired() {
					ctx = setWalletAddress(ctx, payload.GetWalletAddress())
					ctx = setKeyIssuedAt(ctx, payload.GetIssuedAt())
					ctx = setKeyNonce(ctx, payload.GetNonce())
					ctx = setSubject(ctx, payload.GetWalletAddress())
					return next(ctx, r)
				}
			}

			// Fallback: Currently only wallet-based auth is supported.
			return errs.New(errs.Unauthenticated, errors.New("invalid credentials"))
		}

		return h
	}

	return m
}
