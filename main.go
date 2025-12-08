package main

import (
	"context"
	"embed"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/kawai-network/veridium/fantasy"
	"github.com/kawai-network/veridium/fantasy/llamalib"
	"github.com/kawai-network/veridium/fantasy/llamalib/tools/builtin"
	llamaprovider "github.com/kawai-network/veridium/fantasy/providers/llama"
	llamaembed "github.com/kawai-network/veridium/fantasy/providers/llama-embed"
	"github.com/kawai-network/veridium/fantasy/providers/openrouter"
	"github.com/kawai-network/veridium/internal/audio_recorder"
	"github.com/kawai-network/veridium/internal/database"
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

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS

// Dev mode flag - set via environment variable or build tag
var devMode = os.Getenv("VERIDIUM_DEV") == "1" || os.Getenv("DEV") == "1"

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {
	// Check dev mode
	if devMode {
		log.Printf("🔧 Development mode enabled (cache disabled)")
	}

	// Initialize database service
	dbService, err := database.NewService()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbService.Close()

	// Get queries - this is what we'll bind to Wails
	queries := dbService.Queries()

	// Create Search service
	searchService := search.NewService()

	// Initialize TTS service (native OS text-to-speech)
	ttsService, err := tts.NewTTSService()
	if err != nil {
		log.Printf("⚠️  Warning: Failed to initialize TTS service: %v", err)
	} else {
		log.Printf("✅ TTS service initialized successfully")
		log.Printf("   Platform: %s", ttsService.GetPlatformInfo()["platform"])
	}

	// Initialize Whisper STT service (offline, 99 languages, whisper-cpp CLI)
	// Auto-installs whisper-cpp and downloads recommended model in background
	whisperService, err := whisper.NewService()
	if err != nil {
		log.Printf("⚠️  Warning: Failed to initialize Whisper service: %v", err)
		log.Printf("    Speech-to-text features will not be available.")
	}

	defer whisperService.Close()
	log.Printf("✅ Whisper STT service initialized")
	log.Printf("   Models directory: %s", whisperService.GetModelsDirectory())
	log.Printf("   Auto-setup running in background...")

	// Initialize Audio Recorder service (app will be set after creation)
	audioRecorderService := audio_recorder.NewAudioRecorderService(nil)

	// Get user data directory for all services
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		userConfigDir = "."
	}

	// Initialize Llama.cpp Service FIRST (library-based LLM inference)
	// This MUST be initialized before VectorSearchService to load llama.cpp library
	// Auto-installs llama.cpp binaries and downloads models in background
	libService, err := llamalib.NewService()
	if err != nil {
		log.Printf("⚠️  Warning: Failed to initialize Llama service: %v", err)
		log.Printf("    LLM chat features will not be available.")
	} else {
		log.Printf("✅ Llama service initialized")
		log.Printf("   Models directory: %s", libService.GetModelsDirectory())
		log.Printf("   Auto-setup running in background...")
	}

	// Initialize File Service base directory (needed by both FileService and FileProcessor)
	// Use project root directory for easier development access
	fileBaseDir := filepath.Join("files")
	os.MkdirAll(fileBaseDir, 0755)

	// Initialize DuckDB Store (Vector Engine)
	duckDBPath := "data/duckdb.db"
	duckDBStore, err := services.NewDuckDBStore(duckDBPath, 384) // 384 dims for granite-embedding
	if err != nil {
		log.Printf("⚠️  Warning: Failed to initialize DuckDB Store: %v", err)
		log.Printf("    Vector search will fall back to legacy mode.")
	} else {
		log.Printf("✅ DuckDB Store initialized")
		log.Printf("   Database path: %s", duckDBPath)
		defer duckDBStore.Close()
	}

	// Initialize Embedder using pkg/yzma/embedding (custom interface, replaces Eino)
	var embedder llamaembed.Embedder
	var cacheManager *cache.CacheManager
	embeddingModelName := llamalib.GetRecommendedEmbeddingModel()
	embeddingModel, exists := llamalib.GetEmbeddingModel(embeddingModelName)
	if !exists {
		log.Printf("⚠️  Warning: Embedding model not found: %s", embeddingModelName)
	} else {
		installer := llamalib.NewLlamaCppInstaller()
		modelPath := filepath.Join(installer.GetModelsDirectory(), embeddingModel.Filename)
		baseEmbedder, embedErr := llamaembed.NewLlamaEmbedder(&llamaembed.LlamaConfig{
			ModelPath:   modelPath,
			ContextSize: 2048,
		})
		if embedErr != nil {
			log.Printf("⚠️  Warning: Failed to create embedder: %v", embedErr)
		} else {
			log.Printf("✅ Embedder initialized (pkg/yzma/embedding)")
			log.Printf("   Model: %s", embeddingModel.Name)
			log.Printf("   Dimensions: %d", baseEmbedder.Dimensions())
			defer baseEmbedder.Close()

			// Wrap embedder with cache (disabled in dev mode)
			if devMode {
				embedder = baseEmbedder
				log.Printf("⚠️  Cache disabled (dev mode)")
			} else {
				cacheManager = cache.NewCacheManager(baseEmbedder, nil)
				embedder = cacheManager.GetCachedEmbedder(baseEmbedder)
				log.Printf("✅ Cache layer initialized")
				log.Printf("   Embedding cache: max 10000 entries, TTL 24h")
				log.Printf("   LLM cache: max 1000 entries, TTL 1h, similarity threshold 0.95")
			}
		}
	}

	// Initialize Vector Search service (DuckDB + SQLite for semantic search)
	var vectorSearchService *services.VectorSearchService
	if embedder != nil {
		vectorSearchService, err = services.NewVectorSearchService(
			dbService.DB(),
			duckDBStore,
			embedder, // Now using cached embedder
		)
		if err != nil {
			log.Printf("⚠️  Warning: Failed to initialize Vector Search service: %v", err)
			log.Printf("    Semantic search features will not be available.")
		} else {
			log.Printf("✅ Vector Search service initialized (DuckDB + SQLite)")
		}
	} else {
		log.Printf("⚠️  Warning: Embedder not available, Vector Search service disabled")
	}

	// Initialize File Processor service (file parsing + document storage + RAG)
	fileLoader := services.NewFileLoader()
	fileProcessorService := NewFileProcessorService(
		dbService.DB(),
		fileLoader,
		vectorSearchService, // Pass entire service to get embedFunc
		duckDBStore,         // Pass DuckDB store
		libService,          // Pass LibraryService for VL model
		whisperService,      // Pass Whisper service for video transcription
		fileBaseDir,         // Pass file base directory for path resolution
	)
	log.Printf("✅ File Processor service initialized")
	log.Printf("   Handles: file parsing → document storage → RAG processing")
	log.Printf("   Video transcription: ffmpeg + Whisper STT")

	// Initialize Knowledge Base Service (RAG with DuckDB + SQLite)
	var kbService *services.KnowledgeBaseService
	if vectorSearchService != nil && fileProcessorService != nil {
		kbAssetPath := filepath.Join(userConfigDir, "veridium", "kb-assets")

		// Get RAGProcessor from fileProcessorService (we need to expose it)
		// For now, create a new one
		embedder := vectorSearchService.GetEmbedder()
		ragProcessor := services.NewRAGProcessor(dbService.DB(), duckDBStore, fileLoader, embedder)

		kbService, err = services.NewKnowledgeBaseService(dbService, &services.KnowledgeBaseConfig{
			RAGProcessor: ragProcessor,
			VectorSearch: vectorSearchService,
			FileLoader:   fileLoader,
			AssetDir:     kbAssetPath,
		})
		if err != nil {
			log.Printf("⚠️  Warning: Failed to initialize Knowledge Base service: %v", err)
		} else {
			log.Printf("✅ Knowledge Base service initialized (DuckDB + SQLite)")
			log.Printf("   Asset path: %s", kbAssetPath)
		}
	}

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app := application.New(application.Options{
		Name:        "veridium",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			// Database queries - direct sqlc generated code
			application.NewService(queries),
			// Database service - for transaction methods
			application.NewService(dbService),
			// TableViewer service - for database table inspection
			application.NewService(tableviewer.NewService(dbService.DB())),
			// Search service - for web search and crawling
			application.NewService(searchService),
			// TTS service - for text-to-speech (native OS)
			application.NewService(ttsService),
			// Whisper service - for speech-to-text (offline, 99 languages)
			application.NewService(whisperService),
			// Audio recorder service - for native microphone recording
			application.NewService(audioRecorderService),
			// Vector search service - for semantic search using chromem
			application.NewService(vectorSearchService),
			// File processor service - for file parsing + document storage + RAG
			application.NewService(fileProcessorService),
			// Knowledge Base service - for RAG with Chromem + Eino
			application.NewService(kbService),
			// File service - for local file storage
			application.NewService(services.NewFileService(fileBaseDir)),
			// Machine ID service
			application.NewService(&machineid.Service{}),
			// Stable Diffusion service - for image generation
			application.NewService(stablediffusion.New()),
			// Local file system service
			application.NewService(localfs.NewService()),
			// Local file system service
			application.NewService(builtin.NewLocalSystemService()),
			// Native Wails v3 notification service
			application.NewService(notifications.New()),
			// Native Wails v3 notification service
			application.NewService(wailslog.New()),
			// Native Wails v3 sqlite service
			application.NewService(sqlite.New()),
			// User data fileserver (user config directory)
			// Frontend assets are handled by Wails' built-in asset server via embed
			application.NewServiceWithOptions(
				func() *fileserver.FileserverService {
					// Use same base directory as FileService for consistency
					// This ensures uploaded files can be served via /files/ route
					return fileserver.NewWithConfig(&fileserver.Config{
						RootPath: fileBaseDir, // Same as FileService baseDir
					})
				}(),
				application.ServiceOptions{
					Route: "/files",
				},
			),
			// Native Wails v3 kvstore service
			application.NewService(kvstore.New()),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Set app instance for audio recorder service (for event emission)
	audioRecorderService.SetApp(app)

	// Initialize Thread Management Service (needs to be before AgentChatService)
	threadManagementService := services.NewThreadManagementService(app, dbService)
	app.RegisterService(application.NewService(threadManagementService))

	// Initialize Llama Chat Service (OpenAI-compatible chat API)
	// This service provides chat completion functionality using the library service
	var chatModel, titleModel, summaryModel, cleanupModel fantasy.LanguageModel
	if libService != nil {
		// Create llama provider with library service
		ctx := context.Background()
		llamaProvider, err := llamaprovider.New(llamaprovider.WithService(libService))
		if err != nil {
			log.Printf("⚠️  Warning: Failed to create llama provider: %v", err)
		} else {
			// Get the default language model from llama provider
			localModel, err := llamaProvider.LanguageModel(ctx, "")
			if err != nil {
				log.Printf("⚠️  Warning: Failed to get language model from llama provider: %v", err)
			} else {
				// Build chain of models: remote first (if configured), local fallback
				var chainModels []fantasy.LanguageModel

				// [DEV] Hardcoded OpenRouter API key for development testing
				// TODO: Move to config file or environment variable for production
				openRouterKey := "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"
				if openRouterKey != "" {
					openRouterProvider, err := openrouter.New(openrouter.WithAPIKey(openRouterKey))
					if err != nil {
						log.Printf("⚠️  Warning: Failed to create OpenRouter provider: %v", err)
					} else {
						// Use a capable model for tool calling
						remoteModel, err := openRouterProvider.LanguageModel(ctx, "anthropic/claude-3.5-haiku")
						if err != nil {
							log.Printf("⚠️  Warning: Failed to get OpenRouter language model: %v", err)
						} else {
							chainModels = append(chainModels, remoteModel)
							log.Printf("✅ OpenRouter added as primary model (anthropic/claude-3.5-haiku)")
						}
					}
				}

				// Add local llama as fallback
				chainModels = append(chainModels, localModel)

				// Create ChainLanguageModel instances for fallback support
				chatModel, _ = fantasy.NewChain(chainModels)
				titleModel, _ = fantasy.NewChain(chainModels)
				summaryModel, _ = fantasy.NewChain(chainModels)
				cleanupModel, _ = fantasy.NewChain(chainModels)
				log.Printf("✅ ChainLanguageModel instances created (%d models in chain)", len(chainModels))
			}
		}

		// Initialize Agent Chat Service (Yzma-based agent with RAG + DB persistence)
		// Phase 4: Now with Thread Management integration
		if kbService != nil {
			agentChatService := services.NewAgentChatService(
				app,
				dbService,
				libService,
				kbService,
				vectorSearchService,     // File-based RAG (direct attachments)
				threadManagementService, // Phase 4: Thread integration
			)

			// Inject ChainLanguageModel instances for decentralized LLM routing
			if chatModel != nil {
				agentChatService.SetChatModel(chatModel)
			}
			if titleModel != nil {
				agentChatService.SetTitleModel(titleModel)
			}
			if summaryModel != nil {
				agentChatService.SetSummaryModel(summaryModel)
			}

			app.RegisterService(application.NewService(agentChatService))
		}

		// Inject cleanup model to FileProcessorService
		if cleanupModel != nil {
			fileProcessorService.SetLanguageModel(cleanupModel)
			log.Printf("✅ FileProcessorService: LLM cleanup model injected")
		}

		// Add cleanup on shutdown
		app.OnShutdown(func() {
			log.Printf("🧹 Cleaning up Llama Library service...")
			libService.Cleanup()
		})
	}

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	win := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:             "Window 1",
		StartState:        application.WindowStateMaximised,
		EnableDragAndDrop: true, // Enable native drag & drop support
		Mac: application.MacWindow{
			Backdrop: application.MacBackdropTranslucent,
			TitleBar: application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

	// Setup drag & drop event handler
	win.OnWindowEvent(
		events.Common.WindowDropZoneFilesDropped,
		func(event *application.WindowEvent) {
			droppedFiles := event.Context().DroppedFiles()
			details := event.Context().DropZoneDetails()

			log.Printf("[Drag&Drop] Files dropped: %d files", len(droppedFiles))

			// Process files in backend - copy to local storage AND process for RAG
			processedFiles := make([]map[string]interface{}, 0, len(droppedFiles))
			for i, filePath := range droppedFiles {
				log.Printf("[Drag&Drop]   %d. %s", i+1, filePath)

				// Process file: copy to local storage + parse + RAG (all in one)
				// Use DEFAULT_LOBE_CHAT_USER to match the frontend's default user ID
				result, err := fileProcessorService.ProcessFileFromPath(filePath, "DEFAULT_LOBE_CHAT_USER")
				if err != nil {
					log.Printf("[Drag&Drop] Error processing file %s: %v", filePath, err)
					continue
				}

				// Create file info for frontend
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

			if details != nil {
				log.Printf("[Drag&Drop] Drop zone: %s at (%d, %d)", details.ElementID, details.X, details.Y)

				// Emit event to frontend with processed file info
				app.Event.Emit("files:dropped", map[string]interface{}{
					"files":      processedFiles,
					"elementId":  details.ElementID,
					"classList":  details.ClassList,
					"x":          details.X,
					"y":          details.Y,
					"attributes": details.Attributes,
				})
			} else {
				// Drop outside specific zone
				log.Printf("[Drag&Drop] Drop outside specific zone")
				app.Event.Emit("files:dropped", map[string]interface{}{
					"files": processedFiles,
				})
			}
		},
	)

	// Create a goroutine that emits an event containing the current time every second.
	// The frontend can listen to this event and update the UI accordingly.
	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.Event.Emit("time", now)
			time.Sleep(time.Second)
		}
	}()

	// Run the application. This blocks until the application has been exited.
	if err = app.Run(); err != nil {
		log.Fatal(err)
	}
}
