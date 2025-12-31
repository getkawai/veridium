package policy

import (
	"context"

	"github.com/kawai-network/veridium/pkg/mcp-gateway/desktop"
)

func Set(ctx context.Context, data string) error {
	return desktop.NewSecretsClient().SetJfsPolicy(ctx, data)
}
