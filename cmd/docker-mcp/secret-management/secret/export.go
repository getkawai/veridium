package secret

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/kawai-network/veridium/pkg/mcp-gateway/catalog"
	"github.com/kawai-network/veridium/pkg/mcp-gateway/docker"
)

func Export(ctx context.Context, docker docker.Client, serverNames []string) (map[string]string, error) {
	catalog, err := catalog.Get(ctx)
	if err != nil {
		return nil, err
	}

	var secretNames []string
	for _, serverName := range serverNames {
		serverName = strings.TrimSpace(serverName)

		serverSpec, ok := catalog.Servers[serverName]
		if !ok {
			return nil, fmt.Errorf("server %s not found in catalog", serverName)
		}

		for _, s := range serverSpec.Secrets {
			secretNames = append(secretNames, s.Name)
		}
	}

	if len(secretNames) == 0 {
		return map[string]string{}, nil
	}

	sort.Strings(secretNames)

	return docker.ReadSecrets(ctx, secretNames, false)
}
