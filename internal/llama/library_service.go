package llama

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/kawai-network/veridium/pkg/yzma/llama"
)

// LibraryService provides LLM inference using llama.cpp as a library (via yzma)
// This replaces the binary-based approach with direct library calls
type LibraryService struct {
	manager *LlamaCppInstaller

	// Library state
	libPath       string
	isInitialized bool
	initMutex     sync.Mutex

	// Chat model state
	chatModel     llama.Model
	chatContext   llama.Context
	chatVocab     llama.Vocab
	chatSampler   llama.Sampler
	chatModelPath string
	chatMutex     sync.Mutex

	// Embedding model state
	embModel     llama.Model
	embContext   llama.Context
	embVocab     llama.Vocab
	embModelPath string
	embMutex     sync.Mutex

	initOnce sync.Once
}

// NewLibraryService creates a new library-based llama.cpp service
func NewLibraryService() (*LibraryService, error) {
	manager := NewLlamaCppInstaller()

	service := &LibraryService{
		manager: manager,
	}

	log.Printf("📍 [NewLibraryService] Created library-based service instance: %p", service)

	// Start background initialization
	go service.initializeInBackground()

	return service, nil
}

// initializeInBackground handles llama.cpp installation and library loading
func (s *LibraryService) initializeInBackground() {
	log.Printf("🚀 Initializing llama.cpp library in background...")

	// Step 1: Check and install llama.cpp if needed
	if !s.manager.IsLlamaCppInstalled() {
		log.Println("🔧 llama.cpp not found, attempting auto-installation...")

		// Try package manager first
		if err := s.manager.InstallLlamaCpp(); err != nil {
			log.Printf("⚠️  Package manager installation failed: %v", err)
			log.Println("   Falling back to GitHub release download...")

			// Clean up partial downloads
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

			log.Printf("📥 Downloading llama.cpp %s from GitHub...", release.Version)
			if err := s.manager.DownloadRelease(release.Version, nil); err != nil {
				log.Printf("⚠️  Failed to download llama.cpp: %v", err)
				return
			}

			log.Printf("✅ llama.cpp %s installed successfully", release.Version)
		} else {
			log.Println("✅ llama.cpp installed via package manager")
		}
	} else {
		version := s.manager.GetInstalledVersion()
		log.Printf("✅ llama.cpp is installed (version: %s)", version)
	}

	// Step 2: Initialize the library
	if err := s.InitializeLibrary(); err != nil {
		log.Printf("⚠️  Failed to initialize llama.cpp library: %v", err)
		return
	}

	log.Println("✅ llama.cpp library initialized successfully!")

	// Step 3: Auto-download embedding model (in background)
	go func() {
		log.Println("📦 Checking embedding models...")
		if err := s.manager.AutoDownloadRecommendedEmbeddingModel(); err != nil {
			log.Printf("⚠️  Failed to auto-download embedding model: %v", err)
		} else {
			log.Println("✅ Embedding model ready!")
		}
	}()

	// Step 4: Auto-load chat model if available
	s.initOnce.Do(func() {
		models, err := s.GetAvailableModels()
		if err != nil {
			log.Printf("⚠️  Failed to check available models: %v", err)
			return
		}

		if len(models) == 0 {
			log.Println("⚠️  No GGUF models found. Auto-downloading...")
			if err := s.AutoDownloadRecommendedModel(); err != nil {
				log.Printf("⚠️  Failed to auto-download model: %v", err)
				return
			}
			models, _ = s.GetAvailableModels()
		}

		if len(models) > 0 {
			log.Printf("✅ Found %d model(s), auto-loading first model...", len(models))
			if err := s.LoadChatModel(""); err != nil {
				log.Printf("⚠️  Failed to auto-load chat model: %v", err)
			} else {
				log.Println("✅ Chat model loaded and ready!")
			}
		}
	})
}

// InitializeLibrary loads the llama.cpp shared library
func (s *LibraryService) InitializeLibrary() error {
	s.initMutex.Lock()
	defer s.initMutex.Unlock()

	if s.isInitialized {
		return nil
	}

	// Get library path from installer (programmatic, not env var)
	libPath := s.manager.GetLibraryPath()

	// Verify all required libraries exist (libggml, libggml-base, libllama)
	if !s.manager.VerifyAllLibrariesExist() {
		return fmt.Errorf("llama.cpp libraries not found in %s. Please run installer first", libPath)
	}

	s.libPath = libPath
	log.Printf("📚 Loading llama.cpp library from directory: %s", libPath)

	// Log which libraries were found
	requiredPaths := s.manager.GetRequiredLibraryPaths()
	for _, path := range requiredPaths {
		log.Printf("  ✓ Found: %s", filepath.Base(path))
	}

	// Load the library (llama.Load expects a directory path)
	if err := llama.Load(libPath); err != nil {
		return fmt.Errorf("failed to load llama.cpp library from %s: %w", libPath, err)
	}

	// Initialize llama.cpp backend
	llama.Init()

	s.isInitialized = true
	log.Println("✅ llama.cpp library loaded and backend initialized successfully")

	return nil
}

// LoadChatModel loads a chat/generation model
// If modelPath is empty, automatically selects the best available model
func (s *LibraryService) LoadChatModel(modelPath string) error {
	s.chatMutex.Lock()
	defer s.chatMutex.Unlock()

	// Ensure library is initialized
	if !s.isInitialized {
		if err := s.InitializeLibrary(); err != nil {
			return fmt.Errorf("failed to initialize library: %w", err)
		}
	}

	// Auto-select model if not provided
	if modelPath == "" {
		autoModel, err := s.selectBestModel()
		if err != nil {
			return fmt.Errorf("failed to auto-select model: %w", err)
		}
		modelPath = autoModel
		log.Printf("🤖 Auto-selected chat model: %s", modelPath)
	}

	// Verify model file exists
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return fmt.Errorf("model file not found: %s", modelPath)
	}

	// Unload previous model if exists
	if s.chatModel != 0 {
		log.Println("♻️  Unloading previous chat model...")
		if s.chatContext != 0 {
			llama.Free(s.chatContext)
			s.chatContext = 0
		}
		if s.chatModel != 0 {
			llama.ModelFree(s.chatModel)
			s.chatModel = 0
		}
	}

	log.Printf("📥 Loading chat model: %s", filepath.Base(modelPath))

	// Load model with default parameters
	mParams := llama.ModelDefaultParams()
	s.chatModel = llama.ModelLoadFromFile(modelPath, mParams)
	if s.chatModel == 0 {
		return fmt.Errorf("failed to load model from file")
	}

	// Get vocabulary
	s.chatVocab = llama.ModelGetVocab(s.chatModel)

	// Create context with reasonable defaults
	// CRITICAL: Batch size must be >= max tokens processed in a single batch
	// For long prompts, n_tokens_all can exceed n_batch, causing GGML_ASSERT failure
	// Increasing batch size to handle longer prompts safely
	ctxParams := llama.ContextDefaultParams()
	ctxParams.NCtx = 4096   // Context size
	ctxParams.NBatch = 2048 // Batch size - increased from 512 to handle long prompts
	ctxParams.NThreads = int32(runtime.NumCPU())

	s.chatContext = llama.InitFromModel(s.chatModel, ctxParams)
	if s.chatContext == 0 {
		llama.ModelFree(s.chatModel)
		s.chatModel = 0
		return fmt.Errorf("failed to create context")
	}

	// Create sampler chain
	s.chatSampler = llama.SamplerChainInit(llama.SamplerChainDefaultParams())
	llama.SamplerChainAdd(s.chatSampler, llama.SamplerInitTopK(40))
	llama.SamplerChainAdd(s.chatSampler, llama.SamplerInitTopP(0.95, 1))
	llama.SamplerChainAdd(s.chatSampler, llama.SamplerInitTempExt(0.8, 0, 1.0))
	llama.SamplerChainAdd(s.chatSampler, llama.SamplerInitDist(llama.DefaultSeed))

	s.chatModelPath = modelPath
	log.Printf("✅ Chat model loaded successfully: %s", filepath.Base(modelPath))

	return nil
}

// LoadEmbeddingModel loads an embedding model
func (s *LibraryService) LoadEmbeddingModel(modelPath string) error {
	s.embMutex.Lock()
	defer s.embMutex.Unlock()

	// Ensure library is initialized
	if !s.isInitialized {
		if err := s.InitializeLibrary(); err != nil {
			return fmt.Errorf("failed to initialize library: %w", err)
		}
	}

	// Auto-select embedding model if not provided
	if modelPath == "" {
		downloaded := s.manager.GetDownloadedEmbeddingModels()
		if len(downloaded) == 0 {
			return fmt.Errorf("no embedding models available")
		}
		modelPath = filepath.Join(s.manager.ModelsDir, downloaded[0].Filename)
		log.Printf("🤖 Auto-selected embedding model: %s", downloaded[0].Name)
	}

	// Verify model file exists
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return fmt.Errorf("embedding model file not found: %s", modelPath)
	}

	// Unload previous embedding model if exists
	if s.embModel != 0 {
		log.Println("♻️  Unloading previous embedding model...")
		if s.embContext != 0 {
			llama.Free(s.embContext)
			s.embContext = 0
		}
		if s.embModel != 0 {
			llama.ModelFree(s.embModel)
			s.embModel = 0
		}
	}

	log.Printf("📥 Loading embedding model: %s", filepath.Base(modelPath))

	// Load model
	mParams := llama.ModelDefaultParams()
	s.embModel = llama.ModelLoadFromFile(modelPath, mParams)
	if s.embModel == 0 {
		return fmt.Errorf("failed to load embedding model")
	}

	// Get vocabulary
	s.embVocab = llama.ModelGetVocab(s.embModel)

	// Create context for embeddings
	ctxParams := llama.ContextDefaultParams()
	ctxParams.NCtx = 2048
	ctxParams.NBatch = 512
	ctxParams.NThreads = 4
	ctxParams.Embeddings = 1 // Enable embeddings

	s.embContext = llama.InitFromModel(s.embModel, ctxParams)
	if s.embContext == 0 {
		llama.ModelFree(s.embModel)
		s.embModel = 0
		return fmt.Errorf("failed to create embedding context")
	}

	s.embModelPath = modelPath
	log.Printf("✅ Embedding model loaded successfully: %s", filepath.Base(modelPath))

	return nil
}

// Generate generates text from a prompt
func (s *LibraryService) Generate(prompt string, maxTokens int32) (string, error) {
	s.chatMutex.Lock()
	defer s.chatMutex.Unlock()

	if s.chatModel == 0 || s.chatContext == 0 {
		return "", fmt.Errorf("chat model not loaded")
	}

	// Format prompt with chat template if available
	template := llama.ModelChatTemplate(s.chatModel, "")
	if template == "" {
		template = "chatml" // Default template
	}

	// Create a simple chat message
	messages := []llama.ChatMessage{
		llama.NewChatMessage("user", prompt),
	}

	// Apply chat template
	buf := make([]byte, 8192)
	length := llama.ChatApplyTemplate(template, messages, true, buf)
	formattedPrompt := string(buf[:length])

	// Tokenize formatted prompt
	tokens := llama.Tokenize(s.chatVocab, formattedPrompt, true, true)
	if len(tokens) == 0 {
		return "", fmt.Errorf("failed to tokenize prompt")
	}

	// Create batch from tokens
	batch := llama.BatchGetOne(tokens)

	// Handle encoder models
	if llama.ModelHasEncoder(s.chatModel) {
		llama.Encode(s.chatContext, batch)
		start := llama.ModelDecoderStartToken(s.chatModel)
		if start == llama.TokenNull {
			start = llama.VocabBOS(s.chatVocab)
		}
		batch = llama.BatchGetOne([]llama.Token{start})
	}

	// Generate tokens
	var response strings.Builder
	for pos := int32(0); pos < maxTokens; pos += batch.NTokens {
		llama.Decode(s.chatContext, batch)
		token := llama.SamplerSample(s.chatSampler, s.chatContext, -1)

		// Check for end of generation
		if llama.VocabIsEOG(s.chatVocab, token) {
			break
		}

		// Convert token to text
		tokenBuf := make([]byte, 256)
		tokenLength := llama.TokenToPiece(s.chatVocab, token, tokenBuf, 0, false)
		response.Write(tokenBuf[:tokenLength])

		// Prepare next batch
		batch = llama.BatchGetOne([]llama.Token{token})
	}

	return response.String(), nil
}

// GenerateEmbedding generates embeddings for the given text
func (s *LibraryService) GenerateEmbedding(text string) ([]float32, error) {
	s.embMutex.Lock()
	defer s.embMutex.Unlock()

	if s.embModel == 0 || s.embContext == 0 {
		return nil, fmt.Errorf("embedding model not loaded")
	}

	// Tokenize text
	tokens := llama.Tokenize(s.embVocab, text, true, false)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("failed to tokenize text")
	}

	// Create batch
	batch := llama.BatchGetOne(tokens)

	// Decode to get embeddings
	llama.Decode(s.embContext, batch)

	// Get embeddings from context
	nEmbd := llama.ModelNEmbd(s.embModel)
	embeddings := llama.GetEmbeddingsSeq(s.embContext, 0, nEmbd)

	// Copy embeddings to slice
	result := make([]float32, nEmbd)
	for i := int32(0); i < nEmbd; i++ {
		result[i] = embeddings[i]
	}

	return result, nil
}

// selectBestModel automatically selects the best available model
func (s *LibraryService) selectBestModel() (string, error) {
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
		if !strings.HasSuffix(strings.ToLower(name), ".gguf") {
			continue
		}

		// Skip embedding models for chat
		nameLower := strings.ToLower(name)
		if strings.Contains(nameLower, "embedding") || strings.Contains(nameLower, "embed") {
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
		return "", fmt.Errorf("no GGUF models found in %s", modelsDir)
	}

	// Selection strategy: prefer larger, higher quality models
	var bestModel string
	var bestScore int

	for _, model := range models {
		score := 0
		nameLower := strings.ToLower(model.name)

		// Prefer larger models (more parameters = better quality)
		sizeMB := model.size / (1024 * 1024)
		if sizeMB >= 4000 {
			score += 100 // Large models (7B+) get highest priority
		} else if sizeMB >= 2000 {
			score += 50 // Medium models (3B)
		} else if sizeMB >= 1000 {
			score += 25 // Small models (1.5B)
		} else {
			score += 10 // Tiny models (0.5B)
		}

		// Prefer higher quantization quality
		if strings.Contains(nameLower, "q8") {
			score += 50 // Q8 = highest quality
		} else if strings.Contains(nameLower, "q6") {
			score += 40 // Q6 = very high quality
		} else if strings.Contains(nameLower, "q5") {
			score += 30 // Q5 = high quality
		} else if strings.Contains(nameLower, "q4") {
			score += 20 // Q4 = medium quality
		}

		// Prefer instruct/chat models
		if strings.Contains(nameLower, "instruct") || strings.Contains(nameLower, "chat") {
			score += 30
		}

		// Prefer Mistral over Qwen (Mistral generally better quality)
		if strings.Contains(nameLower, "mistral") {
			score += 20
		} else if strings.Contains(nameLower, "qwen") {
			score += 10
		}

		if score > bestScore {
			bestScore = score
			bestModel = model.path
		}
	}

	if bestModel == "" {
		bestModel = models[0].path
	}

	return bestModel, nil
}

// GetAvailableModels returns a list of available GGUF models
func (s *LibraryService) GetAvailableModels() ([]string, error) {
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
func (s *LibraryService) GetModelsDirectory() string {
	return s.manager.GetModelsDirectory()
}

// GetEmbeddingModelsDirectory returns the directory where embedding models are stored
func (s *LibraryService) GetEmbeddingModelsDirectory() string {
	return s.manager.GetModelsDirectory()
}

// IsChatModelLoaded returns true if a chat model is loaded
func (s *LibraryService) IsChatModelLoaded() bool {
	s.chatMutex.Lock()
	defer s.chatMutex.Unlock()
	return s.chatModel != 0 && s.chatContext != 0
}

// IsEmbeddingModelLoaded returns true if an embedding model is loaded
func (s *LibraryService) IsEmbeddingModelLoaded() bool {
	s.embMutex.Lock()
	defer s.embMutex.Unlock()
	return s.embModel != 0 && s.embContext != 0
}

// GetLoadedChatModel returns the path of the currently loaded chat model
func (s *LibraryService) GetLoadedChatModel() string {
	s.chatMutex.Lock()
	defer s.chatMutex.Unlock()
	return s.chatModelPath
}

// GetLoadedEmbeddingModel returns the path of the currently loaded embedding model
func (s *LibraryService) GetLoadedEmbeddingModel() string {
	s.embMutex.Lock()
	defer s.embMutex.Unlock()
	return s.embModelPath
}

// Cleanup releases all loaded models and frees resources
func (s *LibraryService) Cleanup() {
	s.chatMutex.Lock()
	if s.chatContext != 0 {
		llama.Free(s.chatContext)
		s.chatContext = 0
	}
	if s.chatModel != 0 {
		llama.ModelFree(s.chatModel)
		s.chatModel = 0
	}
	s.chatMutex.Unlock()

	s.embMutex.Lock()
	if s.embContext != 0 {
		llama.Free(s.embContext)
		s.embContext = 0
	}
	if s.embModel != 0 {
		llama.ModelFree(s.embModel)
		s.embModel = 0
	}
	s.embMutex.Unlock()

	if s.isInitialized {
		llama.BackendFree()
		s.isInitialized = false
	}

	log.Println("✅ LibraryService cleaned up")
}

// AutoDownloadRecommendedModel downloads the recommended model based on hardware
// Delegates to installer methods
func (s *LibraryService) AutoDownloadRecommendedModel() error {
	return s.manager.AutoDownloadRecommendedChatModel()
}

// Note: Download-related methods have been moved to LlamaCppInstaller.
// Use installer methods:
//   - manager.DownloadChatModel(modelSpec)
//   - manager.DownloadEmbeddingModel(model)
//   - manager.AutoDownloadRecommendedChatModel()
//   - manager.AutoDownloadRecommendedEmbeddingModel()
