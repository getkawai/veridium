package policy

import (
	"context"
	"fmt"

	"github.com/kawai-network/veridium/pkg/mcp-gateway/desktop"
)

func Dump(ctx context.Context) error {
	l, err := desktop.NewSecretsClient().GetJfsPolicy(ctx)
	if err != nil {
		return err
	}

	fmt.Println(l)
	return nil
}
