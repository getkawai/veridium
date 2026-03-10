package app

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/kawai-network/y/logger"

	"github.com/getkawai/database"
	db "github.com/getkawai/database/db"
	"github.com/getkawai/tools"
	"github.com/getkawai/tools/search"
	unillm "github.com/getkawai/unillm"
	llamaembed "github.com/getkawai/unillm/providers/llama-embed"
	"github.com/kawai-network/veridium/internal/audio_recorder"
	"github.com/kawai-network/veridium/internal/services"
	"github.com/kawai-network/veridium/internal/services/cache"
	"github.com/kawai-network/veridium/internal/tts"
	"github.com/kawai-network/x/blockchain"
	"github.com/kawai-network/x/store"
	"go.uber.org/fx"
)

// Configuration constants
const (
	EmbeddingDims = 384
	DefaultUserID = "DEFAULT_LOBE_CHAT_USER"
)

// Context holds all application services and state.
type Context struct {
	// Core Services
	DB      *database.Service
	Queries *db.Queries

	// AI/ML Services
	Embedder     llamaembed.Embedder
	CacheManager *cache.CacheManager

	// Storage Services
	DuckDBStore *services.DuckDBStore
	FileLoader  *services.FileLoader
	KVStore     *store.KVStore

	// Feature Services
	SearchService *search.Service
	TTSService    *tts.TTSService
	AudioRecorder *audio_recorder.AudioRecorderService
	VectorSearch  *services.VectorSearchService
	KBService     *services.KnowledgeBaseService
	RAGProcessor  *services.RAGProcessor
	ToolRegistry  *tools.ToolRegistry
	WalletService *services.WalletService
	DeAIService   *services.DeAIService
	JarvisService *services.JarvisService

	// Blockchain Services
	BlockchainClient   *blockchain.Client
	DepositSyncService *services.DepositSyncService

	// Language Models
	ChatModel    unillm.LanguageModel `json:"-"`
	TitleModel   unillm.LanguageModel `json:"-"`	// Renamed from SummaryModel to TitleModel
	SummaryModel unillm.LanguageModel `json:"-"`
	CleanupModel unillm.LanguageModel `json:"-"`

	// Memory Services
	MemoryIntegration *services.MemoryIntegration
}

// Params holds the dependencies for NewContext
type Params struct {
	fx.In

	DB                 *database.Service
	Queries            *db.Queries
	KVStore            *store.KVStore
	SearchService      *search.Service
	TTSService         *tts.TTSService
	AudioRecorder      *audio_recorder.AudioRecorderService
	DuckDBStore        *services.DuckDBStore
	FileLoader         *services.FileLoader
	VectorSearch       *services.VectorSearchService
	KBService          *services.KnowledgeBaseService
	ToolRegistry       *tools.ToolRegistry
	WalletService      *services.WalletService
	BlockchainClient   *blockchain.Client
	DepositSyncService *services.DepositSyncService
	DeAIService        *services.DeAIService
	JarvisService      *services.JarvisService
	MemoryIntegration  *services.MemoryIntegration

	// Language Models (will be provided by separate fx.Module later)
	ChatModel    unillm.LanguageModel `name:"chatModel"`
	TitleModel   unillm.LanguageModel `name:"titleModel"` // Use specific names for clarity
	SummaryModel unillm.LanguageModel `name:"summaryModel"`
	CleanupModel unillm.LanguageModel `name:"cleanupModel"`

	// Temporarily not provided by fx
	Embedder     llamaembed.Embedder
	CacheManager *cache.CacheManager
	RAGProcessor *services.RAGProcessor
}

// NewContext creates a new application context with all services initialized via dependency injection.
func NewContext(p Params) *Context {
	return &Context{
		DB:                 p.DB,
		Queries:            p.Queries,
		KVStore:            p.KVStore,
		SearchService:      p.SearchService,
		TTSService:         p.TTSService,
		AudioRecorder:      p.AudioRecorder,
		DuckDBStore:        p.DuckDBStore,
		FileLoader:         p.FileLoader,
		VectorSearch:       p.VectorSearch,
		KBService:          p.KBService,
		ToolRegistry:       p.ToolRegistry,
		WalletService:      p.WalletService,
		BlockchainClient:   p.BlockchainClient,
		DepositSyncService: p.DepositSyncService,
		DeAIService:        p.DeAIService,
		JarvisService:      p.JarvisService,
		MemoryIntegration:  p.MemoryIntegration,
		ChatModel:          p.ChatModel,
		TitleModel:         p.TitleModel,
		SummaryModel:       p.SummaryModel,
		CleanupModel:       p.CleanupModel,
		RAGProcessor:       p.RAGProcessor, // Provided by ProvideRAGProcessor
		Embedder:           p.Embedder,     // Provided by ProvideEmbedder
		CacheManager:       p.CacheManager, // Provided by ProvideCacheManager
	}
}

// ProvideSentry initializes Sentry and configures slog to use SentryHandler.
func ProvideSentry(lc fx.Lifecycle) {
	var dsns = []string{
		"https://18511014596d6da4288edc0e714a8c04@o4510629097504768.ingest.us.sentry.io/4510629100650496",
		"https://b66f862d7567c075a44c697757bb8130@o4510618985758720.ingest.us.sentry.io/4510618990804992",
	}

	randDsn := dsns[time.Now().UnixNano()%int64(len(dsns))]

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              randDsn,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		EnableLogs:       true, // Enable Logger API as requested
		BeforeSendLog: func(log *sentry.Log) *sentry.Log {
			// filter all logs below warning
			if log.Severity <= sentry.LogSeverityWarning {
				return nil
			}
			return log
		},
	})
	if err != nil {
		slog.Warn("Sentry initialization failed", "error", err)
		return
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			sentry.Flush(2 * time.Second)
			return nil
		},
	})

	// properties of the logger
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
	})
	// wrap the handler with SentryHandler
	sentryHandler := logger.NewSentryHandler(handler)
	// create a new logger with the SentryHandler
	defaultLogger := slog.New(sentryHandler)
	// set the default logger to the new logger
	slog.SetDefault(defaultLogger)

	slog.Info("Sentry initialized with EnableLogs: true (using slog handler)")
}
