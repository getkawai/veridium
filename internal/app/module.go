package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/getkawai/database"
	db "github.com/getkawai/database/db"
	"github.com/getkawai/tools"
	yzmabuiltin "github.com/getkawai/tools/builtin"
	"github.com/getkawai/tools/search"
	unillm "github.com/getkawai/unillm"
	llamaembed "github.com/getkawai/unillm/providers/llama-embed"
	"github.com/getkawai/unillm/providers/openrouter"
	"github.com/kawai-network/contracts"
	"github.com/kawai-network/veridium/internal/audio_recorder"
	"github.com/kawai-network/veridium/internal/llm"
	"github.com/kawai-network/veridium/internal/services"
	"github.com/kawai-network/veridium/internal/services/cache"
	"github.com/kawai-network/veridium/internal/tts"
	"github.com/kawai-network/x/blockchain"
	"github.com/kawai-network/x/store"
	"github.com/kawai-network/y/paths"
	"go.uber.org/fx"
)

// Module is the collection of all application providers
var Module = fx.Module("app",
	fx.Provide(
		NewContext,
		ProvideDatabase,
		ProvideQueries,
		ProvideKVStore,
		ProvideSearchService,
		ProvideTTSService,
		ProvideAudioRecorder,
		ProvideToolRegistry,
		ProvideDuckDBStore,
		ProvideFileLoader,
		ProvideVectorSearch,
		ProvideKnowledgeBase,
		ProvideWalletService,
		ProvideBlockchainClient,
		ProvideDepositSyncService,
		ProvideDeAIService,
		ProvideJarvisService,
		ProvideMemoryIntegration,
		ProvideEmbedder,
		ProvideCacheManager,
		ProvideRAGProcessor,
		ProvideLocalLanguageModel,
		fx.Annotate(ProvideChatModel, fx.ResultTags(`name:"chatModel"`)),
		fx.Annotate(ProvideTitleModel, fx.ResultTags(`name:"titleModel"`)),
		fx.Annotate(ProvideSummaryModel, fx.ResultTags(`name:"summaryModel"`)),
		fx.Annotate(ProvideCleanupModel, fx.ResultTags(`name:"cleanupModel"`)),
	),
	fx.Invoke(ProvideSentry),
)

// Provider functions with Lifecycle management

func ProvideDatabase(lc fx.Lifecycle) (*database.Service, error) {
	dbService, err := database.NewService()
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Printf("Closing database...")
			return dbService.Close()
		},
	})
	return dbService, nil
}

func ProvideQueries(ds *database.Service) *db.Queries {
	return ds.Queries()
}

func ProvideKVStore() (*store.KVStore, error) {
	return store.NewMultiNamespaceKVStore()
}

func ProvideSearchService() *search.Service {
	return search.NewService()
}

func ProvideTTSService() *tts.TTSService {
	ttsService, err := tts.NewTTSService()
	if err != nil {
		log.Printf("Warning: TTS init failed: %v", err)
		return nil
	}
	return ttsService
}

func ProvideAudioRecorder() *audio_recorder.AudioRecorderService {
	return audio_recorder.NewAudioRecorderService(nil)
}

func ProvideToolRegistry(ds *database.Service) *tools.ToolRegistry {
	registry := tools.NewToolRegistry()
	if err := yzmabuiltin.RegisterAllWithDB(registry, ds.DB()); err != nil {
		log.Printf("Warning: Failed to register builtin tools: %v", err)
	}
	return registry
}

func ProvideDuckDBStore(lc fx.Lifecycle) (*services.DuckDBStore, error) {
	if err := os.MkdirAll(paths.FileBase(), 0755); err != nil {
		return nil, fmt.Errorf("creating base dir %s: %w", paths.FileBase(), err)
	}
	s, err := services.NewDuckDBStore(paths.DuckDB(), EmbeddingDims)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Printf("Closing DuckDB store...")
			return s.Close()
		},
	})
	return s, nil
}

func ProvideFileLoader() *services.FileLoader {
	return services.NewFileLoader()
}

func ProvideVectorSearch(ds *database.Service, duck *services.DuckDBStore) *services.VectorSearchService {
	// Embedder is handled as nil for now to match current state
	vs, err := services.NewVectorSearchService(ds.DB(), duck, nil)
	if err != nil {
		log.Printf("Warning: VectorSearch init failed: %v", err)
		return nil
	}
	return vs
}

func ProvideKnowledgeBase(ds *database.Service, vs *services.VectorSearchService, fl *services.FileLoader, duck *services.DuckDBStore, rag *services.RAGProcessor) *services.KnowledgeBaseService {
	if vs == nil || fl == nil {
		return nil
	}
	kb, err := services.NewKnowledgeBaseService(ds, &services.KnowledgeBaseConfig{
		RAGProcessor: rag,
		VectorSearch: vs,
		FileLoader:   fl,
		AssetDir:     paths.KBAssets(),
	})
	if err != nil {
		log.Printf("Warning: KnowledgeBase init failed: %v", err)
		return nil
	}
	return kb
}

func ProvideWalletService(kv *store.KVStore) *services.WalletService {
	return services.NewWalletService(paths.FileBase(), kv)
}

func ProvideBlockchainClient(kv *store.KVStore) *blockchain.Client {
	config := blockchain.Config{
		RPCUrl:           contracts.MonadRpcUrl,
		TokenAddress:     contracts.KawaiTokenAddress,
		OTCMarketAddress: contracts.OTCMarketAddress,
		USDTAddress:      contracts.StablecoinAddress,
	}

	client, err := blockchain.NewClient(config)
	if err != nil {
		log.Printf("Warning: Failed to initialize blockchain client: %v", err)
		return nil
	}

	if kv != nil {
		kv.SetSupplyQuerier(client)
	}
	return client
}

func ProvideDepositSyncService(lc fx.Lifecycle, kv *store.KVStore) *services.DepositSyncService {
	if kv == nil {
		return nil
	}
	sync, err := services.NewDepositSyncService(kv)
	if err != nil {
		log.Printf("Warning: Failed to initialize DepositSyncService: %v", err)
		return nil
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			sync.Close()
			return nil
		},
	})
	return sync
}

func ProvideDeAIService(ws *services.WalletService, kv *store.KVStore) *services.DeAIService {
	if ws == nil {
		return nil
	}
	return services.NewDeAIService(ws, kv)
}

func ProvideJarvisService(ws *services.WalletService) *services.JarvisService {
	if ws == nil {
		return nil
	}
	return services.NewJarvisService(ws)
}

func ProvideMemoryIntegration(lc fx.Lifecycle) *services.MemoryIntegration {
	muninnDataDir := filepath.Join(paths.Base(), "muninndb", "veridium")
	backend, err := services.NewMuninnMemoryBackend(muninnDataDir, "veridium_default", 10000, false)
	if err != nil {
		log.Printf("Warning: Muninn memory backend init failed: %v", err)
		return nil
	}

	integration, err := services.NewMemoryIntegration(&services.MemoryIntegrationConfig{
		MuninnBackend: backend,
	})
	if err != nil {
		log.Printf("Warning: MemoryIntegration init failed: %v", err)
		_ = backend.Close()
		return nil
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Printf("Closing Muninn memory backend...")
			return backend.Close()
		},
	})

	return integration
}

func ProvideEmbedder() llamaembed.Embedder {
	// TODO: Re-implement embedder initialization after Context.LibService removal.
	// For now, return nil to match original behavior.
	return nil
}

func ProvideCacheManager() *cache.CacheManager {
	// TODO: Implement CacheManager initialization if needed.
	// For now, return nil as it's not explicitly initialized in original context.go.
	return nil
}

func ProvideRAGProcessor(ds *database.Service, duck *services.DuckDBStore, fl *services.FileLoader, embedder llamaembed.Embedder) *services.RAGProcessor {
	// Check for nil dependencies if they are optional
	if ds == nil || duck == nil || fl == nil {
		log.Printf("Warning: RAGProcessor dependencies (DB, DuckDBStore, FileLoader) are nil. RAGProcessor will be nil.")
		return nil
	}
	return services.NewRAGProcessor(ds.DB(), duck, fl, embedder)
}

// Language Models
// This requires a local model as a fallback and then builds a chain.
// We'll create a dummy local model for now as the actual implementation is a TODO.
func ProvideLocalLanguageModel() unillm.LanguageModel {
	// TODO: Re-implement local llama service initialization after Context.LibService removal.
	// For now, return a dummy model.
	return &DummyLanguageModel{}
}

// DummyLanguageModel for placeholders where the real LLM implementation is TODO.
type DummyLanguageModel struct{}

func (d *DummyLanguageModel) Generate(ctx context.Context, call unillm.Call) (*unillm.Response, error) {
	return nil, fmt.Errorf("local model unavailable: fallback dummy invoked")
}

func (d *DummyLanguageModel) Stream(ctx context.Context, call unillm.Call) (unillm.StreamResponse, error) {
	return nil, fmt.Errorf("stream not implemented for dummy model")
}

func (d *DummyLanguageModel) GenerateObject(ctx context.Context, call unillm.ObjectCall) (*unillm.ObjectResponse, error) {
	return nil, fmt.Errorf("generate object not implemented for dummy model")
}

func (d *DummyLanguageModel) StreamObject(ctx context.Context, call unillm.ObjectCall) (unillm.ObjectStreamResponse, error) {
	return nil, fmt.Errorf("stream object not implemented for dummy model")
}

func (d *DummyLanguageModel) Provider() string { return "dummy_local" }
func (d *DummyLanguageModel) Model() string    { return "dummy" }


// ProvideChatModel builds the chain for the chat model
func ProvideChatModel(lc fx.Lifecycle, localModel unillm.LanguageModel) unillm.LanguageModel {
	// We need a background context for buildModelChain
	bgCtx, cancel := context.WithCancel(context.Background())
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			cancel()
			return nil
		},
	})

	chain := llm.BuildModelChain(bgCtx, localModel, openrouter.ModelSelectionCriteria{
		MinContextWindow: 200000,
		RequireReasoning: false,
	}, "Chat Model")
	if len(chain) == 0 {
		return localModel // Fallback to local if chain build fails
	}
	chainModel, err := unillm.NewChain(chain)
	if err != nil {
		log.Printf("Warning: Failed to create chat model chain: %v, falling back to local model", err)
		return localModel
	}
	return chainModel
}

// ProvideTitleModel builds the chain for the title model
func ProvideTitleModel(lc fx.Lifecycle, localModel unillm.LanguageModel) unillm.LanguageModel {
	bgCtx, cancel := context.WithCancel(context.Background())
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			cancel()
			return nil
		},
	})

	chain := llm.BuildModelChain(bgCtx, localModel, openrouter.ModelSelectionCriteria{
		MinContextWindow: 80000,
		RequireReasoning: false,
	}, "Title Model")
	if len(chain) == 0 {
		return localModel
	}
	chainModel, err := unillm.NewChain(chain)
	if err != nil {
		log.Printf("Warning: Failed to create title model chain: %v, falling back to local model", err)
		return localModel
	}
	return chainModel
}

// ProvideSummaryModel builds the chain for the summary model
func ProvideSummaryModel(lc fx.Lifecycle, localModel unillm.LanguageModel) unillm.LanguageModel {
	bgCtx, cancel := context.WithCancel(context.Background())
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			cancel()
			return nil
		},
	})

	chain := llm.BuildModelChain(bgCtx, localModel, openrouter.ModelSelectionCriteria{
		MinContextWindow: 80000,
		RequireReasoning: false,
	}, "Summary Model")
	if len(chain) == 0 {
		return localModel
	}
	chainModel, err := unillm.NewChain(chain)
	if err != nil {
		log.Printf("Warning: Failed to create summary model chain: %v, falling back to local model", err)
		return localModel
	}
	return chainModel
}

// ProvideCleanupModel builds the chain for the cleanup model
func ProvideCleanupModel(lc fx.Lifecycle, localModel unillm.LanguageModel) unillm.LanguageModel {
	bgCtx, cancel := context.WithCancel(context.Background())
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			cancel()
			return nil
		},
	})

	chain := llm.BuildModelChain(bgCtx, localModel, openrouter.ModelSelectionCriteria{
		MinContextWindow: 80000,
		RequireReasoning: false,
	}, "Cleanup Model")
	if len(chain) == 0 {
		return localModel
	}
	chainModel, err := unillm.NewChain(chain)
	if err != nil {
		log.Printf("Warning: Failed to create cleanup model chain: %v, falling back to local model", err)
		return localModel
	}
	return chainModel
}
