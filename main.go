package main

import (
	"context"
	"embed"
	"log"
	"log/slog"
	"os"

	"github.com/getkawai/database"
	"github.com/getkawai/tools"
	"github.com/getkawai/tools/builtin"
	"github.com/getkawai/tools/localfs"
	unillm "github.com/getkawai/unillm"
	"github.com/kawai-network/x/store"
	"github.com/kawai-network/y/config"
	"github.com/kawai-network/y/machineid"
	"github.com/kawai-network/y/paths"
	"github.com/kawai-network/veridium/internal/app"
	"github.com/kawai-network/veridium/internal/image"
	"github.com/kawai-network/veridium/internal/lifecycle"
	"github.com/kawai-network/veridium/internal/services"
	"github.com/kawai-network/veridium/internal/tableviewer"
	"github.com/kawai-network/veridium/internal/topic"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/services/fileserver"
	"github.com/wailsapp/wails/v3/pkg/services/kvstore"
	wailslog "github.com/wailsapp/wails/v3/pkg/services/log"
	"github.com/wailsapp/wails/v3/pkg/services/notifications"
	"github.com/wailsapp/wails/v3/pkg/services/sqlite"
	"go.uber.org/fx"
)

// Version is set at build time
var Version = "0.1.0"

//go:embed all:frontend/dist
var assets embed.FS

// main is the application entry point
func main() {
	// Force local data directory in development mode.
	if os.Getenv("VERIDIUM_DEV") == "1" {
		paths.SetDataDir("data")
	}

	// Initialize environment configuration
	if err := config.Initialize(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// Validate configuration for production
	if err := config.ValidateForProduction(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// Log environment info
	log.Printf("Environment: %s", config.GetEnvironment())
	log.Printf("Network: %s (Chain ID: %d)", config.GetNetworkName(), config.GetChainID())

	// Use fx.App to manage dependencies
	fx.New(
		app.Module,
		fx.Provide(
			NewFileProcessorServiceWrapper,
			NewImageServiceWrapper,
			NewWailsApp,
			NewThreadManagementServiceWrapper,
			NewTopicServiceWrapper,
			NewAgentChatServiceWrapper,
			NewLifecycleManagerWrapper,
		),
		fx.Invoke(RegisterWailsServices),
		fx.Invoke(SetupMainWindow),
	).Run()
}

// FileProcessorServiceWrapper is an fx provider for FileProcessorService
type FileProcessorServiceWrapper struct {
	fx.In
	DB           *database.Service
	FileLoader   *services.FileLoader
	VectorSearch *services.VectorSearchService
	DuckDBStore  *services.DuckDBStore
	CleanupModel unillm.LanguageModel `name:"cleanupModel"`
}

func NewFileProcessorServiceWrapper(p FileProcessorServiceWrapper) *FileProcessorService {
	fileProcessor := NewFileProcessorService(
		p.DB.DB(),
		p.FileLoader,
		p.VectorSearch,
		p.DuckDBStore,
		paths.FileBase(),
	)
	if p.CleanupModel != nil {
		fileProcessor.SetLanguageModel(p.CleanupModel)
		log.Printf("FileProcessor: LLM cleanup model injected")
	}
	return fileProcessor
}

// NewImageServiceWrapper is an fx provider for image.Service
type ImageServiceWrapper struct {
	fx.In
	DB *database.Service
}

func NewImageServiceWrapper(p ImageServiceWrapper) *image.Service {
	return image.NewService(p.DB)
}

// NewWailsApp creates the Wails application instance
func NewWailsApp(lc fx.Lifecycle, ctx *app.Context, fileProcessor *FileProcessorService, sdService *image.Service) *application.App {
	wailsApp := application.New(application.Options{
		Name:        "Kawai",
		Description: "Kawai Network - AI-Powered Blockchain Platform",
		Logger:      slog.Default(),
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	lc.Append(fx.Hook{
		OnStart: func(c context.Context) error {
			log.Printf("Wails App Starting...")
			go func() { // Run Wails app in a goroutine
				if err := wailsApp.Run(); err != nil {
					slog.Error("Application crashed", "error", err)
					os.Exit(1)
				}
			}()
			return nil
		},
		OnStop: func(c context.Context) error {
			log.Printf("Wails App Stopping...")
			wailsApp.Quit()
			return nil
		},
	})

	return wailsApp
}

// RegisterWailsServices defines how services are registered to the Wails app.
// This is an fx.Invoke function as it has side effects (registering services).
func RegisterWailsServices(wailsApp *application.App, ctx *app.Context, fileProcessor *FileProcessorService, sdService *image.Service, threadService *services.ThreadManagementService, topicService *topic.TopicService, agentService *services.AgentChatService, lifecycleManager *lifecycle.Manager) {
	// Register services from buildServiceList
	wailsApp.RegisterService(application.NewService(ctx.Queries))
	wailsApp.RegisterService(application.NewService(ctx.DB))
	wailsApp.RegisterService(application.NewService(tableviewer.NewService(ctx.DB.DB())))

	wailsApp.RegisterService(application.NewService(ctx.SearchService))
	wailsApp.RegisterService(application.NewService(ctx.TTSService))
	wailsApp.RegisterService(application.NewService(ctx.AudioRecorder))

	if ctx.VectorSearch != nil {
		wailsApp.RegisterService(application.NewService(ctx.VectorSearch))
	}
	wailsApp.RegisterService(application.NewService(fileProcessor))
	if ctx.KBService != nil {
		wailsApp.RegisterService(application.NewService(ctx.KBService))
	}

	wailsApp.RegisterService(application.NewService(services.NewFileService(paths.FileBase())))
	wailsApp.RegisterService(application.NewService(localfs.NewService()))
	wailsApp.RegisterService(application.NewService(builtin.NewLocalSystemService()))

	wailsApp.RegisterService(application.NewService(&machineid.Service{}))
	wailsApp.RegisterService(application.NewService(sdService))
	wailsApp.RegisterService(application.NewService(notifications.New()))
	wailsApp.RegisterService(application.NewService(wailslog.New()))
	wailsApp.RegisterService(application.NewService(sqlite.New()))
	wailsApp.RegisterService(application.NewService(kvstore.New()))
	wailsApp.RegisterService(application.NewService(ctx.WalletService))
	wailsApp.RegisterService(application.NewService(ctx.DeAIService))
	wailsApp.RegisterService(application.NewService(ctx.JarvisService))

	wailsApp.RegisterService(application.NewService(ctx.DepositSyncService))
	wailsApp.RegisterService(application.NewService(services.NewMarketplaceService(ctx.KVStore, ctx.BlockchainClient, ctx.WalletService)))
	wailsApp.RegisterService(application.NewService(services.NewReferralService(ctx.KVStore)))
	wailsApp.RegisterService(application.NewService(services.NewCashbackService(ctx.KVStore)))
	wailsApp.RegisterService(application.NewService(&services.ConfigService{}))
	wailsApp.RegisterService(application.NewServiceWithOptions(
		fileserver.NewWithConfig(&fileserver.Config{RootPath: paths.FileBase()}),
		application.ServiceOptions{Route: "/files"},
	))

	// Agent services - AudioRecorder requires deferred App injection to break cycle
	ctx.AudioRecorder.SetApp(wailsApp)
	wailsApp.RegisterService(application.NewService(threadService))

	sdService.SetTopicService(topicService)
	wailsApp.RegisterService(application.NewService(topicService))

	if ctx.KBService != nil {
		wailsApp.RegisterService(application.NewService(agentService))
	}

	// No need for lifecycleManager.RegisterCleanup directly in main.go now,
	// as fx.Lifecycle hooks are handled in provider functions.

	// Register the centralized shutdown handler for Wails
	wailsApp.OnShutdown(lifecycleManager.Shutdown)
}

// NewThreadManagementServiceWrapper is an fx provider for services.ThreadManagementService
type ThreadManagementServiceWrapperParams struct {
	fx.In
	WailsApp *application.App
	DB       *database.Service
}

func NewThreadManagementServiceWrapper(p ThreadManagementServiceWrapperParams) *services.ThreadManagementService {
	return services.NewThreadManagementService(p.WailsApp, p.DB)
}

// NewTopicServiceWrapper is an fx provider for topic.Service
type TopicServiceWrapperParams struct {
	fx.In
	DB       *database.Service
	WailsApp *application.App
}

func NewTopicServiceWrapper(p TopicServiceWrapperParams) *topic.TopicService {
	return topic.NewService(p.DB, p.WailsApp)
}

// NewAgentChatServiceWrapper is an fx provider for services.AgentChatService
type AgentChatServiceWrapperParams struct {
	fx.In
	WailsApp          *application.App
	DB                *database.Service
	KBService         *services.KnowledgeBaseService
	VectorSearch      *services.VectorSearchService
	ThreadService     *services.ThreadManagementService
	TopicService      *topic.TopicService
	ToolRegistry      *tools.ToolRegistry
	KVStore           *store.KVStore
	ChatModel         unillm.LanguageModel `name:"chatModel"`
	TitleModel        unillm.LanguageModel `name:"titleModel"`
	SummaryModel      unillm.LanguageModel `name:"summaryModel"`
	MemoryIntegration *services.MemoryIntegration
}

func NewAgentChatServiceWrapper(p AgentChatServiceWrapperParams) *services.AgentChatService {
	agentService := services.NewAgentChatService(
		p.WailsApp, p.DB, p.KBService, p.VectorSearch, p.ThreadService, p.TopicService, p.ToolRegistry, p.KVStore,
	)
	if p.ChatModel != nil {
		agentService.SetChatModel(p.ChatModel)
	}
	if p.TitleModel != nil {
		agentService.SetTitleModel(p.TitleModel)
	}
	if p.SummaryModel != nil {
		agentService.SetSummaryModel(p.SummaryModel)
	}

	if p.MemoryIntegration != nil {
		if err := agentService.RegisterMemoryTool(p.MemoryIntegration); err != nil {
			log.Printf("⚠️  Failed to register memory tool: %v", err)
		}
	}
	return agentService
}

// NewLifecycleManagerWrapper is an fx provider for lifecycle.Manager
func NewLifecycleManagerWrapper() *lifecycle.Manager {
	return lifecycle.NewManager()
}

// SetupMainWindow is an fx.Invoke function to configure the main window
func SetupMainWindow(wailsApp *application.App, ctx *app.Context, fileProcessor *FileProcessorService) {
	win := wailsApp.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:          "Kawai",
		StartState:     application.WindowStateMaximised,
		EnableFileDrop: true,
		Mac: application.MacWindow{
			Backdrop: application.MacBackdropTranslucent,
			TitleBar: application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

	win.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
		droppedFiles := event.Context().DroppedFiles()
		details := event.Context().DropTargetDetails()

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
