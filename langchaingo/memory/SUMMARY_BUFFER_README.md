# ConversationSummaryBuffer

## Overview

`ConversationSummaryBuffer` is a memory implementation that automatically summarizes old conversations to prevent LLM confusion from very long chat histories. It maintains a balance between context preservation and token efficiency.

## Problem Statement

When conversations grow very long (50+ messages), LLMs can become confused or lose track of important context. Simply keeping all messages leads to:
- **Token limit exhaustion** - Exceeding model context windows
- **Increased costs** - More tokens = higher API costs
- **Degraded performance** - LLMs struggle with very long contexts
- **Lost focus** - Important recent context gets buried

## Solution

`ConversationSummaryBuffer` solves this by:
1. **Summarizing old messages** - Creates a concise summary of older conversation parts
2. **Keeping recent messages** - Maintains the last N conversation pairs in full
3. **Background processing** - Summarization happens asynchronously (non-blocking)
4. **Progressive refinement** - Summary gets updated as conversation grows

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│  Conversation with 20 messages (10 pairs)               │
├─────────────────────────────────────────────────────────┤
│  Messages 1-14:  [SUMMARIZED]                           │
│  → "User asked about X, discussed Y, decided Z..."      │
│                                                          │
│  Messages 15-20: [KEPT IN FULL]                         │
│  → Last 3 conversation pairs for immediate context      │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│  LLM Receives:                                           │
│  1. System: "Previous conversation summary: ..."         │
│  2. Human: "Message 15..."                               │
│  3. AI: "Response 15..."                                 │
│  4. Human: "Message 16..."                               │
│  5. AI: "Response 16..."                                 │
│  6. Human: "Message 17..."                               │
│  7. AI: "Response 17..."                                 │
└─────────────────────────────────────────────────────────┘
```

## Configuration

### Default Settings

```go
const (
    DefaultRecentMessagePairs = 3    // Keep last 3 pairs (6 messages)
    DefaultMaxTokensForRecent = 1500 // ~1.5K tokens for recent messages
    DefaultSummarizeThreshold = 6    // Summarize when > 6 pairs (12 messages)
    SummaryTokenBudget        = 500  // ~500 tokens for summary
)
```

### Why 3 Conversation Pairs?

- **Sufficient context**: 3 pairs (6 messages) provide enough immediate context
- **Token efficient**: ~1000-1500 tokens (safe for all models)
- **Coherent flow**: User can follow the conversation thread
- **Not overwhelming**: Prevents confusion from too many recent messages

## Usage

### Basic Usage

```go
import (
    "context"
    "github.com/kawai-network/veridium/langchaingo/llms/ollama"
    "github.com/kawai-network/veridium/langchaingo/memory"
)

// Create an LLM
llm, _ := ollama.New(
    ollama.WithServerURL("http://localhost:11434"),
    ollama.WithModel("llama2"),
)

// Create summary buffer
summaryBuffer := memory.NewConversationSummaryBuffer(
    llm,
    memory.WithReturnMessages(true),
    memory.WithMemoryKey("history"),
)

// Use it in your chat loop
ctx := context.Background()
summaryBuffer.SaveContext(ctx,
    map[string]any{"input": "User message"},
    map[string]any{"output": "AI response"},
)

// Load memory (returns summary + recent messages)
memoryVars, _ := summaryBuffer.LoadMemoryVariables(ctx, nil)
```

### Advanced Usage with Custom Options

```go
// Custom summary options
summaryOpts := []memory.SummaryBufferOption{
    memory.WithRecentMessagePairs(5),      // Keep last 5 pairs
    memory.WithMaxTokensForRecent(2000),   // Allow up to 2K tokens
    memory.WithSummarizeThreshold(8),      // Summarize after 8 pairs
}

summaryBuffer := memory.NewConversationSummaryBufferWithOptions(
    llm,
    summaryOpts,
    memory.WithChatHistory(history),
    memory.WithReturnMessages(true),
)
```

### Integration with ChatService

```go
// In services/chat_service.go
func (c *ChatService) GetConversationMemoryWithSummary(
    sessionID string, 
    llmModel llms.Model,
) *memory.ConversationSummaryBuffer {
    history := c.sessionMemory.GetChatHistory(sessionID)
    
    summaryOpts := []memory.SummaryBufferOption{
        memory.WithRecentMessagePairs(3),
        memory.WithMaxTokensForRecent(1500),
        memory.WithSummarizeThreshold(6),
    }
    
    return memory.NewConversationSummaryBufferWithOptions(
        llmModel,
        summaryOpts,
        memory.WithChatHistory(history),
        memory.WithReturnMessages(true),
        memory.WithMemoryKey("history"),
    )
}
```

## How It Works

### 1. Message Accumulation
- Messages are stored normally in `ChatHistory`
- No summarization until threshold is reached

### 2. Threshold Detection
- After each `SaveContext()`, checks message count
- If `totalPairs > SummarizeThreshold`, triggers summarization

### 3. Background Summarization
- Runs in a goroutine (non-blocking)
- Prevents concurrent summarization with mutex
- Takes all messages except recent N pairs
- Generates/refines summary using LLM
- Updates `CurrentSummary` field

### 4. Context Loading
- `LoadMemoryVariables()` returns:
  - System message with summary (if exists)
  - Last N conversation pairs in full
- Total context stays within token budget

## Benefits

✅ **Prevents LLM Confusion** - Only sees summary + recent context  
✅ **Token Efficient** - ~2K tokens total (500 summary + 1500 recent)  
✅ **Maintains Coherence** - Recent messages provide immediate context  
✅ **Scalable** - Works for conversations of any length  
✅ **Non-blocking** - Summarization happens in background  
✅ **Configurable** - Can adjust based on model capabilities  
✅ **Progressive** - Summary improves over time  
✅ **Transparent** - Works with existing code patterns  

## Token Budget Breakdown

For a typical configuration:

```
Total Context Budget: ~2000 tokens
├─ Summary:          ~500 tokens (compressed old messages)
└─ Recent Messages:  ~1500 tokens (last 3 pairs in full)

This fits comfortably in even small model contexts (2K-4K tokens)
```

## Comparison with Other Memory Types

| Memory Type | Use Case | Token Usage | Pros | Cons |
|-------------|----------|-------------|------|------|
| `ConversationBuffer` | Short conversations | All messages | Simple, complete history | Grows unbounded |
| `ConversationTokenBuffer` | Token-limited | Prunes old messages | Token-aware | Loses context |
| `ConversationWindowBuffer` | Recent context only | Last N pairs | Predictable size | Loses old context |
| **`ConversationSummaryBuffer`** | **Long conversations** | **Summary + recent** | **Best of both worlds** | **Requires LLM** |

## Performance Considerations

### Summarization Cost
- Runs in background (doesn't block chat)
- Only triggers after threshold (not every message)
- Uses low temperature (0.3) for consistency
- Limited to 500 tokens output

### Memory Usage
- Stores summary in memory (string)
- No additional database tables needed
- Mutex for thread-safety

### Latency
- Zero impact on chat response time
- Summarization happens asynchronously
- Prevents concurrent summarization

## Future Enhancements

Potential improvements (not yet implemented):

1. **Persistent Summary Storage**
   - Save summaries to database
   - Restore on session reload

2. **Configurable Summary Prompts**
   - Custom summarization instructions
   - Domain-specific summaries

3. **Multi-level Summaries**
   - Hierarchical summarization for very long conversations
   - Summary of summaries

4. **Smart Threshold Detection**
   - Dynamic threshold based on message complexity
   - Token-based instead of count-based

5. **Summary Quality Metrics**
   - Track summary compression ratio
   - Measure information retention

## Testing

Run the example tests:

```bash
cd langchaingo/memory
go test -v -run ExampleConversationSummaryBuffer
```

## References

- LangChain Python: [ConversationSummaryBufferMemory](https://python.langchain.com/docs/modules/memory/types/summary_buffer)
- Research: [Lost in the Middle: How Language Models Use Long Contexts](https://arxiv.org/abs/2307.03172)

## License

Same as the parent project.
