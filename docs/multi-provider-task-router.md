# Multi-Provider Task Router

## Overview

TaskRouter adalah sistem routing yang memungkinkan distribusi tugas LLM ke berbagai provider berdasarkan jenis tugas. Ini memungkinkan optimasi biaya dan performa dengan menggunakan provider yang paling sesuai untuk setiap jenis tugas.

## Arsitektur

```
┌─────────────────────────────────────────────────────────────┐
│                    AgentChatService                         │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                    TaskRouter                        │   │
│  │                                                      │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐           │   │
│  │  │   Chat   │  │  Title   │  │ Summary  │  ...      │   │
│  │  │ Provider │  │ Provider │  │ Provider │           │   │
│  │  └────┬─────┘  └────┬─────┘  └────┬─────┘           │   │
│  │       │             │             │                  │   │
│  │       ▼             ▼             ▼                  │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐           │   │
│  │  │OpenRouter│  │  Zhipu   │  │  Local   │           │   │
│  │  │ (Nova)   │  │ (GLM-4)  │  │  Llama   │           │   │
│  │  └──────────┘  └──────────┘  └──────────┘           │   │
│  │                                                      │   │
│  │  ┌──────────────────────────────────────────────┐   │   │
│  │  │            Fallback Provider                  │   │   │
│  │  │              (Local Llama)                    │   │   │
│  │  └──────────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## Task Types

| Task Type | Description | Default Provider |
|-----------|-------------|------------------|
| `chat` | Main conversation/chat | OpenRouter (amazon/nova-2-lite-v1:free) |
| `title` | Topic title generation | Zhipu GLM (glm-4.6) |
| `summary` | History summarization | Local Llama |
| `image_describe` | Image description/analysis | Zhipu GLM (glm-4v-flash) |

## Files

### Core Files

- `internal/llm/task_router.go` - TaskRouter implementation
- `internal/llm/config.go` - Hardcoded development configuration

### Integration

- `internal/services/agent_chat_service.go` - Uses TaskRouter for task distribution
- `internal/services/agent_chat_service_real_stream.go` - Streaming with TaskRouter
- `main.go` - TaskRouter initialization and logging

## Configuration

### Development Configuration (Hardcoded)

```go
// internal/llm/config.go
func GetDefaultDevConfig() *DevConfig {
    return &DevConfig{
        OpenRouterAPIKey: "sk-or-v1-...",
        ZhipuAPIKey:      "...",
        
        ChatProvider:         "openrouter",
        ChatModel:            "amazon/nova-2-lite-v1:free",
        
        TitleProvider:        "zhipu",
        TitleModel:           "glm-4.6",
        
        SummaryProvider:      "local",
        SummaryModel:         "", // auto-detect
        
        ImageDescribeProvider: "zhipu",
        ImageDescribeModel:    "glm-4v-flash",
    }
}
```

### Provider Types

| Provider Type | Endpoint |
|--------------|----------|
| `openrouter` | https://openrouter.ai/api/v1/chat/completions |
| `zhipu` | https://api.z.ai/api/paas/v4/chat/completions |
| `local` | Local Llama via llama.cpp |

## Usage

### Basic Usage

```go
// Get TaskRouter from AgentChatService
taskRouter := service.GetTaskRouter()

// Route chat task
provider := taskRouter.GetProvider(llm.TaskChat)
response, err := provider.Generate(ctx, messages)

// Route title generation
taskRouter.GenerateTitle(ctx, messages)

// Route summary generation
taskRouter.GenerateSummary(ctx, messages)
```

### With Streaming

```go
// Chat with streaming and tools
response, msgs, err := taskRouter.ChatWithTools(
    ctx, 
    messages, 
    toolNames, 
    maxIterations, 
    streamCallback, 
    toolCallback,
)
```

## Fallback Mechanism

1. If a remote provider fails, TaskRouter falls back to Local Llama
2. If no provider is configured for a task, uses fallback provider
3. Fallback provider is set to Local Llama by default

```
Request → Primary Provider → [Success] → Response
                ↓ [Failure]
          Fallback Provider → Response
```

## Logging

TaskRouter logs provider assignments at startup:

```
🔀 TaskRouter: Set fallback provider
🔀 TaskRouter: Local Llama set as fallback provider
🏭 Creating provider: type=openrouter, model=amazon/nova-2-lite-v1:free
🔀 TaskRouter: Set provider for task 'chat'
🔀 TaskRouter: Chat -> openrouter (amazon/nova-2-lite-v1:free)
🏭 Creating provider: type=zhipu, model=glm-4.6
🔀 TaskRouter: Set provider for task 'title'
🔀 TaskRouter: Title -> zhipu (glm-4.6)
```

## Data Flow

### Chat Request Flow

```
1. User sends message
2. AgentChatService.ChatRealStream() receives request
3. TaskRouter.ChatWithTools() routes to chat provider
4. OpenRouter (Nova) processes with streaming
5. Response saved to SQLite database
6. Title generation triggered in background
   └── TaskRouter.GenerateTitle() → Zhipu GLM
       └── If fails → Fallback to Local Llama
```

### Database Tables Used

| Table | Purpose |
|-------|---------|
| `sessions` | Chat sessions |
| `topics` | Conversation topics (with auto-generated titles) |
| `messages` | User, assistant, and tool messages |
| `users` | User accounts |

## Testing

### Integration Tests

```bash
# Run all integration tests
go test -tags=integration -v ./internal/services -run Integration

# Run specific test
go test -tags=integration -v ./internal/services -run TestChatRealStream_Integration
go test -tags=integration -v ./internal/services -run TestMultiProviderRouting_Integration
```

### Test Coverage

- `TestChatRealStream_Integration` - Tests real API calls to OpenRouter
- `TestMultiProviderRouting_Integration` - Verifies task routing configuration
- `TestTitleGeneration_Integration` - Tests title generation with TaskRouter

## Bug Fixes (2025-12-05)

### File Drag & Drop FOREIGN KEY Constraint Fix

**Issue**: Files dropped via drag & drop failed with `FOREIGN KEY constraint failed` error.

**Root Cause**: The drag & drop handler in `main.go` was using `"system"` as the userID:
```go
result, err := fileProcessorService.ProcessFileFromPath(filePath, "system")
```

But the `"system"` user doesn't exist in the database. The `files` table has a foreign key reference to `users(id)`.

**Fix**: Changed to use `"DEFAULT_LOBE_CHAT_USER"` which matches the frontend's default user ID:
```go
// Use DEFAULT_LOBE_CHAT_USER to match the frontend's default user ID
result, err := fileProcessorService.ProcessFileFromPath(filePath, "DEFAULT_LOBE_CHAT_USER")
```

**Affected File**: `main.go` (line ~341)

### Verified Data Flow (After Fix)

1. **Drag & Drop** → File dropped to chat input area
2. **File Copy** → `fileProcessorService.ProcessFileFromPath()` copies to `files/uploads/`
3. **SQLite files table** → File metadata stored with correct `user_id`
4. **SQLite documents table** → Document created with parsed content
5. **SQLite chunks table** → 5 text chunks created from document
6. **DuckDB vectors table** → 5 embeddings stored for semantic search
7. **Chat with attachments** → RAG retrieves relevant chunks
8. **Response** → Assistant describes file content based on RAG context

### Database State After Successful Drop

```sql
-- files table
SELECT id, name, file_type, user_id FROM files;
-- d6edcdf2-...|WhatsApp Chat...txt|txt|DEFAULT_LOBE_CHAT_USER

-- documents table  
SELECT id, title, total_char_count FROM documents;
-- b420b78b-...|...-WhatsApp Chat...txt|2540

-- chunks table (SQLite)
SELECT COUNT(*) FROM chunks WHERE document_id='b420b78b-...';
-- 5

-- vectors table (DuckDB)
SELECT COUNT(*) FROM vectors WHERE file_id='d6edcdf2-...';
-- 5
```

## Future Improvements

1. **Configuration File** - Move API keys from hardcoded to environment variables or config file
2. **Provider Health Checks** - Add health monitoring for remote providers
3. **Load Balancing** - Distribute load across multiple providers of same type
4. **Metrics** - Add provider usage metrics and monitoring
5. **Dynamic Configuration** - Allow runtime configuration updates
