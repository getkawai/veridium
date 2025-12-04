package message

import (
	"reflect"
	"testing"

	"github.com/kawai-network/veridium/types"
)

func TestToolMessage_GetRole(t *testing.T) {
	msg := Tool{Role: "tool", ToolCalls: nil}
	if got := msg.GetRole(); got != "tool" {
		t.Errorf("GetRole() = %q, want %q", got, "tool")
	}
}

func TestToolMessage_GetContent(t *testing.T) {
	msg := Tool{
		Role: "tool",
		ToolCalls: []types.ToolCall{
			{
				Type: "function",
				Function: types.ToolFunction{
					Name: "add",
					Arguments: map[string]string{
						"a": "1",
						"b": "2",
					},
				},
			},
		},
	}
	got := msg.GetContent()
	want := map[string]interface{}{
		"tool_calls": []map[string]interface{}{
			{
				"type": "function",
				"function": map[string]interface{}{
					"name": "add",
					"arguments": map[string]interface{}{
						"a": "1",
						"b": "2",
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetContent() = %v, want %v", got, want)
	}
}

func TestToolResponseMessage_GetRole(t *testing.T) {
	msg := ToolResponse{Role: "tool_response", Name: "result", Content: "42"}
	if got := msg.GetRole(); got != "tool_response" {
		t.Errorf("GetRole() = %q, want %q", got, "tool_response")
	}
}

func TestToolResponseMessage_GetContent(t *testing.T) {
	msg := ToolResponse{Role: "tool_response", Name: "result", Content: "42"}
	got := msg.GetContent()
	want := map[string]interface{}{
		"name":    "result",
		"content": "42",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetContent() = %v, want %v", got, want)
	}
}
