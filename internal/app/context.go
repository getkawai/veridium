// Package app provides the core application context and service initialization.
// This is the SINGLE SOURCE OF TRUTH for service initialization.
// main.go imports and uses this, tests import and use this.
package app

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/kawai-network/veridium/internal/audio_recorder"
	"github.com/kawai-network/veridium/internal/database"
	db "github.com/kawai-network/veridium/internal/database/generated"
	"github.com/kawai-network/veridium/internal/search"
	"github.com/kawai-network/veridium/internal/services"
	"github.com/kawai-network/veridium/internal/services/cache"
	"github.com/kawai-network/veridium/internal/tts"
	"github.com/kawai-network/veridium/internal/whisper"
	"github.com/kawai-network/veridium/pkg/fantasy"
	"github.com/kawai-network/veridium/pkg/fantasy/llamalib"
	llamaprovider "github.com/kawai-network/veridium/pkg/fantasy/providers/llama"
	llamaembed "github.com/kawai-network/veridium/pkg/fantasy/providers/llama-embed"
	"github.com/kawai-network/veridium/pkg/fantasy/providers/openaicompat"
	"github.com/kawai-network/veridium/pkg/fantasy/providers/openrouter"
	"github.com/kawai-network/veridium/pkg/fantasy/tools"
	yzmabuiltin "github.com/kawai-network/veridium/pkg/fantasy/tools/builtin"
)

// Configuration constants
const (
	FileBaseDir   = "files"
	DuckDBPath    = "data/duckdb.db"
	KBAssetPath   = "data/kb-assets"
	EmbeddingDims = 384
	DefaultUserID = "DEFAULT_LOBE_CHAT_USER"
)

// DevMode check
var DevMode = os.Getenv("VERIDIUM_DEV") == "1" || os.Getenv("DEV") == "1"

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

	// Language Models
	ChatModel    fantasy.LanguageModel
	TitleModel   fantasy.LanguageModel
	SummaryModel fantasy.LanguageModel
	CleanupModel fantasy.LanguageModel

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

func (ctx *Context) InitLlamaService() {
	ctx.LibService = llamalib.NewService()
}

func (ctx *Context) InitVectorStore() {
	os.MkdirAll(FileBaseDir, 0755)
	duckDBStore, err := services.NewDuckDBStore(DuckDBPath, EmbeddingDims)
	if err != nil {
		log.Printf("Warning: DuckDB init failed: %v", err)
		return
	}
	ctx.DuckDBStore = duckDBStore
	ctx.AddCleanup(func() { duckDBStore.Close() })
	log.Printf("DuckDB initialized (path: %s)", DuckDBPath)
}

func (ctx *Context) InitEmbedder() {
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

	if DevMode {
		ctx.Embedder = baseEmbedder
	} else {
		ctx.CacheManager = cache.NewCacheManager(baseEmbedder, nil)
		ctx.Embedder = ctx.CacheManager.GetCachedEmbedder(baseEmbedder)
		log.Printf("Cache layer initialized")
	}
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
		AssetDir:     KBAssetPath,
	})
	if err != nil {
		log.Printf("Warning: KnowledgeBase init failed: %v", err)
		return
	}
	ctx.KBService = kbService
	log.Printf("KnowledgeBase initialized")
}

func (ctx *Context) InitWalletService() {
	ctx.WalletService = services.NewWalletService("data")
	log.Printf("WalletService initialized (using jarvis/accounts at ~/.jarvis/)")
}

func (ctx *Context) InitDeAIService() {
	if ctx.WalletService == nil {
		log.Printf("Warning: WalletService not initialized, cannot init DeAIService")
		return
	}
	ctx.DeAIService = services.NewDeAIService(ctx.WalletService)
	log.Printf("DeAIService initialized (BSC Testnet)")
}

func (ctx *Context) InitLanguageModels() {
	if ctx.LibService == nil {
		return
	}

	bgCtx := context.Background()
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

	ctx.ChatModel, _ = fantasy.NewChain(ctx.buildModelChain(bgCtx, localModel, openrouter.ModelSelectionCriteria{
		RequireReasoning: true, RequireAttachments: true, MinContextWindow: 100000,
	}, "ChatModel"), circuitBreaker)

	ctx.TitleModel, _ = fantasy.NewChain(ctx.buildModelChain(bgCtx, localModel, openrouter.ModelSelectionCriteria{}, "TitleModel"), circuitBreaker)

	ctx.SummaryModel, _ = fantasy.NewChain(ctx.buildModelChain(bgCtx, localModel, openrouter.ModelSelectionCriteria{
		MinContextWindow: 50000,
	}, "SummaryModel"), circuitBreaker)

	ctx.CleanupModel, _ = fantasy.NewChain(ctx.buildModelChain(bgCtx, localModel, openrouter.ModelSelectionCriteria{}, "CleanupModel"), circuitBreaker)

	log.Printf("Language models initialized")
}

func (ctx *Context) buildModelChain(bgCtx context.Context, localModel fantasy.LanguageModel, criteria openrouter.ModelSelectionCriteria, taskName string) []fantasy.LanguageModel {
	var chain []fantasy.LanguageModel

	// 1. OpenRouter (free tier)
	orKey := os.Getenv("OPENROUTER_API_KEY")
	if orKey == "" {
		orKey = "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"
	}
	if orKey != "" {
		if provider, err := openrouter.New(openrouter.WithAPIKey(orKey), openrouter.WithModelSelection(criteria)); err == nil {
			if remoteModel, err := provider.LanguageModel(bgCtx, ""); err == nil {
				chain = append(chain, remoteModel)
				catalog := openrouter.GetCatalog()
				if selected := catalog.SelectFreeModel(criteria); selected != nil {
					log.Printf("%s: OpenRouter (%s)", taskName, selected.ID)
				}
			}
		}
	}

	// 2. Pollinations AI (fallback before local)
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

	// 2. ZAI GLM-4.6 (fallback before local)
	zaiKey := os.Getenv("ZAI_API_KEY")
	if zaiKey == "" {
		zaiKey = "a10854167085448cac33753523919ac9.D41CLq6KxXTY7g4u" // dev key
	}
	if zaiKey != "" {
		if provider, err := openaicompat.New(
			openaicompat.WithName("zai"),
			openaicompat.WithBaseURL("https://api.z.ai/api/coding/paas/v4"),
			openaicompat.WithAPIKey(zaiKey),
		); err == nil {
			if zaiModel, err := provider.LanguageModel(bgCtx, "glm-4.6"); err == nil {
				chain = append(chain, zaiModel)
				log.Printf("%s: ZAI (glm-4.6)", taskName)
			}
		}
	}

	// 3. Local model (final fallback)
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
	ctx.InitWalletService()
	ctx.InitDeAIService()
	ctx.InitLanguageModels()
	ctx.InitMemoryServices()
	return nil
}
