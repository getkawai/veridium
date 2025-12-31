package tunnelkit

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/cloudflare/cloudflared/cmd/cloudflared/cliutil"
	"github.com/cloudflare/cloudflared/cmd/cloudflared/tunnel"
	"github.com/cloudflare/cloudflared/connection"
	"github.com/google/uuid"
	"github.com/kawai-network/veridium/pkg/cfapi"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

// Config holds the configuration for creating and running a tunnel
type Config struct {
	AccountID  string
	APIToken   string
	TunnelName string
	Hostname   string
	LocalURL   string
}

// TunnelInfo contains information about a created tunnel
type TunnelInfo struct {
	TunnelID    string
	TunnelToken string
	Hostname    string
	PublicURL   string
}

// GetOrCreateTunnelWithDNS gets an existing tunnel by name or creates a new one if it doesn't exist
func GetOrCreateTunnelWithDNS(cfg Config) (*TunnelInfo, error) {
	// Get zone ID from hostname
	zoneID, err := getZoneID(cfg.APIToken, cfg.Hostname)
	if err != nil {
		return nil, fmt.Errorf("failed to get zone ID: %w", err)
	}

	// Create REST client
	logger := slog.Default()
	client, err := cfapi.NewRESTClient(
		"https://api.cloudflare.com/client/v4",
		cfg.AccountID,
		zoneID,
		cfg.APIToken,
		"tunnelkit",
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	// Try to find existing tunnel by name
	filter := cfapi.NewTunnelFilter()
	filter.ByName(cfg.TunnelName)
	filter.NoDeleted()

	tunnels, err := client.ListTunnels(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list tunnels: %w", err)
	}

	var tunnelID uuid.UUID
	var tunnelToken string

	if len(tunnels) > 0 {
		// Use existing tunnel
		tunnel := tunnels[0]
		tunnelID = tunnel.ID

		// Check if tunnel has active connections
		activeClients, err := client.ListActiveClients(tunnelID)
		if err == nil && len(activeClients) > 0 {
			return nil, fmt.Errorf("tunnel '%s' is currently in use with %d active connection(s)", cfg.TunnelName, len(activeClients))
		}

		// Get tunnel token
		tunnelToken, err = client.GetTunnelToken(tunnelID)
		if err != nil {
			return nil, fmt.Errorf("failed to get tunnel token: %w", err)
		}
	} else {
		// Create new tunnel
		tunnelSecret, err := generateTunnelSecret()
		if err != nil {
			return nil, fmt.Errorf("failed to generate tunnel secret: %w", err)
		}

		tunnelResult, err := client.CreateTunnel(cfg.TunnelName, tunnelSecret)
		if err != nil {
			return nil, fmt.Errorf("failed to create tunnel: %w", err)
		}

		tunnelID = tunnelResult.Tunnel.ID
		tunnelToken = tunnelResult.Token
	}

	// Configure tunnel routing and DNS
	return configureTunnelRouting(client, cfg, tunnelID, tunnelToken)
}

// CreateTunnelWithDNS creates a new Cloudflare Tunnel and sets up DNS routing
func CreateTunnelWithDNS(cfg Config) (*TunnelInfo, error) {
	// Generate tunnel secret
	tunnelSecret, err := generateTunnelSecret()
	if err != nil {
		return nil, fmt.Errorf("failed to generate tunnel secret: %w", err)
	}

	// Get zone ID from hostname
	zoneID, err := getZoneID(cfg.APIToken, cfg.Hostname)
	if err != nil {
		return nil, fmt.Errorf("failed to get zone ID: %w", err)
	}

	// Create REST client
	logger := slog.Default()
	client, err := cfapi.NewRESTClient(
		"https://api.cloudflare.com/client/v4",
		cfg.AccountID,
		zoneID,
		cfg.APIToken,
		"tunnelkit",
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	// Create tunnel
	tunnelResult, err := client.CreateTunnel(cfg.TunnelName, tunnelSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to create tunnel: %w", err)
	}

	// Configure tunnel routing and DNS
	return configureTunnelRouting(client, cfg, tunnelResult.Tunnel.ID, tunnelResult.Token)
}

// GetTunnelByName retrieves an existing tunnel by name
func GetTunnelByName(accountID, apiToken, tunnelName string) (*TunnelInfo, error) {
	// Create REST client (zoneID not needed for listing tunnels)
	logger := slog.Default()
	client, err := cfapi.NewRESTClient(
		"https://api.cloudflare.com/client/v4",
		accountID,
		"", // zoneID not required for tunnel operations
		apiToken,
		"tunnelkit",
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	// Find tunnel by name
	filter := cfapi.NewTunnelFilter()
	filter.ByName(tunnelName)
	filter.NoDeleted()

	tunnels, err := client.ListTunnels(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list tunnels: %w", err)
	}

	if len(tunnels) == 0 {
		return nil, fmt.Errorf("tunnel '%s' not found", tunnelName)
	}

	tunnel := tunnels[0]

	// Get tunnel token
	tunnelToken, err := client.GetTunnelToken(tunnel.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tunnel token: %w", err)
	}

	// Check if tunnel has active connections
	activeClients, err := client.ListActiveClients(tunnel.ID)
	isActive := err == nil && len(activeClients) > 0

	info := &TunnelInfo{
		TunnelID:    tunnel.ID.String(),
		TunnelToken: tunnelToken,
		Hostname:    "", // Will be populated from routes if needed
		PublicURL:   "", // Will be populated from routes if needed
	}

	// Add active status info
	if isActive {
		return nil, fmt.Errorf("tunnel '%s' is currently active with %d connection(s)", tunnelName, len(activeClients))
	}

	return info, nil
}

// ListTunnels lists all tunnels in the account
func ListTunnels(accountID, apiToken string) ([]*TunnelInfo, error) {
	// Create REST client
	logger := slog.Default()
	client, err := cfapi.NewRESTClient(
		"https://api.cloudflare.com/client/v4",
		accountID,
		"",
		apiToken,
		"tunnelkit",
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	// List all non-deleted tunnels
	filter := cfapi.NewTunnelFilter()
	filter.NoDeleted()

	tunnels, err := client.ListTunnels(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list tunnels: %w", err)
	}

	var result []*TunnelInfo
	for _, tunnel := range tunnels {
		// Check if tunnel has active connections
		activeClients, _ := client.ListActiveClients(tunnel.ID)
		isActive := len(activeClients) > 0

		info := &TunnelInfo{
			TunnelID:  tunnel.ID.String(),
			Hostname:  "", // Not available from list API
			PublicURL: "",
		}

		// Add status indicator in hostname field for display
		if isActive {
			info.Hostname = fmt.Sprintf("%s (ACTIVE - %d connections)", tunnel.Name, len(activeClients))
		} else {
			info.Hostname = fmt.Sprintf("%s (inactive)", tunnel.Name)
		}

		result = append(result, info)
	}

	return result, nil
}

// HasActiveConnections checks if a tunnel has active connections
// Returns true if tunnel has active connections, false otherwise
func HasActiveConnections(tunnelID string) (bool, error) {
	// Parse tunnel ID
	id, err := uuid.Parse(tunnelID)
	if err != nil {
		return false, fmt.Errorf("invalid tunnel ID: %w", err)
	}

	// Get credentials from generated tunnels
	// We need accountID and apiToken to check active connections
	// For now, we'll use hardcoded credentials from bulk-tunnels
	// TODO: Make this configurable or read from environment
	accountID := "ceab218751d33cd804878196ad7bef74"
	apiToken := "OP8BZQhyeJxrovCPKt15eUOSC6i5LXTVECGRSMc1"

	// Create REST client
	logger := slog.Default()
	client, err := cfapi.NewRESTClient(
		"https://api.cloudflare.com/client/v4",
		accountID,
		"", // zoneID not required
		apiToken,
		"tunnelkit",
		logger,
	)
	if err != nil {
		return false, fmt.Errorf("failed to create API client: %w", err)
	}

	// Check active connections
	activeClients, err := client.ListActiveClients(id)
	if err != nil {
		return false, fmt.Errorf("failed to list active clients: %w", err)
	}

	return len(activeClients) > 0, nil
}

// RunTunnel runs a Cloudflare Tunnel using the provided token
func RunTunnel(ctx context.Context, tunnelToken string) error {
	// Decode tunnel token
	tokenBytes, err := base64.StdEncoding.DecodeString(tunnelToken)
	if err != nil {
		return fmt.Errorf("failed to decode tunnel token: %w", err)
	}

	var token connection.TunnelToken
	if err := json.Unmarshal(tokenBytes, &token); err != nil {
		return fmt.Errorf("failed to parse tunnel token: %w", err)
	}

	credentials := token.Credentials()
	tunnelProps := &connection.TunnelProperties{
		Credentials: credentials,
	}

	// Setup logger for cloudflared (requires zerolog)
	zerologLogger := zerolog.Nop()

	// Create build info
	buildInfo := &cliutil.BuildInfo{
		GoOS:               "darwin",
		GoVersion:          "go1.21",
		CloudflaredVersion: "dev",
	}

	// Create flag set with required flags
	set := flag.NewFlagSet("tunnel", flag.ContinueOnError)
	set.String("protocol", "http2", "Protocol to use")
	set.String("edge-ip-version", "auto", "Edge IP version")
	set.Duration("rpc-timeout", 5*time.Second, "RPC timeout")
	set.Int("retries", 5, "Number of retries")
	set.Bool("no-tls-verify", true, "Disable TLS verification")
	set.String("origin-ca-pool", "", "Origin CA pool")
	set.String("config", "", "Config file path")

	// Parse empty args
	if err := set.Parse([]string{}); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	// Create cli context
	cliCtx := cli.NewContext(nil, set, nil)

	// Start tunnel server
	return tunnel.StartServer(cliCtx, buildInfo, tunnelProps, &zerologLogger)
}

// RunTunnelWithShutdown runs a tunnel with graceful shutdown handling
func RunTunnelWithShutdown(tunnelToken string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 1)

	// Run tunnel in goroutine
	go func() {
		errChan <- RunTunnel(ctx, tunnelToken)
	}()

	// Wait for signal or error
	select {
	case sig := <-sigChan:
		fmt.Printf("\nReceived signal: %v, shutting down...\n", sig)
		cancel()
		time.Sleep(2 * time.Second)
		return nil
	case err := <-errChan:
		return err
	}
}

// configureTunnelRouting sets up DNS routing and updates the tunnel configuration
func configureTunnelRouting(client *cfapi.RESTClient, cfg Config, tunnelID uuid.UUID, tunnelToken string) (*TunnelInfo, error) {
	// Create or update DNS route
	route := cfapi.NewDNSRoute(cfg.Hostname, true)
	if _, err := client.RouteTunnel(tunnelID, route); err != nil {
		return nil, fmt.Errorf("failed to create DNS route: %w", err)
	}

	// Update tunnel configuration with ingress rules
	if err := updateTunnelConfig(cfg.APIToken, cfg.AccountID, tunnelID, cfg.Hostname, cfg.LocalURL); err != nil {
		return nil, fmt.Errorf("failed to update tunnel config: %w", err)
	}

	return &TunnelInfo{
		TunnelID:    tunnelID.String(),
		TunnelToken: tunnelToken,
		Hostname:    cfg.Hostname,
		PublicURL:   fmt.Sprintf("https://%s", cfg.Hostname),
	}, nil
}

// Helper functions

func generateTunnelSecret() ([]byte, error) {
	secret := make([]byte, 32)
	id := uuid.New()
	copy(secret, id[:])
	return secret, nil
}

func getZoneID(apiToken, hostname string) (string, error) {
	domain := extractDomain(hostname)

	// Use HTTP client directly for zone lookup
	client := &http.Client{Timeout: 15 * time.Second}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.cloudflare.com/client/v4/zones?name=%s", domain), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	type ZoneResponse struct {
		Result []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"result"`
		Success bool `json:"success"`
	}

	var zoneResp ZoneResponse
	if err := json.NewDecoder(resp.Body).Decode(&zoneResp); err != nil {
		return "", err
	}

	if !zoneResp.Success || len(zoneResp.Result) == 0 {
		return "", fmt.Errorf("zone not found for domain: %s", domain)
	}

	return zoneResp.Result[0].ID, nil
}

func extractDomain(hostname string) string {
	parts := strings.Split(hostname, ".")
	if len(parts) >= 2 {
		return strings.Join(parts[len(parts)-2:], ".")
	}
	return hostname
}

func updateTunnelConfig(apiToken, accountID string, tunnelID uuid.UUID, hostname, localURL string) error {
	type IngressRule struct {
		Hostname string `json:"hostname,omitempty"`
		Service  string `json:"service"`
	}

	type TunnelConfig struct {
		Ingress []IngressRule `json:"ingress"`
	}

	type ConfigPayload struct {
		Config TunnelConfig `json:"config"`
	}

	payload := ConfigPayload{
		Config: TunnelConfig{
			Ingress: []IngressRule{
				{
					Hostname: hostname,
					Service:  localURL,
				},
				{
					Service: "http_status:404",
				},
			},
		},
	}

	configBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Use HTTP client directly
	client := &http.Client{Timeout: 15 * time.Second}
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/cfd_tunnel/%s/configurations", accountID, tunnelID)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(configBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update config: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
