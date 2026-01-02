package cliproxy

import (
	"context"

	"github.com/kawai-network/veridium/pkg/cliproxy/sdk/config"
)

// NewFileTokenClientProvider returns the default token-backed client loader.
func NewFileTokenClientProvider() TokenClientProvider {
	return &fileTokenClientProvider{}
}

type fileTokenClientProvider struct{}

func (p *fileTokenClientProvider) Load(ctx context.Context, cfg *config.Config) (*TokenClientResult, error) {
	// Stateless executors handle tokens
	_ = ctx
	_ = cfg
	return &TokenClientResult{SuccessfulAuthed: 0}, nil
}

// NewAPIKeyClientProvider returns the default API key client loader that reuses existing logic.
func NewAPIKeyClientProvider() APIKeyClientProvider {
	return &apiKeyClientProvider{}
}

type apiKeyClientProvider struct{}

func (p *apiKeyClientProvider) Load(ctx context.Context, cfg *config.Config) (*APIKeyClientResult, error) {
	return &APIKeyClientResult{
		GeminiKeyCount:       0,
		VertexCompatKeyCount: 0,
		ClaudeKeyCount:       0,
		CodexKeyCount:        0,
		OpenAICompatCount:    0,
	}, nil
}
