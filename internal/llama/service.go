package llama

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Service provides LLM inference using llama.cpp
type Service struct {
	manager          *LlamaCppReleaseManager
	embeddingManager *EmbeddingManager // NEW: For embedding models
	serverProcess    *exec.Cmd
	serverPort       int
	serverModelPath  string
	serverMutex      sync.Mutex
	initOnce         sync.Once // Ensure initialization happens only once
	autoStarted      bool      // Track if auto-start has been attempted
}

// NewService creates a new llama.cpp service instance
// Automatically installs llama.cpp and downloads recommended model in background
func NewService() (*Service, error) {
	manager := NewLlamaCppReleaseManager()
	embeddingManager := NewEmbeddingManager() // NEW: Initialize embedding manager

	service := &Service{
		manager:          manager,
		embeddingManager: embeddingManager,
		serverPort:       8080, // Default port
	}

	log.Printf("📍 [NewService] Created service instance: %p", service)

	// Start background initialization
	log.Printf("📍 [NewService] Starting goroutine for: %p", service)
	go service.initializeInBackground()

	return service, nil
}

// initializeInBackground handles llama.cpp installation and setup
func (s *Service) initializeInBackground() {
	log.Printf("🚀 Initializing llama.cpp in background... (Service instance: %p)", s)

	// Step 1: Check and install llama.cpp if needed
	if !s.manager.IsLlamaCppInstalled() {
		log.Println("🔧 llama.cpp not found, attempting auto-installation...")

		// Try package manager first (Homebrew on macOS)
		if err := s.manager.InstallLlamaCpp(); err != nil {
			log.Printf("⚠️  Package manager installation failed: %v", err)
			log.Println("   Falling back to GitHub release download...")

			// Fallback to GitHub release download
			// Clean up any partial downloads first
			if err := s.manager.CleanupPartialDownloads(); err != nil {
				log.Printf("⚠️  Cleanup partial downloads: %v", err)
			}

			// Get the latest release
			release, err := s.manager.GetLatestRelease()
			if err != nil {
				log.Printf("⚠️  Failed to get latest llama.cpp release: %v", err)
				log.Printf("   llama.cpp features will not be available")
				return
			}

			log.Printf("📥 Downloading llama.cpp %s from GitHub (this may take a few minutes)...", release.Version)
			if err := s.manager.DownloadRelease(release.Version, nil); err != nil {
				log.Printf("⚠️  Failed to download llama.cpp: %v", err)
				log.Printf("   You can download manually later from the UI")
				return
			}

			log.Printf("✅ llama.cpp %s installed successfully from GitHub", release.Version)
		} else {
			log.Println("✅ llama.cpp installed successfully via package manager")
		}
	} else {
		version := s.manager.GetInstalledVersion()
		log.Printf("✅ llama.cpp is installed (version: %s)", version)
	}

	// Step 2: Check llama-server binary
	serverPath := s.manager.GetServerBinaryPath()
	if _, err := os.Stat(serverPath); err != nil {
		log.Printf("⚠️  llama-server binary not found at: %s", serverPath)
		log.Println("   llama.cpp features will be limited")
		return
	}

	log.Printf("✅ llama-server ready at: %s", serverPath)
	log.Println("🎉 llama.cpp is ready to use!")

	// Step 3: Auto-download embedding model (in background)
	go func() {
		log.Println("📦 Checking embedding models...")
		if err := s.embeddingManager.AutoDownloadRecommendedModel(); err != nil {
			log.Printf("⚠️  Failed to auto-download embedding model: %v", err)
			log.Println("   Embedding features will require manual model download")
		} else {
			log.Println("✅ Embedding model ready!")
		}
	}()

	// Step 4: Auto-start llama-server if not already running
	// Use sync.Once to ensure this only happens once, even if called multiple times
	s.initOnce.Do(func() {
		// Wait a bit for any previous setup to complete
		time.Sleep(2 * time.Second)

		if !s.IsServerRunning() {
			log.Println("🚀 Attempting to auto-start llama-server...")

			// Check if models are available before starting
			models, err := s.GetAvailableModels()
			if err != nil {
				log.Printf("⚠️  Failed to check available models: %v", err)
				log.Println("   Will attempt to auto-download a model...")
				models = []string{} // Treat as no models
			}

			if len(models) == 0 {
				log.Println("⚠️  No GGUF models found. Starting auto-download...")

				// Auto-download recommended model based on hardware
				if err := s.AutoDownloadRecommendedModel(); err != nil {
					log.Printf("⚠️  Failed to auto-download model: %v", err)
					log.Println("   llama-server will not start without a model")
					log.Println("   Please download a model manually to use chat features")
					return
				}

				log.Println("✅ Model downloaded successfully!")
				// Re-check models after download (this will validate them)
				models, err = s.GetAvailableModels()
				if err != nil {
					log.Printf("⚠️  Failed to re-check models after download: %v", err)
				} else if len(models) == 0 {
					log.Println("⚠️  No valid models found after download. Model may be corrupt.")
					log.Println("   Please try downloading again or download manually")
					return
				}
			}

			if len(models) > 0 {
				log.Printf("✅ Found %d model(s), starting llama-server...", len(models))
				if err := s.StartServerAuto(); err != nil {
					log.Printf("❌ Failed to auto-start llama-server: %v", err)
					log.Println("   The server will be started automatically when you use chat features")
				} else {
					// Wait a moment and verify server is running
					time.Sleep(2 * time.Second)
					if s.IsServerRunning() {
						log.Println("✅ llama-server auto-started successfully and is responding")
					} else {
						log.Printf("⚠️  llama-server started but not responding yet")
						log.Println("   It may still be initializing. Try again in a moment.")
					}
				}
			} else {
				log.Println("⚠️  No models available. llama-server will not start.")
				log.Println("   Please download a model to use chat features")
			}
		} else {
			log.Println("✅ llama-server is already running")
		}

		s.autoStarted = true
	})
}

// GetBinaryPath returns the path to the llama.cpp binary directory
func (s *Service) GetBinaryPath() string {
	return s.manager.BinaryPath
}

// IsLlamaCppInstalled checks if llama.cpp is installed
func (s *Service) IsLlamaCppInstalled() bool {
	return s.manager.IsLlamaCppInstalled()
}

// GetInstalledVersion returns the currently installed version
func (s *Service) GetInstalledVersion() string {
	return s.manager.GetInstalledVersion()
}

// StartServer starts the llama-server with automatic model selection
// If modelPath is empty, automatically selects the best available model
func (s *Service) StartServer(modelPath string, port int) error {
	s.serverMutex.Lock()
	defer s.serverMutex.Unlock()

	log.Printf("[StartServer] Called with modelPath=%s, port=%d (current process: %v)", modelPath, port, s.serverProcess != nil)

	// Check if we already have a running process
	if s.serverProcess != nil && s.serverProcess.Process != nil {
		log.Printf("⚠️  llama-server process already exists (PID: %d)", s.serverProcess.Process.Pid)
		return nil
	}

	// Verify llama.cpp is installed
	if !s.manager.IsLlamaCppInstalled() {
		return fmt.Errorf("llama.cpp is not installed")
	}

	// Auto-select model if not provided
	if modelPath == "" {
		autoModel, err := s.selectBestModel()
		if err != nil {
			return fmt.Errorf("failed to auto-select model: %w", err)
		}
		modelPath = autoModel
		log.Printf("🤖 Auto-selected model: %s", modelPath)
	}

	// Verify model file exists
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return fmt.Errorf("model file not found: %s", modelPath)
	}

	serverPath := s.manager.GetServerBinaryPath()
	if _, err := os.Stat(serverPath); err != nil {
		return fmt.Errorf("llama-server binary not found: %w", err)
	}

	// Use provided port or default
	if port > 0 {
		s.serverPort = port
	}

	// Check if port is already in use
	if s.isPortInUse(s.serverPort) {
		log.Printf("⚠️  Port %d is already in use, attempting to find alternative port...", s.serverPort)

		// Try to find an available port
		newPort, err := s.findAvailablePort(s.serverPort)
		if err != nil {
			return fmt.Errorf("port %d is in use and no alternative port found: %w", s.serverPort, err)
		}

		log.Printf("✅ Found available port: %d", newPort)
		s.serverPort = newPort
	}

	// Build command arguments
	args := []string{
		"--host", "127.0.0.1",
		"--port", fmt.Sprintf("%d", s.serverPort),
		"--model", modelPath,
		"--threads", fmt.Sprintf("%d", runtime.NumCPU()),
		"--ctx-size", "4096",
		"--batch-size", "512",
		"--parallel", "4", // Allow 4 concurrent requests
	}

	// Create command
	cmd := exec.Command(serverPath, args...)

	// Capture stdout and stderr for debugging
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the process
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start llama-server: %w", err)
	}

	s.serverProcess = cmd
	s.serverModelPath = modelPath
	log.Printf("✅ llama-server started on port %d (PID: %d)", s.serverPort, cmd.Process.Pid)
	log.Printf("   Model: %s", modelPath)

	// Monitor the process in a goroutine
	go func() {
		if err := cmd.Wait(); err != nil {
			log.Printf("⚠️  llama-server process ended with error: %v", err)
		} else {
			log.Printf("⚠️  llama-server process ended normally")
		}
		s.serverMutex.Lock()
		s.serverProcess = nil
		s.serverModelPath = ""
		s.serverMutex.Unlock()
	}()

	// Wait a moment for server to start
	time.Sleep(2 * time.Second)

	// Verify server is responding
	if !s.IsServerRunning() {
		return fmt.Errorf("llama-server failed to start (not responding on port %d)", s.serverPort)
	}

	return nil
}

// StopServer stops the llama-server if it's running
func (s *Service) StopServer() error {
	s.serverMutex.Lock()
	defer s.serverMutex.Unlock()

	if s.serverProcess == nil || s.serverProcess.Process == nil {
		return fmt.Errorf("llama-server is not running")
	}

	log.Println("Stopping llama-server...")
	if err := s.serverProcess.Process.Kill(); err != nil {
		return fmt.Errorf("failed to stop llama-server: %w", err)
	}

	log.Println("✅ llama-server stopped")
	s.serverProcess = nil
	s.serverModelPath = ""

	return nil
}

// isPortInUse checks if a port is already in use
func (s *Service) isPortInUse(port int) bool {
	address := fmt.Sprintf("127.0.0.1:%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return true // Port is in use
	}
	listener.Close()
	return false
}

// findAvailablePort tries to find an available port starting from the given port
func (s *Service) findAvailablePort(startPort int) (int, error) {
	// Try ports in range: startPort to startPort+100
	for port := startPort + 1; port <= startPort+100; port++ {
		if !s.isPortInUse(port) {
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available port found in range %d-%d", startPort+1, startPort+100)
}

// IsServerRunning checks if llama-server is running by making a health check request
func (s *Service) IsServerRunning() bool {
	// Try to connect to the llama-server health endpoint
	client := &http.Client{Timeout: 2 * time.Second}
	url := fmt.Sprintf("http://127.0.0.1:%d/health", s.serverPort)
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// If we get any response, llama-server is running
	return resp.StatusCode == 200
}

// GetServerStatus returns the current status of llama-server
func (s *Service) GetServerStatus() map[string]interface{} {
	s.serverMutex.Lock()
	defer s.serverMutex.Unlock()

	status := map[string]interface{}{
		"running":    s.IsServerRunning(),
		"port":       s.serverPort,
		"model_path": s.serverModelPath,
		"pid":        nil,
	}

	if s.serverProcess != nil && s.serverProcess.Process != nil {
		status["pid"] = s.serverProcess.Process.Pid
	}

	return status
}

// GetServerURL returns the URL of the running llama-server
func (s *Service) GetServerURL() string {
	return fmt.Sprintf("http://127.0.0.1:%d", s.serverPort)
}

// StartServerAuto starts the llama-server with automatic model selection
// This is the recommended way to start the server for Kawai AI
func (s *Service) StartServerAuto() error {
	return s.StartServer("", 8080) // Empty modelPath triggers auto-selection
}

// selectBestModel automatically selects the best available model
func (s *Service) selectBestModel() (string, error) {
	modelsDir := s.manager.GetModelsDirectory()

	// Ensure models directory exists
	if err := os.MkdirAll(modelsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create models directory: %w", err)
	}

	// Get all GGUF files in models directory
	entries, err := os.ReadDir(modelsDir)
	if err != nil {
		return "", fmt.Errorf("failed to read models directory: %w", err)
	}

	var models []struct {
		path string
		size int64
		name string
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		// Only include GGUF files
		if !strings.HasSuffix(strings.ToLower(name), ".gguf") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		models = append(models, struct {
			path string
			size int64
			name string
		}{
			path: filepath.Join(modelsDir, name),
			size: info.Size(),
			name: name,
		})
	}

	if len(models) == 0 {
		return "", fmt.Errorf("no GGUF models found in %s. Please download a model first", modelsDir)
	}

	// Selection strategy:
	// 1. Prefer models with "qwen" in the name (Kawai's recommended model)
	// 2. Prefer Q4 quantization (good balance of speed/quality)
	// 3. Prefer smaller models for better performance
	// 4. Fall back to any available model

	var bestModel string
	var bestScore int

	for _, model := range models {
		score := 0
		nameLower := strings.ToLower(model.name)

		// Prefer Qwen models
		if strings.Contains(nameLower, "qwen") {
			score += 100
		}

		// Prefer Q4 quantization
		if strings.Contains(nameLower, "q4") {
			score += 50
		} else if strings.Contains(nameLower, "q5") {
			score += 30
		} else if strings.Contains(nameLower, "q8") {
			score += 10
		}

		// Prefer smaller models (better performance)
		// Penalize very large models
		sizeMB := model.size / (1024 * 1024)
		if sizeMB < 5000 { // < 5GB
			score += 20
		} else if sizeMB < 10000 { // < 10GB
			score += 10
		}

		// Prefer models with "instruct" or "chat" in name
		if strings.Contains(nameLower, "instruct") || strings.Contains(nameLower, "chat") {
			score += 30
		}

		if score > bestScore {
			bestScore = score
			bestModel = model.path
		}
	}

	// If no model scored well, just use the first one
	if bestModel == "" {
		bestModel = models[0].path
	}

	return bestModel, nil
}

// validateGGUFFile validates a GGUF file by checking its header structure
func (s *Service) validateGGUFFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// GGUF files start with a magic number: "GGUF" (4 bytes)
	magic := make([]byte, 4)
	if _, err := io.ReadFull(file, magic); err != nil {
		return fmt.Errorf("failed to read magic bytes: %w", err)
	}

	if string(magic) != "GGUF" {
		return fmt.Errorf("invalid GGUF file: wrong magic bytes (expected 'GGUF', got '%s')", string(magic))
	}

	// Read version (4 bytes, little-endian uint32)
	var version uint32
	if err := binary.Read(file, binary.LittleEndian, &version); err != nil {
		return fmt.Errorf("failed to read version: %w", err)
	}

	// Check if version is reasonable (GGUF versions are typically 1, 2, or 3)
	if version == 0 || version > 10 {
		return fmt.Errorf("invalid GGUF file: unreasonable version %d", version)
	}

	// Read tensor count (8 bytes, little-endian uint64)
	var tensorCount uint64
	if err := binary.Read(file, binary.LittleEndian, &tensorCount); err != nil {
		return fmt.Errorf("failed to read tensor count: %w", err)
	}

	// Read metadata key-value count (8 bytes, little-endian uint64)
	var metadataCount uint64
	if err := binary.Read(file, binary.LittleEndian, &metadataCount); err != nil {
		return fmt.Errorf("failed to read metadata count: %w", err)
	}

	// Basic sanity checks
	if tensorCount == 0 {
		return fmt.Errorf("invalid GGUF file: no tensors found")
	}

	if tensorCount > 100000 { // Reasonable upper bound
		return fmt.Errorf("invalid GGUF file: unreasonable tensor count %d", tensorCount)
	}

	if metadataCount > 10000 { // Reasonable upper bound for metadata entries
		return fmt.Errorf("invalid GGUF file: unreasonable metadata count %d", metadataCount)
	}

	log.Printf("GGUF file validation passed: version=%d, tensors=%d, metadata_entries=%d",
		version, tensorCount, metadataCount)
	return nil
}

// validateModelFile validates a model file and removes it if corrupt
func (s *Service) validateModelFile(modelPath string) error {
	// Check if file exists
	fileInfo, err := os.Stat(modelPath)
	if err != nil {
		return fmt.Errorf("model file not found: %w", err)
	}

	// Check minimum file size (GGUF files should be at least a few MB)
	if fileInfo.Size() < 1024*1024 { // Less than 1MB is suspicious
		log.Printf("⚠️  Model file %s is suspiciously small (%d bytes), removing...", modelPath, fileInfo.Size())
		os.Remove(modelPath)
		return fmt.Errorf("model file too small, likely incomplete")
	}

	// Validate GGUF structure
	if err := s.validateGGUFFile(modelPath); err != nil {
		log.Printf("⚠️  Model file %s failed validation: %v, removing...", modelPath, err)
		os.Remove(modelPath)
		return fmt.Errorf("model validation failed: %w", err)
	}

	return nil
}

// GetAvailableModels returns a list of available GGUF models
// Only returns models that pass validation
func (s *Service) GetAvailableModels() ([]string, error) {
	modelsDir := s.manager.GetModelsDirectory()

	entries, err := os.ReadDir(modelsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read models directory: %w", err)
	}

	var models []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasSuffix(strings.ToLower(name), ".gguf") {
			modelPath := filepath.Join(modelsDir, name)
			// Validate each model before including it
			if err := s.validateModelFile(modelPath); err != nil {
				log.Printf("⚠️  Skipping invalid model %s: %v", name, err)
				continue
			}
			models = append(models, modelPath)
		}
	}

	return models, nil
}

// GetModelsDirectory returns the directory where models are stored
func (s *Service) GetModelsDirectory() string {
	return s.manager.GetModelsDirectory()
}

// GetEmbeddingManager returns the embedding manager
func (s *Service) GetEmbeddingManager() *EmbeddingManager {
	return s.embeddingManager
}

// GetEmbeddingModelsDirectory returns the directory where embedding models are stored
func (s *Service) GetEmbeddingModelsDirectory() string {
	return s.embeddingManager.ModelsDir
}

// GetDownloadedEmbeddingModels returns a list of downloaded embedding models
func (s *Service) GetDownloadedEmbeddingModels() []*EmbeddingModel {
	return s.embeddingManager.GetDownloadedModels()
}

// DownloadEmbeddingModel downloads an embedding model with progress callback
func (s *Service) DownloadEmbeddingModel(modelName string, progressCallback func(downloaded, total int64)) error {
	return s.embeddingManager.DownloadModel(modelName, progressCallback)
}

// GetRecommendedEmbeddingModel returns the recommended embedding model
func (s *Service) GetRecommendedEmbeddingModel() string {
	return s.embeddingManager.GetRecommendedModel()
}

// StartEmbeddingServer starts llama-server with an embedding model
func (s *Service) StartEmbeddingServer(port int) error {
	s.serverMutex.Lock()
	defer s.serverMutex.Unlock()

	// Check if server is already running
	if s.IsServerRunning() {
		log.Printf("⚠️  llama-server is already running on port %d", s.serverPort)
		return fmt.Errorf("server already running")
	}

	// Get embedding model path
	embMgr := s.embeddingManager
	downloaded := embMgr.GetDownloadedModels()

	if len(downloaded) == 0 {
		return fmt.Errorf("no embedding models downloaded")
	}

	modelPath, err := embMgr.GetModelPath(downloaded[0].Name)
	if err != nil {
		return fmt.Errorf("failed to get model path: %w", err)
	}

	// Get llama-server binary path
	llamaServer := s.manager.GetServerBinaryPath()
	if _, err := os.Stat(llamaServer); err != nil {
		return fmt.Errorf("llama-server not found: %w", err)
	}

	log.Printf("🚀 Starting llama-server for embeddings...")
	log.Printf("   Model: %s", downloaded[0].Name)
	log.Printf("   Port: %d", port)

	// Build command
	cmd := exec.Command(llamaServer,
		"-m", modelPath,
		"--port", fmt.Sprintf("%d", port),
		"--embedding",
		"--pooling", "mean",
		"--embd-normalize", "2",
		"--ctx-size", "2048",
		"--batch-size", "512",
		"--threads", "4",
	)

	// Set output to logs
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the server
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start embedding server: %w", err)
	}

	s.serverProcess = cmd
	s.serverPort = port
	s.serverModelPath = modelPath

	// Wait a bit for server to start
	time.Sleep(2 * time.Second)

	// Verify server is running
	if !s.IsServerRunning() {
		return fmt.Errorf("embedding server failed to start")
	}

	log.Printf("✅ Embedding server started successfully on port %d", port)
	return nil
}
