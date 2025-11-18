package services

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/schema"
	"github.com/kawai-network/veridium/pkg/contextengine"
)

// ContextEngineBridge bridges the existing ContextEngine with Eino's agent system
type ContextEngineBridge struct {
	contextEngine *contextengine.Engine
}

// NewContextEngineBridge creates a new context engine bridge
func NewContextEngineBridge(contextEngine *contextengine.Engine) *ContextEngineBridge {
	return &ContextEngineBridge{
		contextEngine: contextEngine,
	}
}

// ProcessMessagesForAgent processes messages through the context engine before sending to agent
// This applies all context processors (history, templates, placeholders, etc.)
func (b *ContextEngineBridge) ProcessMessagesForAgent(ctx context.Context, einoMessages []*schema.Message) ([]*schema.Message, error) {
	// Convert Eino messages to ContextEngine messages
	contextMessages := make([]*contextengine.Message, len(einoMessages))
	for i, msg := range einoMessages {
		// Eino schema.Message.Content is string
		contextMessages[i] = &contextengine.Message{
			ID:        "", // Eino messages don't have ID, generate if needed
			Role:      convertEinoRoleToContextRole(msg.Role),
			Content:   msg.Content, // string → interface{}
			CreatedAt: 0,          // Not used in Eino messages
			UpdatedAt: 0,
			Meta:      nil, // TODO: convert if needed
		}
	}

	// Process through context engine
	processedMessages, err := b.contextEngine.Process(ctx, contextMessages)
	if err != nil {
		return nil, fmt.Errorf("context engine processing failed: %w", err)
	}

	// Convert back to Eino messages
	processedEinoMessages := make([]*schema.Message, len(processedMessages))
	for i, msg := range processedMessages {
		// Convert interface{} content back to string
		content := ""
		if msg.Content != nil {
			if strContent, ok := msg.Content.(string); ok {
				content = strContent
			}
		}

		processedEinoMessages[i] = &schema.Message{
			Role:    convertContextRoleToEinoRole(msg.Role),
			Content: content,
		}
	}

	log.Printf("🔄 ContextEngineBridge: Processed %d messages → %d messages", 
		len(einoMessages), len(processedEinoMessages))

	return processedEinoMessages, nil
}

// MergeWithRAGContext merges context-processed messages with RAG context
// This is useful when you want to combine both context engine processing and RAG
func (b *ContextEngineBridge) MergeWithRAGContext(
	processedMessages []*schema.Message,
	ragContext string,
) []*schema.Message {
	if ragContext == "" {
		return processedMessages
	}

	// Inject RAG context as a system message at the beginning
	ragMessage := &schema.Message{
		Role:    schema.System,
		Content: fmt.Sprintf("Knowledge Base Context:\n%s", ragContext),
	}

	// Merge: RAG context first, then processed messages
	mergedMessages := make([]*schema.Message, 0, len(processedMessages)+1)
	mergedMessages = append(mergedMessages, ragMessage)
	mergedMessages = append(mergedMessages, processedMessages...)

	log.Printf("🔗 ContextEngineBridge: Merged RAG context with %d messages", len(processedMessages))

	return mergedMessages
}

// GetEngine returns the underlying context engine for advanced usage
func (b *ContextEngineBridge) GetEngine() *contextengine.Engine {
	return b.contextEngine
}

// convertEinoRoleToContextRole converts Eino role to ContextEngine role
func convertEinoRoleToContextRole(einoRole schema.RoleType) string {
	switch einoRole {
	case schema.System:
		return "system"
	case schema.User:
		return "user"
	case schema.Assistant:
		return "assistant"
	case schema.Tool:
		return "tool"
	default:
		return "user"
	}
}

// convertContextRoleToEinoRole converts ContextEngine role to Eino role
func convertContextRoleToEinoRole(contextRole string) schema.RoleType {
	switch contextRole {
	case "system":
		return schema.System
	case "user":
		return schema.User
	case "assistant":
		return schema.Assistant
	case "tool":
		return schema.Tool
	default:
		return schema.User
	}
}

// ProcessAndMergeRAG is a convenience method that processes messages and merges with RAG context
func (b *ContextEngineBridge) ProcessAndMergeRAG(
	ctx context.Context,
	einoMessages []*schema.Message,
	ragContext string,
) ([]*schema.Message, error) {
	// Process through context engine
	processedMessages, err := b.ProcessMessagesForAgent(ctx, einoMessages)
	if err != nil {
		return nil, err
	}

	// Merge with RAG context
	return b.MergeWithRAGContext(processedMessages, ragContext), nil
}

