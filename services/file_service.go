package services

import (
	"context"

	db "github.com/kawai-network/veridium/internal/database/generated"
)

type FileService struct {
	queries *db.Queries
}

func NewFileService(queries *db.Queries) *FileService {
	return &FileService{
		queries: queries,
	}
}

// File operations

func (s *FileService) GetFile(ctx context.Context, fileID, userID string) (db.File, error) {
	return s.queries.GetFile(ctx, db.GetFileParams{
		ID:     fileID,
		UserID: userID,
	})
}

func (s *FileService) ListFiles(ctx context.Context, userID string, limit, offset int64) ([]db.File, error) {
	return s.queries.ListFiles(ctx, db.ListFilesParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
}

type CreateFileParams struct {
	ID              string
	UserID          string
	FileType        string
	FileHash        *string
	Name            string
	Size            int64
	URL             string
	Source          map[string]interface{}
	ClientID        *string
	Metadata        map[string]interface{}
	ChunkTaskID     *string
	EmbeddingTaskID *string
}

func (s *FileService) CreateFile(ctx context.Context, params CreateFileParams) (db.File, error) {
	now := currentTimestampMs()

	return s.queries.CreateFile(ctx, db.CreateFileParams{
		ID:              params.ID,
		UserID:          params.UserID,
		FileType:        params.FileType,
		FileHash:        toNullString(params.FileHash),
		Name:            params.Name,
		Size:            params.Size,
		Url:             params.URL,
		Source:          toNullJSON(params.Source),
		ClientID:        toNullString(params.ClientID),
		Metadata:        toNullJSON(params.Metadata),
		ChunkTaskID:     toNullString(params.ChunkTaskID),
		EmbeddingTaskID: toNullString(params.EmbeddingTaskID),
		CreatedAt:       now,
		UpdatedAt:       now,
	})
}

type UpdateFileParams struct {
	ID       string
	UserID   string
	Name     string
	Metadata map[string]interface{}
}

func (s *FileService) UpdateFile(ctx context.Context, params UpdateFileParams) (db.File, error) {
	return s.queries.UpdateFile(ctx, db.UpdateFileParams{
		Name:      params.Name,
		Metadata:  toNullJSON(params.Metadata),
		UpdatedAt: currentTimestampMs(),
		ID:        params.ID,
		UserID:    params.UserID,
	})
}

func (s *FileService) DeleteFile(ctx context.Context, fileID, userID string) error {
	return s.queries.DeleteFile(ctx, db.DeleteFileParams{
		ID:     fileID,
		UserID: userID,
	})
}

// Global File operations

func (s *FileService) GetGlobalFile(ctx context.Context, hashID string) (db.GlobalFile, error) {
	return s.queries.GetGlobalFile(ctx, hashID)
}

type CreateGlobalFileParams struct {
	HashID   string
	FileType string
	Size     int64
	URL      string
	Metadata map[string]interface{}
	Creator  string
}

func (s *FileService) CreateGlobalFile(ctx context.Context, params CreateGlobalFileParams) (db.GlobalFile, error) {
	now := currentTimestampMs()

	return s.queries.CreateGlobalFile(ctx, db.CreateGlobalFileParams{
		HashID:     params.HashID,
		FileType:   params.FileType,
		Size:       params.Size,
		Url:        params.URL,
		Metadata:   toNullJSON(params.Metadata),
		Creator:    params.Creator,
		CreatedAt:  now,
		AccessedAt: now,
	})
}

func (s *FileService) UpdateGlobalFileAccess(ctx context.Context, hashID string) error {
	return s.queries.UpdateGlobalFileAccess(ctx, db.UpdateGlobalFileAccessParams{
		AccessedAt: currentTimestampMs(),
		HashID:     hashID,
	})
}

// Knowledge Base operations

func (s *FileService) GetKnowledgeBase(ctx context.Context, kbID, userID string) (db.KnowledgeBase, error) {
	return s.queries.GetKnowledgeBase(ctx, db.GetKnowledgeBaseParams{
		ID:     kbID,
		UserID: userID,
	})
}

func (s *FileService) ListKnowledgeBases(ctx context.Context, userID string) ([]db.KnowledgeBase, error) {
	return s.queries.ListKnowledgeBases(ctx, userID)
}

type CreateKnowledgeBaseParams struct {
	ID          string
	Name        string
	Description *string
	Avatar      *string
	Type        *string
	UserID      string
	ClientID    *string
	IsPublic    bool
	Settings    map[string]interface{}
}

func (s *FileService) CreateKnowledgeBase(ctx context.Context, params CreateKnowledgeBaseParams) (db.KnowledgeBase, error) {
	now := currentTimestampMs()

	return s.queries.CreateKnowledgeBase(ctx, db.CreateKnowledgeBaseParams{
		ID:          params.ID,
		Name:        params.Name,
		Description: toNullString(params.Description),
		Avatar:      toNullString(params.Avatar),
		Type:        toNullString(params.Type),
		UserID:      params.UserID,
		ClientID:    toNullString(params.ClientID),
		IsPublic:    boolToInt(params.IsPublic),
		Settings:    toNullJSON(params.Settings),
		CreatedAt:   now,
		UpdatedAt:   now,
	})
}

type UpdateKnowledgeBaseParams struct {
	ID          string
	UserID      string
	Name        string
	Description *string
	Avatar      *string
	Settings    map[string]interface{}
}

func (s *FileService) UpdateKnowledgeBase(ctx context.Context, params UpdateKnowledgeBaseParams) (db.KnowledgeBase, error) {
	return s.queries.UpdateKnowledgeBase(ctx, db.UpdateKnowledgeBaseParams{
		Name:        params.Name,
		Description: toNullString(params.Description),
		Avatar:      toNullString(params.Avatar),
		Settings:    toNullJSON(params.Settings),
		UpdatedAt:   currentTimestampMs(),
		ID:          params.ID,
		UserID:      params.UserID,
	})
}

func (s *FileService) DeleteKnowledgeBase(ctx context.Context, kbID, userID string) error {
	return s.queries.DeleteKnowledgeBase(ctx, db.DeleteKnowledgeBaseParams{
		ID:     kbID,
		UserID: userID,
	})
}

// Knowledge Base Files operations

func (s *FileService) LinkKnowledgeBaseToFile(ctx context.Context, kbID, fileID, userID string) error {
	return s.queries.LinkKnowledgeBaseToFile(ctx, db.LinkKnowledgeBaseToFileParams{
		KnowledgeBaseID: kbID,
		FileID:          fileID,
		UserID:          userID,
		CreatedAt:       currentTimestampMs(),
	})
}

func (s *FileService) UnlinkKnowledgeBaseFromFile(ctx context.Context, kbID, fileID, userID string) error {
	return s.queries.UnlinkKnowledgeBaseFromFile(ctx, db.UnlinkKnowledgeBaseFromFileParams{
		KnowledgeBaseID: kbID,
		FileID:          fileID,
		UserID:          userID,
	})
}

func (s *FileService) GetKnowledgeBaseFiles(ctx context.Context, kbID, userID string) ([]db.File, error) {
	return s.queries.GetKnowledgeBaseFiles(ctx, db.GetKnowledgeBaseFilesParams{
		KnowledgeBaseID: kbID,
		UserID:          userID,
	})
}

// Session Files operations

func (s *FileService) LinkFileToSession(ctx context.Context, fileID, sessionID, userID string) error {
	return s.queries.LinkFileToSession(ctx, db.LinkFileToSessionParams{
		FileID:    fileID,
		SessionID: sessionID,
		UserID:    userID,
	})
}

func (s *FileService) UnlinkFileFromSession(ctx context.Context, fileID, sessionID, userID string) error {
	return s.queries.UnlinkFileFromSession(ctx, db.UnlinkFileFromSessionParams{
		FileID:    fileID,
		SessionID: sessionID,
		UserID:    userID,
	})
}

func (s *FileService) GetSessionFiles(ctx context.Context, sessionID, userID string) ([]db.File, error) {
	return s.queries.GetSessionFiles(ctx, db.GetSessionFilesParams{
		SessionID: sessionID,
		UserID:    userID,
	})
}
