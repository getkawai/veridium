// Package app provides the core application context and service initialization.
// This is the SINGLE SOURCE OF TRUTH for service initialization.
// main.go imports and uses this, tests import and use this.
package app

import (
	"context"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/kawai-network/veridium/internal/audio_recorder"
	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/internal/database"
	db "github.com/kawai-network/veridium/internal/database/generated"
	"github.com/kawai-network/veridium/internal/paths"
	"github.com/kawai-network/veridium/internal/search"
	"github.com/kawai-network/veridium/internal/services"
	"github.com/kawai-network/veridium/internal/services/cache"
	"github.com/kawai-network/veridium/internal/tts"
	"github.com/kawai-network/veridium/internal/whisper"
	"github.com/kawai-network/veridium/pkg/blockchain"
	"github.com/kawai-network/veridium/pkg/fantasy"
	"github.com/kawai-network/veridium/pkg/fantasy/llamalib"
	googleprovider "github.com/kawai-network/veridium/pkg/fantasy/providers/google"
	llamaprovider "github.com/kawai-network/veridium/pkg/fantasy/providers/llama"
	llamaembed "github.com/kawai-network/veridium/pkg/fantasy/providers/llama-embed"
	"github.com/kawai-network/veridium/pkg/fantasy/providers/openaicompat"
	"github.com/kawai-network/veridium/pkg/fantasy/providers/openrouter"
	"github.com/kawai-network/veridium/pkg/fantasy/providers/pooled"
	"github.com/kawai-network/veridium/pkg/fantasy/tools"
	yzmabuiltin "github.com/kawai-network/veridium/pkg/fantasy/tools/builtin"
	"github.com/kawai-network/veridium/pkg/logger"
	"github.com/kawai-network/veridium/pkg/store"
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
	LibService   *llamalib.Service
	Embedder     llamaembed.Embedder
	CacheManager *cache.CacheManager

	// Storage Services
	DuckDBStore *services.DuckDBStore
	FileLoader  *services.FileLoader
	KVStore     *store.KVStore

	// Feature Services
	SearchService  *search.Service
	TTSService     *tts.TTSService
	WhisperService *whisper.Service
	AudioRecorder  *audio_recorder.AudioRecorderService
	VectorSearch   *services.VectorSearchService
	KBService      *services.KnowledgeBaseService
	RAGProcessor   *services.RAGProcessor
	ToolRegistry   *tools.ToolRegistry
	WalletService  *services.WalletService
	DeAIService    *services.DeAIService
	JarvisService  *services.JarvisService

	// Blockchain Services
	BlockchainClient   *blockchain.Client
	DepositSyncService *services.DepositSyncService

	// Language Models
	ChatModel    fantasy.LanguageModel `json:"-"`
	TitleModel   fantasy.LanguageModel `json:"-"`
	SummaryModel fantasy.LanguageModel `json:"-"`
	CleanupModel fantasy.LanguageModel `json:"-"`

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
	ctx.AddCleanup(func() { dbService.Close() })
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

	if whisperService, err := whisper.NewService(); err != nil {
		log.Printf("Warning: Whisper init failed: %v", err)
	} else {
		ctx.WhisperService = whisperService
		ctx.AddCleanup(func() { whisperService.Close() })
		log.Printf("Whisper initialized")
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
	handler := slog.NewTextHandler(os.Stderr, nil)
	// wrap the handler with SentryHandler
	sentryHandler := logger.NewSentryHandler(handler)
	// create a new logger with the SentryHandler
	logger := slog.New(sentryHandler)
	// set the default logger to the new logger
	slog.SetDefault(logger)

	log.Printf("Sentry initialized with EnableLogs: true (using slog handler)")
}

func (ctx *Context) InitLlamaService() {
	// NewService() automatically starts background initialization
	// which handles: installation check, library loading, and model downloads
	ctx.LibService = llamalib.NewService()
	log.Printf("LlamaService created (initializing in background)")
}

func (ctx *Context) InitVectorStore() {
	os.MkdirAll(paths.FileBase(), 0755)
	duckDBStore, err := services.NewDuckDBStore(paths.DuckDB(), EmbeddingDims)
	if err != nil {
		log.Printf("Warning: DuckDB init failed: %v", err)
		return
	}
	ctx.DuckDBStore = duckDBStore
	ctx.AddCleanup(func() { duckDBStore.Close() })
	log.Printf("DuckDB initialized (path: %s)", paths.DuckDB())
}

func (ctx *Context) InitEmbedder() {
	// Ensure llama library is initialized before creating embedder
	if ctx.LibService == nil {
		log.Printf("Warning: LibService not available, skipping embedder initialization")
		return
	}

	// Wait for library to be ready (with timeout)
	bgCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := ctx.LibService.WaitForInitialization(bgCtx); err != nil {
		log.Printf("Warning: Failed to wait for llama library initialization: %v", err)
		log.Printf("Embedder will not be available")
		return
	}

	modelName := llamalib.GetRecommendedEmbeddingModel()
	model, exists := llamalib.GetEmbeddingModel(modelName)
	if !exists {
		log.Printf("Warning: Embedding model not found: %s", modelName)
		return
	}

	installer := llamalib.NewLlamaCppInstaller()
	modelPath := filepath.Join(installer.GetModelsDirectory(), model.Filename)

	baseEmbedder, err := llamaembed.NewLlamaEmbedder(&llamaembed.LlamaConfig{
		ModelPath:   modelPath,
		ContextSize: 2048,
	})
	if err != nil {
		log.Printf("Warning: Embedder init failed: %v", err)
		return
	}
	ctx.AddCleanup(func() { baseEmbedder.Close() })
	log.Printf("Embedder initialized (model: %s, dims: %d)", model.Name, baseEmbedder.Dimensions())

	// Production-only: always use cache layer
	ctx.CacheManager = cache.NewCacheManager(baseEmbedder, nil)
	ctx.Embedder = ctx.CacheManager.GetCachedEmbedder(baseEmbedder)
	log.Printf("Cache layer initialized")
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
	rpcURL := constant.MonadRpcUrl
	tokenAddress := constant.KawaiTokenAddress
	otcMarketAddress := constant.OTCMarketAddress
	usdtAddress := constant.StablecoinAddress

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
	if ctx.LibService == nil {
		return
	}

	// Use a timeout context to prevent hanging on slow network calls
	bgCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	llamaProvider, err := llamaprovider.New(
		llamaprovider.WithService(ctx.LibService),
		llamaprovider.WithToolRegistry(ctx.ToolRegistry),
	)
	if err != nil {
		log.Printf("Warning: Llama provider failed: %v", err)
		return
	}

	localModel, err := llamaProvider.LanguageModel(bgCtx, "")
	if err != nil {
		log.Printf("Warning: Local LLM failed: %v", err)
		return
	}

	// Circuit breaker: skip rate-limited models until app restart (rate limit is daily, cache is in-memory)
	circuitBreaker := fantasy.WithCircuitBreaker(1, 0)

	// AUTO-DETECT: Use pooled providers if multiple keys available
	// usePooled := len(constant.GetOpenRouterApiKeys()) > 1 || len(constant.GetZaiApiKeys()) > 1
	usePooled := false

	var buildChain func(context.Context, fantasy.LanguageModel, openrouter.ModelSelectionCriteria, string) []fantasy.LanguageModel
	if usePooled {
		buildChain = ctx.buildModelChainV2
		log.Printf("✅ Using POOLED providers (multiple API keys detected)")
	} else {
		buildChain = ctx.buildModelChain
		log.Printf("ℹ️  Using SIMPLE providers (single API key)")
	}

	ctx.ChatModel, err = fantasy.NewChain(buildChain(bgCtx, localModel, openrouter.ModelSelectionCriteria{
		RequireReasoning: true, RequireAttachments: true, MinContextWindow: 100000,
	}, "ChatModel"), circuitBreaker)
	if err != nil {
		log.Printf("Warning: ChatModel chain creation failed: %v", err)
	}

	ctx.TitleModel, err = fantasy.NewChain(buildChain(bgCtx, localModel, openrouter.ModelSelectionCriteria{}, "TitleModel"), circuitBreaker)
	if err != nil {
		log.Printf("Warning: TitleModel chain creation failed: %v", err)
	}

	ctx.SummaryModel, err = fantasy.NewChain(buildChain(bgCtx, localModel, openrouter.ModelSelectionCriteria{
		MinContextWindow: 50000,
	}, "SummaryModel"), circuitBreaker)
	if err != nil {
		log.Printf("Warning: SummaryModel chain creation failed: %v", err)
	}

	ctx.CleanupModel, err = fantasy.NewChain(buildChain(bgCtx, localModel, openrouter.ModelSelectionCriteria{}, "CleanupModel"), circuitBreaker)
	if err != nil {
		log.Printf("Warning: CleanupModel chain creation failed: %v", err)
	}

	log.Printf("Language models initialized with %s", map[bool]string{true: "POOLED providers", false: "SIMPLE providers"}[usePooled])
}

func (ctx *Context) buildModelChain(bgCtx context.Context, localModel fantasy.LanguageModel, criteria openrouter.ModelSelectionCriteria, taskName string) []fantasy.LanguageModel {
	var chain []fantasy.LanguageModel

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

// buildModelChainV2 creates a model chain with pooled providers and automatic rotation.
// This version uses account pooling and smart error handling with metrics.
func (ctx *Context) buildModelChainV2(bgCtx context.Context, localModel fantasy.LanguageModel, criteria openrouter.ModelSelectionCriteria, taskName string) []fantasy.LanguageModel {
	var chain []fantasy.LanguageModel

	// 1. OpenRouter with multiple API keys (pooled with metrics & rotation)
	openRouterKeys := constant.GetOpenRouterApiKeys()
	if len(openRouterKeys) > 0 {
		pooledProvider, err := pooled.New(pooled.Config{
			ProviderName:   "openrouter",
			BaseURL:        "https://openrouter.ai/api/v1",
			ModelName:      "auto", // Will be selected by criteria
			APIKeys:        openRouterKeys,
			EnableMetrics:  true, // Enable metrics tracking
			EnableRotation: true, // Enable auto rotation
			RotationStrategy: &pooled.HealthBasedStrategy{
				MaxConsecutiveFailures: 3,
			},
			CreateClient: func(apiKey string) (fantasy.LanguageModel, error) {
				provider, err := openrouter.New(
					openrouter.WithAPIKey(apiKey),
					openrouter.WithModelSelection(criteria),
				)
				if err != nil {
					return nil, err
				}
				return provider.LanguageModel(bgCtx, "")
			},
		})

		if err == nil {
			chain = append(chain, pooledProvider)
			catalog := openrouter.GetCatalog()
			if selected := catalog.SelectFreeModel(criteria); selected != nil {
				log.Printf("%s: OpenRouter Pooled (%s) with %d keys [Metrics: ON, Rotation: ON]", taskName, selected.ID, len(openRouterKeys))
			}
		} else {
			log.Printf("Warning: Failed to create pooled OpenRouter: %v", err)
		}
	}

	// 2. Pollinations AI (no pooling needed, free service)
	if provider, err := openaicompat.New(
		openaicompat.WithName("pollinations"),
		openaicompat.WithBaseURL("https://text.pollinations.ai/openai"),
		openaicompat.WithAPIKey("dummy"),
	); err == nil {
		if pollinationsModel, err := provider.LanguageModel(bgCtx, "openai"); err == nil {
			chain = append(chain, pollinationsModel)
			log.Printf("%s: Pollinations AI (openai)", taskName)
		}
	}

	// 3. ZAI with multiple API keys (pooled with metrics & rotation)
	zaiKeys := constant.GetZaiApiKeys()
	if len(zaiKeys) > 0 {
		pooledProvider, err := pooled.New(pooled.Config{
			ProviderName:   "zai",
			BaseURL:        "https://api.z.ai/api/coding/paas/v4",
			ModelName:      "glm-4.7",
			APIKeys:        zaiKeys,
			EnableMetrics:  true,
			EnableRotation: true,
			RotationStrategy: &pooled.RoundRobinStrategy{
				RotateAfterRequests: 100, // Rotate after 100 requests
			},
			CreateClient: func(apiKey string) (fantasy.LanguageModel, error) {
				provider, err := openaicompat.New(
					openaicompat.WithName("zai"),
					openaicompat.WithBaseURL("https://api.z.ai/api/coding/paas/v4"),
					openaicompat.WithAPIKey(apiKey),
				)
				if err != nil {
					return nil, err
				}
				return provider.LanguageModel(bgCtx, "glm-4.7")
			},
		})

		if err == nil {
			chain = append(chain, pooledProvider)
			log.Printf("%s: ZAI Pooled (glm-4.7) with %d keys [Metrics: ON, Rotation: ON]", taskName, len(zaiKeys))
		} else {
			log.Printf("Warning: Failed to create pooled ZAI: %v", err)
		}
	}

	// 4. Local model (final fallback)
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

	var llm fantasy.LanguageModel
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
