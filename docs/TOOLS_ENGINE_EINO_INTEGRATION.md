# Tools Engine + Eino Integration

## Executive Summary

**YES, `toolsengine` dapat dan SEBAIKNYA menggunakan `cloudwego/eino`!**

Integrasi ini akan memberikan:
- ✅ **Native Tool Interface**: Implementasi standar untuk tool execution
- ✅ **Rich Tool Ecosystem**: 10+ pre-built tools (web search, wikipedia, command line, dll)
- ✅ **Workflow Integration**: Seamless integration dengan context engine yang sudah pakai Eino
- ✅ **Streaming Support**: Tool execution dengan streaming output
- ✅ **Better Type Safety**: JSON Schema validation built-in

## Current State

### What We Have (toolsengine v1)

```
pkg/toolsengine/
├── types.go          # Custom types
├── registry.go       # Manifest storage
├── engine.go         # Tool generation
└── service.go        # Wails binding
```

**Capabilities:**
- ✅ Manifest management
- ✅ Tool filtering & validation
- ✅ OpenAI-compatible tool generation
- ✅ 96.8% test coverage

**Limitations:**
- ❌ No actual tool execution
- ❌ No streaming support
- ❌ Manual manifest creation
- ❌ Not integrated with context engine workflow

### What Eino Provides

```
cloudwego/eino/
├── components/tool/
│   ├── interface.go       # BaseTool, InvokableTool, StreamableTool
│   └── utils/             # Helper functions for tool creation
├── schema/
│   └── tool.go            # ToolInfo, ParameterInfo, JSONSchema
└── compose/
    └── tool_node.go       # ToolsNode for workflow integration

cloudwego/eino-ext/components/tool/
├── bingsearch/            # Bing search tool
├── duckduckgo/            # DuckDuckGo search
├── wikipedia/             # Wikipedia search
├── httprequest/           # HTTP request tool
├── commandline/           # Command execution
├── browseruse/            # Browser automation
├── googlesearch/          # Google search
├── searxng/               # SearXNG meta-search
├── sequentialthinking/    # Chain-of-thought tool
└── mcp/                   # Model Context Protocol
```

## Architecture Comparison

### Current Architecture (v1)

```
Frontend (TypeScript)
    ↓
ToolsEngineService (Wails)
    ↓
ToolsEngine (Go)
    ↓
ToolRegistry
    ↓
ToolManifest (JSON) → ChatCompletionTool (OpenAI format)
```

**Flow:**
1. Frontend requests tools
2. Service generates OpenAI-compatible tool definitions
3. Frontend sends to LLM
4. LLM returns tool calls
5. **Frontend executes tools** ← Problem!

### Proposed Architecture (v2 with Eino)

```
Frontend (TypeScript)
    ↓
ToolsEngineService (Wails)
    ↓
ToolsEngine (Go)
    ├── ToolRegistry (manifest management)
    └── Eino Tool Interface
        ├── InvokableTool (execution)
        ├── StreamableTool (streaming)
        └── ToolsNode (workflow)
            ↓
Context Engine (Eino Workflow)
    ↓
Eino-ext Tools (10+ pre-built)
```

**Flow:**
1. Frontend requests tools
2. Service generates tool definitions + **registers executors**
3. Frontend sends to LLM
4. LLM returns tool calls
5. **Backend executes tools via Eino** ← Solution!
6. Results flow back through workflow

## Integration Proposal

### Phase 1: Add Eino Tool Interface

**Goal:** Make toolsengine compatible with Eino's tool interface

```go
// pkg/toolsengine/eino_adapter.go
package toolsengine

import (
    "context"
    "encoding/json"
    
    "github.com/cloudwego/eino/components/tool"
    "github.com/cloudwego/eino/schema"
)

// EinoTool wraps a ToolManifest as an Eino InvokableTool
type EinoTool struct {
    manifest  *ToolManifest
    apiDef    *APIDefinition
    executor  ToolExecutor // New: function to execute the tool
}

// ToolExecutor is a function that executes a tool
type ToolExecutor func(ctx context.Context, args map[string]interface{}) (interface{}, error)

// Info implements tool.BaseTool
func (t *EinoTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
    // Convert APIDefinition to schema.ToolInfo
    params := convertToParameterInfo(t.apiDef.Parameters)
    
    return &schema.ToolInfo{
        Name:        t.manifest.Identifier + "." + t.apiDef.Name,
        Desc:        t.apiDef.Description,
        ParamsOneOf: schema.NewParamsOneOfByParams(params),
    }, nil
}

// InvokableRun implements tool.InvokableTool
func (t *EinoTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
    // Parse arguments
    var args map[string]interface{}
    if err := json.Unmarshal([]byte(argumentsInJSON), &args); err != nil {
        return "", err
    }
    
    // Execute tool
    result, err := t.executor(ctx, args)
    if err != nil {
        return "", err
    }
    
    // Convert result to JSON string
    resultJSON, err := json.Marshal(result)
    if err != nil {
        return "", err
    }
    
    return string(resultJSON), nil
}

// NewEinoTool creates an Eino-compatible tool from a manifest
func NewEinoTool(manifest *ToolManifest, apiDef *APIDefinition, executor ToolExecutor) tool.InvokableTool {
    return &EinoTool{
        manifest: manifest,
        apiDef:   apiDef,
        executor: executor,
    }
}
```

### Phase 2: Tool Executor Registry

**Goal:** Register and manage tool executors

```go
// pkg/toolsengine/executor_registry.go
package toolsengine

import (
    "context"
    "fmt"
    "sync"
)

// ExecutorRegistry manages tool executors
type ExecutorRegistry struct {
    executors map[string]ToolExecutor
    mu        sync.RWMutex
}

// NewExecutorRegistry creates a new executor registry
func NewExecutorRegistry() *ExecutorRegistry {
    return &ExecutorRegistry{
        executors: make(map[string]ToolExecutor),
    }
}

// Register registers a tool executor
func (r *ExecutorRegistry) Register(toolID, apiName string, executor ToolExecutor) {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    key := fmt.Sprintf("%s.%s", toolID, apiName)
    r.executors[key] = executor
}

// Get retrieves a tool executor
func (r *ExecutorRegistry) Get(toolID, apiName string) (ToolExecutor, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    key := fmt.Sprintf("%s.%s", toolID, apiName)
    executor, exists := r.executors[key]
    return executor, exists
}

// Execute executes a tool
func (r *ExecutorRegistry) Execute(ctx context.Context, toolID, apiName string, args map[string]interface{}) (interface{}, error) {
    executor, exists := r.Get(toolID, apiName)
    if !exists {
        return nil, fmt.Errorf("executor not found: %s.%s", toolID, apiName)
    }
    
    return executor(ctx, args)
}
```

### Phase 3: Enhanced ToolsEngine

**Goal:** Integrate Eino tools into engine

```go
// pkg/toolsengine/engine.go (enhanced)
package toolsengine

import (
    "github.com/cloudwego/eino/components/tool"
)

// ToolsEngine (enhanced)
type ToolsEngine struct {
    registry              *ToolRegistry
    executorRegistry      *ExecutorRegistry  // New
    defaultToolIDs        []string
    enableChecker         EnableChecker
    functionCallChecker   FunctionCallChecker
    toolNameGenerator     ToolNameGenerator
}

// GenerateEinoTools generates Eino-compatible tools
func (e *ToolsEngine) GenerateEinoTools(params GenerateToolsParams) ([]tool.InvokableTool, error) {
    // Get enabled manifests
    enabledManifests, _ := e.filterEnabledTools(allToolIDs, params)
    
    // Convert to Eino tools
    einoTools := make([]tool.InvokableTool, 0)
    
    for _, manifest := range enabledManifests {
        for _, api := range manifest.API {
            // Get executor
            executor, exists := e.executorRegistry.Get(manifest.Identifier, api.Name)
            if !exists {
                // Skip if no executor registered
                continue
            }
            
            // Create Eino tool
            einoTool := NewEinoTool(&manifest, &api, executor)
            einoTools = append(einoTools, einoTool)
        }
    }
    
    return einoTools, nil
}

// RegisterToolExecutor registers a tool executor
func (e *ToolsEngine) RegisterToolExecutor(toolID, apiName string, executor ToolExecutor) {
    e.executorRegistry.Register(toolID, apiName, executor)
}
```

### Phase 4: Pre-built Tool Integration

**Goal:** Use Eino-ext tools out of the box

```go
// pkg/toolsengine/builtin/web_search.go
package builtin

import (
    "context"
    
    "github.com/kawai-network/veridium/cloudwego/eino-ext/components/tool/duckduckgo"
    "github.com/kawai-network/veridium/pkg/toolsengine"
)

// RegisterWebSearchTools registers web search tools
func RegisterWebSearchTools(engine *toolsengine.ToolsEngine) error {
    // Create DuckDuckGo tool
    ddgTool, err := duckduckgo.NewTool(context.Background())
    if err != nil {
        return err
    }
    
    // Register executor
    engine.RegisterToolExecutor("web-search", "search", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
        query := args["query"].(string)
        
        // Execute via Eino tool
        resultJSON, err := ddgTool.InvokableRun(ctx, `{"query": "`+query+`"}`)
        if err != nil {
            return nil, err
        }
        
        return resultJSON, nil
    })
    
    return nil
}
```

### Phase 5: Context Engine Integration

**Goal:** Use tools in context engineering workflow

```go
// pkg/contextengine/eino/graph.go (enhanced)
package eino

import (
    "github.com/cloudwego/eino/compose"
    "github.com/kawai-network/veridium/pkg/toolsengine"
)

// BuildContextEngineeringGraph (enhanced)
func (gb *GraphBuilder) BuildContextEngineeringGraph(toolsEngine *toolsengine.ToolsEngine) (*compose.Workflow[MessageInput, MessageOutput], error) {
    wf := compose.NewWorkflow[MessageInput, MessageOutput]()
    
    // ... existing processors ...
    
    // Add ToolsNode for tool execution
    if gb.config.Tools != nil && len(gb.config.Tools) > 0 {
        // Generate Eino tools
        einoTools, err := toolsEngine.GenerateEinoTools(toolsengine.GenerateToolsParams{
            ToolIDs:  extractToolIDs(gb.config.Tools),
            Model:    gb.config.Model,
            Provider: gb.config.Provider,
        })
        if err != nil {
            return nil, err
        }
        
        // Create ToolsNode
        toolsNode := compose.NewToolsNode(einoTools)
        
        // Add to workflow
        wf.AddNode("toolExecution", toolsNode).AddInput("toolCall")
        wf.AddNode("afterToolExecution", afterToolLambda).AddInput("toolExecution")
    }
    
    // ... rest of workflow ...
    
    return wf, nil
}
```

## Benefits of Integration

### 1. **Unified Tool Ecosystem**

**Before:**
- Frontend tools (TypeScript)
- Backend manifests (JSON)
- No execution capability

**After:**
- Unified Eino tool interface
- 10+ pre-built tools ready to use
- Backend execution with streaming

### 2. **Better Architecture**

```
┌─────────────────────────────────────────────────┐
│           Frontend (TypeScript)                  │
│  - UI for tool selection                        │
│  - Display tool results                         │
└─────────────────┬───────────────────────────────┘
                  │ Wails IPC
┌─────────────────▼───────────────────────────────┐
│        ToolsEngineService (Wails)               │
│  - GenerateTools()                              │
│  - ExecuteTool()  ← New!                        │
│  - StreamTool()   ← New!                        │
└─────────────────┬───────────────────────────────┘
                  │
┌─────────────────▼───────────────────────────────┐
│           ToolsEngine (Go)                      │
│  ┌──────────────────────────────────────────┐  │
│  │ ToolRegistry (manifest management)       │  │
│  └──────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────┐  │
│  │ ExecutorRegistry (tool execution)        │  │
│  └──────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────┐  │
│  │ Eino Adapter (interface compatibility)   │  │
│  └──────────────────────────────────────────┘  │
└─────────────────┬───────────────────────────────┘
                  │
┌─────────────────▼───────────────────────────────┐
│         Eino Tool Interface                     │
│  - InvokableTool                                │
│  - StreamableTool                               │
│  - ToolsNode (workflow integration)             │
└─────────────────┬───────────────────────────────┘
                  │
┌─────────────────▼───────────────────────────────┐
│       Eino-ext Pre-built Tools                  │
│  - Web Search (DuckDuckGo, Bing, Google)       │
│  - Wikipedia                                    │
│  - HTTP Request                                 │
│  - Command Line                                 │
│  - Browser Automation                           │
│  - Sequential Thinking                          │
│  - MCP (Model Context Protocol)                 │
└─────────────────────────────────────────────────┘
```

### 3. **Workflow Integration**

```go
// Example: Chat with tool execution in workflow
messages := []*schema.Message{
    {Role: "user", Content: "Search for latest AI news"},
}

// Context engineering with tools
result, err := contextEngine.Process(ctx, messages, contextengine.Config{
    Tools: []contextengine.Tool{
        {ID: "web-search", Enabled: true},
        {ID: "wikipedia", Enabled: true},
    },
    Model:    "gpt-4",
    Provider: "openai",
})

// Workflow automatically:
// 1. Generates tool definitions
// 2. Sends to LLM
// 3. Executes tool calls
// 4. Injects results back
// 5. Continues conversation
```

### 4. **Streaming Support**

```go
// Stream tool execution results
stream, err := toolsEngine.StreamTool(ctx, "web-search", "search", map[string]interface{}{
    "query": "latest AI news",
})

for {
    chunk, err := stream.Recv()
    if err == io.EOF {
        break
    }
    // Send chunk to frontend
    sendToFrontend(chunk)
}
```

### 5. **Type Safety**

```go
// Before: Manual JSON schema
api := APIDefinition{
    Parameters: map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "query": map[string]interface{}{
                "type": "string",
            },
        },
    },
}

// After: Type-safe schema
params := map[string]*schema.ParameterInfo{
    "query": {
        Type:     schema.String,
        Desc:     "Search query",
        Required: true,
    },
}
```

## Implementation Roadmap

### Week 1: Foundation
- [ ] Create `pkg/toolsengine/eino_adapter.go`
- [ ] Create `pkg/toolsengine/executor_registry.go`
- [ ] Add tests for adapters
- [ ] Update `ToolsEngine` with executor support

### Week 2: Integration
- [ ] Integrate Eino tools into engine
- [ ] Add `GenerateEinoTools()` method
- [ ] Add `RegisterToolExecutor()` method
- [ ] Update service with execution methods

### Week 3: Pre-built Tools
- [ ] Create `pkg/toolsengine/builtin/` package
- [ ] Integrate DuckDuckGo search
- [ ] Integrate Wikipedia
- [ ] Integrate HTTP request tool
- [ ] Add tool discovery/registration

### Week 4: Context Engine Integration
- [ ] Update context engine to use tools
- [ ] Add ToolsNode to workflow
- [ ] Add tool execution in pipeline
- [ ] Add streaming support

### Week 5: Frontend Integration
- [ ] Update Wails service
- [ ] Generate TypeScript bindings
- [ ] Create frontend tool execution UI
- [ ] Add streaming UI components

### Week 6: Testing & Documentation
- [ ] Comprehensive integration tests
- [ ] Performance benchmarks
- [ ] Update documentation
- [ ] Create examples

## Migration Path

### For Existing Users

**No breaking changes!** The current API remains:

```go
// v1 API (still works)
tools, err := engine.GenerateTools(params)

// v2 API (new, optional)
einoTools, err := engine.GenerateEinoTools(params)
```

### For New Features

```go
// 1. Register executors
engine.RegisterToolExecutor("web-search", "search", webSearchExecutor)

// 2. Generate tools with execution
einoTools, err := engine.GenerateEinoTools(params)

// 3. Use in workflow
toolsNode := compose.NewToolsNode(einoTools)
```

## Comparison Table

| Feature | Current (v1) | With Eino (v2) |
|---------|-------------|----------------|
| Manifest Management | ✅ | ✅ |
| Tool Generation | ✅ | ✅ |
| Tool Execution | ❌ | ✅ |
| Streaming | ❌ | ✅ |
| Pre-built Tools | ❌ | ✅ (10+) |
| Workflow Integration | ❌ | ✅ |
| Type Safety | Partial | ✅ Full |
| JSON Schema | Manual | ✅ Built-in |
| Context Engine Integration | ❌ | ✅ |
| Test Coverage | 96.8% | 96.8%+ |

## Conclusion

**Recommendation: INTEGRATE with Eino!**

### Why?

1. **Already using Eino**: Context engine uses Eino, so integration is natural
2. **Rich ecosystem**: 10+ pre-built tools ready to use
3. **Better architecture**: Unified tool interface, execution, and workflow
4. **No breaking changes**: Existing API remains, new features are additive
5. **Future-proof**: Eino is actively maintained by CloudWeGo

### Next Steps

1. ✅ **Approve this proposal**
2. 🔄 **Start Phase 1**: Create Eino adapter
3. 🔄 **Integrate pre-built tools**: Start with web search
4. 🔄 **Update context engine**: Add tool execution to workflow
5. 🔄 **Update frontend**: Add tool execution UI

### Timeline

- **Phase 1-2**: 2 weeks (foundation)
- **Phase 3-4**: 2 weeks (integration)
- **Phase 5-6**: 2 weeks (frontend & testing)
- **Total**: 6 weeks to full integration

---

**Status**: ✅ Ready for implementation
**Priority**: High
**Complexity**: Medium
**Impact**: High

