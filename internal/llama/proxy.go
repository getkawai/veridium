package llama

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ProxyService provides HTTP proxy via Wails bindings
// This allows WebView to call localhost llama-server through Go
type ProxyService struct {
	service    *Service
	httpClient *http.Client
}

// NewProxyService creates a new proxy service
func NewProxyService(llamaService *Service) *ProxyService {
	return &ProxyService{
		service: llamaService,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// ProxyRequest represents a generic HTTP proxy request
type ProxyRequest struct {
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

// ProxyResponse represents a generic HTTP proxy response
type ProxyResponse struct {
	Status     int               `json:"status"`
	StatusText string            `json:"statusText"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

// Fetch proxies HTTP requests to llama-server
// This mimics browser fetch() API but goes through Go
func (p *ProxyService) Fetch(ctx context.Context, request ProxyRequest) (*ProxyResponse, error) {
	// Ensure llama-server is running
	if !p.service.IsServerRunning() {
		return nil, fmt.Errorf("llama-server is not running")
	}

	// Build full URL
	baseURL := p.service.GetServerURL()
	url := baseURL + request.Path

	// Create HTTP request
	var bodyReader io.Reader
	if request.Body != "" {
		bodyReader = bytes.NewBufferString(request.Body)
	}

	req, err := http.NewRequestWithContext(ctx, request.Method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}

	// Execute request
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Build response headers
	responseHeaders := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			responseHeaders[key] = values[0]
		}
	}

	return &ProxyResponse{
		Status:     resp.StatusCode,
		StatusText: resp.Status,
		Headers:    responseHeaders,
		Body:       string(bodyBytes),
	}, nil
}
