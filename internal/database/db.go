package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"

	"github.com/kawai-network/veridium/internal/database/generated"
)

// Service provides database operations
type Service struct {
	db      *sql.DB
	queries *db.Queries
}

// NewService creates a new database service
func NewService() (*Service, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		userConfigDir = "."
	}
	appDataDir := filepath.Join(userConfigDir, "veridium")

	if err := os.MkdirAll(appDataDir, 0o755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(appDataDir, "veridium.db")
	database, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Enable foreign keys
	if _, err := database.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, err
	}

	// Enable WAL mode for better concurrency
	if _, err := database.Exec("PRAGMA journal_mode = WAL"); err != nil {
		return nil, err
	}

	queries := db.New(database)

	return &Service{
		db:      database,
		queries: queries,
	}, nil
}

// Close closes the database connection
func (s *Service) Close() error {
	return s.db.Close()
}

// Queries returns the generated queries interface
func (s *Service) Queries() *db.Queries {
	return s.queries
}

// DB returns the underlying database connection
func (s *Service) DB() *sql.DB {
	return s.db
}

// WithTx executes a function within a transaction
func (s *Service) WithTx(ctx context.Context, fn func(*db.Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	qtx := s.queries.WithTx(tx)

	if err := fn(qtx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

