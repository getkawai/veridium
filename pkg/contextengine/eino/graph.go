package eino

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	proc "github.com/kawai-network/veridium/pkg/contextengine/processors"
	prov "github.com/kawai-network/veridium/pkg/contextengine/providers"
	"github.com/kawai-network/veridium/pkg/contextengine/types"
)

// MessageInput wraps messages for workflow input
type MessageInput struct {
	Messages []*schema.Message
}

// MessageOutput wraps messages for workflow output
type MessageOutput struct {
	Messages []*schema.Message
}

// GraphBuilder builds Eino workflow graphs for context engineering
type GraphBuilder struct {
	config *types.ContextConfig
}

// NewGraphBuilder creates a new graph builder
func NewGraphBuilder(config *types.ContextConfig) *GraphBuilder {
	return &GraphBuilder{
		config: config,
	}
}

// BuildContextEngineeringGraph builds the main context engineering workflow
func (gb *GraphBuilder) BuildContextEngineeringGraph() (*compose.Workflow[MessageInput, MessageOutput], error) {
	wf := compose.NewWorkflow[MessageInput, MessageOutput]()

	// Import processors and providers
	processors := gb.buildProcessors()
	providers := gb.buildProviders()

	// Helper function to wrap lambda nodes that work with []*schema.Message
	wrapLambda := func(lambda *compose.Lambda) *compose.Lambda {
		return compose.InvokableLambda(func(ctx context.Context, input MessageOutput) (MessageOutput, error) {
			// Create a temporary workflow to invoke the lambda
			tempWf := compose.NewWorkflow[[]*schema.Message, []*schema.Message]()
			tempWf.AddLambdaNode("temp", lambda).AddInput(compose.START)
			tempWf.End().AddInput("temp")

			compiled, err := tempWf.Compile(ctx)
			if err != nil {
				return MessageOutput{}, err
			}

			result, err := compiled.Invoke(ctx, input.Messages)
			if err != nil {
				return MessageOutput{}, err
			}
			return MessageOutput{Messages: result}, nil
		})
	}

	// 0. Group Message Flatten (MUST be first, to normalize group messages)
	// Special handling for first node - it receives MessageInput from START
	if processors.GroupMessageFlatten != nil {
		firstNodeLambda := compose.InvokableLambda(func(ctx context.Context, input MessageInput) (MessageOutput, error) {
			// Create a temporary workflow to invoke the lambda
			tempWf := compose.NewWorkflow[[]*schema.Message, []*schema.Message]()
			tempWf.AddLambdaNode("temp", processors.GroupMessageFlatten).AddInput(compose.START)
			tempWf.End().AddInput("temp")

			compiled, err := tempWf.Compile(ctx)
			if err != nil {
				return MessageOutput{}, err
			}

			result, err := compiled.Invoke(ctx, input.Messages)
			if err != nil {
				return MessageOutput{}, err
			}
			return MessageOutput{Messages: result}, nil
		})
		wf.AddLambdaNode("groupFlatten", firstNodeLambda).AddInput(compose.START)
	} else {
		// Passthrough wrapper
		passthroughLambda := compose.InvokableLambda(func(ctx context.Context, input MessageInput) (MessageOutput, error) {
			return MessageOutput{Messages: input.Messages}, nil
		})
		wf.AddLambdaNode("groupFlatten", passthroughLambda).AddInput(compose.START)
	}

	lastNode := "groupFlatten"

	// 1. History Truncate (after group flatten, before any message injection)
	if processors.HistoryTruncate != nil {
		wf.AddLambdaNode("truncate", wrapLambda(processors.HistoryTruncate)).AddInput(lastNode)
		lastNode = "truncate"
	}

	// 2. System Role Injection
	if providers.SystemRole != nil {
		wf.AddLambdaNode("systemRole", wrapLambda(providers.SystemRole)).
			AddInput(lastNode)
		lastNode = "systemRole"
	}

	// 3. Inbox Guide Provider
	if providers.InboxGuide != nil {
		wf.AddLambdaNode("inboxGuide", wrapLambda(providers.InboxGuide)).
			AddInput(lastNode)
		lastNode = "inboxGuide"
	}

	// 4. Tool System Role Provider
	if providers.ToolSystemRole != nil {
		wf.AddLambdaNode("toolSystemRole", wrapLambda(providers.ToolSystemRole)).
			AddInput(lastNode)
		lastNode = "toolSystemRole"
	}

	// 5. History Summary Provider
	if providers.HistorySummary != nil {
		wf.AddLambdaNode("historySummary", wrapLambda(providers.HistorySummary)).
			AddInput(lastNode)
		lastNode = "historySummary"
	}

	// 6. Input Template Processing
	if processors.InputTemplate != nil {
		wf.AddLambdaNode("template", wrapLambda(processors.InputTemplate)).
			AddInput(lastNode)
		lastNode = "template"
	}

	// 7. Placeholder Variables
	if processors.PlaceholderVariables != nil {
		wf.AddLambdaNode("placeholders", wrapLambda(processors.PlaceholderVariables)).
			AddInput(lastNode)
		lastNode = "placeholders"
	}

	// 8. Message Content Processing
	if processors.MessageContent != nil {
		wf.AddLambdaNode("content", wrapLambda(processors.MessageContent)).
			AddInput(lastNode)
		lastNode = "content"
	}

	// 9. Tool Call Processing
	if processors.ToolCall != nil {
		wf.AddLambdaNode("toolCalls", wrapLambda(processors.ToolCall)).
			AddInput(lastNode)
		lastNode = "toolCalls"
	}

	// 10. Tool Message Reorder
	if processors.ToolMessageReorder != nil {
		wf.AddLambdaNode("reorder", wrapLambda(processors.ToolMessageReorder)).
			AddInput(lastNode)
		lastNode = "reorder"
	}

	// 11. Message Cleanup
	if processors.MessageCleanup != nil {
		wf.AddLambdaNode("cleanup", wrapLambda(processors.MessageCleanup)).
			AddInput(lastNode)
		lastNode = "cleanup"
	}

	wf.End().AddInput(lastNode)

	return wf, nil
}

// processors holds all processor lambdas
type processors struct {
	GroupMessageFlatten  *compose.Lambda
	HistoryTruncate      *compose.Lambda
	InputTemplate        *compose.Lambda
	PlaceholderVariables *compose.Lambda
	MessageContent       *compose.Lambda
	ToolCall             *compose.Lambda
	ToolMessageReorder   *compose.Lambda
	MessageCleanup       *compose.Lambda
}

// providers holds all provider lambdas
type providers struct {
	SystemRole     *compose.Lambda
	InboxGuide     *compose.Lambda
	ToolSystemRole *compose.Lambda
	HistorySummary *compose.Lambda
}

// buildProcessors builds all processor lambdas
func (gb *GraphBuilder) buildProcessors() *processors {
	p := &processors{}

	// Import processor packages
	procConfig := gb.config

	// Group Message Flatten (always enabled)
	p.GroupMessageFlatten = proc.NewGroupMessageFlattenLambda(proc.GroupMessageFlattenConfig{})

	// History Truncate
	if procConfig.EnableHistoryCount {
		p.HistoryTruncate = proc.NewHistoryTruncateLambda(proc.HistoryTruncateConfig{
			EnableHistoryCount: procConfig.EnableHistoryCount,
			HistoryCount:       procConfig.HistoryCount,
		})
	}

	// Input Template
	if procConfig.InputTemplate != "" {
		p.InputTemplate = proc.NewInputTemplateLambda(proc.InputTemplateConfig{
			InputTemplate: procConfig.InputTemplate,
		})
	}

	// Placeholder Variables
	if len(procConfig.Variables) > 0 {
		generators := make(map[string]func() string)
		for k, v := range procConfig.Variables {
			// Convert interface{} to string generator
			if strVal, ok := v.(string); ok {
				val := strVal // Capture for closure
				generators[k] = func() string { return val }
			} else if fn, ok := v.(func() string); ok {
				generators[k] = fn
			}
		}
		if len(generators) > 0 {
			p.PlaceholderVariables = proc.NewPlaceholderVariablesLambda(proc.PlaceholderVariablesConfig{
				VariableGenerators: generators,
				Depth:              2,
			})
		}
	}

	// Message Content
	p.MessageContent = proc.NewMessageContentLambda(proc.MessageContentConfig{
		FileContext: proc.FileContextConfig{
			Enabled:        procConfig.MessageContent.FileContext.Enabled,
			IncludeFileURL: procConfig.MessageContent.FileContext.IncludeFileURL,
		},
		IsCanUseVideo:  procConfig.MessageContent.IsCanUseVideo,
		IsCanUseVision: procConfig.MessageContent.IsCanUseVision,
		Model:          procConfig.Model,
		Provider:       procConfig.Provider,
	})

	// Tool Call
	if len(procConfig.Tools) > 0 {
		p.ToolCall = proc.NewToolCallLambda(proc.ToolCallConfig{
			IsCanUseFC: procConfig.MessageContent.IsCanUseVision, // Placeholder
			Model:      procConfig.Model,
			Provider:   procConfig.Provider,
		})
	}

	// Tool Message Reorder
	p.ToolMessageReorder = proc.ToolMessageReorderLambda()

	// Message Cleanup
	p.MessageCleanup = proc.MessageCleanupLambda()

	return p
}

// buildProviders builds all provider lambdas
func (gb *GraphBuilder) buildProviders() *providers {
	providersList := &providers{}
	config := gb.config

	// System Role Injector
	if config.SystemRole != "" {
		providersList.SystemRole = prov.NewSystemRoleInjectorLambda(prov.SystemRoleInjectorConfig{
			SystemRole: config.SystemRole,
		})
	}

	// Inbox Guide Provider
	if config.IsWelcomeQuestion && config.SessionID != "" {
		// Note: InboxSessionID and InboxGuideSystemRole should come from config
		// For now, skip if not provided
		providersList.InboxGuide = nil // Will be set when config is complete
	}

	// Tool System Role Provider
	if len(config.Tools) > 0 {
		providersList.ToolSystemRole = prov.NewToolSystemRoleProviderLambda(prov.ToolSystemRoleConfig{
			GetToolSystemRoles: func(tools []interface{}) string {
				// This should be provided by caller
				return ""
			},
			IsCanUseFC: func(model, provider string) bool {
				// This should be provided by caller
				return false
			},
			Model:    config.Model,
			Provider: config.Provider,
			Tools:    convertToolsToInterface(config.Tools),
		})
	}

	// History Summary Provider
	if config.HistorySummary != "" {
		providersList.HistorySummary = prov.NewHistorySummaryProviderLambda(prov.HistorySummaryConfig{
			HistorySummary: config.HistorySummary,
		})
	}

	return providersList
}

// convertToolsToInterface converts Tool slice to []interface{}
func convertToolsToInterface(tools []types.Tool) []interface{} {
	result := make([]interface{}, len(tools))
	for i, tool := range tools {
		result[i] = tool
	}
	return result
}

// ProcessMessages processes messages through the context engineering pipeline
func (gb *GraphBuilder) ProcessMessages(ctx context.Context, messages []*types.Message) ([]*types.Message, error) {
	// Convert types.Message to schema.Message
	schemaMessages := convertToSchemaMessages(messages)

	// Build graph
	graph, err := gb.BuildContextEngineeringGraph()
	if err != nil {
		return nil, err
	}

	// Compile graph
	compiled, err := graph.Compile(ctx)
	if err != nil {
		return nil, err
	}

	// Invoke graph with MessageInput wrapper
	input := MessageInput{Messages: schemaMessages}
	result, err := compiled.Invoke(ctx, input)
	if err != nil {
		return nil, err
	}

	// Convert back to types.Message
	return convertFromSchemaMessages(result.Messages), nil
}

// convertToSchemaMessages converts internal Message type to Eino schema.Message
func convertToSchemaMessages(msgs []*types.Message) []*schema.Message {
	result := make([]*schema.Message, len(msgs))
	for i, msg := range msgs {
		result[i] = convertToSchemaMessage(msg)
	}
	return result
}

// convertToSchemaMessage converts a single Message to schema.Message
func convertToSchemaMessage(msg *types.Message) *schema.Message {
	schemaMsg := &schema.Message{
		Role: convertRoleType(msg.Role),
	}

	// Handle content conversion
	switch v := msg.Content.(type) {
	case string:
		schemaMsg.Content = v
	case []types.ContentPart:
		// Convert to UserInputMultiContent for user messages
		if msg.Role == "user" {
			parts := make([]schema.MessageInputPart, 0, len(v))
			for _, part := range v {
				if inputPart := convertToMessageInputPart(part); inputPart != nil {
					parts = append(parts, *inputPart)
				}
			}
			if len(parts) > 0 {
				schemaMsg.UserInputMultiContent = parts
			}
		} else {
			// For assistant messages, use Content field with text concatenation
			textParts := make([]string, 0)
			for _, part := range v {
				if part.Type == "text" && part.Text != nil {
					textParts = append(textParts, *part.Text)
				}
			}
			if len(textParts) > 0 {
				schemaMsg.Content = textParts[0] // Use first text part
			}
		}
	default:
		if str, ok := v.(string); ok {
			schemaMsg.Content = str
		}
	}

	return schemaMsg
}

// convertFromSchemaMessages converts Eino schema.Message to internal Message type
func convertFromSchemaMessages(msgs []*schema.Message) []*types.Message {
	result := make([]*types.Message, len(msgs))
	for i, msg := range msgs {
		result[i] = convertFromSchemaMessage(msg)
	}
	return result
}

// convertFromSchemaMessage converts a single schema.Message to Message
func convertFromSchemaMessage(msg *schema.Message) *types.Message {
	result := &types.Message{
		Role:      string(msg.Role),
		CreatedAt: 0,
		UpdatedAt: 0,
		Meta:      make(map[string]interface{}),
	}

	// Handle content conversion
	if len(msg.UserInputMultiContent) > 0 {
		// Convert from UserInputMultiContent
		parts := make([]types.ContentPart, 0, len(msg.UserInputMultiContent))
		for _, part := range msg.UserInputMultiContent {
			if cp := convertFromMessageInputPart(part); cp != nil {
				parts = append(parts, *cp)
			}
		}
		if len(parts) > 0 {
			result.Content = parts
		} else if msg.Content != "" {
			result.Content = msg.Content
		}
	} else if msg.Content != "" {
		result.Content = msg.Content
	}

	return result
}

// convertRoleType converts string role to schema.RoleType
func convertRoleType(role string) schema.RoleType {
	switch role {
	case "user":
		return schema.User
	case "assistant":
		return schema.Assistant
	case "system":
		return schema.System
	case "tool":
		return schema.Tool
	default:
		return schema.User // Default to user
	}
}

// convertToMessageInputPart converts ContentPart to MessageInputPart
func convertToMessageInputPart(part types.ContentPart) *schema.MessageInputPart {
	switch part.Type {
	case "text":
		if part.Text != nil {
			return &schema.MessageInputPart{
				Type: schema.ChatMessagePartTypeText,
				Text: *part.Text,
			}
		}
	case "image_url":
		if part.ImageURL != nil {
			detail := schema.ImageURLDetailAuto
			if part.ImageURL.Detail != "" {
				detail = schema.ImageURLDetail(part.ImageURL.Detail)
			}
			return &schema.MessageInputPart{
				Type: schema.ChatMessagePartTypeImageURL,
				Image: &schema.MessageInputImage{
					MessagePartCommon: schema.MessagePartCommon{
						URL: &part.ImageURL.URL,
					},
					Detail: detail,
				},
			}
		}
	case "video_url":
		if part.VideoURL != nil {
			return &schema.MessageInputPart{
				Type: schema.ChatMessagePartTypeVideoURL,
				Video: &schema.MessageInputVideo{
					MessagePartCommon: schema.MessagePartCommon{
						URL: &part.VideoURL.URL,
					},
				},
			}
		}
	}
	return nil
}

// convertFromMessageInputPart converts MessageInputPart to ContentPart
func convertFromMessageInputPart(part schema.MessageInputPart) *types.ContentPart {
	switch part.Type {
	case schema.ChatMessagePartTypeText:
		if part.Text != "" {
			text := part.Text
			return &types.ContentPart{
				Type: "text",
				Text: &text,
			}
		}
	case schema.ChatMessagePartTypeImageURL:
		if part.Image != nil && part.Image.URL != nil {
			cp := &types.ContentPart{
				Type: "image_url",
				ImageURL: &types.ImageURL{
					URL:    *part.Image.URL,
					Detail: string(part.Image.Detail),
				},
			}
			return cp
		}
	case schema.ChatMessagePartTypeVideoURL:
		if part.Video != nil && part.Video.URL != nil {
			return &types.ContentPart{
				Type: "video_url",
				VideoURL: &types.VideoURL{
					URL: *part.Video.URL,
				},
			}
		}
	}
	return nil
}
