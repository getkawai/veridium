package llama

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/kawai-network/veridium/pkg/yzma/llama"
	"github.com/kawai-network/veridium/pkg/yzma/mtmd"
)

// LibraryService provides LLM inference using llama.cpp as a library (via yzma)
// This replaces the binary-based approach with direct library calls
type LibraryService struct {
	installer *LlamaCppInstaller

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

	// Vision-Language (VL) model state
	vlModel     llama.Model
	vlContext   llama.Context
	vlVocab     llama.Vocab
	vlSampler   llama.Sampler
	vlMTMDCtx   mtmd.Context
	vlModelPath string
	vlMutex     sync.Mutex

	initOnce sync.Once
	initChan chan struct{} // Closed when library is initialized
}

// NewLibraryService creates a new library-based llama.cpp service
func NewLibraryService() (*LibraryService, error) {
	installer := NewLlamaCppInstaller()

	service := &LibraryService{
		installer: installer,
		initChan:  make(chan struct{}),
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
	if !s.installer.IsLlamaCppInstalled() {
		log.Println("🔧 llama.cpp not found, attempting auto-installation...")

		// Use InstallLlamaCpp which now uses download.InstallLibraries
		// This handles version management, auto-upgrade, and fallback automatically
		if err := s.installer.InstallLlamaCpp(); err != nil {
			log.Printf("⚠️  Failed to install llama.cpp: %v", err)
			log.Printf("   llama.cpp features will not be available")
			return
		}

		log.Println("✅ llama.cpp installed successfully")
	} else {
		log.Println("✅ llama.cpp is already installed")
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
		if err := s.installer.AutoDownloadRecommendedEmbeddingModel(); err != nil {
			log.Printf("⚠️  Failed to auto-download embedding model: %v", err)
		} else {
			log.Println("✅ Embedding model ready!")
		}
	}()

	// Step 4: Auto-load chat model if available (in BACKGROUND to avoid blocking UI)
	go func() {
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
	}()

	// Step 5: Auto-download VL model (in background)
	go func() {
		log.Println("📦 Checking VL models...")
		if err := s.AutoDownloadRecommendedVLModel(); err != nil {
			log.Printf("⚠️  Failed to auto-download VL model: %v", err)
		} else {
			log.Println("✅ VL model ready!")
		}
	}()
}

// InitializeLibrary loads the llama.cpp shared library
func (s *LibraryService) InitializeLibrary() error {
	s.initMutex.Lock()
	defer s.initMutex.Unlock()

	if s.isInitialized {
		return nil
	}

	// Get library path from installer (programmatic, not env var)
	libPath := s.installer.GetLibraryPath()

	// Verify all required libraries exist (libggml, libggml-base, libllama)
	if !s.installer.VerifyAllLibrariesExist() {
		return fmt.Errorf("llama.cpp libraries not found in %s. Please run installer first", libPath)
	}

	s.libPath = libPath
	log.Printf("📚 Loading llama.cpp library from directory: %s", libPath)

	// Log which libraries were found
	requiredPaths := s.installer.GetRequiredLibraryPaths()
	for _, path := range requiredPaths {
		log.Printf("  ✓ Found: %s", filepath.Base(path))
	}

	// Load the library (llama.Load expects a directory path)
	if err := llama.Load(libPath); err != nil {
		return fmt.Errorf("failed to load llama.cpp library from %s: %w", libPath, err)
	}

	// Load the mtmd library for multimodal/VL support
	if err := mtmd.Load(libPath); err != nil {
		log.Printf("⚠️  Failed to load mtmd library (VL features may not work): %v", err)
		// Don't fail completely - VL is optional
	} else {
		log.Println("✅ mtmd library loaded (VL support enabled)")
	}

	// Initialize llama.cpp backend
	llama.Init()

	s.isInitialized = true

	// Signal initialization complete (idempotent close)
	select {
	case <-s.initChan:
	default:
		close(s.initChan)
	}

	log.Println("✅ llama.cpp library loaded and backend initialized successfully")

	return nil
}

// WaitForInitialization waits for the library to be initialized
func (s *LibraryService) WaitForInitialization(ctx context.Context) error {
	select {
	case <-s.initChan:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
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
	var err error
	s.chatModel, err = llama.ModelLoadFromFile(modelPath, mParams)
	if err != nil || s.chatModel == 0 {
		return fmt.Errorf("failed to load model from file: %w", err)
	}

	// Get vocabulary
	s.chatVocab = llama.ModelGetVocab(s.chatModel)

	// Create context with reasonable defaults
	// CRITICAL: Batch size must be >= max tokens processed in a single batch
	// For long prompts, n_tokens_all can exceed n_batch, causing GGML_ASSERT failure
	// Increasing batch size to handle longer prompts safely
	ctxParams := llama.ContextDefaultParams()
	ctxParams.NCtx = 16384  // Context size - increased from 4096 to support long conversations
	ctxParams.NBatch = 2048 // Batch size - increased from 512 to handle long prompts
	ctxParams.NThreads = int32(runtime.NumCPU())

	s.chatContext, err = llama.InitFromModel(s.chatModel, ctxParams)
	if err != nil || s.chatContext == 0 {
		llama.ModelFree(s.chatModel)
		s.chatModel = 0
		return fmt.Errorf("failed to create context: %w", err)
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
		downloaded := s.installer.GetDownloadedEmbeddingModels()
		if len(downloaded) == 0 {
			return fmt.Errorf("no embedding models available")
		}
		modelPath = filepath.Join(s.installer.ModelsDir, downloaded[0].Filename)
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
	var err error
	s.embModel, err = llama.ModelLoadFromFile(modelPath, mParams)
	if err != nil || s.embModel == 0 {
		return fmt.Errorf("failed to load embedding model: %w", err)
	}

	// Get vocabulary
	s.embVocab = llama.ModelGetVocab(s.embModel)

	// Create context for embeddings
	ctxParams := llama.ContextDefaultParams()
	ctxParams.NCtx = 2048
	ctxParams.NBatch = 512
	ctxParams.NThreads = 4
	ctxParams.Embeddings = 1 // Enable embeddings

	s.embContext, err = llama.InitFromModel(s.embModel, ctxParams)
	if err != nil || s.embContext == 0 {
		llama.ModelFree(s.embModel)
		s.embModel = 0
		return fmt.Errorf("failed to create embedding context: %w", err)
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
	// Increased buffer size to handle long conversation histories
	buf := make([]byte, 65536) // 64KB buffer for long conversations
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

// LoadVLModel loads a Vision-Language (VL) model
// If modelPath is empty, automatically selects the best available VL model
func (s *LibraryService) LoadVLModel(modelPath string) error {
	s.vlMutex.Lock()
	defer s.vlMutex.Unlock()

	// Ensure library is initialized
	if !s.isInitialized {
		if err := s.InitializeLibrary(); err != nil {
			return fmt.Errorf("failed to initialize library: %w", err)
		}
	}

	// Auto-select VL model if not provided
	if modelPath == "" {
		autoModel, err := s.selectBestVLModel()
		if err != nil {
			return fmt.Errorf("failed to auto-select VL model: %w", err)
		}
		modelPath = autoModel
		log.Printf("🤖 Auto-selected VL model: %s", modelPath)
	}

	// Verify model file exists
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return fmt.Errorf("VL model file not found: %s", modelPath)
	}

	// Unload previous VL model if exists
	if s.vlMTMDCtx != 0 {
		log.Println("♻️  Unloading previous VL multimodal context...")
		mtmd.Free(s.vlMTMDCtx)
		s.vlMTMDCtx = 0
	}
	if s.vlContext != 0 {
		llama.Free(s.vlContext)
		s.vlContext = 0
	}
	if s.vlModel != 0 {
		llama.ModelFree(s.vlModel)
		s.vlModel = 0
	}

	log.Printf("📥 Loading VL model: %s", filepath.Base(modelPath))

	// Load VL model
	mParams := llama.ModelDefaultParams()
	var err error
	s.vlModel, err = llama.ModelLoadFromFile(modelPath, mParams)
	if err != nil || s.vlModel == 0 {
		return fmt.Errorf("failed to load VL model from file: %w", err)
	}

	// Get vocabulary
	s.vlVocab = llama.ModelGetVocab(s.vlModel)

	// Find the MMTRoj projector file (mmproj-xxx.gguf)
	modelsDir := s.installer.GetModelsDirectory()
	projectorPath, err := s.findProjectorForModel(modelPath, modelsDir)
	if err != nil {
		llama.ModelFree(s.vlModel)
		s.vlModel = 0
		return fmt.Errorf("failed to find projector file: %w", err)
	}

	log.Printf("📱 Found projector: %s", filepath.Base(projectorPath))

	// Initialize MTMD multimodal context
	mtmdParams := mtmd.ContextParamsDefault()
	mtmdParams.UseGPU = true // Enable GPU acceleration

	// Set logging level using mtmd.LogSet instead of Verbosity field
	// mtmd.LogSet(llama.LogNormal) // Use LogNormal for standard logging, or LogSilent() to disable

	s.vlMTMDCtx, err = mtmd.InitFromFile(projectorPath, s.vlModel, mtmdParams)
	if err != nil || s.vlMTMDCtx == 0 {
		llama.ModelFree(s.vlModel)
		s.vlModel = 0
		return fmt.Errorf("failed to initialize MTMD context: %w", err)
	}

	// Check if model supports vision
	if !mtmd.SupportVision(s.vlMTMDCtx) {
		mtmd.Free(s.vlMTMDCtx)
		llama.ModelFree(s.vlModel)
		s.vlMTMDCtx = 0
		s.vlModel = 0
		return fmt.Errorf("model does not support vision")
	}

	// Create llama context for text generation
	ctxParams := llama.ContextDefaultParams()
	ctxParams.NCtx = 16384  // Context size - increased from 4096 to support long conversations
	ctxParams.NBatch = 2048 // Batch size
	ctxParams.NThreads = int32(runtime.NumCPU())

	s.vlContext, err = llama.InitFromModel(s.vlModel, ctxParams)
	if err != nil || s.vlContext == 0 {
		mtmd.Free(s.vlMTMDCtx)
		llama.ModelFree(s.vlModel)
		s.vlMTMDCtx = 0
		s.vlModel = 0
		return fmt.Errorf("failed to create VL context: %w", err)
	}

	// Create sampler chain for VL responses
	s.vlSampler = llama.SamplerChainInit(llama.SamplerChainDefaultParams())
	llama.SamplerChainAdd(s.vlSampler, llama.SamplerInitTopK(40))
	llama.SamplerChainAdd(s.vlSampler, llama.SamplerInitTopP(0.95, 1))
	llama.SamplerChainAdd(s.vlSampler, llama.SamplerInitTempExt(0.1, 0, 1.0)) // Lower temp for VL
	llama.SamplerChainAdd(s.vlSampler, llama.SamplerInitDist(llama.DefaultSeed))

	s.vlModelPath = modelPath
	log.Printf("✅ VL model loaded successfully: %s", filepath.Base(modelPath))

	return nil
}

// ProcessImageWithText processes an image with accompanying text using VL model
func (s *LibraryService) ProcessImageWithText(imagePath, prompt string, maxTokens int32) (string, error) {
	s.vlMutex.Lock()
	defer s.vlMutex.Unlock()

	if s.vlModel == 0 || s.vlContext == 0 || s.vlMTMDCtx == 0 {
		return "", fmt.Errorf("VL model not loaded")
	}

	// Verify image file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return "", fmt.Errorf("image file not found: %s", imagePath)
	}

	log.Printf("🖼️  Processing image: %s", filepath.Base(imagePath))
	log.Printf("💬 Prompt: %s", prompt)

	// Load image using MTMD
	bitmap := mtmd.BitmapInitFromFile(s.vlMTMDCtx, imagePath)
	if bitmap == 0 {
		return "", fmt.Errorf("failed to load image")
	}
	defer mtmd.BitmapFree(bitmap)

	// Create input text with image placeholder
	imagePlaceholder := mtmd.DefaultMarker() // Usually "<__media__>"
	fullPrompt := "Here is an image: " + imagePlaceholder + "\n\n" + prompt + "\n\nDescribe what you see in detail."

	inputText := mtmd.NewInputText(fullPrompt, true, true)

	// Initialize output chunks container (must be initialized before Tokenize)
	outputChunks := mtmd.InputChunksInit()
	defer mtmd.InputChunksFree(outputChunks)

	// Tokenize input (text + image)
	bitmaps := []mtmd.Bitmap{bitmap}
	if mtmd.Tokenize(s.vlMTMDCtx, outputChunks, inputText, bitmaps) != 0 {
		return "", fmt.Errorf("failed to tokenize input with image")
	}

	// Get new n_past after processing chunks
	var newNPast llama.Pos
	nPast := llama.Pos(0)
	seqID := llama.SeqId(0)

	// Process multimodal input using helper function
	// logitsLast must be true to compute logits for the last token (needed for sampling)
	nBatch := int32(512)
	logitsLast := true
	if mtmd.HelperEvalChunks(s.vlMTMDCtx, s.vlContext, outputChunks, nPast, seqID, nBatch, logitsLast, &newNPast) != 0 {
		return "", fmt.Errorf("failed to process multimodal input")
	}

	// Generate response
	var response strings.Builder
	pos := int32(0)
	for pos < maxTokens {
		token := llama.SamplerSample(s.vlSampler, s.vlContext, -1)

		// Check for end of generation
		if llama.VocabIsEOG(s.vlVocab, token) {
			break
		}

		// Convert token to text
		tokenBuf := make([]byte, 256)
		tokenLength := llama.TokenToPiece(s.vlVocab, token, tokenBuf, 0, false)
		response.Write(tokenBuf[:tokenLength])

		// Prepare next batch
		batch := llama.BatchGetOne([]llama.Token{token})

		// Decode next token
		llama.Decode(s.vlContext, batch)
		pos += batch.NTokens
	}

	result := response.String()
	log.Printf("✅ Image processed successfully")

	return result, nil
}

// findProjectorForModel finds the corresponding MMTRoj projector file for a VL model
func (s *LibraryService) findProjectorForModel(modelPath, modelsDir string) (string, error) {
	modelName := strings.TrimSuffix(filepath.Base(modelPath), ".gguf")

	entries, err := os.ReadDir(modelsDir)
	if err != nil {
		return "", fmt.Errorf("failed to read models directory: %w", err)
	}

	// Look for projector files that match the model
	possiblePrefixes := []string{
		"mmproj-" + modelName,
		"mmproj",
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".gguf") {
			continue
		}

		for _, prefix := range possiblePrefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				return filepath.Join(modelsDir, name), nil
			}
		}
	}

	return "", fmt.Errorf("no matching projector file found for model %s", modelName)
}

// selectBestVLModel automatically selects the best available VL model
func (s *LibraryService) selectBestVLModel() (string, error) {
	modelsDir := s.installer.GetModelsDirectory()

	VLModels, err := s.installer.GetAvailableVLModels()
	if err != nil {
		return "", fmt.Errorf("failed to get available VL models: %w", err)
	}

	if len(VLModels) == 0 {
		return "", fmt.Errorf("no VL models found in %s. Run AutoDownloadRecommendedVLModel() first", modelsDir)
	}

	// For now, just return the first (largest) model
	// TODO: Implement scoring based on RAM and quality like selectBestModel
	return VLModels[0], nil
}

// IsVLModelLoaded returns true if a VL model is loaded
func (s *LibraryService) IsVLModelLoaded() bool {
	s.vlMutex.Lock()
	defer s.vlMutex.Unlock()
	return s.vlModel != 0 && s.vlContext != 0 && s.vlMTMDCtx != 0
}

// GetLoadedVLModel returns the path of the currently loaded VL model
func (s *LibraryService) GetLoadedVLModel() string {
	s.vlMutex.Lock()
	defer s.vlMutex.Unlock()
	return s.vlModelPath
}

// AutoDownloadRecommendedVLModel downloads the recommended VL model with hardware detection
// Delegates to installer methods
func (s *LibraryService) AutoDownloadRecommendedVLModel() error {
	log.Println("📦 Auto-downloading VL model...")
	if err := s.installer.AutoDownloadRecommendedVLModel(); err != nil {
		return fmt.Errorf("failed to download VL model: %w", err)
	}

	// Also check if projector file exists, if not it will be downloaded with the model
	return s.AutoDownloadRecommendedVLProjector()
}

// AutoDownloadRecommendedVLProjector ensures projector files are available for VL models
func (s *LibraryService) AutoDownloadRecommendedVLProjector() error {
	// VL models typically come with their own projector files, so this is usually handled
	// by the model download. This method can be extended if separate projector downloads
	// are needed in the future.
	log.Println("✅ VL projector files should be available with the model")
	return nil
}

// selectBestModel automatically selects the best available model
// For non-reasoning mode: prefer Llama 3.2 (non-reasoning models)
// For reasoning mode: prefer Qwen3 (reasoning models)
func (s *LibraryService) selectBestModel() (string, error) {
	modelsDir := s.installer.GetModelsDirectory()

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

	// Selection strategy: prefer non-reasoning models (Llama, Mistral) over reasoning models (Qwen)
	// Reasoning models should only be used when explicitly requested
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

		// CRITICAL: Prefer non-reasoning models (Llama, Mistral) by default
		// Reasoning models (Qwen) generate <think> tags which should be avoided in default mode
		if strings.Contains(nameLower, "llama") {
			score += 100 // Llama is best for non-reasoning (no think tags)
		} else if strings.Contains(nameLower, "mistral") {
			score += 80 // Mistral is also good for non-reasoning
		} else if strings.Contains(nameLower, "qwen") {
			score -= 50 // Penalize Qwen (reasoning model) - only use if explicitly requested
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
	modelsDir := s.installer.GetModelsDirectory()

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
	return s.installer.GetModelsDirectory()
}

// GetEmbeddingModelsDirectory returns the directory where embedding models are stored
func (s *LibraryService) GetEmbeddingModelsDirectory() string {
	return s.installer.GetModelsDirectory()
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
	// Cleanup VL model
	s.vlMutex.Lock()
	if s.vlMTMDCtx != 0 {
		mtmd.Free(s.vlMTMDCtx)
		s.vlMTMDCtx = 0
	}
	if s.vlContext != 0 {
		llama.Free(s.vlContext)
		s.vlContext = 0
	}
	if s.vlModel != 0 {
		llama.ModelFree(s.vlModel)
		s.vlModel = 0
	}
	s.vlMutex.Unlock()

	// Cleanup chat model
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

	// Cleanup embedding model
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
	return s.installer.AutoDownloadRecommendedChatModel()
}

// Note: Download-related methods have been moved to LlamaCppInstaller.
// Use installer methods:
//   - manager.DownloadChatModel(modelSpec)
//   - manager.DownloadEmbeddingModel(model)
//   - manager.AutoDownloadRecommendedChatModel()
//   - manager.AutoDownloadRecommendedEmbeddingModel()

// GetHardwareSpecs returns the detected hardware specifications
// This is used for validating if the system can handle reasoning models
func (s *LibraryService) GetHardwareSpecs() *HardwareSpecs {
	if s.installer == nil {
		return nil
	}
	return s.installer.HardwareSpecs
}
