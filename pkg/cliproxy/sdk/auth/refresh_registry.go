package auth

import (
	"time"

	cliproxyauth "github.com/kawai-network/veridium/pkg/cliproxy/sdk/cliproxy/auth"
)

func registerRefreshLead(provider string, factory func() Authenticator) {
	cliproxyauth.RegisterRefreshLeadProvider(provider, func() *time.Duration {
		if factory == nil {
			return nil
		}
		auth := factory()
		if auth == nil {
			return nil
		}
		return auth.RefreshLead()
	})
}
