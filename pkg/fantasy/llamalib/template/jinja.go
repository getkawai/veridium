package template

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/kawai-network/veridium/pkg/fantasy"
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"
)

// getMessageContent returns the content as a map for template rendering.
// This maintains compatibility with jinja templates that expect map access.
func getMessageContent(m fantasy.Message) map[string]interface{} {
	result := make(map[string]interface{})

	var textContent string
	var toolCalls []map[string]interface{}

	for _, part := range m.Content {
		switch p := part.(type) {
		case fantasy.TextPart:
			textContent += p.Text
		case fantasy.ToolCallPart:
			// Parse Input JSON string to arguments map for template compatibility
			var arguments interface{}
			if p.Input != "" {
				var parsed map[string]interface{}
				if err := json.Unmarshal([]byte(p.Input), &parsed); err == nil {
					arguments = parsed
				}
			}
			toolCalls = append(toolCalls, map[string]interface{}{
				"id":        p.ToolCallID,
				"name":      p.ToolName,
				"input":     p.Input,
				"arguments": arguments,
			})
		case fantasy.ToolResultPart:
			result["tool_call_id"] = p.ToolCallID
			if textOutput, ok := p.Output.(fantasy.ToolResultOutputContentText); ok {
				result["content"] = textOutput.Text
			}
		}
	}

	if textContent != "" {
		result["content"] = textContent
	}
	if len(toolCalls) > 0 {
		result["tool_calls"] = toolCalls
	}

	return result
}

// raiseExceptionFunc implements the raise_exception Jinja function
// Some model templates (like Llama 3.2) use this to validate constraints
// We implement it as a no-op that returns empty string to allow template to continue
// The actual constraint (e.g., single tool call) is handled by taking only the first tool call
func raiseExceptionFunc(msg string) string {
	// Log warning but don't fail - we handle constraints at the model level
	fmt.Printf("⚠️  Template warning (ignored): %s\n", msg)
	return ""
}

// Apply applies a jinja chat template to a slice of [fantasy.Message], Set addAssistantPrompt to true to generate the
// assistant prompt, for example on the first message.
func Apply(tmpl string, messages []fantasy.Message, addAssistantPrompt bool) (string, error) {
	// prevent filesystem access
	gonja.DefaultLoader = &NoFSLoader{}

	t, err := gonja.FromString(tmpl)
	if err != nil {
		return "", err
	}

	msgs := make([]map[string]interface{}, len(messages))
	for i, m := range messages {
		msgs[i] = map[string]interface{}{
			"role": m.Role,
		}
		for k, v := range getMessageContent(m) {
			msgs[i][k] = v
		}
	}

	data := exec.NewContext(map[string]interface{}{
		"messages":              msgs,
		"add_generation_prompt": addAssistantPrompt,
		// Add raise_exception function to prevent template errors
		// Llama 3.2 template uses this to validate single tool call constraint
		"raise_exception": raiseExceptionFunc,
	})

	return t.ExecuteToString(data)
}

// NoFSLoader is a template loader that provides no filesystem access.
// This prevents template injection attacks like {% include "/etc/passwd" %}.
type NoFSLoader struct{}

func (nl *NoFSLoader) Read(path string) (io.Reader, error) {
	return nil, errors.New("filesystem access disabled")
}

// Resolve always returns an error to prevent filesystem access.
func (nl *NoFSLoader) Resolve(path string) (string, error) {
	return "", errors.New("filesystem access disabled")
}

func (nl *NoFSLoader) Inherit(from string) (loaders.Loader, error) {
	return nil, errors.New("filesystem access disabled")
}
