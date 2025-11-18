package eino

import (
	"context"

	"github.com/kawai-network/veridium/context-engine/internal/types"
)

// ContextState holds shared state for the context engineering pipeline
type ContextState struct {
	Messages    []*types.Message
	Metadata    map[string]interface{}
	IsAborted   bool
	AbortReason string
}

// NewContextState creates a new context state
func NewContextState(messages []*types.Message) *ContextState {
	return &ContextState{
		Messages: messages,
		Metadata: make(map[string]interface{}),
	}
}

// StatePreHandler is a function that handles state before node execution
type StatePreHandler func(ctx context.Context, state *ContextState) context.Context

// StatePostHandler is a function that handles state after node execution
type StatePostHandler func(ctx context.Context, state *ContextState, output interface{}) context.Context

// GetMetadata retrieves metadata value
func (cs *ContextState) GetMetadata(key string) (interface{}, bool) {
	if cs.Metadata == nil {
		return nil, false
	}
	val, ok := cs.Metadata[key]
	return val, ok
}

// SetMetadata sets metadata value
func (cs *ContextState) SetMetadata(key string, value interface{}) {
	if cs.Metadata == nil {
		cs.Metadata = make(map[string]interface{})
	}
	cs.Metadata[key] = value
}

// Abort marks the state as aborted
func (cs *ContextState) Abort(reason string) {
	cs.IsAborted = true
	cs.AbortReason = reason
}
