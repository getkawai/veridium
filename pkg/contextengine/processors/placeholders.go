package processors

import (
	"context"
	"regexp"
	"strings"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// PlaceholderVariablesConfig holds configuration for placeholder variable processing
type PlaceholderVariablesConfig struct {
	VariableGenerators map[string]func() string
	Depth              int
}

const defaultPlaceholderDepth = 2

var placeholderRegex = regexp.MustCompile(`{{(.*?)}}`)

// NewPlaceholderVariablesLambda creates a lambda node for placeholder variable processing
func NewPlaceholderVariablesLambda(config PlaceholderVariablesConfig) *compose.Lambda {
	depth := config.Depth
	if depth == 0 {
		depth = defaultPlaceholderDepth
	}

	return compose.InvokableLambda(func(ctx context.Context, msgs []*schema.Message) ([]*schema.Message, error) {
		if len(config.VariableGenerators) == 0 {
			return msgs, nil
		}

		result := make([]*schema.Message, len(msgs))
		for i, msg := range msgs {
			result[i] = processMessagePlaceholders(msg, config.VariableGenerators, depth)
		}

		return result, nil
	})
}

// processMessagePlaceholders processes placeholder variables in a message
func processMessagePlaceholders(msg *schema.Message, generators map[string]func() string, depth int) *schema.Message {
	result := &schema.Message{
		Role:                  msg.Role,
		Content:               msg.Content,
		Name:                  msg.Name,
		ToolCalls:             msg.ToolCalls,
		ToolCallID:            msg.ToolCallID,
		ToolName:              msg.ToolName,
		ResponseMeta:          msg.ResponseMeta,
		UserInputMultiContent: msg.UserInputMultiContent,
		AssistantGenMultiContent: msg.AssistantGenMultiContent,
	}

	// Process string content
	if msg.Content != "" {
		result.Content = parsePlaceholderVariables(msg.Content, generators, depth)
	}

	// Process UserInputMultiContent
	if len(msg.UserInputMultiContent) > 0 {
		processedParts := make([]schema.MessageInputPart, len(msg.UserInputMultiContent))
		for i, part := range msg.UserInputMultiContent {
			processedParts[i] = part
			if part.Type == schema.ChatMessagePartTypeText && part.Text != "" {
				processedText := parsePlaceholderVariables(part.Text, generators, depth)
				processedParts[i].Text = processedText
			}
		}
		result.UserInputMultiContent = processedParts
	}

	return result
}

// parsePlaceholderVariables replaces template variables with actual values
func parsePlaceholderVariables(text string, generators map[string]func() string, depth int) string {
	result := text

	// Recursive parsing to handle nested variables
	for i := 0; i < depth; i++ {
		matches := placeholderRegex.FindAllStringSubmatch(result, -1)
		if len(matches) == 0 {
			break
		}

		// Extract variable names
		variables := make(map[string]string)
		for _, match := range matches {
			if len(match) > 1 {
				varName := strings.TrimSpace(match[1])
				if generator, exists := generators[varName]; exists {
					value := generator()
					if value != "" {
						variables[varName] = value
					}
				}
			}
		}

		if len(variables) == 0 {
			break
		}

		// Replace variables one by one
		tempResult := result
		for varName, value := range variables {
			// Escape special regex characters
			escapedVarName := regexp.QuoteMeta(varName)
			pattern := regexp.MustCompile(`{{\s*` + escapedVarName + `\s*}}`)
			tempResult = pattern.ReplaceAllString(tempResult, value)
		}

		if tempResult == result {
			break
		}
		result = tempResult
	}

	return result
}

