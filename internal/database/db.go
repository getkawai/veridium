package database

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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

// NewService creates a new database service
func NewService() (*Service, error) {
	// Use project directory for database storage
	appDataDir := "./data"

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

	// Initialize schema if needed (check if users table exists)
	var tableExists int
	err = database.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='users'").Scan(&tableExists)
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

	// Ensure default user and inbox session exist (for desktop single-user app)
	if err := service.ensureDefaultUserAndInbox(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ensure default user and inbox: %w", err)
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
func (s *Service) CreateMessageWithRelations(ctx context.Context, params CreateMessageWithRelationsParams, userId string) (db.Message, error) {
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
				UserID:    userId,
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
				UserID:     userId,
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
func (s *Service) UpdateMessageWithImages(ctx context.Context, params UpdateMessageWithImagesParams, userId string) error {
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
				UserID:    userId,
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
	UserId     string
}

// DeleteMessageWithRelated deletes a message and its related tool messages in a transaction
func (s *Service) DeleteMessageWithRelated(ctx context.Context, toolCallIdsJson string, messageIds []string, userId string) error {
	return s.WithTx(ctx, func(q *db.Queries) error {
		// Get related tool messages if tool call IDs provided
		if toolCallIdsJson != "" && toolCallIdsJson != "[]" {
			var toolCallIds []string
			if err := json.Unmarshal([]byte(toolCallIdsJson), &toolCallIds); err == nil {
				// Fetch each tool message (batch operation in Go)
				for _, toolCallId := range toolCallIds {
					msgId, err := q.GetMessageByToolCallId(ctx, db.GetMessageByToolCallIdParams{
						ToolCallID: sql.NullString{String: toolCallId, Valid: true},
						UserID:     userId,
					})
					if err == nil {
						messageIds = append(messageIds, msgId)
					}
					// Ignore not found errors
				}
			}
		}

		// Delete all messages
		if err := q.BatchDeleteMessages(ctx, db.BatchDeleteMessagesParams{
			UserID: userId,
			Ids:    messageIds,
		}); err != nil {
			return fmt.Errorf("failed to delete messages: %w", err)
		}

		return nil
	})
}

// GetMessagesByToolCallIds fetches messages by tool call IDs (batch operation)
func (s *Service) GetMessagesByToolCallIds(ctx context.Context, toolCallIdsJson string, userId string) ([]string, error) {
	var toolCallIds []string
	if err := json.Unmarshal([]byte(toolCallIdsJson), &toolCallIds); err != nil {
		return nil, fmt.Errorf("failed to parse tool call IDs: %w", err)
	}

	results := make([]string, 0, len(toolCallIds))
	for _, toolCallId := range toolCallIds {
		msgId, err := s.queries.GetMessageByToolCallId(ctx, db.GetMessageByToolCallIdParams{
			ToolCallID: sql.NullString{String: toolCallId, Valid: true},
			UserID:     userId,
		})
		if err == nil {
			results = append(results, msgId)
		}
		// Ignore not found errors
	}

	return results, nil
}

// GetDocumentsByFileIds fetches documents by file IDs (batch operation)
func (s *Service) GetDocumentsByFileIds(ctx context.Context, fileIdsJson string, userId string) ([]db.GetDocumentByFileIdRow, error) {
	var fileIds []string
	if err := json.Unmarshal([]byte(fileIdsJson), &fileIds); err != nil {
		return nil, fmt.Errorf("failed to parse file IDs: %w", err)
	}

	results := make([]db.GetDocumentByFileIdRow, 0, len(fileIds))
	for _, fileId := range fileIds {
		doc, err := s.queries.GetDocumentByFileId(ctx, db.GetDocumentByFileIdParams{
			FileID: sql.NullString{String: fileId, Valid: true},
			UserID: userId,
		})
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
				UserID:          params.File.UserID,
				CreatedAt:       params.File.CreatedAt,
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
	UserID           string
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

		// 2. Delete embeddings for each chunk
		for _, chunkId := range chunkIds {
			if chunkId.Valid {
				// Try to delete embedding (may not exist)
				_ = q.DeleteEmbedding(ctx, db.DeleteEmbeddingParams{
					ID:     chunkId.String,
					UserID: sql.NullString{String: params.UserID, Valid: true},
				})
			}
		}

		// 3. Delete chunks
		for _, chunkId := range chunkIds {
			if chunkId.Valid {
				err := q.DeleteChunk(ctx, db.DeleteChunkParams{
					ID:     chunkId.String,
					UserID: sql.NullString{String: params.UserID, Valid: true},
				})
				if err != nil && err != sql.ErrNoRows {
					return fmt.Errorf("failed to delete chunk: %w", err)
				}
			}
		}

		// 4. Delete file record
		err = q.DeleteFile(ctx, db.DeleteFileParams{
			ID:     params.FileID,
			UserID: params.UserID,
		})
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
func (s *Service) DeleteAIProviderWithModels(ctx context.Context, providerID string, userID string) error {
	return s.WithTx(ctx, func(q *db.Queries) error {
		// 1. Delete all models of the provider
		err := q.DeleteModelsByProvider(ctx, db.DeleteModelsByProviderParams{
			ProviderID: providerID,
			UserID:     userID,
		})
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("failed to delete models: %w", err)
		}

		// 2. Delete the provider
		err = q.DeleteAIProvider(ctx, db.DeleteAIProviderParams{
			ID:     providerID,
			UserID: userID,
		})
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
// DEFAULT USER AND INBOX INITIALIZATION (Desktop Single-User App)
// ============================================================================

const defaultUserID = "DEFAULT_LOBE_CHAT_USER"

// ensureDefaultUserAndInbox ensures the default user and inbox session exist
// This is called during database initialization for desktop single-user apps
func (s *Service) ensureDefaultUserAndInbox(ctx context.Context) error {
	now := int64(1000) // Use a fixed timestamp for default user

	// 1. Ensure default user exists
	err := s.queries.EnsureUserExists(ctx, db.EnsureUserExistsParams{
		ID:        defaultUserID,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return fmt.Errorf("failed to ensure default user: %w", err)
	}

	// 2. Check if inbox session already exists
	_, err = s.queries.GetSessionBySlug(ctx, db.GetSessionBySlugParams{
		Slug:   "inbox",
		UserID: defaultUserID,
	})

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

	return nil
}

// createDefaultInboxSession creates the default inbox session with agent
func (s *Service) createDefaultInboxSession(ctx context.Context) error {
	return s.WithTx(ctx, func(q *db.Queries) error {
		now := int64(1000)
		sessionID := "default-inbox-session"
		agentID := "default-inbox-agent"

		// 1. Create session
		_, err := q.CreateSession(ctx, db.CreateSessionParams{
			ID:              sessionID,
			UserID:          defaultUserID,
			Slug:            "inbox",
			Title:           sql.NullString{Valid: false},
			Description:     sql.NullString{Valid: false},
			Avatar:          sql.NullString{Valid: false},
			BackgroundColor: sql.NullString{Valid: false},
			Type:            sql.NullString{String: "agent", Valid: true},
			GroupID:         sql.NullString{Valid: false},
			ClientID:        sql.NullString{Valid: false},
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
			UserID:           defaultUserID,
			Slug:             sql.NullString{Valid: false},
			Title:            sql.NullString{Valid: false},
			Description:      sql.NullString{Valid: false},
			Tags:             sql.NullString{String: "[]", Valid: true},
			Avatar:           sql.NullString{Valid: false},
			BackgroundColor:  sql.NullString{Valid: false},
			Plugins:          sql.NullString{String: "[]", Valid: true},
			ClientID:         sql.NullString{Valid: false},
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
			UserID:    defaultUserID,
		})
		if err != nil {
			return fmt.Errorf("failed to link agent to session: %w", err)
		}

		return nil
	})
}
