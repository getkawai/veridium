package llama

import (
	"fmt"
	"log"
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
	manager         *LlamaCppReleaseManager
	serverProcess   *exec.Cmd
	serverPort      int
	serverModelPath string
	serverMutex     sync.Mutex
}

// NewService creates a new llama.cpp service instance
// Automatically installs llama.cpp and downloads recommended model in background
func NewService() (*Service, error) {
	manager := NewLlamaCppReleaseManager()

	service := &Service{
		manager:    manager,
		serverPort: 8080, // Default port
	}

	// Start background initialization
	go service.initializeInBackground()

	return service, nil
}

// initializeInBackground handles llama.cpp installation and setup
func (s *Service) initializeInBackground() {
	log.Println("🚀 Initializing llama.cpp in background...")

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

	// Step 3: Auto-start llama-server if not already running
	// Wait a bit for any previous setup to complete
	time.Sleep(1 * time.Second)

	if !s.IsServerRunning() {
		// Check if models are available before starting
		models, err := s.GetAvailableModels()
		if err != nil {
			log.Printf("⚠️  Failed to check available models: %v", err)
			log.Println("   llama-server will not auto-start. Please download a model first.")
			return
		}

		if len(models) == 0 {
			log.Println("⚠️  No GGUF models found. Starting auto-download...")
			
			// Auto-download recommended model based on hardware
			if err := s.AutoDownloadRecommendedModel(); err != nil {
				log.Printf("⚠️  Failed to auto-download model: %v", err)
				log.Println("   You can download a model manually later")
			return
			}
			
			log.Println("✅ Model downloaded successfully!")
		}

		log.Println("🚀 Auto-starting llama-server...")
		if err := s.StartServerAuto(); err != nil {
			log.Printf("⚠️  Failed to auto-start llama-server: %v", err)
			log.Println("   You can start it manually later")
		} else {
			log.Println("✅ llama-server auto-started successfully")
		}
	} else {
		log.Println("✅ llama-server is already running")
	}
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

	// Check if server is already running
	if s.IsServerRunning() {
		log.Println("llama-server is already running")
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

	// Build command arguments
	args := []string{
		"--host", "127.0.0.1",
		"--port", fmt.Sprintf("%d", s.serverPort),
		"--model", modelPath,
		"--threads", fmt.Sprintf("%d", runtime.NumCPU()),
		"--ctx-size", "4096",
		"--batch-size", "512",
	}

	// Create command
	cmd := exec.Command(serverPath, args...)

	// Start the process
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start llama-server: %w", err)
	}

	s.serverProcess = cmd
	s.serverModelPath = modelPath
	log.Printf("✅ llama-server started on port %d (PID: %d)", s.serverPort, cmd.Process.Pid)

	// Monitor the process in a goroutine
	go func() {
		if err := cmd.Wait(); err != nil {
			log.Printf("⚠️  llama-server process ended: %v", err)
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

// GetAvailableModels returns a list of available GGUF models
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
			models = append(models, filepath.Join(modelsDir, name))
		}
	}

	return models, nil
}

// GetModelsDirectory returns the directory where models are stored
func (s *Service) GetModelsDirectory() string {
	return s.manager.GetModelsDirectory()
}
