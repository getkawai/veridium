package llama

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// ProxyService provides HTTP proxy via Wails bindings
// This allows WebView to call localhost llama-server through Go
type ProxyService struct {
	service    *Service
	httpClient *http.Client
	app        *application.App
}

// NewProxyService creates a new proxy service
func NewProxyService(llamaService *Service, app *application.App) *ProxyService {
	log.Printf("📍 [proxy.go] Creating ProxyService with llamaService: %p", llamaService)
	return &ProxyService{
		service: llamaService,
		app:     app,
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

// replaceKawaiAutoModel replaces "kawai-auto" model name with actual model path
func (p *ProxyService) replaceKawaiAutoModel(body string) (string, error) {
	if body == "" {
		return body, nil
	}

	// Parse JSON body
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		// If not JSON, return as-is
		log.Printf("[Proxy] Failed to parse JSON body: %v", err)
		return body, nil
	}

	// Check if model field exists and is "kawai-auto"
	if model, ok := data["model"].(string); ok {
		log.Printf("[Proxy] Request model: %s", model)
		if model == "kawai-auto" {
			// Get the actual model path from the running server
			log.Printf("[Proxy] serverModelPath: %s", p.service.serverModelPath)
			if p.service.serverModelPath != "" {
				log.Printf("[Proxy] Replacing 'kawai-auto' with actual model: %s", p.service.serverModelPath)
				data["model"] = p.service.serverModelPath

				// Marshal back to JSON
				modifiedBody, err := json.Marshal(data)
				if err != nil {
					return body, fmt.Errorf("failed to marshal modified body: %w", err)
				}
				return string(modifiedBody), nil
			} else {
				log.Printf("[Proxy] ⚠️  serverModelPath is empty, cannot replace kawai-auto")
			}
		}
	} else {
		log.Printf("[Proxy] No model field found in request body")
	}

	return body, nil
}

// Fetch proxies HTTP requests to llama-server
// This mimics browser fetch() API but goes through Go
func (p *ProxyService) Fetch(ctx context.Context, request ProxyRequest) (*ProxyResponse, error) {
	log.Printf("📍 [Fetch] Called with method=%s, path=%s, body_len=%d", request.Method, request.Path, len(request.Body))

	// Modify request body: replace kawai-auto model and force stream:false
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(request.Body), &data); err == nil {
		modified := false
		
		// Replace kawai-auto model name if present
		if model, ok := data["model"].(string); ok {
			log.Printf("[Proxy] Request model: %s", model)
			if model == "kawai-auto" {
				log.Printf("[Proxy] serverModelPath: %s", p.service.serverModelPath)
				if p.service.serverModelPath != "" {
					log.Printf("[Proxy] Replacing 'kawai-auto' with actual model: %s", p.service.serverModelPath)
					data["model"] = p.service.serverModelPath
					modified = true
				} else {
					log.Printf("[Proxy] ⚠️  serverModelPath is empty, cannot replace kawai-auto")
				}
			}
		}
		
		// WORKAROUND: Fetch() cannot handle streaming responses
		// Force stream: false if present in request body
		if stream, ok := data["stream"].(bool); ok && stream {
			log.Printf("[Fetch] ⚠️  Detected stream:true in request, forcing stream:false for Fetch()")
			data["stream"] = false
			modified = true
			
			// IMPORTANT: Remove stream_options when stream:false
			// llama-server hangs if stream_options is present with stream:false
			if _, hasStreamOptions := data["stream_options"]; hasStreamOptions {
				log.Printf("[Fetch] Removing stream_options (incompatible with stream:false)")
				delete(data, "stream_options")
			}
		}
		
		// Add default max_tokens if not present to prevent very long responses
		if _, hasMaxTokens := data["max_tokens"]; !hasMaxTokens {
			log.Printf("[Fetch] Adding default max_tokens=2048 to prevent timeout")
			data["max_tokens"] = 2048
			modified = true
		}
		
		// Marshal back to JSON if modified
		if modified {
			if modifiedJSON, err := json.Marshal(data); err == nil {
				request.Body = string(modifiedJSON)
				log.Printf("[Fetch] Request body modified successfully")
			} else {
				log.Printf("[Fetch] ⚠️  Failed to marshal modified body: %v", err)
			}
		}
	} else {
		log.Printf("[Fetch] ⚠️  Failed to parse request body as JSON: %v", err)
	}

	// Ensure llama-server is running, attempt auto-start if not
	if !p.service.IsServerRunning() {
		log.Printf("[Proxy] llama-server is not running, attempting auto-start...")
		if err := p.service.StartServerAuto(); err != nil {
			return nil, fmt.Errorf("llama-server is not running and failed to auto-start: %w", err)
		}
		// Wait a moment for server to be ready
		time.Sleep(2 * time.Second)
		// Verify server is now running
		if !p.service.IsServerRunning() {
			return nil, fmt.Errorf("llama-server auto-start completed but server is not responding")
		}
		log.Printf("[Proxy] llama-server auto-started successfully")
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
	log.Printf("[Fetch] Executing request to: %s", url)
	log.Printf("[Fetch] Request method: %s, body length: %d", request.Method, len(request.Body))
	if len(request.Body) > 0 && len(request.Body) < 1000 {
		log.Printf("[Fetch] Final request body: %s", request.Body)
	} else if len(request.Body) > 0 {
		log.Printf("[Fetch] Final request body (first 500 chars): %s...", request.Body[:500])
	}
	
	// Add explicit timeout for the request
	reqCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	
	log.Printf("[Fetch] About to call httpClient.Do()...")
	resp, err := p.httpClient.Do(req)
	if err != nil {
		log.Printf("[Fetch] ❌ Request failed: %v", err)
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("[Fetch] ✅ Response received! Status: %d %s", resp.StatusCode, resp.Status)

	// Read response body with timeout
	log.Printf("[Fetch] About to read response body...")

	// Create a channel to read body with timeout
	type result struct {
		data []byte
		err  error
	}
	resultChan := make(chan result, 1)

	go func() {
		data, err := io.ReadAll(resp.Body)
		resultChan <- result{data, err}
	}()

	// Wait for result or timeout
	var bodyBytes []byte
	select {
	case res := <-resultChan:
		if res.err != nil {
			log.Printf("[Fetch] ❌ Failed to read response body: %v", res.err)
			return nil, fmt.Errorf("failed to read response body: %w", res.err)
		}
		bodyBytes = res.data
	case <-time.After(10 * time.Second):
		log.Printf("[Fetch] ❌ Timeout reading response body after 10 seconds")
		return nil, fmt.Errorf("timeout reading response body")
	}

	log.Printf("[Fetch] ✅ Response body read successfully, length: %d bytes", len(bodyBytes))
	if len(bodyBytes) > 0 && len(bodyBytes) < 500 {
		log.Printf("[Fetch] Response body preview: %s", string(bodyBytes))
	} else if len(bodyBytes) > 0 {
		log.Printf("[Fetch] Response body preview (first 200 chars): %s...", string(bodyBytes[:min(200, len(bodyBytes))]))
	}

	// Build response headers
	responseHeaders := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			responseHeaders[key] = values[0]
		}
	}

	log.Printf("[Fetch] ✅ Returning response to frontend")
	return &ProxyResponse{
		Status:     resp.StatusCode,
		StatusText: resp.Status,
		Headers:    responseHeaders,
		Body:       string(bodyBytes),
	}, nil
}

// StreamFetch proxies HTTP requests and streams response via Wails events
// This enables real-time streaming for SSE responses
func (p *ProxyService) StreamFetch(ctx context.Context, requestID string, request ProxyRequest) error {
	log.Printf("[StreamFetch] Starting stream for request ID: %s", requestID)

	// Replace kawai-auto model name if present
	modifiedBody, err := p.replaceKawaiAutoModel(request.Body)
	if err != nil {
		log.Printf("⚠️  Failed to replace kawai-auto model: %v", err)
	} else {
		request.Body = modifiedBody
	}

	// Ensure llama-server is running, attempt auto-start if not
	if !p.service.IsServerRunning() {
		log.Printf("[StreamFetch] llama-server is not running, attempting auto-start...")
		if err := p.service.StartServerAuto(); err != nil {
			return fmt.Errorf("llama-server is not running and failed to auto-start: %w", err)
		}
		// Wait a moment for server to be ready
		time.Sleep(2 * time.Second)
		// Verify server is now running
		if !p.service.IsServerRunning() {
			return fmt.Errorf("llama-server auto-start completed but server is not responding")
		}
		log.Printf("[StreamFetch] llama-server auto-started successfully")
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
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}

	// Execute request
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Build response headers
	responseHeaders := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			responseHeaders[key] = values[0]
		}
	}

	// Emit response metadata
	log.Printf("[StreamFetch] Emitting metadata for request ID: %s", requestID)
	p.app.Event.Emit(fmt.Sprintf("stream:%s:meta", requestID), map[string]interface{}{
		"status":     resp.StatusCode,
		"statusText": resp.Status,
		"headers":    responseHeaders,
	})

	// Stream response body line by line (for SSE format)
	reader := bufio.NewReader(resp.Body)
	lineCount := 0
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("[StreamFetch] Error reading response: %v", err)
			return fmt.Errorf("failed to read response: %w", err)
		}

		lineCount++
		lineStr := string(line)
		// Debug: log first few lines
		if lineCount <= 3 {
			log.Printf("[StreamFetch] Emitting line %d: %q", lineCount, lineStr)
		}
		// Emit each line as event
		p.app.Event.Emit(fmt.Sprintf("stream:%s:data", requestID), lineStr)
	}

	// Emit end event
	log.Printf("[StreamFetch] Stream ended for request ID: %s (sent %d lines)", requestID, lineCount)
	p.app.Event.Emit(fmt.Sprintf("stream:%s:end", requestID))

	return nil
}
