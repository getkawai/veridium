package template

import (
	"os"
	"strings"
	"testing"

	"github.com/kawai-network/veridium/fantasy"
	"github.com/kawai-network/veridium/types"
)

func TestChatMLTemplate(t *testing.T) {
	// Load the chatml.jinja template
	tmplPath := "./prompts/chatml.jinja"
	tmplBytes, err := os.ReadFile(tmplPath)
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}
	tmpl := string(tmplBytes)

	// Prepare chat messages
	messages := fantasy.Prompt{
		fantasy.NewUserMessage("Hello, how are you?"),
		fantasy.Message{
			Role:    fantasy.MessageRoleAssistant,
			Content: []fantasy.MessagePart{fantasy.TextPart{Text: "I'm fine, thank you!"}},
		},
	}

	result, err := Apply(tmpl, messages, true)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	t.Logf("Apply output:\n%s", result)
	// Basic checks
	if !strings.Contains(result, "<|im_start|>user") || !strings.Contains(result, "<|im_start|>assistant") {
		t.Error("Output does not contain expected role markers")
	}
	if !strings.Contains(result, "Hello, how are you?") || !strings.Contains(result, "I'm fine, thank you!") {
		t.Error("Output does not contain expected message content")
	}
}

func TestQwen25InstructTemplateWithToolCall(t *testing.T) {
	// Load the qwen2.5-instruct.jinja template
	tmplPath := "./prompts/qwen2.5-instruct.jinja"
	tmplBytes, err := os.ReadFile(tmplPath)
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}
	tmpl := string(tmplBytes)

	// Prepare messages with a tool call
	messages := fantasy.Prompt{
		fantasy.NewUserMessage("What is 2 + 3?"),
		types.NewToolCallMessage([]types.ToolCall{
			{
				Type: "function",
				Function: types.ToolFunction{
					Name: "add",
					Arguments: map[string]string{
						"a": "2",
						"b": "3",
					},
				},
			},
		}),
		types.NewToolResultMessage("", "add", "5"),
	}

	result, err := Apply(tmpl, messages, true)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	t.Logf("Apply output:\n%s", result)

	// Basic checks
	if !strings.Contains(result, "<|im_start|>user") {
		t.Error("Output does not contain expected user role marker")
	}
	if !strings.Contains(result, "What is 2 + 3?") {
		t.Error("Output does not contain expected user message content")
	}
	if !strings.Contains(result, "<tool_call>") {
		t.Error("Output does not contain expected tool_call marker")
	}
	if !strings.Contains(result, "add") {
		t.Error("Output does not contain expected function name 'add'")
	}
	if !strings.Contains(result, "<tool_response>") {
		t.Error("Output does not contain expected tool_response marker")
	}
	if !strings.Contains(result, "5") {
		t.Error("Output does not contain expected tool response content '5'")
	}
}

func TestApplyJinjaTemplateWithToolMessage(t *testing.T) {
	// Load the qwen2.5-instruct.jinja template (supports tool calls)
	tmplPath := "./prompts/qwen2.5-instruct.jinja"
	tmplBytes, err := os.ReadFile(tmplPath)
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}
	tmpl := string(tmplBytes)

	// Prepare messages with ToolCallMessage
	messages := fantasy.Prompt{
		fantasy.NewUserMessage("Call the calculator function"),
		types.NewToolCallMessage([]types.ToolCall{
			{
				Type: "function",
				Function: types.ToolFunction{
					Name: "calculator",
					Arguments: map[string]string{
						"operation": "add",
						"x":         "10",
						"y":         "20",
					},
				},
			},
		}),
	}

	result, err := Apply(tmpl, messages, true)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	t.Logf("Apply output:\n%s", result)

	// Basic checks
	if !strings.Contains(result, "<|im_start|>user") {
		t.Error("Output does not contain expected user role marker")
	}
	if !strings.Contains(result, "Call the calculator function") {
		t.Error("Output does not contain expected user message content")
	}
	if !strings.Contains(result, "<|im_start|>assistant") {
		t.Error("Output does not contain expected assistant role marker")
	}
	if !strings.Contains(result, "<tool_call>") {
		t.Error("Output does not contain expected tool_call marker")
	}
	if !strings.Contains(result, "calculator") {
		t.Error("Output does not contain expected function name 'calculator'")
	}
}

func TestApplyJinjaTemplateWithToolResponseMessage(t *testing.T) {
	// Load the qwen2.5-instruct.jinja template (supports tool calls)
	tmplPath := "./prompts/qwen2.5-instruct.jinja"
	tmplBytes, err := os.ReadFile(tmplPath)
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}
	tmpl := string(tmplBytes)

	// Prepare messages with ToolResultMessage
	messages := fantasy.Prompt{
		fantasy.NewUserMessage("What is the result?"),
		types.NewToolCallMessage([]types.ToolCall{
			{
				Type: "function",
				Function: types.ToolFunction{
					Name: "calculator",
					Arguments: map[string]string{
						"x": "10",
						"y": "20",
					},
				},
			},
		}),
		types.NewToolResultMessage("", "calculator", "30"),
	}

	result, err := Apply(tmpl, messages, true)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	t.Logf("Apply output:\n%s", result)

	// Basic checks
	if !strings.Contains(result, "<|im_start|>user") {
		t.Error("Output does not contain expected user role marker")
	}
	if !strings.Contains(result, "What is the result?") {
		t.Error("Output does not contain expected user message content")
	}
	if !strings.Contains(result, "<tool_response>") {
		t.Error("Output does not contain expected tool_response marker")
	}
	if !strings.Contains(result, "30") {
		t.Error("Output does not contain expected tool response content '30'")
	}
}

func TestApplyJinjaTemplateWithMultipleToolCalls(t *testing.T) {
	// Load the qwen2.5-instruct.jinja template
	tmplPath := "./prompts/qwen2.5-instruct.jinja"
	tmplBytes, err := os.ReadFile(tmplPath)
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}
	tmpl := string(tmplBytes)

	// Prepare messages with multiple tool calls
	messages := fantasy.Prompt{
		fantasy.NewUserMessage("Calculate 2+3 and 5*7"),
		types.NewToolCallMessage([]types.ToolCall{
			{
				Type: "function",
				Function: types.ToolFunction{
					Name: "add",
					Arguments: map[string]string{
						"a": "2",
						"b": "3",
					},
				},
			},
			{
				Type: "function",
				Function: types.ToolFunction{
					Name: "multiply",
					Arguments: map[string]string{
						"a": "5",
						"b": "7",
					},
				},
			},
		}),
		types.NewToolResultMessage("", "add", "5"),
		types.NewToolResultMessage("", "multiply", "35"),
	}

	result, err := Apply(tmpl, messages, true)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	t.Logf("Apply output:\n%s", result)

	// Basic checks
	if !strings.Contains(result, "add") {
		t.Error("Output does not contain expected function name 'add'")
	}
	if !strings.Contains(result, "multiply") {
		t.Error("Output does not contain expected function name 'multiply'")
	}
	if !strings.Contains(result, "5") {
		t.Error("Output does not contain expected tool response '5'")
	}
	if !strings.Contains(result, "35") {
		t.Error("Output does not contain expected tool response '35'")
	}
}
