package llamacpp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

// LlamaClient represents a client for interacting with llama-server or llama-cli
type LlamaClient struct {
	serverURL    string
	httpClient   *http.Client
	serverCmd    *exec.Cmd
	modelPath    string
	embeddingDim int
}

// EmbeddingRequest represents the request structure for embedding API
type EmbeddingRequest struct {
	Content string `json:"content"`
}

// EmbeddingResponse represents the response structure for embedding API
type EmbeddingResponse struct {
	Index     int         `json:"index"`
	Embedding [][]float32 `json:"embedding"`
}

// NewLlamaClient creates a new llama client that can use either llama-server or llama-cli
func NewLlamaClient(modelPath string, gpuLayers int) (*LlamaClient, error) {
	client := &LlamaClient{
		modelPath:  modelPath,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	// First try to start llama-server
	if err := client.startServer(gpuLayers); err != nil {
		return nil, fmt.Errorf("failed to start llama-server: %w", err)
	}

	// Wait for server to be ready and get embedding dimension
	if err := client.waitForServer(); err != nil {
		client.Close()
		return nil, fmt.Errorf("server failed to start properly: %w", err)
	}

	// Get embedding dimension by making a test request
	if err := client.getEmbeddingDimension(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to determine embedding dimension: %w", err)
	}

	return client, nil
}

// startServer starts the llama-server process
func (c *LlamaClient) startServer(gpuLayers int) error {
	// Find llama-server binary
	serverPath, err := c.findLlamaServer()
	if err != nil {
		return err
	}

	// Find an available port
	port, err := c.findAvailablePort()
	if err != nil {
		return err
	}

	c.serverURL = fmt.Sprintf("http://localhost:%d", port)

	// Build command arguments
	args := []string{
		"-m", c.modelPath,
		"--port", strconv.Itoa(port),
		"--embedding",
		"--log-disable",
		"-ngl", strconv.Itoa(gpuLayers),
	}

	// Start the server
	c.serverCmd = exec.Command(serverPath, args...)
	c.serverCmd.Stdout = nil // Suppress output
	c.serverCmd.Stderr = nil // Suppress output

	if err := c.serverCmd.Start(); err != nil {
		return fmt.Errorf("failed to start llama-server: %w", err)
	}

	return nil
}

// findLlamaServer finds the llama-server binary
func (c *LlamaClient) findLlamaServer() (string, error) {
	// Check if llama-server is in PATH
	if path, err := exec.LookPath("llama-server"); err == nil {
		return path, nil
	}

	// Check common installation paths
	commonPaths := []string{
		"/usr/local/bin/llama-server",
		"/usr/bin/llama-server",
		"./llama-server",
		"../llama.cpp/llama-server",
		"../llama.cpp/build/bin/llama-server",
	}

	for _, path := range commonPaths {
		if absPath, err := filepath.Abs(path); err == nil {
			if _, err := os.Stat(absPath); err == nil {
				return absPath, nil
			}
		}
	}

	return "", fmt.Errorf("llama-server binary not found in PATH or common locations")
}

// findAvailablePort finds an available port starting from 8080
func (c *LlamaClient) findAvailablePort() (int, error) {
	for port := 8080; port < 8180; port++ {
		if c.isPortAvailable(port) {
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available ports found")
}

// isPortAvailable checks if a port is available
func (c *LlamaClient) isPortAvailable(port int) bool {
	conn, err := http.Get(fmt.Sprintf("http://localhost:%d", port))
	if err != nil {
		return true // Port is available
	}
	conn.Body.Close()
	return false // Port is in use
}

// waitForServer waits for the server to be ready
func (c *LlamaClient) waitForServer() error {
	maxAttempts := 30
	for i := 0; i < maxAttempts; i++ {
		resp, err := c.httpClient.Get(c.serverURL + "/health")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == 200 {
				return nil
			}
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("server did not become ready within 30 seconds")
}

// getEmbeddingDimension determines the embedding dimension by making a test request
func (c *LlamaClient) getEmbeddingDimension() error {
	embedding, err := c.EmbedText("test")
	if err != nil {
		return err
	}
	c.embeddingDim = len(embedding)
	return nil
}

// EmbedText generates embeddings for the given text
func (c *LlamaClient) EmbedText(text string) ([]float32, error) {
	reqBody := EmbeddingRequest{Content: text}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(
		c.serverURL+"/embedding",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to make embedding request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	var embeddingRespArray []EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embeddingRespArray); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(embeddingRespArray) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}

	if len(embeddingRespArray[0].Embedding) == 0 {
		return nil, fmt.Errorf("empty embedding returned")
	}

	return embeddingRespArray[0].Embedding[0], nil
}

// EmbeddingDimension returns the dimension of embeddings produced by this model
func (c *LlamaClient) EmbeddingDimension() int {
	return c.embeddingDim
}

// Close stops the server and cleans up resources
func (c *LlamaClient) Close() error {
	if c.serverCmd != nil && c.serverCmd.Process != nil {
		// Try graceful shutdown first
		if err := c.serverCmd.Process.Signal(os.Interrupt); err != nil {
			// Force kill if graceful shutdown fails
			c.serverCmd.Process.Kill()
		}
		c.serverCmd.Wait()
		c.serverCmd = nil
	}
	return nil
}
