package toolsengine

import (
	"context"
	"testing"

	"github.com/cloudwego/eino/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Example 1: Simple map-based tool (backward compatible)
func TestNewTool_Simple(t *testing.T) {
	tool := NewTool(
		"calculator",
		"calculator",
		"Simple calculator",
		map[string]*schema.ParameterInfo{
			"expression": {
				Type:     schema.String,
				Desc:     "Math expression",
				Required: true,
			},
		},
		func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
			return map[string]interface{}{
				"result": 42,
			}, nil
		},
	)

	require.NotNil(t, tool)
	assert.Equal(t, "calculator", tool.ID)

	// Execute
	result, err := tool.InvokableRun(context.Background(), `{"expression":"2+2"}`)
	require.NoError(t, err)
	assert.Contains(t, result, "result")
}

// Example 2: Type-safe tool with utils.NewTool
func TestNewTypedTool(t *testing.T) {
	// Define types
	type CalcRequest struct {
		Expression string `json:"expression"`
	}
	type CalcResponse struct {
		Result float64 `json:"result"`
	}

	// Create function
	fn := func(ctx context.Context, req *CalcRequest) (*CalcResponse, error) {
		return &CalcResponse{Result: 42}, nil
	}

	// Create tool info
	toolInfo := &schema.ToolInfo{
		Name: "calculator",
		Desc: "Calculate expressions",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"expression": {
				Type:     schema.String,
				Desc:     "Math expression",
				Required: true,
			},
		}),
	}

	// Create tool
	tool := NewTypedTool("calc", toolInfo, fn)

	require.NotNil(t, tool)
	assert.Equal(t, "calc", tool.ID)

	// Execute
	result, err := tool.InvokableRun(context.Background(), `{"expression":"2+2"}`)
	require.NoError(t, err)
	assert.Contains(t, result, "result")
}

// Example 3: Auto-inferred tool (most convenient)
func TestInferTool(t *testing.T) {
	// Define types with JSON schema tags
	type SearchRequest struct {
		Query      string `json:"query" jsonschema_description:"Search query"`
		MaxResults int    `json:"max_results" jsonschema_description:"Maximum number of results"`
	}
	type SearchResponse struct {
		Results []string `json:"results" jsonschema_description:"Search results"`
	}

	// Create function
	fn := func(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
		return &SearchResponse{
			Results: []string{"result1", "result2"},
		}, nil
	}

	// Auto-infer tool (schema generated automatically!)
	tool, err := InferTool("search", "web_search", "Search the web", fn)

	require.NoError(t, err)
	require.NotNil(t, tool)
	assert.Equal(t, "search", tool.ID)

	// Verify tool info
	info, err := tool.Info(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "web_search", info.Name)
	assert.Equal(t, "Search the web", info.Desc)

	// Execute
	result, err := tool.InvokableRun(context.Background(), `{"query":"golang","max_results":5}`)
	require.NoError(t, err)
	assert.Contains(t, result, "results")
}

// Example 4: Builder pattern (fluent API)
func TestToolBuilder(t *testing.T) {
	tool, err := NewToolBuilder("weather", "get_weather").
		WithDescription("Get weather information").
		WithParameter("city", schema.String, "City name", true).
		WithParameter("unit", schema.String, "Temperature unit (C or F)", false).
		WithCategory("utility").
		WithVersion("1.0.0").
		WithExecutor(func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
			return map[string]interface{}{
				"temperature": 25.5,
				"condition":   "Sunny",
			}, nil
		}).
		Build()

	require.NoError(t, err)
	require.NotNil(t, tool)
	assert.Equal(t, "weather", tool.ID)
	assert.Equal(t, "utility", tool.Category)
	assert.Equal(t, "1.0.0", tool.Version)
}

// Benchmark: Compare different approaches
func BenchmarkToolExecution_Simple(b *testing.B) {
	tool := NewTool("test", "test", "Test", nil, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return "result", nil
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = tool.InvokableRun(context.Background(), `{"test":"value"}`)
	}
}

func BenchmarkToolExecution_Typed(b *testing.B) {
	type Req struct{ Test string }
	type Resp struct{ Result string }

	fn := func(ctx context.Context, req *Req) (*Resp, error) {
		return &Resp{Result: "result"}, nil
	}

	toolInfo := &schema.ToolInfo{Name: "test", Desc: "Test"}
	tool := NewTypedTool("test", toolInfo, fn)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = tool.InvokableRun(context.Background(), `{"test":"value"}`)
	}
}

func BenchmarkToolExecution_Inferred(b *testing.B) {
	type Req struct {
		Test string `json:"test" jsonschema_description:"Test value"`
	}
	type Resp struct {
		Result string `json:"result"`
	}

	fn := func(ctx context.Context, req *Req) (*Resp, error) {
		return &Resp{Result: "result"}, nil
	}

	tool, _ := InferTool("test", "test", "Test tool", fn)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = tool.InvokableRun(context.Background(), `{"test":"value"}`)
	}
}

