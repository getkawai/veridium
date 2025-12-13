package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	db "github.com/kawai-network/veridium/internal/database/generated"
)

const (
	// RemoteAgentIndexURL is the URL to fetch the agent index from
	RemoteAgentIndexURL = "https://registry.npmmirror.com/@lobehub/agents-index/1.42.0/files/public/index.json"
)

// AgentMeta matches the JSON structure of the remote index
type AgentMeta struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Avatar      string   `json:"avatar"`
}

// AgentIndexItem matches the JSON structure of the remote index
type AgentIndexItem struct {
	Identifier string    `json:"identifier"`
	Category   string    `json:"category"`
	Author     string    `json:"author"`
	Meta       AgentMeta `json:"meta"`
	CreatedAt  string    `json:"createdAt"`
	Homepage   string    `json:"homepage"`
	SystemRole string    `json:"systemRole"`
}

type AgentIndexResponse struct {
	Agents []AgentIndexItem `json:"agents"`
}

// SeedAvailableAgents fetches agents from remote and seeds them into the database
func (s *Service) SeedAvailableAgents(ctx context.Context) error {
	// 1. Check if we already have a significant number of agents
	var count int
	// We no longer filter by user_id
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM agents").Scan(&count)
	if err != nil {
		log.Printf("Failed to check existing agents: %v", err)
	}

	// If we have > 100 agents, assume already seeded.
	if count > 100 {
		log.Println("Agents already seeded (count > 100). Skipping remote fetch.")
		return nil
	}

	// Fetch remote JSON
	log.Printf("Seeding agents from %s...", RemoteAgentIndexURL)
	resp, err := http.Get(RemoteAgentIndexURL)
	if err != nil {
		log.Printf("Failed to fetch agent index: %v. Skipping seeding.", err)
		return nil // Don't crash the app, just skip seeding
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to fetch agent index: status %d. Skipping.", resp.StatusCode)
		return nil
	}

	var indexResponse AgentIndexResponse
	if err := json.NewDecoder(resp.Body).Decode(&indexResponse); err != nil {
		return fmt.Errorf("failed to decode agent index: %w", err)
	}

	log.Printf("Found %d agents to seed.", len(indexResponse.Agents))

	now := time.Now().UnixMilli()
	count = 0

	// Use a transaction for batch insert
	return s.WithTx(ctx, func(q *db.Queries) error {
		for _, agent := range indexResponse.Agents {
			// Skip invalid agents
			if agent.Identifier == "" {
				continue
			}

			// Prepare tags JSON
			tagsJson, _ := json.Marshal(agent.Meta.Tags)

			// Prepare default config matchin `createDefaultInboxSession`
			// We use a safe default config
			defaultAgentConfig := `{"autoCreateTopicThreshold":2,"displayMode":"chat","enableAutoCreateTopic":true,"enableCompressHistory":true,"enableHistoryCount":true,"enableReasoning":false,"enableStreaming":true,"historyCount":20,"reasoningBudgetToken":1024,"searchFCModel":{"model":"kawai-auto","provider":"kawai"},"searchMode":"off"}`
			defaultParams := `{"frequency_penalty":0,"presence_penalty":0,"temperature":0.6,"top_p":1}`

			// 1. Insert Agent
			// Note: SystemRole is missing from index, so we might leave it empty
			// or use description as a placeholder.
			_, err := q.CreateAgent(ctx, db.CreateAgentParams{
				ID:               agent.Identifier,
				Slug:             sql.NullString{String: agent.Identifier, Valid: true},
				Title:            sql.NullString{String: agent.Meta.Title, Valid: true},
				Description:      sql.NullString{String: agent.Meta.Description, Valid: true},
				Tags:             sql.NullString{String: string(tagsJson), Valid: true},
				Avatar:           sql.NullString{String: agent.Meta.Avatar, Valid: true},
				BackgroundColor:  sql.NullString{Valid: false},
				Plugins:          sql.NullString{String: "[]", Valid: true},
				ChatConfig:       sql.NullString{String: defaultAgentConfig, Valid: true},
				FewShots:         sql.NullString{Valid: false},
				Model:            sql.NullString{String: "kawai-auto", Valid: true}, // Default to our local model
				Params:           sql.NullString{String: defaultParams, Valid: true},
				Provider:         sql.NullString{String: "kawai", Valid: true},
				SystemRole:       sql.NullString{String: agent.Meta.Description, Valid: true}, // Fallback: use description
				Tts:              sql.NullString{Valid: false},
				Virtual:          0,
				OpeningMessage:   sql.NullString{Valid: false},
				OpeningQuestions: sql.NullString{String: "[]", Valid: true},
				CreatedAt:        now,
				UpdatedAt:        now,
			})

			if err != nil {
				// Assume duplicate, skip
				continue
			}

			// 2. Create Session for the agent
			_, err = q.CreateSession(ctx, db.CreateSessionParams{
				ID:              agent.Identifier,
				Slug:            agent.Identifier, // Session slug often matches agent slug
				Title:           sql.NullString{String: agent.Meta.Title, Valid: true},
				Description:     sql.NullString{String: agent.Meta.Description, Valid: true},
				Avatar:          sql.NullString{String: agent.Meta.Avatar, Valid: true},
				BackgroundColor: sql.NullString{Valid: false},
				Type:            sql.NullString{String: "agent", Valid: true},
				GroupID:         sql.NullString{Valid: false}, // Default group
				Pinned:          0,
				CreatedAt:       now,
				UpdatedAt:       now,
			})

			if err != nil {
				log.Printf("Failed to create session for agent %s: %v", agent.Identifier, err)
				continue
			}

			// 3. Link Agent to Session
			err = q.LinkAgentToSession(ctx, db.LinkAgentToSessionParams{
				AgentID:   agent.Identifier,
				SessionID: agent.Identifier,
			})

			if err != nil {
				log.Printf("Failed to link agent %s to session: %v", agent.Identifier, err)
				continue
			}

			count++
		}
		log.Printf("Successfully seeded %d new agents and sessions.", count)
		return nil
	})
}
