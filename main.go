package main

import (
	"embed"
	_ "embed"
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
	"github.com/wailsapp/wails/v3/pkg/application"
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

	// Initialize Llama service (LLM inference using llama.cpp)
	// Auto-installs llama.cpp in background
	llamaService, err := llama.NewService()
	if err != nil {
		log.Printf("⚠️  Warning: Failed to initialize Llama service: %v", err)
		log.Printf("    LLM features will not be available.")
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
	}

	// Initialize File Processor service (file parsing + document storage + RAG)
	loadFileService := &LoadFileService{}
	fileProcessorService := NewFileProcessorService(
		dbService.DB(),
		loadFileService,
		vectorSearchService.GetChromemDB(),
	)
	log.Printf("✅ File Processor service initialized")
	log.Printf("   Handles: file parsing → document storage → RAG processing")

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app := application.New(application.Options{
		Name:        "veridium",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(&GreetService{}),
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
			// Llama service - for LLM inference using llama.cpp
			application.NewService(llamaService),
			// Audio recorder service - for native microphone recording
			application.NewService(audioRecorderService),
			// Vector search service - for semantic search using chromem
			application.NewService(vectorSearchService),
			// File processor service - for file parsing + document storage + RAG
			application.NewService(fileProcessorService),
			// Machine ID service
			application.NewService(&machineid.Service{}),
			// Temp file service
			application.NewService(&TempFileService{}),
			// Node.js equivalent services
			application.NewService(&NodeFsService{}),
			application.NewService(&NodePathService{}),
			application.NewService(&NodeBufferService{}),
			application.NewService(&NodeExecService{}),
			application.NewService(&NodeOsService{}),
			application.NewService(&ZipService{}),
			application.NewService(&LoadFileService{}),
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
					// Get user data directory
					userConfigDir, err := os.UserConfigDir()
					if err != nil {
						// Fallback to current directory
						userConfigDir = "."
					}
					appDataDir := filepath.Join(userConfigDir, "veridium")

					return fileserver.NewWithConfig(&fileserver.Config{
						RootPath:               appDataDir,
						EnableDirectoryListing: true, // Enable for user data access
						EnableCORS:             true, // Enable CORS for web access
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

	// Initialize Llama proxy service (for WebView to call llama-server via streaming)
	llamaProxyService := llama.NewProxyService(llamaService, app)
	app.RegisterService(application.NewService(llamaProxyService))

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:      "Window 1",
		StartState: application.WindowStateMaximised,
		Mac: application.MacWindow{
			Backdrop: application.MacBackdropTranslucent,
			TitleBar: application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

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
