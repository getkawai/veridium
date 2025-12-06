package types

import (
	"reflect"
	"testing"
)

func TestMessage_GetRole(t *testing.T) {
	tests := []struct {
		name string
		msg  Message
		want string
	}{
		{"user message", NewUserMessage("Hello"), "user"},
		{"system message", NewSystemMessage("You are helpful"), "system"},
		{"assistant message", NewAssistantMessage("Hi there!"), "assistant"},
		{"tool result", NewToolResultMessage("call_123", "search", "results"), "tool"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.GetRole(); got != tt.want {
				t.Errorf("GetRole() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMessage_GetContent(t *testing.T) {
	t.Run("text message", func(t *testing.T) {
		msg := NewAssistantMessage("Hi there!")
		content := msg.GetContent()
		if content["content"] != "Hi there!" {
			t.Errorf("GetContent() content = %v, want %v", content["content"], "Hi there!")
		}
	})

	t.Run("tool call message", func(t *testing.T) {
		toolCalls := []ToolCall{
			{
				ID:   "call_123",
				Type: "function",
				Function: ToolFunction{
					Name: "add",
					Arguments: map[string]string{
						"a": "1",
						"b": "2",
					},
				},
			},
		}
		msg := NewToolCallMessage(toolCalls)
		content := msg.GetContent()

		calls, ok := content["tool_calls"].([]map[string]interface{})
		if !ok || len(calls) != 1 {
			t.Errorf("GetContent() tool_calls not found or wrong length")
			return
		}

		fn := calls[0]["function"].(map[string]interface{})
		if fn["name"] != "add" {
			t.Errorf("tool call name = %v, want %v", fn["name"], "add")
		}
	})

	t.Run("tool result message", func(t *testing.T) {
		msg := NewToolResultMessage("call_123", "search", "found 5 results")
		content := msg.GetContent()

		if content["name"] != "search" {
			t.Errorf("GetContent() name = %v, want %v", content["name"], "search")
		}
		if content["content"] != "found 5 results" {
			t.Errorf("GetContent() content = %v, want %v", content["content"], "found 5 results")
		}
		if content["tool_call_id"] != "call_123" {
			t.Errorf("GetContent() tool_call_id = %v, want %v", content["tool_call_id"], "call_123")
		}
	})
}

func TestMessage_GetText(t *testing.T) {
	msg := NewUserMessage("Hello world")
	if got := msg.GetText(); got != "Hello world" {
		t.Errorf("GetText() = %q, want %q", got, "Hello world")
	}
}

func TestMessage_GetToolCalls(t *testing.T) {
	toolCalls := []ToolCall{
		{ID: "call_1", Type: "function", Function: ToolFunction{Name: "func1"}},
		{ID: "call_2", Type: "function", Function: ToolFunction{Name: "func2"}},
	}
	msg := NewToolCallMessage(toolCalls)

	got := msg.GetToolCalls()
	if len(got) != 2 {
		t.Errorf("GetToolCalls() length = %d, want %d", len(got), 2)
	}
	if got[0].ID != "call_1" || got[1].ID != "call_2" {
		t.Errorf("GetToolCalls() IDs mismatch")
	}
}

func TestMessage_HasToolCalls(t *testing.T) {
	t.Run("message with tool calls", func(t *testing.T) {
		msg := NewToolCallMessage([]ToolCall{{ID: "1"}})
		if !msg.HasToolCalls() {
			t.Error("HasToolCalls() = false, want true")
		}
	})

	t.Run("message without tool calls", func(t *testing.T) {
		msg := NewAssistantMessage("Hello")
		if msg.HasToolCalls() {
			t.Error("HasToolCalls() = true, want false")
		}
	})
}

func TestNewUserMessage(t *testing.T) {
	t.Run("text only", func(t *testing.T) {
		msg := NewUserMessage("Hello")
		if msg.Role != MessageRoleUser {
			t.Errorf("Role = %v, want %v", msg.Role, MessageRoleUser)
		}
		if len(msg.Content) != 1 {
			t.Errorf("Content length = %d, want %d", len(msg.Content), 1)
		}
	})

	t.Run("with files", func(t *testing.T) {
		file := FilePart{
			Filename:  "test.png",
			Data:      []byte{0x89, 0x50},
			MediaType: "image/png",
		}
		msg := NewUserMessage("Check this", file)
		if len(msg.Content) != 2 {
			t.Errorf("Content length = %d, want %d", len(msg.Content), 2)
		}
	})
}

func TestNewToolResultMessage(t *testing.T) {
	msg := NewToolResultMessage("call_123", "search", "results")

	if msg.Role != MessageRoleTool {
		t.Errorf("Role = %v, want %v", msg.Role, MessageRoleTool)
	}

	if len(msg.Content) != 1 {
		t.Fatalf("Content length = %d, want 1", len(msg.Content))
	}

	part, ok := msg.Content[0].(ToolResultPart)
	if !ok {
		t.Fatalf("Content[0] is not ToolResultPart")
	}

	if part.ToolCallID != "call_123" {
		t.Errorf("ToolCallID = %v, want %v", part.ToolCallID, "call_123")
	}
	if part.ToolName != "search" {
		t.Errorf("ToolName = %v, want %v", part.ToolName, "search")
	}
	if part.Content != "results" {
		t.Errorf("Content = %v, want %v", part.Content, "results")
	}
	if part.IsError {
		t.Error("IsError = true, want false")
	}
}

func TestNewToolErrorMessage(t *testing.T) {
	msg := NewToolErrorMessage("call_456", "fetch", "network error")

	part := msg.Content[0].(ToolResultPart)
	if !part.IsError {
		t.Error("IsError = false, want true")
	}
}

func TestPartTypes(t *testing.T) {
	tests := []struct {
		name string
		part MessagePart
		want MessageContentType
	}{
		{"TextPart", TextPart{Text: "hello"}, MessageContentTypeText},
		{"ReasoningPart", ReasoningPart{Text: "thinking..."}, MessageContentTypeReasoning},
		{"FilePart", FilePart{Filename: "test.txt"}, MessageContentTypeFile},
		{"ToolCallPart", ToolCallPart{}, MessageContentTypeToolCall},
		{"ToolResultPart", ToolResultPart{}, MessageContentTypeToolResult},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.part.GetType(); got != tt.want {
				t.Errorf("GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrompt(t *testing.T) {
	prompt := Prompt{
		NewSystemMessage("You are helpful"),
		NewUserMessage("Hello"),
		NewAssistantMessage("Hi!"),
	}

	if len(prompt) != 3 {
		t.Errorf("Prompt length = %d, want 3", len(prompt))
	}

	if prompt[0].GetRole() != "system" {
		t.Errorf("prompt[0] role = %v, want system", prompt[0].GetRole())
	}
}

func TestNewTextMessage(t *testing.T) {
	tests := []struct {
		role MessageRole
		text string
	}{
		{MessageRoleUser, "user message"},
		{MessageRoleAssistant, "assistant message"},
		{MessageRoleSystem, "system message"},
	}

	for _, tt := range tests {
		t.Run(string(tt.role), func(t *testing.T) {
			msg := NewTextMessage(tt.role, tt.text)
			if msg.Role != tt.role {
				t.Errorf("Role = %v, want %v", msg.Role, tt.role)
			}
			if msg.GetText() != tt.text {
				t.Errorf("GetText() = %v, want %v", msg.GetText(), tt.text)
			}
		})
	}
}

func TestGetContent_Empty(t *testing.T) {
	msg := Message{Role: MessageRoleUser, Content: []MessagePart{}}
	content := msg.GetContent()

	if len(content) != 0 {
		t.Errorf("GetContent() should be empty for message with no parts, got %v", content)
	}
}

func TestFilePart_GetType(t *testing.T) {
	f := FilePart{
		Filename:  "image.png",
		Data:      []byte{0x89, 0x50, 0x4E, 0x47},
		MediaType: "image/png",
	}

	if got := f.GetType(); got != MessageContentTypeFile {
		t.Errorf("GetType() = %v, want %v", got, MessageContentTypeFile)
	}

	if !reflect.DeepEqual(f.Data, []byte{0x89, 0x50, 0x4E, 0x47}) {
		t.Error("FilePart data mismatch")
	}
}
