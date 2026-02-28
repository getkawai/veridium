package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	db "github.com/getkawai/database/db"
	_ "modernc.org/sqlite"
)

// toNullString helper
func toNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

const (
	RemoteAgentIndexURL          = "https://registry.npmmirror.com/@lobehub/agents-index/1.42.0/files/public/index.json"
	GitHubAgentDetailURLTemplate = "https://raw.githubusercontent.com/lobehub/lobe-chat-agents/main/locales/%s/index.json"
	MaxConcurrent                = 10 // Concurrent fetches
	FetchTimeout                 = 10 * time.Second
)

type AgentMeta struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Avatar      string   `json:"avatar"`
}

type AgentConfig struct {
	SystemRole       string   `json:"systemRole"`
	OpeningMessage   string   `json:"openingMessage"`
	OpeningQuestions []string `json:"openingQuestions"`
}

type AgentIndexItem struct {
	Identifier string    `json:"identifier"`
	Category   string    `json:"category"`
	Author     string    `json:"author"`
	Meta       AgentMeta `json:"meta"`
}

type AgentIndexResponse struct {
	Agents []AgentIndexItem `json:"agents"`
}

type AgentDetail struct {
	Config AgentConfig `json:"config"`
	Meta   AgentMeta   `json:"meta"`
}

type AgentData struct {
	Agent  AgentIndexItem
	Detail *AgentDetail
	Error  error
}

func fetchAgentDetails(identifier string) (*AgentDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), FetchTimeout)
	defer cancel()

	url := fmt.Sprintf(GitHubAgentDetailURLTemplate, identifier)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	var detail AgentDetail
	if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
		return nil, err
	}

	return &detail, nil
}

func main() {
	log.Println("🌱 Starting agent seeding process...")

	// 1. Fetch agent index
	log.Println("📥 Fetching agent index from NPM registry...")
	resp, err := http.Get(RemoteAgentIndexURL)
	if err != nil {
		log.Fatalf("Failed to fetch agent index: %v", err)
	}
	defer resp.Body.Close()

	var indexResponse AgentIndexResponse
	if err := json.NewDecoder(resp.Body).Decode(&indexResponse); err != nil {
		log.Fatalf("Failed to decode agent index: %v", err)
	}

	log.Printf("✅ Found %d agents to seed\n", len(indexResponse.Agents))

	// 2. Fetch details concurrently
	log.Println("🚀 Fetching agent details from GitHub (concurrent)...")

	agentChan := make(chan AgentIndexItem, len(indexResponse.Agents))
	resultChan := make(chan AgentData, len(indexResponse.Agents))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < MaxConcurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for agent := range agentChan {
				detail, err := fetchAgentDetails(agent.Identifier)
				resultChan <- AgentData{
					Agent:  agent,
					Detail: detail,
					Error:  err,
				}
			}
		}()
	}

	// Send agents to workers
	go func() {
		for _, agent := range indexResponse.Agents {
			if agent.Identifier != "" {
				agentChan <- agent
			}
		}
		close(agentChan)
	}()

	// Wait for workers
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	var agentDataList []AgentData
	processed := 0
	failed := 0
	for data := range resultChan {
		processed++
		if data.Error != nil {
			failed++
			if processed%50 == 0 {
				log.Printf("⏳ Progress: %d/%d (failed: %d)", processed, len(indexResponse.Agents), failed)
			}
		} else {
			if processed%50 == 0 {
				log.Printf("⏳ Progress: %d/%d (failed: %d)", processed, len(indexResponse.Agents), failed)
			}
		}
		agentDataList = append(agentDataList, data)
	}

	log.Printf("✅ Fetched %d agents (%d failed)\n", len(agentDataList)-failed, failed)

	// 3. Create database and insert
	dbPath := "data/seed_agents.db"
	os.Remove(dbPath) // Clean start

	log.Printf("💾 Creating database: %s\n", dbPath)
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer conn.Close()

	// Create schema from schema.sql file
	log.Println("📋 Loading schema from schema.sql...")
	schemaBytes, err := os.ReadFile("internal/database/schema/schema.sql")
	if err != nil {
		log.Fatalf("Failed to read schema file: %v", err)
	}

	if _, err := conn.Exec(string(schemaBytes)); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	// Initialize queries
	queries := db.New(conn)
	ctx := context.Background()

	// Insert agents
	log.Println("📝 Inserting agents into database...")
	now := time.Now().UnixMilli()
	defaultAgentConfig := `{"autoCreateTopicThreshold":2,"displayMode":"chat","enableAutoCreateTopic":true,"enableCompressHistory":true,"enableHistoryCount":true,"enableStreaming":true,"historyCount":20,"searchFCModel":{"model":"kawai-auto","provider":"kawai"},"searchMode":"off"}`
	defaultParams := `{"frequency_penalty":0,"presence_penalty":0,"temperature":0.6,"top_p":1}`

	inserted := 0
	for _, data := range agentDataList {
		agent := data.Agent

		// Prepare data
		tagsJSON, _ := json.Marshal(agent.Meta.Tags)
		systemRole := agent.Meta.Description // Fallback
		openingMessage := ""
		openingQuestions := "[]"

		if data.Detail != nil && data.Error == nil {
			if data.Detail.Config.SystemRole != "" {
				systemRole = data.Detail.Config.SystemRole
			}
			openingMessage = data.Detail.Config.OpeningMessage
			if len(data.Detail.Config.OpeningQuestions) > 0 {
				questionsJSON, _ := json.Marshal(data.Detail.Config.OpeningQuestions)
				openingQuestions = string(questionsJSON)
			}
		}

		// Insert agent
		_, err := queries.CreateAgent(ctx, db.CreateAgentParams{
			ID:               agent.Identifier,
			Title:            toNullString(agent.Meta.Title),
			Description:      toNullString(agent.Meta.Description),
			Tags:             toNullString(string(tagsJSON)),
			Avatar:           toNullString(agent.Meta.Avatar),
			BackgroundColor:  sql.NullString{Valid: false}, // NULL
			Plugins:          toNullString("[]"),
			ChatConfig:       toNullString(defaultAgentConfig),
			FewShots:         sql.NullString{Valid: false}, // NULL
			Model:            toNullString("kawai-auto"),
			Params:           toNullString(defaultParams),
			Provider:         toNullString("kawai"),
			SystemRole:       toNullString(systemRole),
			Tts:              sql.NullString{Valid: false}, // NULL
			Virtual:          0,
			OpeningMessage:   toNullString(openingMessage),
			OpeningQuestions: toNullString(openingQuestions),
			CreatedAt:        now,
			UpdatedAt:        now,
		})

		if err != nil {
			log.Printf("⚠️  Failed to insert agent %s: %v", agent.Identifier, err)
			continue
		}

		// Insert session
		_, err = queries.CreateSession(ctx, db.CreateSessionParams{
			ID:              agent.Identifier,
			Title:           toNullString(agent.Meta.Title),
			Description:     toNullString(agent.Meta.Description),
			Avatar:          toNullString(agent.Meta.Avatar),
			BackgroundColor: sql.NullString{Valid: false}, // NULL
			Type:            toNullString("agent"),
			GroupID:         toNullString("default"),
			Pinned:          0,
			CreatedAt:       now,
			UpdatedAt:       now,
		})

		if err != nil {
			log.Printf("⚠️  Failed to insert session for %s: %v", agent.Identifier, err)
			continue
		}

		// Link agent to session
		err = queries.LinkAgentToSession(ctx, db.LinkAgentToSessionParams{
			AgentID:   agent.Identifier,
			SessionID: agent.Identifier,
		})

		if err != nil {
			log.Printf("⚠️  Failed to link agent %s to session: %v", agent.Identifier, err)
		}

		inserted++
	}

	log.Printf("✅ Inserted %d agents and sessions\n", inserted)

	// 4. Generate SQL dump
	log.Println("📦 Generating SQL dump...")
	dumpPath := "internal/database/seed_data.sql"

	dumpFile, err := os.Create(dumpPath)
	if err != nil {
		log.Fatalf("Failed to create dump file: %v", err)
	}
	defer dumpFile.Close()

	// Write header
	fmt.Fprintln(dumpFile, "-- Auto-generated agent seed data")
	fmt.Fprintf(dumpFile, "-- Generated at: %s\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(dumpFile, "-- Total agents: %d\n\n", inserted)
	fmt.Fprintln(dumpFile, "PRAGMA foreign_keys=OFF;")
	fmt.Fprintln(dumpFile, "BEGIN TRANSACTION;")

	// Dump agents
	// Dump agents
	rows, err := conn.Query("SELECT id, title, description, tags, avatar, background_color, plugins, chat_config, few_shots, model, params, provider, system_role, tts, virtual, opening_message, opening_questions, created_at, updated_at FROM agents")
	if err != nil {
		log.Fatalf("Failed to query agents: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, title, description, tags, avatar, bgColor, plugins, chatConfig, fewShots, model, params, provider, systemRole, tts, openingMsg, openingQ sql.NullString
		var virtual int
		var createdAt, updatedAt int64

		err := rows.Scan(&id, &title, &description, &tags, &avatar, &bgColor, &plugins, &chatConfig, &fewShots, &model, &params, &provider, &systemRole, &tts, &virtual, &openingMsg, &openingQ, &createdAt, &updatedAt)
		if err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}

		fmt.Fprintf(dumpFile, "INSERT INTO agents (id, title, description, tags, avatar, background_color, plugins, chat_config, few_shots, model, params, provider, system_role, tts, virtual, opening_message, opening_questions, created_at, updated_at) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %d, %s, %s, %d, %d);\n",
			sqlQuote(id), sqlQuote(title), sqlQuote(description), sqlQuote(tags), sqlQuote(avatar), sqlQuote(bgColor), sqlQuote(plugins), sqlQuote(chatConfig), sqlQuote(fewShots), sqlQuote(model), sqlQuote(params), sqlQuote(provider), sqlQuote(systemRole), sqlQuote(tts), virtual, sqlQuote(openingMsg), sqlQuote(openingQ), createdAt, updatedAt)

	}

	// Dump sessions
	sessionRows, err := conn.Query("SELECT id, title, description, avatar, background_color, type, group_id, pinned, created_at, updated_at FROM sessions")
	if err != nil {
		log.Fatalf("Failed to query sessions: %v", err)
	}
	defer sessionRows.Close()

	for sessionRows.Next() {
		var id, title, description, avatar, bgColor, typ, groupID sql.NullString
		var pinned int
		var createdAt, updatedAt int64

		err := sessionRows.Scan(&id, &title, &description, &avatar, &bgColor, &typ, &groupID, &pinned, &createdAt, &updatedAt)
		if err != nil {
			log.Printf("Failed to scan session row: %v", err)
			continue
		}

		fmt.Fprintf(dumpFile, "INSERT INTO sessions (id, title, description, avatar, background_color, type, group_id, pinned, created_at, updated_at) VALUES (%s, %s, %s, %s, %s, %s, %s, %d, %d, %d);\n",
			sqlQuote(id), sqlQuote(title), sqlQuote(description), sqlQuote(avatar), sqlQuote(bgColor), sqlQuote(typ), sqlQuote(groupID), pinned, createdAt, updatedAt)
	}

	// Dump agents_to_sessions
	linkRows, err := conn.Query("SELECT agent_id, session_id FROM agents_to_sessions")
	if err != nil {
		log.Fatalf("Failed to query agents_to_sessions: %v", err)
	}
	defer linkRows.Close()

	for linkRows.Next() {
		var agentID, sessionID string
		err := linkRows.Scan(&agentID, &sessionID)
		if err != nil {
			log.Printf("Failed to scan link row: %v", err)
			continue
		}

		fmt.Fprintf(dumpFile, "INSERT INTO agents_to_sessions (agent_id, session_id) VALUES ('%s', '%s');\n", agentID, sessionID)
	}

	fmt.Fprintln(dumpFile, "COMMIT;")
	fmt.Fprintln(dumpFile, "PRAGMA foreign_keys=ON;")

	log.Printf("✅ SQL dump saved to: %s\n", dumpPath)
	log.Println("🎉 Seeding complete!")
	log.Printf("\n💡 Next steps:\n")
	log.Printf("   1. Review the dump file: %s\n", dumpPath)
	log.Printf("   2. Restart the app - it will load from the dump\n")
	log.Printf("   3. Commit the dump file to git\n")
}

func sqlQuote(ns sql.NullString) string {
	if !ns.Valid {
		return "NULL"
	}
	// Escape single quotes
	escaped := ""
	for _, c := range ns.String {
		if c == '\'' {
			escaped += "''"
		} else {
			escaped += string(c)
		}
	}
	return "'" + escaped + "'"
}
