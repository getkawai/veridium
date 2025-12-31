package interceptors

import (
	"context"
	"fmt"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TelemetryMiddleware tracks list operations and other gateway operations
func TelemetryMiddleware() mcp.Middleware {
	return func(next mcp.MethodHandler) mcp.MethodHandler {
		return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
			// Debug log all methods if debug is enabled
			if os.Getenv("DOCKER_MCP_TELEMETRY_DEBUG") != "" {
				fmt.Fprintf(os.Stderr, "[MCP-MIDDLEWARE] Method called: %s\n", method)
			}

			// Call the next handler
			result, err := next(ctx, method, req)

			return result, err
		}
	}
}
