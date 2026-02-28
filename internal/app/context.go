// Package app provides the core application context and service initialization.
// This is the SINGLE SOURCE OF TRUTH for service initialization.
// main.go imports and uses this, tests import and use this.
package app

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/getkawai/database"
	db "github.com/getkawai/database/db"
	"github.com/getkawai/tools"
	yzmabuiltin "github.com/getkawai/tools/builtin"
	"github.com/getkawai/tools/search"
	unillm "github.com/getkawai/unillm"
	googleprovider "github.com/getkawai/unillm/providers/google"
	llamaembed "github.com/getkawai/unillm/providers/llama-embed"
	"github.com/getkawai/unillm/providers/openaicompat"
	"github.com/getkawai/unillm/providers/openrouter"
	"github.com/getsentry/sentry-go"
	"github.com/kawai-network/contracts"
	"github.com/kawai-network/veridium/internal/audio_recorder"
	"github.com/kawai-network/veridium/internal/services"
	"github.com/kawai-network/veridium/internal/services/cache"
	"github.com/kawai-network/veridium/internal/tts"
	"github.com/kawai-network/x/blockchain"
	"github.com/kawai-network/x/constant"
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

	// Memory Services (MemGPT-style)
	MemoryService     *services.MemoryService
	MemoryEnrichment  *services.MemoryEnrichmentService
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
	log.Printf("DuckDB initialized (path: %s)", paths.DuckDB())
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
	// TODO: Re-implement language model chain initialization after Context.LibService removal.
}

func (ctx *Context) buildModelChain(bgCtx context.Context, localModel unillm.LanguageModel, criteria openrouter.ModelSelectionCriteria, taskName string) []unillm.LanguageModel {
	var chain []unillm.LanguageModel

	// 1. Google Gemini 2.5 Flash-Lite (free tier with highest limits: 15 RPM, 1000 RPD)
	if apiKey := constant.GetRandomGeminiApiKey(); apiKey != "" {
		log.Printf("🔍 %s: Initializing Google Gemini 2.5 Flash-Lite...", taskName)
		if provider, err := googleprovider.New(googleprovider.WithGeminiAPIKey(apiKey)); err == nil {
			if geminiModel, err := provider.LanguageModel(bgCtx, "gemini-2.5-flash-lite"); err == nil {
				chain = append(chain, geminiModel)
				log.Printf("✅ %s: Added Google Gemini (gemini-2.5-flash-lite) to chain [15 RPM, 1000 RPD]", taskName)
			} else {
				log.Printf("❌ %s: Google Gemini provider initialized but failed to get model: %v", taskName, err)
			}
		} else {
			log.Printf("❌ %s: Failed to initialize Google Gemini provider: %v", taskName, err)
		}
	} else {
		log.Printf("ℹ️  %s: Skipping Google Gemini (no API key)", taskName)
	}

	// 2. OpenRouter (free tier)
	if apiKey := constant.GetRandomOpenRouterApiKey(); apiKey != "" {
		log.Printf("🔍 %s: Initializing OpenRouter...", taskName)
		if provider, err := openrouter.New(openrouter.WithAPIKey(apiKey), openrouter.WithModelSelection(criteria)); err == nil {
			if remoteModel, err := provider.LanguageModel(bgCtx, ""); err == nil {
				chain = append(chain, remoteModel)
				catalog := openrouter.GetCatalog()
				if selected := catalog.SelectFreeModel(criteria); selected != nil {
					log.Printf("✅ %s: Added OpenRouter (%s) to chain", taskName, selected.ID)
				} else {
					log.Printf("⚠️  %s: OpenRouter initialized but no free model matched criteria", taskName)
				}
			} else {
				log.Printf("❌ %s: OpenRouter provider initialized but failed to get model: %v", taskName, err)
			}
		} else {
			log.Printf("❌ %s: Failed to initialize OpenRouter provider: %v", taskName, err)
		}
	} else {
		log.Printf("ℹ️  %s: Skipping OpenRouter (no API key)", taskName)
	}

	// 3. Pollinations AI (fallback before local)
	if provider, err := openaicompat.New(
		openaicompat.WithName("pollinations"),
		openaicompat.WithBaseURL("https://text.pollinations.ai/openai"),
		openaicompat.WithAPIKey("dummy"), // Pollinations doesn't require API key, but SDK needs one
	); err == nil {
		if pollinationsModel, err := provider.LanguageModel(bgCtx, "openai"); err == nil {
			chain = append(chain, pollinationsModel)
			log.Printf("%s: Pollinations AI (openai)", taskName)
		}
	}

	// 4. ZAI GLM-4.6 (fallback before local)
	if provider, err := openaicompat.New(
		openaicompat.WithName("zai"),
		openaicompat.WithBaseURL("https://api.z.ai/api/coding/paas/v4"),
		openaicompat.WithAPIKey(constant.GetRandomZaiApiKey()),
	); err == nil {
		if zaiModel, err := provider.LanguageModel(bgCtx, "glm-4.7"); err == nil {
			chain = append(chain, zaiModel)
			log.Printf("%s: ZAI (glm-4.7)", taskName)
		}
	}

	// 5. Local model (final fallback)
	chain = append(chain, localModel)
	log.Printf("%s: Chain created with %d models (fallback: %s/%s)", taskName, len(chain), localModel.Provider(), localModel.Model())
	return chain
}

func (ctx *Context) InitMemoryServices() {
	if ctx.DB == nil || ctx.DuckDBStore == nil || ctx.Embedder == nil {
		log.Printf("Warning: Prerequisites not met for Memory services")
		return
	}

	memService, err := services.NewMemoryService(ctx.DB, &services.MemoryServiceConfig{
		DuckDB: ctx.DuckDBStore, Embedder: ctx.Embedder, EmbeddingDim: EmbeddingDims,
	})
	if err != nil {
		log.Printf("Warning: MemoryService init failed: %v", err)
		return
	}
	ctx.MemoryService = memService

	var llm unillm.LanguageModel
	if ctx.ChatModel != nil {
		llm = ctx.ChatModel
		log.Printf("Memory enrichment: using LLM")
	} else {
		log.Printf("Memory enrichment: rule-based fallback")
	}

	enrichService, err := services.NewMemoryEnrichmentService(&services.MemoryEnrichmentConfig{
		MemoryService: memService, LLM: llm,
	})
	if err != nil {
		log.Printf("Warning: MemoryEnrichment init failed: %v", err)
		return
	}
	ctx.MemoryEnrichment = enrichService

	integration, err := services.NewMemoryIntegration(&services.MemoryIntegrationConfig{
		MemoryService: memService, EnrichmentService: enrichService,
	})
	if err != nil {
		log.Printf("Warning: MemoryIntegration init failed: %v", err)
		return
	}
	ctx.MemoryIntegration = integration

	log.Printf("Memory services initialized (MemGPT-style)")
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
