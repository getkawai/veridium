// Package main is the entry point for the Veridium desktop application.
// Veridium is a Wails v3 application that provides AI-powered features including:
//   - Local LLM inference via llama.cpp
//   - Vector search and RAG (Retrieval-Augmented Generation)
//   - Speech-to-text (Whisper) and text-to-speech
//   - Knowledge base management
//   - File processing with automatic embedding
package main

import (
	"context"
	"embed"
	"log"
	"os"
	"path/filepath"

	"github.com/kawai-network/veridium/fantasy"
	"github.com/kawai-network/veridium/fantasy/llamalib"
	"github.com/kawai-network/veridium/fantasy/llamalib/tools/builtin"
	llamaprovider "github.com/kawai-network/veridium/fantasy/providers/llama"
	llamaembed "github.com/kawai-network/veridium/fantasy/providers/llama-embed"
	"github.com/kawai-network/veridium/fantasy/providers/openrouter"
	"github.com/kawai-network/veridium/internal/audio_recorder"
	"github.com/kawai-network/veridium/internal/database"
	db "github.com/kawai-network/veridium/internal/database/generated"
	"github.com/kawai-network/veridium/internal/machineid"
	"github.com/kawai-network/veridium/internal/search"
	"github.com/kawai-network/veridium/internal/services"
	"github.com/kawai-network/veridium/internal/services/cache"
	"github.com/kawai-network/veridium/internal/stablediffusion"
	"github.com/kawai-network/veridium/internal/tableviewer"
	"github.com/kawai-network/veridium/internal/tts"
	"github.com/kawai-network/veridium/internal/whisper"
	"github.com/kawai-network/veridium/pkg/localfs"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/services/fileserver"
	"github.com/wailsapp/wails/v3/pkg/services/kvstore"
	wailslog "github.com/wailsapp/wails/v3/pkg/services/log"
	"github.com/wailsapp/wails/v3/pkg/services/notifications"
	"github.com/wailsapp/wails/v3/pkg/services/sqlite"
)

// ============================================================================
// EMBEDDED ASSETS
// ============================================================================

// assets embeds the frontend build output (frontend/dist) into the binary.
// This allows the application to serve the frontend without external files.
//
//go:embed all:frontend/dist
var assets embed.FS

// ============================================================================
// CONFIGURATION
// ============================================================================

// devMode enables development mode which disables caching for faster iteration.
// Set VERIDIUM_DEV=1 or DEV=1 environment variable to enable.
var devMode = os.Getenv("VERIDIUM_DEV") == "1" || os.Getenv("DEV") == "1"

const (
	// fileBaseDir is the directory for storing uploaded/processed files.
	fileBaseDir = "files"

	// duckDBPath is the path to the DuckDB database file for vector storage.
	duckDBPath = "data/duckdb.db"

	// embeddingDims is the dimension size for the embedding model (granite-embedding uses 384).
	embeddingDims = 384

	// defaultUserID is the default user ID for single-user desktop mode.
	defaultUserID = "DEFAULT_LOBE_CHAT_USER"
)

// ============================================================================
// APPLICATION CONTEXT
// ============================================================================

// AppContext holds all application services and state.
// It provides a centralized place to manage service lifecycle and dependencies.
type AppContext struct {
	// Core Services
	DB            *database.Service // SQLite database service
	Queries       *db.Queries       // Generated sqlc queries
	UserConfigDir string            // User's config directory path

	// AI/ML Services
	LibService   *llamalib.Service   // llama.cpp library service for LLM inference
	Embedder     llamaembed.Embedder // Text embedding service
	CacheManager *cache.CacheManager // Cache manager for embeddings and LLM responses

	// Storage Services
	DuckDBStore *services.DuckDBStore // DuckDB vector store for embeddings
	FileLoader  *services.FileLoader  // File loading and parsing service

	// Feature Services
	SearchService  *search.Service                      // Web search service
	TTSService     *tts.TTSService                      // Text-to-speech service (native OS)
	WhisperService *whisper.Service                     // Speech-to-text service (Whisper)
	AudioRecorder  *audio_recorder.AudioRecorderService // Microphone recording service
	VectorSearch   *services.VectorSearchService        // Semantic vector search service
	FileProcessor  *FileProcessorService                // File processing and RAG pipeline
	KBService      *services.KnowledgeBaseService       // Knowledge base management

	// Language Models - ChainLanguageModel instances for fallback support
	ChatModel    fantasy.LanguageModel // Primary chat completion model
	TitleModel   fantasy.LanguageModel // Model for generating chat titles
	SummaryModel fantasy.LanguageModel // Model for summarization tasks
	CleanupModel fantasy.LanguageModel // Model for text cleanup/formatting after OCR/Whisper

	// cleanupFuncs stores cleanup functions to be called on shutdown (LIFO order)
	cleanupFuncs []func()
}

// NewAppContext creates a new AppContext with initialized user config directory.
func NewAppContext() *AppContext {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		userConfigDir = "."
	}
	return &AppContext{
		UserConfigDir: userConfigDir,
		cleanupFuncs:  make([]func(), 0),
	}
}

// AddCleanup registers a cleanup function to be called on application shutdown.
// Functions are called in reverse order (LIFO) during Cleanup().
func (ctx *AppContext) AddCleanup(fn func()) {
	ctx.cleanupFuncs = append(ctx.cleanupFuncs, fn)
}

// Cleanup executes all registered cleanup functions in reverse order.
// This ensures proper resource cleanup (e.g., closing database connections).
func (ctx *AppContext) Cleanup() {
	for i := len(ctx.cleanupFuncs) - 1; i >= 0; i-- {
		ctx.cleanupFuncs[i]()
	}
}

// ============================================================================
// SERVICE INITIALIZATION
// ============================================================================

// initDatabase initializes the SQLite database service and sqlc queries.
// Returns error if database initialization fails (fatal for the application).
func (ctx *AppContext) initDatabase() error {
	dbService, err := database.NewService()
	if err != nil {
		return err
	}
	ctx.DB = dbService
	ctx.Queries = dbService.Queries()
	ctx.AddCleanup(func() { dbService.Close() })
	return nil
}

// initBasicServices initializes non-critical services (search, TTS, Whisper, audio).
// These services are optional - failures are logged but don't stop the application.
func (ctx *AppContext) initBasicServices() {
	ctx.SearchService = search.NewService()

	if ttsService, err := tts.NewTTSService(); err != nil {
		log.Printf("Warning: Failed to initialize TTS service: %v", err)
	} else {
		ctx.TTSService = ttsService
		log.Printf("TTS service initialized (platform: %s)", ttsService.GetPlatformInfo()["platform"])
	}

	if whisperService, err := whisper.NewService(); err != nil {
		log.Printf("Warning: Failed to initialize Whisper service: %v", err)
	} else {
		ctx.WhisperService = whisperService
		ctx.AddCleanup(func() { whisperService.Close() })
		log.Printf("Whisper STT service initialized (models: %s)", whisperService.GetModelsDirectory())
	}

	ctx.AudioRecorder = audio_recorder.NewAudioRecorderService(nil)
}

// initLlamaService initializes the llama.cpp library service.
// This MUST be called before initVectorStore and initEmbedder as they depend on llama.cpp.
func (ctx *AppContext) initLlamaService() {
	ctx.LibService = llamalib.NewService()
}

// initVectorStore initializes DuckDB for vector storage.
// Creates the files directory and sets up the vector database for embedding storage.
func (ctx *AppContext) initVectorStore() {
	os.MkdirAll(fileBaseDir, 0755)

	duckDBStore, err := services.NewDuckDBStore(duckDBPath, embeddingDims)
	if err != nil {
		log.Printf("Warning: Failed to initialize DuckDB Store: %v", err)
		return
	}
	ctx.DuckDBStore = duckDBStore
	ctx.AddCleanup(func() { duckDBStore.Close() })
	log.Printf("DuckDB Store initialized (path: %s)", duckDBPath)
}

// initEmbedder initializes the text embedding service using llama.cpp.
// In production mode, embeddings are cached to improve performance.
// In dev mode, caching is disabled for faster iteration.
func (ctx *AppContext) initEmbedder() {
	embeddingModelName := llamalib.GetRecommendedEmbeddingModel()
	embeddingModel, exists := llamalib.GetEmbeddingModel(embeddingModelName)
	if !exists {
		log.Printf("Warning: Embedding model not found: %s", embeddingModelName)
		return
	}

	installer := llamalib.NewLlamaCppInstaller()
	modelPath := filepath.Join(installer.GetModelsDirectory(), embeddingModel.Filename)

	baseEmbedder, err := llamaembed.NewLlamaEmbedder(&llamaembed.LlamaConfig{
		ModelPath:   modelPath,
		ContextSize: 2048,
	})
	if err != nil {
		log.Printf("Warning: Failed to create embedder: %v", err)
		return
	}
	ctx.AddCleanup(func() { baseEmbedder.Close() })
	log.Printf("Embedder initialized (model: %s, dims: %d)", embeddingModel.Name, baseEmbedder.Dimensions())

	if devMode {
		ctx.Embedder = baseEmbedder
		log.Printf("Cache disabled (dev mode)")
	} else {
		ctx.CacheManager = cache.NewCacheManager(baseEmbedder, nil)
		ctx.Embedder = ctx.CacheManager.GetCachedEmbedder(baseEmbedder)
		log.Printf("Cache layer initialized")
	}
}

// initVectorSearch initializes the semantic vector search service.
// Requires Embedder to be initialized first. Combines DuckDB and SQLite for search.
func (ctx *AppContext) initVectorSearch() {
	if ctx.Embedder == nil {
		log.Printf("Warning: Embedder not available, Vector Search service disabled")
		return
	}

	vectorSearch, err := services.NewVectorSearchService(ctx.DB.DB(), ctx.DuckDBStore, ctx.Embedder)
	if err != nil {
		log.Printf("Warning: Failed to initialize Vector Search service: %v", err)
		return
	}
	ctx.VectorSearch = vectorSearch
	log.Printf("Vector Search service initialized (DuckDB + SQLite)")
}

// initFileProcessor initializes the file processing service.
// Handles file parsing, document storage, RAG processing, and video transcription.
func (ctx *AppContext) initFileProcessor() {
	ctx.FileLoader = services.NewFileLoader()
	ctx.FileProcessor = NewFileProcessorService(
		ctx.DB.DB(),
		ctx.FileLoader,
		ctx.VectorSearch,
		ctx.DuckDBStore,
		ctx.LibService,
		ctx.WhisperService,
		fileBaseDir,
	)
	log.Printf("File Processor service initialized")
}

// initKnowledgeBase initializes the knowledge base service for RAG.
// Requires VectorSearch and FileProcessor to be initialized first.
func (ctx *AppContext) initKnowledgeBase() {
	if ctx.VectorSearch == nil || ctx.FileProcessor == nil {
		return
	}

	kbAssetPath := filepath.Join(ctx.UserConfigDir, "veridium", "kb-assets")
	embedder := ctx.VectorSearch.GetEmbedder()
	ragProcessor := services.NewRAGProcessor(ctx.DB.DB(), ctx.DuckDBStore, ctx.FileLoader, embedder)

	kbService, err := services.NewKnowledgeBaseService(ctx.DB, &services.KnowledgeBaseConfig{
		RAGProcessor: ragProcessor,
		VectorSearch: ctx.VectorSearch,
		FileLoader:   ctx.FileLoader,
		AssetDir:     kbAssetPath,
	})
	if err != nil {
		log.Printf("Warning: Failed to initialize Knowledge Base service: %v", err)
		return
	}
	ctx.KBService = kbService
	log.Printf("Knowledge Base service initialized (asset path: %s)", kbAssetPath)
}

// initLanguageModels initializes the language models for chat, title generation, etc.
// Creates separate model chains with appropriate criteria for each task.
func (ctx *AppContext) initLanguageModels() {
	if ctx.LibService == nil {
		return
	}

	bgCtx := context.Background()

	llamaProvider, err := llamaprovider.New(llamaprovider.WithService(ctx.LibService))
	if err != nil {
		log.Printf("Warning: Failed to create llama provider: %v", err)
		return
	}

	localModel, err := llamaProvider.LanguageModel(bgCtx, "")
	if err != nil {
		log.Printf("Warning: Failed to get language model from llama provider: %v", err)
		return
	}

	// ChatModel: needs reasoning, attachments, large context for conversation
	chatChain := ctx.buildModelChainWithCriteria(bgCtx, localModel, openrouter.ModelSelectionCriteria{
		RequireReasoning:   true,
		RequireAttachments: true,
		MinContextWindow:   100000,
	}, "ChatModel")
	ctx.ChatModel, _ = fantasy.NewChain(chatChain)

	// TitleModel: simple task, no special requirements
	titleChain := ctx.buildModelChainWithCriteria(bgCtx, localModel, openrouter.ModelSelectionCriteria{}, "TitleModel")
	ctx.TitleModel, _ = fantasy.NewChain(titleChain)

	// SummaryModel: needs large context for long input, no reasoning needed
	summaryChain := ctx.buildModelChainWithCriteria(bgCtx, localModel, openrouter.ModelSelectionCriteria{
		MinContextWindow: 50000,
	}, "SummaryModel")
	ctx.SummaryModel, _ = fantasy.NewChain(summaryChain)

	// CleanupModel: simple text cleanup, no special requirements
	cleanupChain := ctx.buildModelChainWithCriteria(bgCtx, localModel, openrouter.ModelSelectionCriteria{}, "CleanupModel")
	ctx.CleanupModel, _ = fantasy.NewChain(cleanupChain)

	log.Printf("Language models initialized with task-specific criteria")
}

// buildModelChainWithCriteria creates a chain with OpenRouter (auto-selected) + local fallback.
func (ctx *AppContext) buildModelChainWithCriteria(bgCtx context.Context, localModel fantasy.LanguageModel, criteria openrouter.ModelSelectionCriteria, taskName string) []fantasy.LanguageModel {
	var chainModels []fantasy.LanguageModel

	// TODO: Move to config file or environment variable for production
	openRouterKey := "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"
	if openRouterKey != "" {
		if provider, err := openrouter.New(openrouter.WithAPIKey(openRouterKey), openrouter.WithModelSelection(criteria)); err == nil {
			if remoteModel, err := provider.LanguageModel(bgCtx, ""); err == nil {
				catalog := openrouter.GetCatalog()
				selected := catalog.SelectFreeModel(criteria)
				modelName := "none found"
				if selected != nil {
					modelName = selected.ID
				}
				chainModels = append(chainModels, remoteModel)
				log.Printf("%s: OpenRouter (%s)", taskName, modelName)
			}
		}
	}

	chainModels = append(chainModels, localModel)
	return chainModels
}

// ============================================================================
// WAILS APPLICATION SETUP
// ============================================================================

// createWailsApp creates and configures the Wails v3 application instance.
// Sets up all services, asset handling, and macOS-specific options.
func createWailsApp(ctx *AppContext) *application.App {
	return application.New(application.Options{
		Name:        "veridium",
		Description: "A demo of using raw HTML & CSS",
		Services:    buildServiceList(ctx),
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})
}

// buildServiceList creates the list of services to register with Wails.
// Services are grouped by category: Database, Core Features, AI/ML, Storage, Utilities.
func buildServiceList(ctx *AppContext) []application.Service {
	return []application.Service{
		// Database
		application.NewService(ctx.Queries),
		application.NewService(ctx.DB),
		application.NewService(tableviewer.NewService(ctx.DB.DB())),

		// Core Features
		application.NewService(ctx.SearchService),
		application.NewService(ctx.TTSService),
		application.NewService(ctx.WhisperService),
		application.NewService(ctx.AudioRecorder),

		// AI/ML
		application.NewService(ctx.VectorSearch),
		application.NewService(ctx.FileProcessor),
		application.NewService(ctx.KBService),

		// File & Storage
		application.NewService(services.NewFileService(fileBaseDir)),
		application.NewService(localfs.NewService()),
		application.NewService(builtin.NewLocalSystemService()),

		// Utilities
		application.NewService(&machineid.Service{}),
		application.NewService(stablediffusion.New()),

		// Wails Native Services
		application.NewService(notifications.New()),
		application.NewService(wailslog.New()),
		application.NewService(sqlite.New()),
		application.NewService(kvstore.New()),

		// File Server
		application.NewServiceWithOptions(
			fileserver.NewWithConfig(&fileserver.Config{RootPath: fileBaseDir}),
			application.ServiceOptions{Route: "/files"},
		),
	}
}

// registerAgentServices registers services that need the Wails app instance.
// This includes ThreadManagement, AgentChat, and sets up shutdown handlers.
func registerAgentServices(app *application.App, ctx *AppContext) {
	ctx.AudioRecorder.SetApp(app)

	threadService := services.NewThreadManagementService(app, ctx.DB)
	app.RegisterService(application.NewService(threadService))

	if ctx.LibService != nil && ctx.KBService != nil {
		agentService := services.NewAgentChatService(
			app, ctx.DB, ctx.LibService, ctx.KBService, ctx.VectorSearch, threadService,
		)

		if ctx.ChatModel != nil {
			agentService.SetChatModel(ctx.ChatModel)
		}
		if ctx.TitleModel != nil {
			agentService.SetTitleModel(ctx.TitleModel)
		}
		if ctx.SummaryModel != nil {
			agentService.SetSummaryModel(ctx.SummaryModel)
		}

		app.RegisterService(application.NewService(agentService))
	}

	if ctx.CleanupModel != nil {
		ctx.FileProcessor.SetLanguageModel(ctx.CleanupModel)
		log.Printf("FileProcessorService: LLM cleanup model injected")
	}

	if ctx.LibService != nil {
		app.OnShutdown(func() {
			log.Printf("Cleaning up Llama Library service...")
			ctx.LibService.Cleanup()
		})
	}
}

// ============================================================================
// WINDOW SETUP
// ============================================================================

// createMainWindow creates and configures the main application window.
// Sets up window properties, drag-and-drop support, and macOS styling.
func createMainWindow(app *application.App, ctx *AppContext) {
	win := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:             "Veridium",
		StartState:        application.WindowStateMaximised,
		EnableDragAndDrop: true,
		Mac: application.MacWindow{
			Backdrop: application.MacBackdropTranslucent,
			TitleBar: application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

	setupDragDropHandler(win, app, ctx)
}

// setupDragDropHandler configures the drag-and-drop event handler for the window.
// Processes dropped files through FileProcessor and emits events to frontend.
func setupDragDropHandler(win *application.WebviewWindow, app *application.App, ctx *AppContext) {
	win.OnWindowEvent(events.Common.WindowDropZoneFilesDropped, func(event *application.WindowEvent) {
		droppedFiles := event.Context().DroppedFiles()
		details := event.Context().DropZoneDetails()

		log.Printf("[Drag&Drop] Files dropped: %d files", len(droppedFiles))

		processedFiles := processDroppedFiles(droppedFiles, ctx)
		emitDropEvent(app, processedFiles, details)
	})
}

// processDroppedFiles processes each dropped file through the FileProcessor.
// Returns a slice of file info maps containing URLs, IDs, and processing status.
func processDroppedFiles(files []string, ctx *AppContext) []map[string]interface{} {
	processedFiles := make([]map[string]interface{}, 0, len(files))

	for i, filePath := range files {
		log.Printf("[Drag&Drop]   %d. %s", i+1, filePath)

		result, err := ctx.FileProcessor.ProcessFileFromPath(filePath, defaultUserID)
		if err != nil {
			log.Printf("[Drag&Drop] Error processing file %s: %v", filePath, err)
			continue
		}

		fileInfo := map[string]interface{}{
			"originalPath": filePath,
			"url":          result.RelativeURL,
			"name":         result.Filename,
			"fileId":       result.FileID,
			"documentId":   result.DocumentID,
			"processing":   result.Processing,
		}
		processedFiles = append(processedFiles, fileInfo)
		log.Printf("[Drag&Drop] File processed: %s -> %s (fileId: %s)", result.Filename, result.RelativeURL, result.FileID)
	}

	return processedFiles
}

// emitDropEvent emits the "files:dropped" event to the frontend with processed file info.
// Includes drop zone details (element ID, coordinates) if available.
func emitDropEvent(app *application.App, files []map[string]interface{}, details *application.DropZoneDetails) {
	eventData := map[string]interface{}{"files": files}

	if details != nil {
		log.Printf("[Drag&Drop] Drop zone: %s at (%d, %d)", details.ElementID, details.X, details.Y)
		eventData["elementId"] = details.ElementID
		eventData["classList"] = details.ClassList
		eventData["x"] = details.X
		eventData["y"] = details.Y
		eventData["attributes"] = details.Attributes
	} else {
		log.Printf("[Drag&Drop] Drop outside specific zone")
	}

	app.Event.Emit("files:dropped", eventData)
}

// ============================================================================
// MAIN ENTRY POINT
// ============================================================================

// main is the application entry point.
// Initializes all services, creates the Wails app, and starts the event loop.
func main() {
	if devMode {
		log.Printf("Development mode enabled (cache disabled)")
	}

	ctx := NewAppContext()
	defer ctx.Cleanup()

	// Initialize all services
	if err := ctx.initDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	ctx.initBasicServices()
	ctx.initLlamaService()
	ctx.initVectorStore()
	ctx.initEmbedder()
	ctx.initVectorSearch()
	ctx.initFileProcessor()
	ctx.initKnowledgeBase()
	ctx.initLanguageModels()

	// Create and configure Wails app
	app := createWailsApp(ctx)
	registerAgentServices(app, ctx)
	createMainWindow(app, ctx)

	// Run application
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
