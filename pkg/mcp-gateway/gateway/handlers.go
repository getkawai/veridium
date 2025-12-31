package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/kawai-network/veridium/pkg/mcp-gateway/catalog"
)

func getClientConfig(readOnlyHint *bool, ss *mcp.ServerSession, server *mcp.Server) *clientConfig {
	return &clientConfig{readOnly: readOnlyHint, serverSession: ss, server: server}
}

// inferServerType determines the type of MCP server based on its configuration
func inferServerType(serverConfig *catalog.ServerConfig) string {
	if serverConfig.Spec.Remote.Transport == "http" {
		return "streaming"
	}

	if serverConfig.Spec.Remote.Transport == "sse" {
		return "sse"
	}

	// Check for Docker image
	if serverConfig.Spec.Image != "" {
		return "docker"
	}

	// Unknown type
	return "unknown"
}

func (g *Gateway) mcpToolHandler(tool catalog.Tool) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Convert CallToolParamsRaw to CallToolParams
		var args any
		if len(req.Params.Arguments) > 0 {
			if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
				return nil, fmt.Errorf("failed to unmarshal arguments: %w", err)
			}
		}
		params := &mcp.CallToolParams{
			Meta:      req.Params.Meta,
			Name:      req.Params.Name,
			Arguments: args,
		}
		return g.clientPool.runToolContainer(ctx, tool, params)
	}
}

func (g *Gateway) mcpServerToolHandler(serverName string, server *mcp.Server, annotations *mcp.ToolAnnotations, originalToolName string) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Look up server configuration
		serverConfig, _, ok := g.configuration.Find(serverName)
		if !ok {
			return nil, fmt.Errorf("server %q not found in configuration", serverName)
		}

		// Debug logging to stderr
		if os.Getenv("DOCKER_MCP_TELEMETRY_DEBUG") != "" {
			fmt.Fprintf(os.Stderr, "[MCP-HANDLER] Tool call received: %s from server: %s\n", req.Params.Name, serverConfig.Name)
		}

		var readOnlyHint *bool
		if annotations != nil && annotations.ReadOnlyHint {
			readOnlyHint = &annotations.ReadOnlyHint
		}

		client, err := g.clientPool.AcquireClient(ctx, serverConfig, getClientConfig(readOnlyHint, req.Session, server))
		if err != nil {
			return nil, err
		}
		defer g.clientPool.ReleaseClient(client)

		// Convert CallToolParamsRaw to CallToolParams
		var args any
		if len(req.Params.Arguments) > 0 {
			if jsonErr := json.Unmarshal(req.Params.Arguments, &args); jsonErr != nil {
				return nil, fmt.Errorf("failed to unmarshal arguments: %w", jsonErr)
			}
		}
		params := &mcp.CallToolParams{
			Meta:      req.Params.Meta,
			Name:      originalToolName,
			Arguments: args,
		}

		// Execute the tool call
		result, err := client.Session().CallTool(ctx, params)

		if err != nil {
			return nil, err
		}

		return result, nil
	}
}

func (g *Gateway) mcpServerPromptHandler(serverName string, server *mcp.Server) mcp.PromptHandler {
	return func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		// Look up server configuration
		serverConfig, _, ok := g.configuration.Find(serverName)
		if !ok {
			return nil, fmt.Errorf("server %q not found in configuration", serverName)
		}

		// Debug logging to stderr
		if os.Getenv("DOCKER_MCP_TELEMETRY_DEBUG") != "" {
			fmt.Fprintf(os.Stderr, "[MCP-HANDLER] Prompt get received: %s from server: %s\n", req.Params.Name, serverConfig.Name)
		}

		client, err := g.clientPool.AcquireClient(ctx, serverConfig, getClientConfig(nil, req.Session, server))
		if err != nil {
			return nil, err
		}
		defer g.clientPool.ReleaseClient(client)

		result, err := client.Session().GetPrompt(ctx, req.Params)
		if err != nil {
			return nil, err
		}

		return result, nil
	}
}

func (g *Gateway) mcpServerResourceHandler(serverName string, server *mcp.Server) mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		// Look up server configuration
		serverConfig, _, ok := g.configuration.Find(serverName)
		if !ok {
			return nil, fmt.Errorf("server %q not found in configuration", serverName)
		}

		// Debug logging to stderr
		if os.Getenv("DOCKER_MCP_TELEMETRY_DEBUG") != "" {
			fmt.Fprintf(os.Stderr, "[MCP-HANDLER] Resource read received: %s from server: %s\n", req.Params.URI, serverConfig.Name)
		}

		client, err := g.clientPool.AcquireClient(ctx, serverConfig, getClientConfig(nil, req.Session, server))
		if err != nil {
			return nil, err
		}
		defer g.clientPool.ReleaseClient(client)

		result, err := client.Session().ReadResource(ctx, req.Params)
		if err != nil {
			return nil, err
		}

		return result, nil
	}
}
