package database

import (
	"context"
	"fmt"
	"log"
	"os"
)

// SeedAvailableAgents loads pre-seeded agents from SQL dump
func (s *Service) SeedAvailableAgents(ctx context.Context) error {
	// 1. Check if we already have agents
	var count int
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM agents WHERE id != 'default-inbox-agent'").Scan(&count)
	if err != nil {
		log.Printf("Failed to check existing agents: %v", err)
	}

	// Skip if we already have agents (more than just inbox)
	if count > 0 {
		log.Printf("Agents already seeded (%d agents found). Skipping.", count)
		return nil
	}

	log.Println("Loading pre-seeded agents from SQL dump...")

	// 2. Load SQL dump from embedded file
	seedDataPath := "internal/database/seed_data.sql"
	seedData, err := os.ReadFile(seedDataPath)
	if err != nil {
		log.Printf("Failed to read seed data file: %v. Skipping agent seeding.", err)
		log.Printf("To generate seed data, run: go run cmd/seed-agents/main.go")
		return nil // Don't crash the app
	}

	// 3. Execute SQL dump
	_, err = s.db.ExecContext(ctx, string(seedData))
	if err != nil {
		return fmt.Errorf("failed to execute seed data: %w", err)
	}

	// 4. Count inserted agents
	err = s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM agents WHERE id != 'default-inbox-agent'").Scan(&count)
	if err == nil {
		log.Printf("✅ Successfully loaded %d agents from seed data", count)
	}

	return nil
}
