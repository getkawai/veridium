// Package mid provides app level middleware support.
package mid

import (
	"context"
	"time"

	"github.com/kawai-network/veridium/cmd/server/foundation/web"
)

func checkIsError(e web.Encoder) error {
	err, hasError := e.(error)
	if hasError {
		return err
	}

	return nil
}

// =============================================================================

type ctxKey int

const (
	subjectKey ctxKey = iota + 1
	walletAddressKey
	keyIssuedAtKey
	keyNonceKey
)

func setSubject(ctx context.Context, subject string) context.Context {
	return context.WithValue(ctx, subjectKey, subject)
}

// GetSubject returns the subject from the context.
func GetSubject(ctx context.Context) string {
	v, ok := ctx.Value(subjectKey).(string)
	if !ok {
		return ""
	}
	return v
}

// =============================================================================
// Wallet Authentication Context

func setWalletAddress(ctx context.Context, address string) context.Context {
	return context.WithValue(ctx, walletAddressKey, address)
}

// GetWalletAddress returns the wallet address from the context.
// This is set by WalletAuthenticate middleware.
func GetWalletAddress(ctx context.Context) string {
	v, ok := ctx.Value(walletAddressKey).(string)
	if !ok {
		return ""
	}
	return v
}

func setKeyIssuedAt(ctx context.Context, issuedAt time.Time) context.Context {
	return context.WithValue(ctx, keyIssuedAtKey, issuedAt)
}

// GetKeyIssuedAt returns the API key issuance timestamp from the context.
func GetKeyIssuedAt(ctx context.Context) time.Time {
	v, ok := ctx.Value(keyIssuedAtKey).(time.Time)
	if !ok {
		return time.Time{}
	}
	return v
}

func setKeyNonce(ctx context.Context, nonce string) context.Context {
	return context.WithValue(ctx, keyNonceKey, nonce)
}

// GetKeyNonce returns the API key nonce from the context.
// This can be used for individual key revocation.
func GetKeyNonce(ctx context.Context) string {
	v, ok := ctx.Value(keyNonceKey).(string)
	if !ok {
		return ""
	}
	return v
}
