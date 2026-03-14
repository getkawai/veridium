// Package app provides the core application context and service initialization.
// This is the SINGLE SOURCE OF TRUTH for service initialization.
// main.go imports and uses this, tests import and use this.
package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/getkawai/database"
	db "github.com/getkawai/database/db"
	"github.com/getkawai/tools"
	yzmabuiltin "github.com/getkawai/tools/builtin"
	"github.com/getkawai/tools/search"
	unillm "github.com/getkawai/unillm"
	llamaembed "github.com/getkawai/unillm/providers/llama-embed"
	"github.com/getkawai/unillm/providers/openrouter"
	"github.com/getsentry/sentry-go"
	"github.com/kawai-network/contracts"
	"github.com/kawai-network/veridium/internal/audio_recorder"
	"github.com/kawai-network/veridium/internal/services"
	"github.com/kawai-network/veridium/internal/services/cache"
	"github.com/kawai-network/veridium/internal/tts"
	"github.com/kawai-network/x/blockchain"
	"github.com/kawai-network/x/llm"
	"github.com/kawai-network/x/store"
	"github.com/kawai-network/y/logger"
	"github.com/kawai-network/y/paths"
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
	TitleModel   unillm.LanguageModel `json:"-"`
	SummaryModel unillm.LanguageModel `json:"-"`
	CleanupModel unillm.LanguageModel `json:"-"`

	// Memory Services
	MemoryIntegration *services.MemoryIntegration

	cleanupFuncs []func()
}

func NewContext() *Context {
	return &Context{
		cleanupFuncs: make([]func(), 0),
	}
}

func (ctx *Context) AddCleanup(fn func()) {
	ctx.cleanupFuncs = append(ctx.cleanupFuncs, fn)
}

func (ctx *Context) Cleanup() {
	for i := len(ctx.cleanupFuncs) - 1; i >= 0; i-- {
		ctx.cleanupFuncs[i]()
	}
}

func (ctx *Context) InitDatabase() error {
	dbService, err := database.NewService()
	if err != nil {
		return err
	}
	ctx.DB = dbService
	ctx.Queries = dbService.Queries()
	ctx.AddCleanup(func() { _ = dbService.Close() })
	return nil
}

func (ctx *Context) InitBasicServices() {
	ctx.SearchService = search.NewService()

	if ttsService, err := tts.NewTTSService(); err != nil {
		log.Printf("Warning: TTS init failed: %v", err)
	} else {
		ctx.TTSService = ttsService
		log.Printf("TTS initialized")
	}

	ctx.AudioRecorder = audio_recorder.NewAudioRecorderService(nil)

	// Initialize Tool Registry
	ctx.ToolRegistry = tools.NewToolRegistry()
	if err := yzmabuiltin.RegisterAllWithDB(ctx.ToolRegistry, ctx.DB.DB()); err != nil {
		log.Printf("Warning: Failed to register builtin tools: %v", err)
	} else {
		log.Printf("Tool Registry initialized with builtin tools")
	}
}

var dsns = []string{
	"https://18511014596d6da4288edc0e714a8c04@o4510629097504768.ingest.us.sentry.io/4510629100650496",
	"https://b66f862d7567c075a44c697757bb8130@o4510618985758720.ingest.us.sentry.io/4510618990804992",
}

func getRandomDsn() string {
	return dsns[time.Now().UnixNano()%int64(len(dsns))]
}

func (ctx *Context) InitSentry() {
	// Production-only build: Sentry always enabled

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              getRandomDsn(),
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
		log.Printf("Sentry initialization failed: %v\n", err)
		return
	}

	// Flush buffered events before the program terminates.
	// Set the timeout to the maximum duration the program can afford to wait.
	ctx.AddCleanup(func() {
		sentry.Flush(2 * time.Second)
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

	log.Printf("Sentry initialized with EnableLogs: true (using slog handler)")
}

func (ctx *Context) InitLlamaService() {
	// TODO: Re-implement local llama service initialization after Context.LibService removal.
}

func (ctx *Context) InitVectorStore() {
	_ = os.MkdirAll(paths.FileBase(), 0755)
	duckDBStore, err := services.NewDuckDBStore(paths.DuckDB(), EmbeddingDims)
	if err != nil {
		log.Printf("Warning: DuckDB init failed: %v", err)
		return
	}
	ctx.DuckDBStore = duckDBStore
	ctx.AddCleanup(func() { _ = duckDBStore.Close() })
	log.Printf("DuckDB initialized for RAG/vector retrieval (path: %s)", paths.DuckDB())
}

func (ctx *Context) InitEmbedder() {
	// TODO: Re-implement embedder initialization after Context.LibService removal.
}

func (ctx *Context) InitVectorSearch() {
	if ctx.Embedder == nil {
		log.Printf("Warning: Embedder not available, VectorSearch disabled")
		return
	}
	vectorSearch, err := services.NewVectorSearchService(ctx.DB.DB(), ctx.DuckDBStore, ctx.Embedder)
	if err != nil {
		log.Printf("Warning: VectorSearch init failed: %v", err)
		return
	}
	ctx.VectorSearch = vectorSearch
	log.Printf("VectorSearch initialized")
}

func (ctx *Context) InitFileLoader() {
	ctx.FileLoader = services.NewFileLoader()
}

func (ctx *Context) InitKnowledgeBase() {
	if ctx.VectorSearch == nil || ctx.FileLoader == nil {
		return
	}
	embedder := ctx.VectorSearch.GetEmbedder()
	ragProcessor := services.NewRAGProcessor(ctx.DB.DB(), ctx.DuckDBStore, ctx.FileLoader, embedder)
	ctx.RAGProcessor = ragProcessor

	kbService, err := services.NewKnowledgeBaseService(ctx.DB, &services.KnowledgeBaseConfig{
		RAGProcessor: ragProcessor,
		VectorSearch: ctx.VectorSearch,
		FileLoader:   ctx.FileLoader,
		AssetDir:     paths.KBAssets(),
	})
	if err != nil {
		log.Printf("Warning: KnowledgeBase init failed: %v", err)
		return
	}
	ctx.KBService = kbService
	log.Printf("KnowledgeBase initialized")
}

func (ctx *Context) InitWalletService() {
	// Dependent on KVStore for API Key generation
	if ctx.KVStore == nil {
		// Try init KVStore first if nil?
		// Actually InitAll order matters.
		// For now assuming KVStore is ready or we pass nil and handle it?
		// Better to reorder InitAll.
	}
	ctx.WalletService = services.NewWalletService(paths.FileBase(), ctx.KVStore)
	log.Printf("WalletService initialized")
}

func (ctx *Context) InitDeAIService() {
	if ctx.WalletService == nil {
		log.Printf("Warning: WalletService not initialized, cannot init DeAIService")
		return
	}
	ctx.DeAIService = services.NewDeAIService(ctx.WalletService, ctx.KVStore)
	log.Printf("DeAIService initialized (Monad Testnet)")
}

func (ctx *Context) InitJarvisService() {
	if ctx.WalletService == nil {
		log.Printf("Warning: WalletService not initialized, cannot init JarvisService")
		return
	}
	ctx.JarvisService = services.NewJarvisService(ctx.WalletService)
	log.Printf("JarvisService initialized (multi-chain support)")
}

func (ctx *Context) InitBlockchainClient() {
	// Monad configuration from constants for desktop app
	rpcURL := contracts.MonadRpcUrl
	tokenAddress := contracts.KawaiTokenAddress
	otcMarketAddress := contracts.OTCMarketAddress
	usdtAddress := contracts.StablecoinAddress

	// Initialize blockchain client
	config := blockchain.Config{
		RPCUrl:           rpcURL,
		TokenAddress:     tokenAddress,
		OTCMarketAddress: otcMarketAddress,
		USDTAddress:      usdtAddress,
	}

	client, err := blockchain.NewClient(config)
	if err != nil {
		log.Printf("Warning: Failed to initialize blockchain client: %v", err)
		log.Printf("  RPC URL: %s", rpcURL)
		log.Printf("  Marketplace will work in local-only mode")
		return
	}

	ctx.BlockchainClient = client
	log.Printf("✅ Blockchain client initialized")

	// Inject as supply querier for KVStore halving logic if KVStore is ready
	if ctx.KVStore != nil {
		ctx.KVStore.SetSupplyQuerier(client)
		log.Printf("Blockchain client injected into KVStore for halving logic")
	}

	log.Printf("  RPC URL: %s", rpcURL)
	log.Printf("  Chain ID: %s", client.ChainID.String())
	log.Printf("  Token Address: %s", tokenAddress)
	log.Printf("  Escrow Address: %s", otcMarketAddress)
	log.Printf("  USDT Address: %s", usdtAddress)

	// Initialize DepositSyncService for manual deposit sync
	if ctx.KVStore != nil {
		syncService, err := services.NewDepositSyncService(ctx.KVStore)
		if err != nil {
			log.Printf("Warning: Failed to initialize DepositSyncService: %v", err)
		} else {
			ctx.DepositSyncService = syncService
			ctx.AddCleanup(func() { syncService.Close() })
			log.Printf("✅ DepositSyncService initialized")
		}
	}
}

func (ctx *Context) InitKVStore() {
	kvStore, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		log.Printf("Warning: Failed to initialize KVStore: %v", err)
		return
	}
	ctx.KVStore = kvStore
	log.Printf("KVStore initialized with Cloudflare KV")
}

func (ctx *Context) InitLanguageModels() {
	bgCtx := context.Background()

	log.Printf("🔍 InitLanguageModels: building model chains...")

	var localModel unillm.LanguageModel

	buildChain := func(taskName string, criteria openrouter.ModelSelectionCriteria) unillm.LanguageModel {
		return llm.BuildChain(llm.ModelChainConfig{
			Context:    bgCtx,
			LocalModel: localModel,
			TaskName:   taskName,
			Criteria:   criteria,
		})
	}

	// Chat model: general purpose
	ctx.ChatModel = buildChain("chat", openrouter.ModelSelectionCriteria{
		RequireReasoning:   false,
		RequireAttachments: false,
		MinContextWindow:   8192,
	})

	// Title model: short text generation
	ctx.TitleModel = buildChain("title", openrouter.ModelSelectionCriteria{
		RequireReasoning:   false,
		RequireAttachments: false,
		MinContextWindow:   4096,
	})

	// Summary model: conversation summarization
	ctx.SummaryModel = buildChain("summary", openrouter.ModelSelectionCriteria{
		RequireReasoning:   false,
		RequireAttachments: false,
		MinContextWindow:   8192,
	})

	// Cleanup model: OCR/transcript cleanup
	ctx.CleanupModel = buildChain("cleanup", openrouter.ModelSelectionCriteria{
		RequireReasoning:   false,
		RequireAttachments: false,
		MinContextWindow:   4096,
	})
}

func buildMuninnMemoryIntegration(backend *services.MuninnMemoryBackend) (*services.MemoryIntegration, error) {
	if backend == nil {
		return nil, fmt.Errorf("muninn backend is required")
	}

	integration, err := services.NewMemoryIntegration(&services.MemoryIntegrationConfig{
		MuninnBackend: backend,
	})
	if err != nil {
		return nil, err
	}
	if !integration.UsesMuninnBackend() {
		return nil, fmt.Errorf("memory integration is not in muninn mode")
	}
	return integration, nil
}

func (ctx *Context) InitMemoryServices() {
	muninnDataDir := filepath.Join(paths.Base(), "muninndb", "veridium")
	backend, err := services.NewMuninnMemoryBackend(muninnDataDir, "veridium_default", 10000, false)
	if err != nil {
		log.Printf("Warning: Muninn memory backend init failed: %v", err)
		return
	}
	ctx.AddCleanup(func() {
		if err := backend.Close(); err != nil {
			log.Printf("Warning: failed to close Muninn memory backend: %v", err)
		}
	})

	integration, err := buildMuninnMemoryIntegration(backend)
	if err != nil {
		log.Printf("Warning: MemoryIntegration init failed: %v", err)
		return
	}

	ctx.MemoryIntegration = integration

	log.Printf("Memory services initialized (MuninnDB active, connection=veridium-memory, vault strategy=user+session scoped)")
}

// InitAll initializes all services in the correct order
func (ctx *Context) InitAll() error {
	ctx.InitSentry() // Initialize Sentry first to capture any startup errors

	if err := ctx.InitDatabase(); err != nil {
		return err
	}
	ctx.InitBasicServices()
	ctx.InitLlamaService()
	ctx.InitVectorStore()
	ctx.InitEmbedder()
	ctx.InitVectorSearch()
	ctx.InitFileLoader()
	ctx.InitKnowledgeBase()
	ctx.InitKVStore()          // MOVED UP
	ctx.InitWalletService()    // Depends on KVStore
	ctx.InitBlockchainClient() // Initialize blockchain client after wallet service
	ctx.InitDeAIService()
	ctx.InitJarvisService()
	// ctx.InitAPIKeyService() // Removed
	// ctx.InitAuthService() // Removed

	log.Printf("🚀 Starting InitLanguageModels()...")
	ctx.InitLanguageModels()
	log.Printf("✅ InitLanguageModels() completed")

	log.Printf("🚀 Starting InitMemoryServices()...")
	ctx.InitMemoryServices()
	log.Printf("✅ InitMemoryServices() completed")

	return nil
}
