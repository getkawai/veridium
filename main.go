// Package main is the entry point for the Veridium desktop application.
package main

import (
	"embed"
	"log"
	"log/slog"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/kawai-network/veridium/internal/app"
	"github.com/kawai-network/veridium/internal/image"
	"github.com/kawai-network/veridium/internal/machineid"
	"github.com/kawai-network/veridium/internal/paths"
	"github.com/kawai-network/veridium/internal/services"
	"github.com/kawai-network/veridium/internal/tableviewer"
	"github.com/kawai-network/veridium/internal/topic"
	"github.com/kawai-network/veridium/pkg/fantasy/tools/builtin"
	"github.com/kawai-network/veridium/pkg/localfs"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/services/fileserver"
	"github.com/wailsapp/wails/v3/pkg/services/kvstore"
	wailslog "github.com/wailsapp/wails/v3/pkg/services/log"
	"github.com/wailsapp/wails/v3/pkg/services/notifications"
	"github.com/wailsapp/wails/v3/pkg/services/sqlite"
)

//go:embed all:frontend/dist
var assets embed.FS

// main is the application entry point
func main() {
	if app.DevMode {
		log.Printf("Development mode enabled")
	}

	// Initialize core services using internal/app (SINGLE SOURCE OF TRUTH)
	ctx := app.NewContext()
	defer ctx.Cleanup()

	if err := ctx.InitAll(); err != nil {
		sentry.CaptureException(err)
		log.Fatalf("Failed to initialize: %v", err)
	}

	fileProcessor := NewFileProcessorService(
		ctx.DB.DB(),
		ctx.FileLoader,
		ctx.VectorSearch,
		ctx.DuckDBStore,
		ctx.LibService,
		ctx.WhisperService,
		paths.FileBase(),
	)

	// Initialize Stable Diffusion in background
	sdEngine := image.NewEngine()
	sdEngine.InitializeInBackground()

	// Create Service wrapper (with DB)
	sdService := image.NewService(ctx.DB, sdEngine)

	// Create Wails app
	wailsApp := createWailsApp(ctx, fileProcessor, sdService)
	registerAgentServices(wailsApp, ctx, fileProcessor, sdService)
	createMainWindow(wailsApp, ctx, fileProcessor)

	if err := wailsApp.Run(); err != nil {
		// Use slog.Error which will be captured by SentryHandler
		slog.Error("Application crashed", "error", err)
		os.Exit(1)
	}
}

func createWailsApp(ctx *app.Context, fileProcessor *FileProcessorService, sdService *image.Service) *application.App {
	return application.New(application.Options{
		Name:        "veridium",
		Description: "Veridium AI Assistant",
		Logger:      slog.Default(),
		Services:    buildServiceList(ctx, fileProcessor, sdService),
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})
}

func buildServiceList(ctx *app.Context, fileProcessor *FileProcessorService, sdService *image.Service) []application.Service {
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
		application.NewService(fileProcessor),
		application.NewService(ctx.KBService),

		// File & Storage
		application.NewService(services.NewFileService(paths.FileBase())),
		application.NewService(localfs.NewService()),
		application.NewService(builtin.NewLocalSystemService()),

		// Utilities
		application.NewService(&machineid.Service{}),
		application.NewService(sdService), // Use the initialized sdService

		// Wails Native Services
		application.NewService(notifications.New()),
		application.NewService(wailslog.New()),
		application.NewService(sqlite.New()),
		application.NewService(kvstore.New()),
		application.NewService(ctx.WalletService),
		application.NewService(ctx.DeAIService),
		application.NewService(ctx.JarvisService),

		// Blockchain Services
		application.NewService(ctx.DepositSyncService),

		// Marketplace
		application.NewService(services.NewMarketplaceService(ctx.KVStore, ctx.BlockchainClient, ctx.WalletService)),

		// File Server
		application.NewServiceWithOptions(
			fileserver.NewWithConfig(&fileserver.Config{RootPath: paths.FileBase()}),
			application.ServiceOptions{Route: "/files"},
		),
	}
}

func registerAgentServices(wailsApp *application.App, ctx *app.Context, fileProcessor *FileProcessorService, sdService *image.Service) {
	ctx.AudioRecorder.SetApp(wailsApp)

	threadService := services.NewThreadManagementService(wailsApp, ctx.DB)
	wailsApp.RegisterService(application.NewService(threadService))

	// Initialize TopicService (for title generation)
	topicService := topic.NewService(ctx.DB, wailsApp)
	// Inject TopicService into StableDiffusion
	sdService.SetTopicService(topicService)
	// Register TopicService
	wailsApp.RegisterService(application.NewService(topicService))

	if ctx.LibService != nil && ctx.KBService != nil {
		agentService := services.NewAgentChatService(
			wailsApp, ctx.DB, ctx.LibService, ctx.KBService, ctx.VectorSearch, threadService, topicService, ctx.ToolRegistry, ctx.KVStore,
		)

		if ctx.ChatModel != nil {
			agentService.SetChatModel(ctx.ChatModel)
		}
		if ctx.TitleModel != nil {
			agentService.SetTitleModel(ctx.TitleModel) // This will also set it on topicService
		}
		if ctx.SummaryModel != nil {
			agentService.SetSummaryModel(ctx.SummaryModel)
		}

		// Register memory tool for recalling stored memories
		if ctx.MemoryIntegration != nil {
			if err := agentService.RegisterMemoryTool(ctx.MemoryIntegration); err != nil {
				log.Printf("⚠️  Failed to register memory tool: %v", err)
			}
		}

		wailsApp.RegisterService(application.NewService(agentService))
	}

	if ctx.CleanupModel != nil {
		fileProcessor.SetLanguageModel(ctx.CleanupModel)
		log.Printf("FileProcessor: LLM cleanup model injected")
	}

	if ctx.LibService != nil {
		wailsApp.OnShutdown(func() {
			log.Printf("Cleaning up Llama Library...")
			ctx.LibService.Cleanup()
		})
	}

	// Cleanup Stable Diffusion processes on shutdown
	wailsApp.OnShutdown(func() {
		sdService.Cleanup() // Cleanup is available via embedded Engine
	})
}

func createMainWindow(wailsApp *application.App, ctx *app.Context, fileProcessor *FileProcessorService) {
	win := wailsApp.Window.NewWithOptions(application.WebviewWindowOptions{
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

	win.OnWindowEvent(events.Common.WindowDropZoneFilesDropped, func(event *application.WindowEvent) {
		droppedFiles := event.Context().DroppedFiles()
		details := event.Context().DropZoneDetails()

		log.Printf("[Drag&Drop] %d files dropped", len(droppedFiles))

		processedFiles := make([]map[string]interface{}, 0, len(droppedFiles))
		for i, filePath := range droppedFiles {
			log.Printf("[Drag&Drop] %d. %s", i+1, filePath)

			result, err := fileProcessor.ProcessFileFromPath(filePath)
			if err != nil {
				log.Printf("[Drag&Drop] Error: %v", err)
				continue
			}

			processedFiles = append(processedFiles, map[string]interface{}{
				"originalPath": filePath,
				"url":          result.RelativeURL,
				"name":         result.Filename,
				"fileId":       result.FileID,
				"documentId":   result.DocumentID,
				"processing":   result.Processing,
			})
		}

		eventData := map[string]interface{}{"files": processedFiles}
		if details != nil {
			eventData["elementId"] = details.ElementID
			eventData["classList"] = details.ClassList
			eventData["x"] = details.X
			eventData["y"] = details.Y
			eventData["attributes"] = details.Attributes
		}

		wailsApp.Event.Emit("files:dropped", eventData)
	})
}
