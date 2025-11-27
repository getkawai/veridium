package main

import (
	"embed"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/kawai-network/veridium/internal/audio_recorder"
	"github.com/kawai-network/veridium/internal/database"
	"github.com/kawai-network/veridium/internal/llama"
	"github.com/kawai-network/veridium/internal/machineid"
	"github.com/kawai-network/veridium/internal/search"
	"github.com/kawai-network/veridium/internal/services"
	"github.com/kawai-network/veridium/internal/tableviewer"
	"github.com/kawai-network/veridium/internal/tts"
	"github.com/kawai-network/veridium/internal/whisper"
	"github.com/kawai-network/veridium/pkg/chromem"
	"github.com/kawai-network/veridium/pkg/contextengine"
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

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {
	// Initialize database service
	dbService, err := database.NewService()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbService.Close()

	// Get queries - this is what we'll bind to Wails
	queries := dbService.Queries()

	// Create TableViewer service
	tableViewerService := tableviewer.NewService(dbService.DB())

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

	// Initialize Vector Search service (chromem for semantic search)
	// Get user data directory for vector DB persistence
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		userConfigDir = "."
	}
	vectorDBPath := filepath.Join(userConfigDir, "veridium", "vector-db")
	vectorSearchService, err := services.NewVectorSearchService(vectorDBPath, "llama", "http://localhost:8080")
	if err != nil {
		log.Printf("⚠️  Warning: Failed to initialize Vector Search service: %v", err)
		log.Printf("    Semantic search features will use fallback mode.")
	} else {
		log.Printf("✅ Vector Search service initialized (chromem)")
		log.Printf("   Database path: %s", vectorDBPath)
		log.Printf("   Embedding provider: llama.cpp (llama-server)")
		log.Printf("   Embedding endpoint: http://localhost:8080")
		log.Printf("   Note: Embedding models auto-download in background")
		log.Printf("   Note: Use llamaService.StartEmbeddingServer(8080) to start embedding server")
	}

	// Initialize File Service base directory (needed by both FileService and FileProcessor)
	fileBaseDir := filepath.Join(userConfigDir, "veridium", "files")
	os.MkdirAll(fileBaseDir, 0755)

	// Initialize File Processor service (file parsing + document storage + RAG)
	fileLoader := services.NewFileLoader()
	fileProcessorService := NewFileProcessorService(
		dbService.DB(),
		fileLoader,
		vectorSearchService, // Pass entire service to get chromemDB and embedFunc
		fileBaseDir,         // Pass file base directory for path resolution
	)
	log.Printf("✅ File Processor service initialized")
	log.Printf("   Handles: file parsing → document storage → RAG processing")

	// Initialize File Service (local storage for desktop)

	fileStorage := services.NewLocalFileStorage(fileBaseDir)
	fileSvc := services.NewFileService("system", fileStorage)
	log.Printf("✅ File Service initialized")

	// Initialize Llama.cpp Library Service (library-based LLM inference)
	// Auto-installs llama.cpp binaries and downloads models in background
	libService, err := llama.NewLibraryService()
	if err != nil {
		log.Printf("⚠️  Warning: Failed to initialize Llama Library service: %v", err)
		log.Printf("    LLM chat features will not be available.")
	} else {
		log.Printf("✅ Llama Library service initialized")
		log.Printf("   Models directory: %s", libService.GetModelsDirectory())
		log.Printf("   Auto-setup running in background...")
	}

	// Initialize Knowledge Base Service (RAG with Chromem + Eino)
	var kbService *services.KnowledgeBaseService
	if libService != nil {
		// Initialize embedding model for KB
		embeddingModelPath := filepath.Join(libService.GetModelsDirectory(),
			"granite-embedding-107m-multilingual-Q6_K_L.gguf")
		embedFunc := chromem.NewEmbeddingFuncLlamaWithPreloadedLibrary(embeddingModelPath)

		kbPath := filepath.Join(userConfigDir, "veridium", "knowledge-bases")
		kbAssetPath := filepath.Join(userConfigDir, "veridium", "kb-assets")

		kbService, err = services.NewKnowledgeBaseService(dbService, &services.KnowledgeBaseConfig{
			ChromemPath:   kbPath,
			EmbeddingFunc: embedFunc,
			AssetDir:      kbAssetPath,
		})
		if err != nil {
			log.Printf("⚠️  Warning: Failed to initialize Knowledge Base service: %v", err)
		} else {
			log.Printf("✅ Knowledge Base service initialized")
			log.Printf("   Vector DB path: %s", kbPath)
			log.Printf("   Asset path: %s", kbAssetPath)
			log.Printf("   Embedding model: granite-embedding-107m-multilingual")
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
			application.NewService(tableViewerService),
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
			// Context Engine service - for message context engineering
			application.NewService(contextengine.NewContextEngineService()),
			// Tools Engine service - for tool/plugin management
			application.NewService(NewToolsEngineService()),
			// Knowledge Base service - for RAG with Chromem + Eino
			application.NewService(kbService),
			// Machine ID service
			application.NewService(&machineid.Service{}),
			// File service
			application.NewService(fileSvc),
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
						RootPath:               fileBaseDir, // Same as FileService baseDir
						EnableDirectoryListing: true,        // Enable for user data access
						EnableCORS:             true,        // Enable CORS for web access
						IndexFile:              "index.html",
						AllowedExtensions:      []string{".html", ".css", ".js", ".png", ".jpg", ".jpeg", ".gif", ".svg", ".pdf", ".docx", ".txt", ".md"}, // User files + images
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
	log.Printf("✅ Thread Management service registered")
	log.Printf("   Supports: conversation branching, thread creation, thread switching")

	// Initialize Llama Chat Service (OpenAI-compatible chat API)
	// This service provides chat completion functionality using the library service
	if libService != nil {
		llamaChatService := llama.NewLibraryChatService(libService, app)
		app.RegisterService(application.NewService(llamaChatService))
		log.Printf("✅ Llama Chat service registered")
		log.Printf("   OpenAI-compatible API: ChatCompletion, ChatCompletionStream")
		log.Printf("   Supports: temperature, top_p, top_k, max_tokens")

		// Phase 3: Initialize integration bridges
		// These bridge existing engines to Eino agent system
		var toolsBridge *services.ToolsEngineBridge
		var contextBridge *services.ContextEngineBridge

		// Initialize tools engine bridge (from existing ToolsEngineService)
		toolsEngineService := NewToolsEngineService()
		if toolsEngineService.engine != nil {
			toolsBridge = services.NewToolsEngineBridge(toolsEngineService.engine)
			log.Printf("🔧 Tools Engine Bridge initialized")
		}

		// Initialize context engine bridge (from existing context engine)
		contextEngine := contextengine.New(contextengine.Config{
			SystemRole:         "",
			EnableHistoryCount: true,
			HistoryCount:       20, // Keep last 20 messages
		})
		contextBridge = services.NewContextEngineBridge(contextEngine)
		log.Printf("🔄 Context Engine Bridge initialized")

		// Initialize Agent Chat Service (Eino-based agent with RAG + DB persistence + Bridges)
		// Phase 4: Now with Thread Management integration
		if kbService != nil {
			agentChatService := services.NewAgentChatService(
				app,
				dbService,
				libService,
				kbService,
				toolsBridge,             // Phase 3: Tools integration
				contextBridge,           // Phase 3: Context processing
				threadManagementService, // Phase 4: Thread integration
			)
			app.RegisterService(application.NewService(agentChatService))
			log.Printf("✅ Agent Chat service registered")
			log.Printf("   Eino-based agent with RAG capabilities")
			log.Printf("   Supports: tool calling, knowledge base search, multi-turn conversations")
			log.Printf("   Session persistence: SQLite (messages + metadata)")
			log.Printf("   Phase 3: Tools Engine Bridge integrated")
			log.Printf("   Phase 3: Context Engine Bridge integrated")
			log.Printf("   Phase 4: Thread Management integrated")
			log.Printf("   Phase 4: Auto Topic & Thread support")
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

			// Process files in backend - copy to local storage
			processedFiles := make([]map[string]interface{}, 0, len(droppedFiles))
			for i, filePath := range droppedFiles {
				log.Printf("[Drag&Drop]   %d. %s", i+1, filePath)

				// Copy file to local storage
				relativeKey, err := fileSvc.CopyFileFromAbsolutePath(filePath)
				if err != nil {
					log.Printf("[Drag&Drop] Error copying file %s: %v", filePath, err)
					continue
				}

				// Create file info for frontend
				fileInfo := map[string]interface{}{
					"originalPath": filePath,
					"savedKey":     relativeKey,
					"url":          "/files/" + relativeKey,
					"name":         filepath.Base(filePath),
				}
				processedFiles = append(processedFiles, fileInfo)
				log.Printf("[Drag&Drop] File saved: %s -> %s", filepath.Base(filePath), relativeKey)
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
