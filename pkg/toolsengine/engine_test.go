package toolsengine

import (
	"context"
	"testing"

	"github.com/cloudwego/eino/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewToolsEngine(t *testing.T) {
	engine, err := NewToolsEngine(Config{})
	require.NoError(t, err)
	assert.NotNil(t, engine)
}

func TestRegisterAndGetTool(t *testing.T) {
	engine, _ := NewToolsEngine(Config{})

	tool := NewTool("test", "test_tool", "Test tool", nil, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return "result", nil
	})

	err := engine.RegisterTool(tool)
	require.NoError(t, err)

	retrieved, exists := engine.GetTool("test")
	assert.True(t, exists)
	assert.Equal(t, "test", retrieved.ID)
}

func TestGenerateTools(t *testing.T) {
	engine, _ := NewToolsEngine(Config{})

	params := map[string]*schema.ParameterInfo{
		"query": {
			Type:     schema.String,
			Desc:     "Search query",
			Required: true,
		},
	}

	tool := NewTool("search", "search_tool", "Search tool", params, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return "result", nil
	})

	engine.RegisterTool(tool)

	tools, err := engine.GenerateTools(GenerateToolsParams{
		ToolIDs:  []string{"search"},
		Model:    "gpt-4",
		Provider: "openai",
	})

	require.NoError(t, err)
	assert.Equal(t, 1, len(tools))
	assert.Equal(t, "function", tools[0].Type)
	assert.Equal(t, "search_tool", tools[0].Function.Name)
}

func TestExecuteTool(t *testing.T) {
	engine, _ := NewToolsEngine(Config{})

	called := false
	tool := NewTool("test", "test_tool", "Test", nil, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		called = true
		return map[string]interface{}{
			"success": true,
			"input":   args["query"],
		}, nil
	})

	engine.RegisterTool(tool)

	result, err := engine.ExecuteTool(context.Background(), "test", `{"query":"test"}`)

	require.NoError(t, err)
	assert.True(t, called)
	assert.Contains(t, result, "success")
}
