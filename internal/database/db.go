package database

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"

	db "github.com/kawai-network/veridium/internal/database/generated"
)

//go:embed schema/schema.sql
var schemaSQL string

// Service provides database operations
type Service struct {
	db      *sql.DB
	queries *db.Queries
}

// NewService creates a new database service using default path (./data/veridium.db)
func NewService() (*Service, error) {
	return NewServiceWithPath("")
}

// NewServiceWithPath creates a new database service with custom database path
// If dbPath is empty, uses default path (./data/veridium.db)
func NewServiceWithPath(dbPath string) (*Service, error) {
	if dbPath == "" {
		// Use project directory for database storage
		appDataDir := "./data"
		if err := os.MkdirAll(appDataDir, 0o755); err != nil {
			return nil, err
		}
		dbPath = filepath.Join(appDataDir, "veridium.db")
	} else {
		// Ensure parent directory exists for custom path
		if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
			return nil, err
		}
	}

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

	// Initialize schema if needed (check if sessions table exists - users table is gone)
	var tableExists int
	err = database.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='sessions'").Scan(&tableExists)
	if err != nil {
		return nil, fmt.Errorf("failed to check schema: %w", err)
	}

	if tableExists == 0 {
		// Schema doesn't exist, initialize it
		fmt.Println("Initializing database schema...")
		if _, err := database.Exec(schemaSQL); err != nil {
			return nil, fmt.Errorf("failed to initialize schema: %w", err)
		}
		fmt.Println("✅ Database schema initialized successfully")
	} else {
		fmt.Println("✅ Database schema already initialized")
	}

	queries := db.New(database)

	service := &Service{
		db:      database,
		queries: queries,
	}

	// Ensure default user settings and inbox session exist (for desktop single-user app)
	if err := service.ensureDefaultData(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ensure default data: %w", err)
	}

	return service, nil
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

// CreateMessageWithRelationsParams contains all data needed to create a message with its relations
type CreateMessageWithRelationsParams struct {
	Message    db.CreateMessageParams
	Plugin     *db.CreateMessagePluginParams
	FileIds    []string
	FileChunks []struct {
		ChunkId    string
		QueryId    string
		Similarity sql.NullInt64
	}
}

// CreateMessageWithRelations creates a message and all its related data in a single transaction
func (s *Service) CreateMessageWithRelations(ctx context.Context, params CreateMessageWithRelationsParams) (db.Message, error) {
	var result db.Message

	err := s.WithTx(ctx, func(q *db.Queries) error {
		// 1. Create message
		msg, err := q.CreateMessage(ctx, params.Message)
		if err != nil {
			return fmt.Errorf("failed to create message: %w", err)
		}
		result = msg

		// 2. Create plugin if needed
		if params.Plugin != nil {
			if _, err := q.CreateMessagePlugin(ctx, *params.Plugin); err != nil {
				return fmt.Errorf("failed to create message plugin: %w", err)
			}
		}

		// 3. Link files
		for _, fileId := range params.FileIds {
			if err := q.LinkMessageToFile(ctx, db.LinkMessageToFileParams{
				FileID:    fileId,
				MessageID: msg.ID,
			}); err != nil {
				return fmt.Errorf("failed to link file %s: %w", fileId, err)
			}
		}

		// 4. Link file chunks
		for _, chunk := range params.FileChunks {
			if err := q.LinkMessageQueryToChunk(ctx, db.LinkMessageQueryToChunkParams{
				MessageID:  sql.NullString{String: msg.ID, Valid: true},
				QueryID:    sql.NullString{String: chunk.QueryId, Valid: true},
				ChunkID:    sql.NullString{String: chunk.ChunkId, Valid: true},
				Similarity: chunk.Similarity,
			}); err != nil {
				return fmt.Errorf("failed to link chunk %s: %w", chunk.ChunkId, err)
			}
		}

		return nil
	})

	return result, err
}

// UpdateMessageWithImagesParams contains data for updating a message with images
type UpdateMessageWithImagesParams struct {
	MessageId string
	Message   db.UpdateMessageParams
	ImageIds  []string
}

// UpdateMessageWithImages updates a message and links images in a single transaction
func (s *Service) UpdateMessageWithImages(ctx context.Context, params UpdateMessageWithImagesParams) error {
	return s.WithTx(ctx, func(q *db.Queries) error {
		// 1. Update message
		if _, err := q.UpdateMessage(ctx, params.Message); err != nil {
			return fmt.Errorf("failed to update message: %w", err)
		}

		// 2. Link images
		for _, imageId := range params.ImageIds {
			if err := q.LinkMessageToFile(ctx, db.LinkMessageToFileParams{
				FileID:    imageId,
				MessageID: params.MessageId,
			}); err != nil {
				return fmt.Errorf("failed to link image %s: %w", imageId, err)
			}
		}

		return nil
	})
}

// DeleteMessageWithRelatedParams contains IDs for batch deletion
type DeleteMessageWithRelatedParams struct {
	MessageIds []string
}

// DeleteMessageWithRelated deletes a message and its related tool messages in a transaction
func (s *Service) DeleteMessageWithRelated(ctx context.Context, toolCallIdsJson string, messageIds []string) error {
	return s.WithTx(ctx, func(q *db.Queries) error {
		// Get related tool messages if tool call IDs provided
		if toolCallIdsJson != "" && toolCallIdsJson != "[]" {
			var toolCallIds []string
			if err := json.Unmarshal([]byte(toolCallIdsJson), &toolCallIds); err == nil {
				// Fetch each tool message (batch operation in Go)
				for _, toolCallId := range toolCallIds {
					// GetMessageByToolCallId no longer needs UserID
					msgId, err := q.GetMessageByToolCallId(ctx, sql.NullString{String: toolCallId, Valid: true})
					if err == nil {
						messageIds = append(messageIds, msgId)
					}
					// Ignore not found errors
				}
			}
		}

		// Delete all messages
		// BatchDeleteMessages no longer takes UserID, just the slice of IDs
		if err := q.BatchDeleteMessages(ctx, messageIds); err != nil {
			return fmt.Errorf("failed to delete messages: %w", err)
		}

		return nil
	})
}

// GetMessagesByToolCallIds fetches messages by tool call IDs (batch operation)
func (s *Service) GetMessagesByToolCallIds(ctx context.Context, toolCallIdsJson string) ([]string, error) {
	var toolCallIds []string
	if err := json.Unmarshal([]byte(toolCallIdsJson), &toolCallIds); err != nil {
		return nil, fmt.Errorf("failed to parse tool call IDs: %w", err)
	}

	results := make([]string, 0, len(toolCallIds))
	for _, toolCallId := range toolCallIds {
		// GetMessageByToolCallId now expects just the toolCallID (or params if more complex, but we removed userID)
		// Based on messages.sql, it takes `tool_call_id`.
		msgId, err := s.queries.GetMessageByToolCallId(ctx, sql.NullString{String: toolCallId, Valid: true})
		if err == nil {
			results = append(results, msgId)
		}
		// Ignore not found errors
	}

	return results, nil
}

// GetDocumentsByFileIds fetches documents by file IDs (batch operation)
func (s *Service) GetDocumentsByFileIds(ctx context.Context, fileIdsJson string) ([]db.GetDocumentByFileIdRow, error) {
	var fileIds []string
	if err := json.Unmarshal([]byte(fileIdsJson), &fileIds); err != nil {
		return nil, fmt.Errorf("failed to parse file IDs: %w", err)
	}

	results := make([]db.GetDocumentByFileIdRow, 0, len(fileIds))
	for _, fileId := range fileIds {
		// GetDocumentByFileID now takes just fileID
		doc, err := s.queries.GetDocumentByFileId(ctx, sql.NullString{String: fileId, Valid: true})
		if err == nil {
			results = append(results, doc)
		}
		// Ignore not found errors
	}

	return results, nil
}

// ============================================================================
// FILE OPERATIONS WITH TRANSACTIONS
// ============================================================================

// CreateFileWithLinksParams contains all data needed to create a file with its links
type CreateFileWithLinksParams struct {
	File          db.CreateFileParams
	GlobalFile    *db.CreateGlobalFileParams
	KnowledgeBase *string // Knowledge base ID to link to
}

// CreateFileWithLinks creates a file, optionally creates a global file, and links to knowledge base
// All operations are atomic within a transaction
func (s *Service) CreateFileWithLinks(ctx context.Context, params CreateFileWithLinksParams) (*db.File, error) {
	var result *db.File

	err := s.WithTx(ctx, func(q *db.Queries) error {
		// 1. Create global file if provided
		if params.GlobalFile != nil {
			if _, err := q.CreateGlobalFile(ctx, *params.GlobalFile); err != nil {
				// Ignore if already exists
				if err.Error() != "UNIQUE constraint failed" {
					return fmt.Errorf("failed to create global file: %w", err)
				}
			}
		}

		// 2. Create file
		file, err := q.CreateFile(ctx, params.File)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}

		// 3. Link to knowledge base if provided
		if params.KnowledgeBase != nil {
			err = q.LinkKnowledgeBaseToFile(ctx, db.LinkKnowledgeBaseToFileParams{
				KnowledgeBaseID: *params.KnowledgeBase,
				FileID:          file.ID,
			})
			if err != nil {
				return fmt.Errorf("failed to link to knowledge base: %w", err)
			}
		}

		result = &file
		return nil
	})

	return result, err
}

// DeleteFileWithCascadeParams contains data needed to delete a file with all related data
type DeleteFileWithCascadeParams struct {
	FileID           string
	RemoveGlobalFile bool
	FileHash         string
}

// DeleteFileWithCascade deletes a file and all its related chunks and embeddings
// All operations are atomic within a transaction
func (s *Service) DeleteFileWithCascade(ctx context.Context, params DeleteFileWithCascadeParams) error {
	return s.WithTx(ctx, func(q *db.Queries) error {
		// 1. Get chunk IDs for this file
		chunkIds, err := q.GetFileChunkIds(ctx, sql.NullString{String: params.FileID, Valid: true})
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("failed to get chunk IDs: %w", err)
		}

		// 2. Delete chunks (embeddings are stored in DuckDB, not SQLite)
		for _, chunkId := range chunkIds {
			if chunkId.Valid {
				// DeleteChunk now just takes ID
				err := q.DeleteChunk(ctx, chunkId.String)
				if err != nil && err != sql.ErrNoRows {
					return fmt.Errorf("failed to delete chunk: %w", err)
				}
			}
		}

		// 4. Delete file record
		// DeleteFile now just takes ID
		err = q.DeleteFile(ctx, params.FileID)
		if err != nil {
			return fmt.Errorf("failed to delete file: %w", err)
		}

		// 5. Check if global file should be deleted
		if params.RemoveGlobalFile && params.FileHash != "" {
			countResult, err := q.CountFilesByHash(ctx, sql.NullString{String: params.FileHash, Valid: true})
			if err != nil {
				return fmt.Errorf("failed to count files by hash: %w", err)
			}

			// Delete global file if no other files use it
			if countResult == 0 {
				_ = q.DeleteGlobalFile(ctx, params.FileHash)
			}
		}

		return nil
	})
}

// ============================================================================
// AI PROVIDER OPERATIONS WITH TRANSACTIONS
// ============================================================================

// DeleteAIProviderWithModels deletes an AI provider and all its models atomically
func (s *Service) DeleteAIProviderWithModels(ctx context.Context, providerID string) error {
	return s.WithTx(ctx, func(q *db.Queries) error {
		// 1. Delete all models of the provider
		// DeleteModelsByProvider now just takes providerID
		err := q.DeleteModelsByProvider(ctx, providerID)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("failed to delete models: %w", err)
		}

		// 2. Delete the provider
		// DeleteAIProvider now just takes providerID
		err = q.DeleteAIProvider(ctx, providerID)
		if err != nil {
			return fmt.Errorf("failed to delete provider: %w", err)
		}

		return nil
	})
}

// BatchInsertAIModelsParams contains data for batch inserting AI models
type BatchInsertAIModelsParams struct {
	Models []db.CreateAIModelParams
}

// BatchInsertAIModels inserts multiple AI models atomically
// Ignores conflicts (models that already exist)
func (s *Service) BatchInsertAIModels(ctx context.Context, models []db.CreateAIModelParams) ([]db.AiModel, error) {
	var results []db.AiModel

	err := s.WithTx(ctx, func(q *db.Queries) error {
		for _, model := range models {
			result, err := q.CreateAIModel(ctx, model)
			if err != nil {
				// Ignore UNIQUE constraint failures
				continue
			}
			results = append(results, result)
		}
		return nil
	})

	return results, err
}

// ============================================================================
// DEFAULT DATA INITIALIZATION (Desktop Single-User App)
// ============================================================================

const defaultUserID = "DEFAULT_LOBE_CHAT_USER"
const DefaultInboxAgentID = "inbox"

// ensureDefaultData ensures the default settings and inbox session exist
// This is called during database initialization for desktop single-user apps
func (s *Service) ensureDefaultData(ctx context.Context) error {
	// 1. Ensure default user settings exist
	if err := s.ensureDefaultUserSettings(ctx); err != nil {
		return fmt.Errorf("failed to ensure default user settings: %w", err)
	}

	// 2. Ensure default AI provider (Kawai) exists
	if err := s.ensureDefaultAIProvider(ctx); err != nil {
		return fmt.Errorf("failed to ensure default AI provider: %w", err)
	}

	// 3. Check if inbox session already exists
	// GetSession no longer needs UserID, and we use fixed ID "inbox"
	_, err := s.queries.GetSession(ctx, "inbox")

	if err == sql.ErrNoRows {
		// Inbox doesn't exist, create it
		if err := s.createDefaultInboxSession(ctx); err != nil {
			return fmt.Errorf("failed to create inbox session: %w", err)
		}
		fmt.Println("✅ Default inbox session created")
	} else if err != nil {
		return fmt.Errorf("failed to check inbox session: %w", err)
	} else {
		// Inbox already exists
		fmt.Println("✅ Default inbox session already exists")
	}

	// 4. Seed available agents from SQL dump (fast, synchronous)
	seedCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := s.SeedAvailableAgents(seedCtx); err != nil {
		fmt.Printf("⚠️  Failed to seed agents: %v\n", err)
	}

	return nil
}

// createDefaultInboxSession creates the default inbox session with agent
func (s *Service) createDefaultInboxSession(ctx context.Context) error {
	return s.WithTx(ctx, func(q *db.Queries) error {
		now := int64(1000)
		sessionID := DefaultInboxAgentID
		agentID := DefaultInboxAgentID

		// 1. Create session
		_, err := q.CreateSession(ctx, db.CreateSessionParams{
			ID:              sessionID,
			Title:           sql.NullString{Valid: false},
			Description:     sql.NullString{Valid: false},
			Avatar:          sql.NullString{Valid: false},
			BackgroundColor: sql.NullString{Valid: false},
			Type:            sql.NullString{String: "agent", Valid: true},
			GroupID:         sql.NullString{Valid: false},
			Pinned:          0,
			CreatedAt:       now,
			UpdatedAt:       now,
		})
		if err != nil {
			return fmt.Errorf("failed to create session: %w", err)
		}

		// 2. Create default agent with kawai-auto model
		defaultAgentConfig := `{"autoCreateTopicThreshold":2,"displayMode":"chat","enableAutoCreateTopic":true,"enableCompressHistory":true,"enableHistoryCount":true,"enableReasoning":false,"enableStreaming":true,"historyCount":20,"reasoningBudgetToken":1024,"searchFCModel":{"model":"kawai-auto","provider":"kawai"},"searchMode":"off"}`
		defaultParams := `{"frequency_penalty":0,"presence_penalty":0,"temperature":1,"top_p":1}`

		_, err = q.CreateAgent(ctx, db.CreateAgentParams{
			ID:               agentID,
			Title:            sql.NullString{Valid: false},
			Description:      sql.NullString{Valid: false},
			Tags:             sql.NullString{String: "[]", Valid: true},
			Avatar:           sql.NullString{Valid: false},
			BackgroundColor:  sql.NullString{Valid: false},
			Plugins:          sql.NullString{String: "[]", Valid: true},
			ChatConfig:       sql.NullString{String: defaultAgentConfig, Valid: true},
			FewShots:         sql.NullString{Valid: false},
			Model:            sql.NullString{String: "kawai-auto", Valid: true},
			Params:           sql.NullString{String: defaultParams, Valid: true},
			Provider:         sql.NullString{String: "kawai", Valid: true},
			SystemRole:       sql.NullString{Valid: false},
			Tts:              sql.NullString{Valid: false},
			Virtual:          0,
			OpeningMessage:   sql.NullString{Valid: false},
			OpeningQuestions: sql.NullString{String: "[]", Valid: true},
			CreatedAt:        now,
			UpdatedAt:        now,
		})
		if err != nil {
			return fmt.Errorf("failed to create agent: %w", err)
		}

		// 3. Link agent to session
		err = q.LinkAgentToSession(ctx, db.LinkAgentToSessionParams{
			AgentID:   agentID,
			SessionID: sessionID,
		})
		if err != nil {
			return fmt.Errorf("failed to link agent to session: %w", err)
		}

		return nil
	})
}

// ensureDefaultUserSettings ensures default user settings exist
func (s *Service) ensureDefaultUserSettings(ctx context.Context) error {
	// Check if user settings already exist
	_, err := s.queries.GetUserSettings(ctx, defaultUserID)
	if err == nil {
		// Settings already exist
		fmt.Println("✅ Default user settings already exist")
		return nil
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("failed to check user settings: %w", err)
	}

	// Create default settings matching frontend defaults
	defaultGeneral := `{"fontSize":14,"animationMode":"agile","transitionMode":"fadeIn","highlighterTheme":"lobe-theme","mermaidTheme":"lobe-theme"}`
	defaultLanguageModel := `{"kawai":{"enabled":true}}`
	defaultHotkey := `{}`
	defaultTool := `{}`
	defaultImage := `{}`

	_, err = s.queries.UpsertUserSettings(ctx, db.UpsertUserSettingsParams{
		ID:            defaultUserID,
		General:       sql.NullString{String: defaultGeneral, Valid: true},
		LanguageModel: sql.NullString{String: defaultLanguageModel, Valid: true},
		Hotkey:        sql.NullString{String: defaultHotkey, Valid: true},
		Tool:          sql.NullString{String: defaultTool, Valid: true},
		Image:         sql.NullString{String: defaultImage, Valid: true},
		KeyVaults:     sql.NullString{String: "{}", Valid: true},
		Tts:           sql.NullString{String: "", Valid: false},
		SystemAgent:   sql.NullString{String: "", Valid: false},
		DefaultAgent:  sql.NullString{String: "", Valid: false},
	})
	if err != nil {
		return fmt.Errorf("failed to create default user settings: %w", err)
	}

	fmt.Println("✅ Default user settings created")
	return nil
}

// ensureDefaultAIProvider ensures the default Kawai AI provider exists
func (s *Service) ensureDefaultAIProvider(ctx context.Context) error {
	now := int64(1000) // Use a fixed timestamp for default entries

	// Check if Kawai provider already exists
	_, err := s.queries.GetAIProvider(ctx, "kawai")
	if err == nil {
		// Provider already exists, ensure model abilities are up to date
		fmt.Println("✅ Default Kawai AI provider already exists")
		if err := s.updateKawaiAutoModelAbilities(ctx); err != nil {
			return fmt.Errorf("failed to update kawai-auto model abilities: %w", err)
		}
		return nil
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("failed to check AI provider: %w", err)
	}

	// Create default Kawai provider
	providerSettings := `{"defaultShowBrowserRequest":true,"proxyUrl":{"placeholder":"http://127.0.0.1:8080/v1"},"responseAnimation":{"speed":2,"text":"smooth"},"showApiKey":false,"showModelFetcher":false}`
	providerConfig := `{}`

	_, err = s.queries.CreateAIProvider(ctx, db.CreateAIProviderParams{
		ID:            "kawai",
		Name:          sql.NullString{String: "Kawai", Valid: true},
		Sort:          sql.NullInt64{Int64: 0, Valid: true},
		Enabled:       sql.NullInt64{Int64: 1, Valid: true}, // enabled = true
		FetchOnClient: sql.NullInt64{Int64: 0, Valid: true}, // fetchOnClient = false
		CheckModel:    sql.NullString{String: "", Valid: false},
		Logo:          sql.NullString{String: "", Valid: false},
		Description:   sql.NullString{String: "Kawai AI - Local LLM inference powered by llama.cpp", Valid: true},
		KeyVaults:     sql.NullString{String: "", Valid: false},
		Source:        sql.NullString{String: "builtin", Valid: true},
		Settings:      sql.NullString{String: providerSettings, Valid: true},
		Config:        sql.NullString{String: providerConfig, Valid: true},
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	if err != nil {
		return fmt.Errorf("failed to create Kawai AI provider: %w", err)
	}

	// Create default kawai-auto model
	modelAbilities := `{"functionCall":true,"vision":true,"files":true}`
	modelParams := `{}`
	modelConfig := `{}`

	_, err = s.queries.CreateAIModel(ctx, db.CreateAIModelParams{
		ID:                  "kawai-auto",
		DisplayName:         sql.NullString{String: "Kawai Auto", Valid: true},
		Description:         sql.NullString{String: "Automatically selects the best local model for your hardware", Valid: true},
		Organization:        sql.NullString{String: "", Valid: false},
		Enabled:             sql.NullInt64{Int64: 1, Valid: true},
		ProviderID:          "kawai",
		Type:                "chat",
		Sort:                sql.NullInt64{Int64: 0, Valid: true},
		Pricing:             sql.NullString{String: "", Valid: false},
		Parameters:          sql.NullString{String: modelParams, Valid: true},
		Config:              sql.NullString{String: modelConfig, Valid: true},
		Abilities:           sql.NullString{String: modelAbilities, Valid: true},
		ContextWindowTokens: sql.NullInt64{Int64: 128000, Valid: true},
		Source:              sql.NullString{String: "builtin", Valid: true},
		ReleasedAt:          sql.NullString{String: "", Valid: false},
		CreatedAt:           now,
		UpdatedAt:           now,
	})
	if err != nil {
		return fmt.Errorf("failed to create kawai-auto model: %w", err)
	}

	fmt.Println("✅ Default Kawai AI provider and model created")
	return nil
}

// updateKawaiAutoModelAbilities updates the abilities of the kawai-auto model
// This is used to migrate existing databases to the latest model capabilities
func (s *Service) updateKawaiAutoModelAbilities(ctx context.Context) error {
	// Get current model
	model, err := s.queries.GetAIModel(ctx, db.GetAIModelParams{
		ID:         "kawai-auto",
		ProviderID: "kawai",
	})
	if err != nil {
		if err == sql.ErrNoRows {
			// Model doesn't exist, skip update
			return nil
		}
		return fmt.Errorf("failed to get kawai-auto model: %w", err)
	}

	// Check if abilities need updating
	currentAbilities := model.Abilities.String
	expectedAbilities := `{"functionCall":true,"vision":true,"files":true}`

	if currentAbilities == expectedAbilities {
		// Already up to date
		return nil
	}

	// Update model abilities
	_, err = s.queries.UpdateAIModel(ctx, db.UpdateAIModelParams{
		ID:         "kawai-auto",
		ProviderID: "kawai",
		Abilities:  sql.NullString{String: expectedAbilities, Valid: true},
		UpdatedAt:  int64(1000),
	})
	if err != nil {
		return fmt.Errorf("failed to update kawai-auto model abilities: %w", err)
	}

	fmt.Println("✅ Updated kawai-auto model abilities to include file support")
	return nil
}
