# Context Engine Go

Go implementation of LobeChat Context Engine using Eino framework for powerful graph orchestration.

## Features

- **Complete Processor Suite**: All 9 processors + 4 providers (feature parity with TypeScript version)
- **Pipeline Management**: Add, remove, disable/enable processors dynamically
- **Validation & Statistics**: Built-in pipeline validation and comprehensive stats
- **Graph Orchestration**: Eino framework with potential parallelism
- **Strong Type Safety**: Compile-time type checking
- **Streaming Support**: Built-in streaming (4 paradigms)
- **Flexible Configuration**: Deep cloning, custom processors, and more

## Installation

```bash
go get github.com/kawai-network/veridium/context-engine
```

## Quick Start

```go
package main

import (
    "context"
    "github.com/kawai-network/veridium/pkg/contextenginepkg/contextengine"
)

func main() {
    config := contextengine.Config{
        SystemRole: "You are a helpful assistant",
        HistoryCount: 10,
    }
    
    engine := contextengine.New(config)
    
    messages := []*contextengine.Message{
        {Role: "user", Content: "Hello"},
    }
    
    result, err := engine.Process(context.Background(), messages)
    if err != nil {
        panic(err)
    }
    
    // Use processed messages
    _ = result
}
```

## Architecture

This package implements the same functionality as `@lobechat/context-engine` TypeScript package but uses Eino framework for graph orchestration.

### Components

- **Processors** (9 total):
  - `GroupMessageFlatten` - Flatten group messages into assistant + tool sequences
  - `HistoryTruncate` - Limit message history
  - `InputTemplate` - Apply input templates
  - `PlaceholderVariables` - Replace placeholder variables
  - `MessageContent` - Process message content (images, videos, files)
  - `ToolCall` - Handle tool/function calls
  - `ToolMessageReorder` - Reorder tool messages
  - `MessageCleanup` - Clean up message metadata

- **Providers** (4 total):
  - `SystemRoleInjector` - Inject system role messages
  - `InboxGuide` - Add inbox guide for welcome messages
  - `ToolSystemRole` - Add tool system roles
  - `HistorySummary` - Inject history summaries

- **Graph**: Eino workflow orchestration with sequential execution

### Pipeline Management

```go
engine := contextengine.New(config)

// Add custom processor
engine.AddProcessor(contextengine.CustomProcessor{
    Name: "MyProcessor",
    Process: func(messages []*contextengine.Message) ([]*contextengine.Message, error) {
        // Your custom logic
        return messages, nil
    },
    Order: 5,
})

// Disable built-in processor
engine.DisableProcessor("HistoryTruncate")

// Validate pipeline
result := engine.Validate()
if !result.Valid {
    fmt.Println("Validation errors:", result.Errors)
}

// Get statistics
stats := engine.GetStats()
fmt.Printf("Total processors: %d\n", stats.ProcessorCount)

// Clone engine
clonedEngine := engine.Clone()
```

## License

MIT

