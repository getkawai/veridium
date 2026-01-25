package version

import "github.com/kawai-network/veridium/pkg/mcp-gateway/version"

// Re-export from pkg for backward compatibility
var Version = version.Version

func UserAgent() string {
	return version.UserAgent()
}
