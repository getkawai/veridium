package services

import (
	"context"
	"database/sql"

	db "github.com/kawai-network/veridium/internal/database/generated"
)

type AgentService struct {
	queries *db.Queries
}

func NewAgentService(queries *db.Queries) *AgentService {
	return &AgentService{
		queries: queries,
	}
}

// Agent operations

func (s *AgentService) GetAgent(ctx context.Context, agentID, userID string) (db.Agent, error) {
	return s.queries.GetAgent(ctx, db.GetAgentParams{
		ID:     agentID,
		UserID: userID,
	})
}

func (s *AgentService) GetAgentBySlug(ctx context.Context, slug, userID string) (db.Agent, error) {
	return s.queries.GetAgentBySlug(ctx, db.GetAgentBySlugParams{
		Slug:   sql.NullString{String: slug, Valid: true},
		UserID: userID,
	})
}

func (s *AgentService) ListAgents(ctx context.Context, userID string, limit, offset int64) ([]db.Agent, error) {
	return s.queries.ListAgents(ctx, db.ListAgentsParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
}

func (s *AgentService) SearchAgents(ctx context.Context, userID, query string, limit int64) ([]db.Agent, error) {
	searchPattern := "%" + query + "%"
	return s.queries.SearchAgents(ctx, db.SearchAgentsParams{
		UserID:  userID,
		Title:   sql.NullString{String: searchPattern, Valid: true},
		Title_2: sql.NullString{String: searchPattern, Valid: true},
		Limit:   limit,
	})
}

type CreateAgentParams struct {
	ID               string
	Slug             *string
	Title            *string
	Description      *string
	Tags             []string
	Avatar           *string
	BackgroundColor  *string
	Plugins          []string
	ClientID         *string
	UserID           string
	ChatConfig       map[string]interface{}
	FewShots         map[string]interface{}
	Model            *string
	Params           map[string]interface{}
	Provider         *string
	SystemRole       *string
	TTS              map[string]interface{}
	Virtual          bool
	OpeningMessage   *string
	OpeningQuestions []string
}

func (s *AgentService) CreateAgent(ctx context.Context, params CreateAgentParams) (db.Agent, error) {
	now := currentTimestampMs()

	return s.queries.CreateAgent(ctx, db.CreateAgentParams{
		ID:               params.ID,
		Slug:             toNullString(params.Slug),
		Title:            toNullString(params.Title),
		Description:      toNullString(params.Description),
		Tags:             toNullJSONArray(params.Tags),
		Avatar:           toNullString(params.Avatar),
		BackgroundColor:  toNullString(params.BackgroundColor),
		Plugins:          toNullJSONArray(params.Plugins),
		ClientID:         toNullString(params.ClientID),
		UserID:           params.UserID,
		ChatConfig:       toNullJSON(params.ChatConfig),
		FewShots:         toNullJSON(params.FewShots),
		Model:            toNullString(params.Model),
		Params:           toNullJSON(params.Params),
		Provider:         toNullString(params.Provider),
		SystemRole:       toNullString(params.SystemRole),
		Tts:              toNullJSON(params.TTS),
		Virtual:          boolToInt(params.Virtual),
		OpeningMessage:   toNullString(params.OpeningMessage),
		OpeningQuestions: toNullJSONArray(params.OpeningQuestions),
		CreatedAt:        now,
		UpdatedAt:        now,
	})
}

type UpdateAgentParams struct {
	ID               string
	UserID           string
	Title            *string
	Description      *string
	Tags             []string
	Avatar           *string
	BackgroundColor  *string
	Plugins          []string
	ChatConfig       map[string]interface{}
	FewShots         map[string]interface{}
	Model            *string
	Params           map[string]interface{}
	Provider         *string
	SystemRole       *string
	TTS              map[string]interface{}
	OpeningMessage   *string
	OpeningQuestions []string
}

func (s *AgentService) UpdateAgent(ctx context.Context, params UpdateAgentParams) (db.Agent, error) {
	return s.queries.UpdateAgent(ctx, db.UpdateAgentParams{
		Title:            toNullString(params.Title),
		Description:      toNullString(params.Description),
		Tags:             toNullJSONArray(params.Tags),
		Avatar:           toNullString(params.Avatar),
		BackgroundColor:  toNullString(params.BackgroundColor),
		Plugins:          toNullJSONArray(params.Plugins),
		ChatConfig:       toNullJSON(params.ChatConfig),
		FewShots:         toNullJSON(params.FewShots),
		Model:            toNullString(params.Model),
		Params:           toNullJSON(params.Params),
		Provider:         toNullString(params.Provider),
		SystemRole:       toNullString(params.SystemRole),
		Tts:              toNullJSON(params.TTS),
		OpeningMessage:   toNullString(params.OpeningMessage),
		OpeningQuestions: toNullJSONArray(params.OpeningQuestions),
		UpdatedAt:        currentTimestampMs(),
		ID:               params.ID,
		UserID:           params.UserID,
	})
}

func (s *AgentService) DeleteAgent(ctx context.Context, agentID, userID string) error {
	return s.queries.DeleteAgent(ctx, db.DeleteAgentParams{
		ID:     agentID,
		UserID: userID,
	})
}

// Agent relationship operations

func (s *AgentService) LinkAgentToSession(ctx context.Context, agentID, sessionID, userID string) error {
	return s.queries.LinkAgentToSession(ctx, db.LinkAgentToSessionParams{
		AgentID:   agentID,
		SessionID: sessionID,
		UserID:    userID,
	})
}

func (s *AgentService) UnlinkAgentFromSession(ctx context.Context, agentID, sessionID, userID string) error {
	return s.queries.UnlinkAgentFromSession(ctx, db.UnlinkAgentFromSessionParams{
		AgentID:   agentID,
		SessionID: sessionID,
		UserID:    userID,
	})
}

func (s *AgentService) GetSessionAgents(ctx context.Context, sessionID, userID string) ([]db.Agent, error) {
	return s.queries.GetSessionAgents(ctx, db.GetSessionAgentsParams{
		SessionID: sessionID,
		UserID:    userID,
	})
}

func (s *AgentService) LinkAgentToFile(ctx context.Context, fileID, agentID, userID string, enabled bool) error {
	now := currentTimestampMs()
	return s.queries.LinkAgentToFile(ctx, db.LinkAgentToFileParams{
		FileID:    fileID,
		AgentID:   agentID,
		Enabled:   boolToInt(enabled),
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	})
}

func (s *AgentService) UnlinkAgentFromFile(ctx context.Context, fileID, agentID, userID string) error {
	return s.queries.UnlinkAgentFromFile(ctx, db.UnlinkAgentFromFileParams{
		FileID:  fileID,
		AgentID: agentID,
		UserID:  userID,
	})
}

func (s *AgentService) GetAgentFiles(ctx context.Context, agentID, userID string) ([]db.File, error) {
	return s.queries.GetAgentFiles(ctx, db.GetAgentFilesParams{
		AgentID: agentID,
		UserID:  userID,
	})
}

func (s *AgentService) LinkAgentToKnowledgeBase(ctx context.Context, agentID, knowledgeBaseID, userID string, enabled bool) error {
	now := currentTimestampMs()
	return s.queries.LinkAgentToKnowledgeBase(ctx, db.LinkAgentToKnowledgeBaseParams{
		AgentID:         agentID,
		KnowledgeBaseID: knowledgeBaseID,
		UserID:          userID,
		Enabled:         boolToInt(enabled),
		CreatedAt:       now,
		UpdatedAt:       now,
	})
}

func (s *AgentService) UnlinkAgentFromKnowledgeBase(ctx context.Context, agentID, knowledgeBaseID, userID string) error {
	return s.queries.UnlinkAgentFromKnowledgeBase(ctx, db.UnlinkAgentFromKnowledgeBaseParams{
		AgentID:         agentID,
		KnowledgeBaseID: knowledgeBaseID,
		UserID:          userID,
	})
}

func (s *AgentService) GetAgentKnowledgeBases(ctx context.Context, agentID, userID string) ([]db.KnowledgeBase, error) {
	return s.queries.GetAgentKnowledgeBases(ctx, db.GetAgentKnowledgeBasesParams{
		AgentID: agentID,
		UserID:  userID,
	})
}
