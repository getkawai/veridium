# ToolsEngine - Eino Native

Simple, clean implementation of tools engine using **CloudWeGo Eino** directly.

## Architecture

```
ToolsEngine (Simple wrapper)
    ↓
Eino Tool Interface (Direct)
    ↓
Eino utils.InferTool (Auto schema generation)
```

## Key Features

✅ **Eino-Native**: Direct use of `tool.InvokableTool` interface  
✅ **Auto Schema**: Use `utils.InferTool` for automatic schema generation  
✅ **Type-Safe**: Go generics for type safety  
✅ **Simple**: Minimal wrapper, maximum Eino power  

## Quick Start

### 1. Create Tool with Auto Schema (Recommended)

```go
import (
    "github.com/cloudwego/eino/components/tool/utils"
    "github.com/kawai-network/veridium/pkg/toolsengine"
)

// Define request/response types
type SearchRequest struct {
    Query      string `json:"query" jsonschema_description:"Search query"`
    MaxResults int    `json:"max_results" jsonschema_description:"Max results"`
}

type SearchResponse struct {
    Results []string `json:"results"`
}

// Create function
func search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
    // Your search logic
    return &SearchResponse{Results: []string{"result1", "result2"}}, nil
}

// Create tool with auto schema
einoTool, _ := utils.InferTool("web_search", "Search the web", search)

// Wrap and register
tool := toolsengine.WrapEinoTool("web-search", einoTool)
engine.RegisterTool(tool)
```

### 2. Create Tool Manually (Simple)

```go
tool := toolsengine.NewTool(
    "calculator",
    "calculator",
    "Perform calculations",
    map[string]*schema.ParameterInfo{
        "expression": {
            Type:     schema.String,
            Desc:     "Math expression",
            Required: true,
        },
    },
    func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
        expr := args["expression"].(string)
        // Calculate
        return map[string]interface{}{"result": 42}, nil
    },
)

engine.RegisterTool(tool)
```

### 3. Use Eino Built-in Tools

```go
import (
    "github.com/cloudwego/eino-ext/components/tool/duckduckgo/v2"
)

// Create DuckDuckGo search tool
ddgTool, _ := duckduckgo.NewTextSearchTool(ctx, &duckduckgo.Config{
    ToolName: "web_search",
    ToolDesc: "Search the web",
})

// Wrap and register
tool := toolsengine.WrapEinoTool("web-search", ddgTool)
engine.RegisterTool(tool)
```

## Usage

### Initialize Engine

```go
engine, err := toolsengine.NewToolsEngine(toolsengine.Config{})
```

### Register Tools

```go
// Register builtin tools
builtin.RegisterAllBuiltinTools(engine)

// Register custom tool
engine.RegisterTool(myTool)
```

### Generate OpenAI Tools

```go
tools, err := engine.GenerateTools(toolsengine.GenerateToolsParams{
    ToolIDs:  []string{"web-search", "calculator"},
    Model:    "gpt-4",
    Provider: "openai",
})
```

### Execute Tool

```go
result, err := engine.ExecuteTool(ctx, "calculator", `{"expression":"2+2"}`)
```

### Get Eino Tools (for Context Engine)

```go
einoTools := engine.GetEinoTools([]string{"web-search"})

// Use in context engine
contextEngine.Process(contextengine.ProcessRequest{
    Messages: messages,
    Config: contextengine.Config{
        Tools: einoTools,
    },
})
```

## Built-in Tools

### Calculator
- **ID**: `calculator`
- **Description**: Perform mathematical calculations
- **Supports**: `+`, `-`, `*`, `/`, `sqrt()`, `sin()`, `cos()`, `tan()`, `pow()`, `pi`, `e`

```go
result, _ := engine.ExecuteTool(ctx, "calculator", `{"expression":"sqrt(16) + pow(2,3)"}`)
// {"expression":"sqrt(16) + pow(2,3)","result":12}
```

### Web Search (Mock)
- **ID**: `web-search`
- **Description**: Search the web
- **Status**: Mock implementation (TODO: integrate real API)

```go
result, _ := engine.ExecuteTool(ctx, "web-search", `{"query":"golang","max_results":5}`)
```

## Eino Integration Examples

### Using Eino's InferTool

```go
// 1. Define types with JSON schema tags
type WeatherRequest struct {
    City string `json:"city" jsonschema_description:"City name"`
    Unit string `json:"unit" jsonschema_description:"Temperature unit (C or F)"`
}

type WeatherResponse struct {
    Temperature float64 `json:"temperature"`
    Condition   string  `json:"condition"`
}

// 2. Create function
func getWeather(ctx context.Context, req *WeatherRequest) (*WeatherResponse, error) {
    return &WeatherResponse{
        Temperature: 25.5,
        Condition:   "Sunny",
    }, nil
}

// 3. Auto-generate tool
einoTool, _ := utils.InferTool("get_weather", "Get weather information", getWeather)

// 4. Wrap and register
tool := toolsengine.WrapEinoTool("weather", einoTool)
engine.RegisterTool(tool)
```

### Using Eino-Ext Tools

```go
import (
    ddg "github.com/cloudwego/eino-ext/components/tool/duckduckgo/v2"
    wikipedia "github.com/cloudwego/eino-ext/components/tool/wikipedia"
)

// DuckDuckGo Search
ddgTool, _ := ddg.NewTextSearchTool(ctx, &ddg.Config{
    ToolName:   "duckduckgo_search",
    ToolDesc:   "Search the web using DuckDuckGo",
    MaxResults: 10,
})
engine.RegisterTool(toolsengine.WrapEinoTool("ddg-search", ddgTool))

// Wikipedia
wikiTool, _ := wikipedia.NewTool(ctx, &wikipedia.Config{
    ToolName: "wikipedia",
    ToolDesc: "Search Wikipedia",
})
engine.RegisterTool(toolsengine.WrapEinoTool("wikipedia", wikiTool))
```

## API Reference

### Tool

```go
type Tool struct {
    tool.InvokableTool  // Eino interface
    
    ID       string
    Category string
    Version  string
    Enabled  bool
    Metadata map[string]interface{}
}
```

### ToolsEngine

```go
type ToolsEngine struct {
    registry *ToolRegistry
}

// Methods
func NewToolsEngine(config Config) (*ToolsEngine, error)
func (e *ToolsEngine) RegisterTool(t *Tool) error
func (e *ToolsEngine) GetTool(id string) (*Tool, bool)
func (e *ToolsEngine) GenerateTools(params GenerateToolsParams) ([]ChatCompletionTool, error)
func (e *ToolsEngine) GetEinoTools(toolIDs []string) []tool.InvokableTool
func (e *ToolsEngine) ExecuteTool(ctx context.Context, toolID string, argsJSON string) (string, error)
```

### Helper Functions

```go
// Wrap existing Eino tool
func WrapEinoTool(id string, einoTool tool.InvokableTool) *Tool

// Create tool manually
func NewTool(id, name, desc string, params map[string]*schema.ParameterInfo, executor ToolExecutor) *Tool

// Builder pattern
func NewToolBuilder(id, name string) *ToolBuilder

// Convert to OpenAI format
func ConvertToOpenAI(ctx context.Context, einoTools []tool.InvokableTool) ([]ChatCompletionTool, error)
```

## Testing

```bash
cd pkg/toolsengine
go test -v -cover
```

**Coverage**: 40.3%

## Comparison with Frontend

| Feature | Frontend (TypeScript) | Backend (Go + Eino) |
|---------|----------------------|---------------------|
| Tool Definition | Manual JSON manifest | Auto from Go types |
| Type Safety | Runtime validation | Compile-time |
| Schema Generation | Manual | Automatic |
| Execution | N/A | Built-in |
| Integration | Limited | Full Eino ecosystem |

## Eino Resources

- **Eino Core**: `/cloudwego/eino/components/tool/`
- **Eino Extensions**: `/cloudwego/eino-ext/components/tool/`
- **Examples**: `/cloudwego/eino-examples/components/tool/`

### Available Eino-Ext Tools

- ✅ **DuckDuckGo Search** - Web search
- ✅ **Bing Search** - Web search
- ✅ **Google Search** - Web search
- ✅ **Wikipedia** - Encyclopedia
- ✅ **HTTP Request** - GET/POST/PUT/DELETE
- ✅ **Command Line** - Execute shell commands
- ✅ **Python Executor** - Run Python code
- ✅ **Browser Use** - Browser automation
- ✅ **MCP** - Model Context Protocol
- ✅ **SearXNG** - Meta search engine

## Next Steps

1. **Integrate Real Search APIs**
   - Replace mock web search with DuckDuckGo
   - Add Bing/Google search options

2. **Add More Built-in Tools**
   - File operations
   - HTTP requests
   - Database queries
   - Image processing

3. **Use Eino's InferTool**
   - Migrate calculator to use `utils.InferTool`
   - Auto-generate schemas from Go structs

4. **Integrate Eino-Ext Tools**
   - Add DuckDuckGo search
   - Add Wikipedia
   - Add HTTP request tool

## License

Apache License 2.0

