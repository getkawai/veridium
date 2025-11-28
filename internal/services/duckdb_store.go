package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kawai-network/veridium/pkg/xlog"
	_ "github.com/marcboeker/go-duckdb"
)

// DuckDBStore handles vector storage using DuckDB
type DuckDBStore struct {
	db *sql.DB
}

// NewDuckDBStore creates a new DuckDB store
func NewDuckDBStore(path string, embeddingDim int) (*DuckDBStore, error) {
	// Open DuckDB connection
	// If path is empty, it uses in-memory DB (useful for testing)
	dsn := path
	if dsn == "" {
		dsn = ":memory:"
	}

	db, err := sql.Open("duckdb", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open duckdb: %w", err)
	}

	store := &DuckDBStore{db: db}

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

	// 2. Enable experimental persistence (required for disk-backed DB)
	if _, err := s.db.Exec("SET hnsw_enable_experimental_persistence = true"); err != nil {
		xlog.Warn("Failed to enable experimental persistence (might be in-memory DB)", "error", err)
	}

	// 3. Create vectors table
	// We use DOUBLE[] for embeddings (64-bit float, consistent with Eino)
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS vectors (
			id TEXT PRIMARY KEY,
			file_id TEXT,
			embedding DOUBLE[%d]
		)
	`, dim)
	if _, err := s.db.Exec(query); err != nil {
		return fmt.Errorf("failed to create vectors table: %w", err)
	}

	// 4. Create HNSW index
	// We check if index exists by trying to create it with IF NOT EXISTS (if supported)
	// or catching the error. DuckDB supports CREATE INDEX IF NOT EXISTS.
	if _, err := s.db.Exec("CREATE INDEX IF NOT EXISTS vec_idx ON vectors USING HNSW (embedding)"); err != nil {
		return fmt.Errorf("failed to create HNSW index: %w", err)
	}

	return nil
}

// UpsertVector inserts or updates a vector
func (s *DuckDBStore) UpsertVector(ctx context.Context, id string, fileID string, embedding []float64) error {
	// DuckDB supports INSERT OR REPLACE for simple ID replacement.
	// go-duckdb driver supports passing []float64 for DOUBLE[] columns.

	query := `INSERT OR REPLACE INTO vectors (id, file_id, embedding) VALUES (?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, id, fileID, embedding)
	if err != nil {
		return fmt.Errorf("failed to upsert vector: %w", err)
	}

	return nil
}

// SearchVectors searches for similar vectors
func (s *DuckDBStore) SearchVectors(ctx context.Context, embedding []float64, limit int) ([]string, error) {
	// Query using array_distance (Euclidean distance by default for HNSW)
	// We want the closest ones, so ORDER BY distance ASC
	query := `
		SELECT id 
		FROM vectors 
		ORDER BY array_distance(embedding, ?) ASC 
		LIMIT ?
	`

	rows, err := s.db.QueryContext(ctx, query, embedding, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search vectors: %w", err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return ids, nil
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
