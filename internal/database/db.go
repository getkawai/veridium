package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"

	db "github.com/kawai-network/veridium/internal/database/generated"
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
