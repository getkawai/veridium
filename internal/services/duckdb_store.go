package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/kawai-network/veridium/pkg/xlog"
	_ "github.com/marcboeker/go-duckdb"
)

// DuckDBStore handles vector storage using DuckDB
type DuckDBStore struct {
	db     *sql.DB
	config *HNSWConfig
}

// HNSWConfig contains configuration options for HNSW index
// Reference: https://duckdb.org/docs/stable/core_extensions/vss
type HNSWConfig struct {
	// Metric determines the distance function to use
	// Options: "l2sq" (Euclidean/L2), "cosine", "ip" (inner product)
	// Default: "l2sq"
	Metric string

	// EfConstruction controls the size of the dynamic candidate list during index construction
	// Higher values = better quality index but slower construction
	// Default: 128, Range: typically 100-500
	EfConstruction int

	// EfSearch controls the size of the dynamic candidate list during search
	// Higher values = better recall but slower search
	// Default: 64, Range: typically 10-500
	// Can be overridden at runtime with SET hnsw_ef_search
	EfSearch int

	// M is the number of bi-directional links created for each node
	// Higher values = better recall but more memory usage
	// Default: 16, Range: typically 5-48
	M int
}

// DefaultHNSWConfig returns default HNSW configuration
func DefaultHNSWConfig() *HNSWConfig {
	return &HNSWConfig{
		Metric:         "l2sq",
		EfConstruction: 128,
		EfSearch:       64,
		M:              16,
	}
}

// NewDuckDBStore creates a new DuckDB store with default HNSW configuration
func NewDuckDBStore(path string, embeddingDim int) (*DuckDBStore, error) {
	return NewDuckDBStoreWithConfig(path, embeddingDim, DefaultHNSWConfig())
}

// NewDuckDBStoreWithConfig creates a new DuckDB store with custom HNSW configuration
func NewDuckDBStoreWithConfig(path string, embeddingDim int, config *HNSWConfig) (*DuckDBStore, error) {
	// Open DuckDB connection
	// If path is empty, it uses in-memory DB (useful for testing)
	dsn := path
	if dsn == "" {
		dsn = ":memory:"
	}

	// Add access_mode=READ_WRITE to ensure we can write
	// Note: We don't use ?access_mode=READ_WRITE because it might interfere with path
	// DuckDB will use default READ_WRITE mode
	db, err := sql.Open("duckdb", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open duckdb: %w", err)
	}

	store := &DuckDBStore{
		db:     db,
		config: config,
	}

	// Initialize VSS extension and schema
	if err := store.init(embeddingDim); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize duckdb store: %w", err)
	}

	return store, nil
}

// Close closes the database connection
func (s *DuckDBStore) Close() error {
	return s.db.Close()
}

// init initializes the database schema and extensions
func (s *DuckDBStore) init(dim int) error {
	// 1. Install and Load VSS extension
	// We try to load first, if it fails, we try to install
	if _, err := s.db.Exec("LOAD vss"); err != nil {
		xlog.Info("VSS extension not loaded, attempting to install...")
		if _, err := s.db.Exec("INSTALL vss"); err != nil {
			return fmt.Errorf("failed to install vss extension: %w", err)
		}
		if _, err := s.db.Exec("LOAD vss"); err != nil {
			return fmt.Errorf("failed to load vss extension after install: %w", err)
		}
	}

	// 2. Checkpoint WAL to avoid replay issues with old schema
	if _, err := s.db.Exec("CHECKPOINT"); err != nil {
		xlog.Warn("Failed to checkpoint WAL", "error", err)
	}

	// 3. Enable experimental persistence (required for disk-backed DB)
	if _, err := s.db.Exec("SET hnsw_enable_experimental_persistence = true"); err != nil {
		xlog.Warn("Failed to enable experimental persistence (might be in-memory DB)", "error", err)
	}

	// 4. Create vectors table
	// We use FLOAT[] for embeddings (32-bit float, required by HNSW index)
	// Note: Go's float64 will be converted to float32 when inserting

	// Drop existing table if it exists (for schema migration from DOUBLE to FLOAT)
	if _, err := s.db.Exec("DROP TABLE IF EXISTS vectors"); err != nil {
		xlog.Warn("Failed to drop existing vectors table", "error", err)
	}

	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS vectors (
			id TEXT PRIMARY KEY,
			file_id TEXT,
			embedding FLOAT[%d]
		)
	`, dim)
	xlog.Info("Creating vectors table", "query", query)
	if _, err := s.db.Exec(query); err != nil {
		return fmt.Errorf("failed to create vectors table: %w", err)
	}

	// 4. Create HNSW index with configuration
	// Reference: https://duckdb.org/docs/stable/core_extensions/vss
	indexQuery := fmt.Sprintf(`
		CREATE INDEX IF NOT EXISTS vec_idx ON vectors 
		USING HNSW (embedding) 
		WITH (
			metric = '%s',
			ef_construction = %d,
			ef_search = %d,
			M = %d
		)
	`, s.config.Metric, s.config.EfConstruction, s.config.EfSearch, s.config.M)

	xlog.Info("Creating HNSW index", "query", indexQuery, "config", s.config)
	if _, err := s.db.Exec(indexQuery); err != nil {
		return fmt.Errorf("failed to create HNSW index: %w", err)
	}

	return nil
}

// UpsertVector inserts or updates a vector
func (s *DuckDBStore) UpsertVector(ctx context.Context, id string, fileID string, embedding []float32) error {
	// DuckDB supports INSERT OR REPLACE for simple ID replacement.
	// Note: go-duckdb doesn't support []float32 as parameter, so we build the array literal
	embeddingStr := "["
	for i, v := range embedding {
		if i > 0 {
			embeddingStr += ", "
		}
		embeddingStr += fmt.Sprintf("%f", v)
	}
	embeddingStr += "]"

	dim := len(embedding)
	query := fmt.Sprintf(`INSERT OR REPLACE INTO vectors (id, file_id, embedding) VALUES (?, ?, %s::FLOAT[%d])`, embeddingStr, dim)
	_, err := s.db.ExecContext(ctx, query, id, fileID)
	if err != nil {
		return fmt.Errorf("failed to upsert vector: %w", err)
	}

	return nil
}

// VectorSearchResult represents a search result with similarity score
type VectorSearchResult struct {
	ID         string
	Similarity float64
}

// DistanceMetric represents the distance metric to use for vector search
type DistanceMetric string

const (
	// DistanceEuclidean uses Euclidean distance (L2)
	DistanceEuclidean DistanceMetric = "euclidean"
	// DistanceCosine uses cosine distance (1 - cosine_similarity)
	DistanceCosine DistanceMetric = "cosine"
	// DistanceInnerProduct uses negative inner product
	DistanceInnerProduct DistanceMetric = "inner_product"
)

// SearchVectors searches for similar vectors and returns IDs with similarity scores
// Uses Euclidean distance (L2) by default, which is accelerated by HNSW index
// Reference: https://duckdb.org/2024/10/23/whats-new-in-the-vss-extension
func (s *DuckDBStore) SearchVectors(ctx context.Context, embedding []float32, limit int) ([]VectorSearchResult, error) {
	return s.SearchVectorsWithMetric(ctx, embedding, limit, DistanceEuclidean)
}

// SearchVectorsWithMetric searches for similar vectors using specified distance metric
// Supports: euclidean, cosine, inner_product
// All metrics are accelerated by HNSW index when available
// Reference: https://duckdb.org/2024/10/23/whats-new-in-the-vss-extension
func (s *DuckDBStore) SearchVectorsWithMetric(ctx context.Context, embedding []float32, limit int, metric DistanceMetric) ([]VectorSearchResult, error) {
	// Note: go-duckdb doesn't support Array types ([]float32, []float64, []int) as query parameters
	// We build the array literal as a string instead
	embeddingStr := "["
	for i, v := range embedding {
		if i > 0 {
			embeddingStr += ", "
		}
		embeddingStr += fmt.Sprintf("%f", v)
	}
	embeddingStr += "]"

	// Get embedding dimension for type cast
	dim := len(embedding)

	// Choose distance function based on metric
	var distanceFunc string
	switch metric {
	case DistanceCosine:
		distanceFunc = "array_cosine_distance"
	case DistanceInnerProduct:
		distanceFunc = "array_negative_inner_product"
	case DistanceEuclidean:
		fallthrough
	default:
		distanceFunc = "array_distance"
	}

	query := fmt.Sprintf(`
		SELECT id, %s(embedding, %s::FLOAT[%d]) as distance
		FROM vectors 
		ORDER BY distance ASC 
		LIMIT ?
	`, distanceFunc, embeddingStr, dim)

	rows, err := s.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search vectors: %w", err)
	}
	defer rows.Close()

	var results []VectorSearchResult
	for rows.Next() {
		var id string
		var distance float64
		if err := rows.Scan(&id, &distance); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		// Convert distance to similarity (1 / (1 + distance))
		// This normalization works well for all distance metrics
		similarity := 1.0 / (1.0 + distance)
		results = append(results, VectorSearchResult{
			ID:         id,
			Similarity: similarity,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return results, nil
}

// SearchVectorsIDs is a convenience method that returns only IDs (for backward compatibility)
func (s *DuckDBStore) SearchVectorsIDs(ctx context.Context, embedding []float32, limit int) ([]string, error) {
	results, err := s.SearchVectors(ctx, embedding, limit)
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(results))
	for i, r := range results {
		ids[i] = r.ID
	}
	return ids, nil
}

// BatchSearchRequest represents a single search request in a batch
type BatchSearchRequest struct {
	QueryID   string
	Embedding []float32
}

// BatchSearchResult represents results for a single query in a batch
type BatchSearchResult struct {
	QueryID string
	Results []VectorSearchResult
}

// BatchSearchVectors performs multiple vector searches in a single query using LATERAL joins
// This is significantly faster than running individual searches (up to 66× speedup)
// Reference: https://duckdb.org/2024/10/23/whats-new-in-the-vss-extension
//
// Example: Search 1000 queries against 10000 vectors with limit=5
// - Individual searches: ~10 seconds
// - Batch search: ~0.15 seconds (66× faster!)
func (s *DuckDBStore) BatchSearchVectors(ctx context.Context, queries []BatchSearchRequest, limit int) ([]BatchSearchResult, error) {
	if len(queries) == 0 {
		return []BatchSearchResult{}, nil
	}

	// Create a temporary table for query embeddings
	// Note: We use a CTE (Common Table Expression) to avoid creating actual tables
	dim := len(queries[0].Embedding)

	// Build VALUES clause for all query embeddings
	var valuesBuilder strings.Builder
	valuesBuilder.WriteString("WITH queries AS (SELECT * FROM (VALUES ")

	for i, q := range queries {
		if i > 0 {
			valuesBuilder.WriteString(", ")
		}

		// Build embedding array literal
		embeddingStr := "["
		for j, v := range q.Embedding {
			if j > 0 {
				embeddingStr += ", "
			}
			embeddingStr += fmt.Sprintf("%f", v)
		}
		embeddingStr += "]"

		valuesBuilder.WriteString(fmt.Sprintf("('%s', %s::FLOAT[%d])", q.QueryID, embeddingStr, dim))
	}

	valuesBuilder.WriteString(") AS t(query_id, embedding)) ")

	// Build LATERAL join query
	// This will use HNSW index for each query efficiently
	query := valuesBuilder.String() + fmt.Sprintf(`
		SELECT 
			queries.query_id,
			items.id,
			items.distance
		FROM queries, LATERAL (
			SELECT
				vectors.id,
				array_distance(queries.embedding, vectors.embedding) AS distance
			FROM vectors
			ORDER BY distance ASC
			LIMIT %d
		) AS items
		ORDER BY queries.query_id, items.distance
	`, limit)

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute batch search: %w", err)
	}
	defer rows.Close()

	// Group results by query_id
	resultsMap := make(map[string][]VectorSearchResult)
	for rows.Next() {
		var queryID, vectorID string
		var distance float64
		if err := rows.Scan(&queryID, &vectorID, &distance); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		similarity := 1.0 / (1.0 + distance)
		resultsMap[queryID] = append(resultsMap[queryID], VectorSearchResult{
			ID:         vectorID,
			Similarity: similarity,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	// Convert map to slice, maintaining query order
	results := make([]BatchSearchResult, 0, len(queries))
	for _, q := range queries {
		results = append(results, BatchSearchResult{
			QueryID: q.QueryID,
			Results: resultsMap[q.QueryID],
		})
	}

	return results, nil
}

// SetEfSearch overrides the ef_search parameter at runtime for this connection
// Higher values = better recall but slower search
// Reference: https://duckdb.org/docs/stable/core_extensions/vss#index-options
func (s *DuckDBStore) SetEfSearch(efSearch int) error {
	query := fmt.Sprintf("SET hnsw_ef_search = %d", efSearch)
	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to set ef_search: %w", err)
	}
	xlog.Info("Updated ef_search parameter", "ef_search", efSearch)
	return nil
}

// ResetEfSearch resets the ef_search parameter to default (from index creation)
func (s *DuckDBStore) ResetEfSearch() error {
	_, err := s.db.Exec("RESET hnsw_ef_search")
	if err != nil {
		return fmt.Errorf("failed to reset ef_search: %w", err)
	}
	xlog.Info("Reset ef_search to default")
	return nil
}

// CompactIndex triggers re-compaction of the HNSW index to prune deleted items
// Call this after significant number of updates/deletes to maintain query quality
// Reference: https://duckdb.org/docs/stable/core_extensions/vss#inserts-updates-deletes-and-re-compaction
func (s *DuckDBStore) CompactIndex(indexName string) error {
	query := fmt.Sprintf("PRAGMA hnsw_compact_index('%s')", indexName)
	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to compact index: %w", err)
	}
	xlog.Info("Compacted HNSW index", "index_name", indexName)
	return nil
}

// GetIndexStats returns statistics about the HNSW index
func (s *DuckDBStore) GetIndexStats() (map[string]interface{}, error) {
	// Query to get index information
	query := `
		SELECT 
			COUNT(*) as total_vectors,
			COUNT(CASE WHEN id IS NOT NULL THEN 1 END) as active_vectors
		FROM vectors
	`

	row := s.db.QueryRow(query)
	var totalVectors, activeVectors int64
	if err := row.Scan(&totalVectors, &activeVectors); err != nil {
		return nil, fmt.Errorf("failed to get index stats: %w", err)
	}

	stats := map[string]interface{}{
		"total_vectors":  totalVectors,
		"active_vectors": activeVectors,
		"deleted_ratio":  float64(totalVectors-activeVectors) / float64(totalVectors),
		"config":         s.config,
	}

	return stats, nil
}

// DeleteVector deletes a vector by ID
func (s *DuckDBStore) DeleteVector(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM vectors WHERE id = ?", id)
	return err
}

// DeleteVectorsByFileID deletes all vectors for a file
func (s *DuckDBStore) DeleteVectorsByFileID(ctx context.Context, fileID string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM vectors WHERE file_id = ?", fileID)
	return err
}
