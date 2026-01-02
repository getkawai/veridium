package constant

import "time"

// Timeout constants for various operations
// These can be overridden via environment variables in the future
const (
	// Blockchain operation timeouts
	BlockchainCallTimeout    = 10 * time.Second
	BlockchainReceiptTimeout = 30 * time.Second
	BlockchainMaxRetries     = 3
	BlockchainInitialBackoff = 500 * time.Millisecond
	BlockchainMaxBackoff     = 5 * time.Second

	// LLM operation timeouts
	LLMGenerateTimeout        = 120 * time.Second
	LLMCleanupTimeout         = 60 * time.Second
	LLMTitleGenerationTimeout = 3 * time.Minute // Background operation, can be longer

	// Database operation timeouts
	DatabaseQueryTimeout     = 5 * time.Second
	DatabaseMigrationTimeout = 30 * time.Second
	DatabaseSeedTimeout      = 10 * time.Second

	// File processing timeouts
	FileProcessingTimeout  = 5 * time.Minute
	FileUploadTimeout      = 2 * time.Minute
	ImageProcessingTimeout = 30 * time.Second

	// Network operation timeouts
	HTTPRequestTimeout  = 30 * time.Second
	HTTPDownloadTimeout = 5 * time.Minute
	WebSocketTimeout    = 60 * time.Second

	// Cache and storage timeouts
	CacheOperationTimeout = 2 * time.Second
	KVStoreTimeout        = 5 * time.Second
	VectorSearchTimeout   = 10 * time.Second
)

// Retry configuration constants
const (
	// Default retry configuration
	DefaultMaxRetries        = 3
	DefaultInitialBackoff    = 100 * time.Millisecond
	DefaultMaxBackoff        = 2 * time.Second
	DefaultBackoffMultiplier = 2.0

	// KV Store retry configuration
	KVStoreMaxRetries     = 3
	KVStoreInitialBackoff = 100 * time.Millisecond
	KVStoreMaxBackoff     = 2 * time.Second
)
