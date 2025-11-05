package main

import (
	"embed"
	_ "embed"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/kawai-network/veridium/internal/database"
	"github.com/kawai-network/veridium/internal/services/search"
	"github.com/kawai-network/veridium/internal/services/tableviewer"
	"github.com/kawai-network/veridium/services"
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
	ttsService, err := services.NewTTSService()
	if err != nil {
		log.Printf("⚠️  Warning: Failed to initialize TTS service: %v", err)
	} else {
		log.Printf("✅ TTS service initialized successfully")
		log.Printf("   Platform: %s", ttsService.GetPlatformInfo()["platform"])
	}

	// Initialize Whisper STT service (offline, cross-platform)
	whisperService, err := services.NewWhisperService()
	if err != nil {
		log.Printf("⚠️  Warning: Failed to initialize Whisper service: %v", err)
	} else {
		defer whisperService.Close()
		log.Printf("✅ Whisper service initialized successfully")
		log.Printf("   Models directory: %s", whisperService.GetModelsDirectory())
	}

	// Initialize Hybrid STT service (Native + Whisper fallback)
	hybridSTTService, err := services.NewHybridSTTService("en-US")
	if err != nil {
		log.Printf("⚠️  Warning: Failed to initialize Hybrid STT service: %v", err)
		log.Printf("    Speech-to-text features will not be available.")
	} else {
		defer hybridSTTService.Close()
		log.Printf("✅ Hybrid STT service initialized successfully")
		engines := hybridSTTService.GetAvailableEngines()
		log.Printf("   Available engines: %v", engines)
		log.Printf("   Current engine: %s", hybridSTTService.GetCurrentEngine())
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
			// Whisper service - for offline STT
			application.NewService(whisperService),
			// Hybrid STT service - for speech-to-text (Native + Whisper)
			application.NewService(hybridSTTService),
			// Machine ID service
			application.NewService(&MachineIDService{}),
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
