# Agent Chat Service Documentation

## Daftar Isi
- [Overview](#overview)
- [Arsitektur](#arsitektur)
- [Komponen Utama](#komponen-utama)
- [Alur Kerja](#alur-kerja)
- [API Reference](#api-reference)
- [Fitur-Fitur](#fitur-fitur)
- [Konfigurasi](#konfigurasi)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

---

## Overview

**Agent Chat Service** adalah layanan chat berbasis AI yang menggunakan framework Eino (CloudWeGo) untuk menyediakan kemampuan conversational AI dengan fitur-fitur canggih seperti:

- 🤖 **Agent-based Architecture**: Menggunakan ADK (Agent Development Kit) dari Eino
- 📚 **RAG (Retrieval-Augmented Generation)**: Integrasi dengan Knowledge Base
- 🛠️ **Tool Integration**: Dukungan untuk tools eksternal
- 🧠 **Reasoning Modes**: Tiga mode reasoning (Disabled, Enabled, Verbose)
- 💾 **Persistent Storage**: Hybrid DB + in-memory caching
- 📋 **Auto Summarization**: Kompresi history otomatis
- 🔄 **Streaming Support**: Token-by-token streaming
- 🌳 **Thread Management**: Dukungan untuk conversation branching

### File Location
```
internal/services/agent_chat_service.go
```

### Dependencies
```go
- github.com/cloudwego/eino/adk          // Agent Development Kit
- github.com/cloudwego/eino/components/tool
- github.com/cloudwego/eino/compose
- github.com/cloudwego/eino/schema
- internal/database                      // Database layer
- internal/llama                         // LLM integration
- github.com/wailsapp/wails/v3          // Desktop app framework
```

---

## Arsitektur

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Frontend (TypeScript)                     │
│                  - React Components                          │
│                  - Event Listeners                           │
└────────────────────────┬────────────────────────────────────┘
                         │ Wails Bridge
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                  AgentChatService (Go)                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   Session    │  │   Context    │  │    Tools     │     │
│  │  Management  │  │    Engine    │  │    Engine    │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │     RAG      │  │   History    │  │   Thread     │     │
│  │   Workflow   │  │   Summary    │  │  Management  │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    Eino Framework                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │  ADK Agent   │  │  ChatModel   │  │  Tool System │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                   LLM Layer (llama.cpp)                      │
│  - Model Loading & Management                                │
│  - Token Generation                                          │
│  - Streaming Support                                         │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                  Database (SQLite)                           │
│  - Sessions, Messages, Topics, Threads                       │
│  - Knowledge Base                                            │
└─────────────────────────────────────────────────────────────┘
```

### Component Diagram

```
AgentChatService
├── Session Management
│   ├── In-Memory Cache (map[string]*AgentSession)
│   └── Database Persistence
│
├── Agent Creation & Execution
│   ├── ADK Agent (Eino)
│   ├── LlamaEinoModel (LLM Adapter)
│   └── Tool Integration
│
├── RAG Workflow
│   ├── Knowledge Base Service
│   └── Vector Search
│
├── Context Processing
│   ├── Context Engine Bridge
│   └── Message Processing
│
├── History Management
│   ├── Auto Summarization
│   ├── Incremental Summary
│   └── Context Compression
│
├── Thread Management
│   ├── Topic Creation
│   ├── Thread Branching
│   └── Parent-Child Relationships
│
└── Utility Models
    ├── Title Generation (Small Model)
    └── Summary Generation (Small Model)
```

---

## Komponen Utama

### 1. AgentChatService

**Struct Definition:**
```go
type AgentChatService struct {
    app           *application.App
    db            *database.Service
    libService    *llama.LibraryService
    llamaModel    *llama.LlamaEinoModel
    kbService     *KnowledgeBaseService
    ragWorkflow   *RAGWorkflow
    
    // Bridges
    toolsBridge   *ToolsEngineBridge
    contextBridge *ContextEngineBridge
    threadService *ThreadManagementService
    
    // Configuration
    reasoningConfig ReasoningConfig
    
    // Utility Models
    titleModelPath   string
    summaryModelPath string
    
    // Session Cache
    sessions      map[string]*AgentSession
    sessionsMutex sync.RWMutex
}
```

**Responsibilities:**
- Mengelola lifecycle chat sessions
- Orchestrasi antara berbagai komponen (RAG, Tools, Context)
- Menangani streaming dan non-streaming responses
- Auto-summarization dan context management
- Model selection dan optimization

### 2. AgentSession

**Struct Definition:**
```go
type AgentSession struct {
    SessionID       string
    UserID          string
    Agent           adk.Agent
    Messages        []*schema.Message
    KnowledgeBaseID string
    Tools           []tool.BaseTool
    Context         map[string]any
    CreatedAt       int64
    UpdatedAt       int64
    DBSession       *db.Session
    
    // Topic & Thread
    TopicID  string
    ThreadID string
}
```

**Lifecycle:**
1. **Creation**: Saat user mengirim pesan pertama
2. **Caching**: Disimpan di memory untuk akses cepat
3. **Persistence**: Disinkronkan ke database
4. **Restoration**: Dimuat dari DB jika cache miss
5. **Cleanup**: Dihapus dari cache saat tidak aktif

### 3. ChatRequest & ChatResponse

**Request Structure:**
```go
type ChatRequest struct {
    // Identity
    SessionID string
    UserID    string
    
    // Content
    Message string
    
    // Context
    TopicID  string  // Auto-created if empty
    ThreadID string  // For branching
    ParentID string  // Parent message
    
    // Configuration
    KnowledgeBaseID string
    Tools           []string
    Context         map[string]any
    Temperature     float32
    MaxTokens       int
    Stream          bool
}
```

**Response Structure:**
```go
type ChatResponse struct {
    // IDs
    MessageID string
    SessionID string
    TopicID   string
    ThreadID  string
    
    // Content
    Message      string
    ToolCalls    []schema.ToolCall
    Sources      []*schema.Document
    FinishReason string
    Usage        *schema.TokenUsage
    
    // Metadata
    CreatedAt int64
    Error     string
}
```

---

## Alur Kerja

### 1. Chat Request Flow

```
┌─────────────────────────────────────────────────────────────┐
│ 1. User sends message                                        │
│    - ChatRequest with message, sessionID, userID            │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 2. Get or Create Session                                     │
│    a. Check in-memory cache                                  │
│    b. If not found, load from DB                            │
│    c. If still not found, create new session                │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 3. Load Context (if exists)                                  │
│    - Load history summary from topic                         │
│    - Load thread messages (if threadID provided)            │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 4. Validate & Auto-Switch Model                             │
│    - Check if current model matches reasoning mode          │
│    - Auto-switch to recommended model if needed             │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 5. Create Topic (if first message)                          │
│    - Auto-create topic with placeholder title               │
│    - Get topicID for message saving                         │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 6. Save User Message to DB                                   │
│    - Store with sessionID, topicID, threadID                │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 7. Process Messages through Context Engine                   │
│    - Apply context transformations                           │
│    - Inject history summary                                  │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 8. Generate Response                                         │
│    ┌──────────────────┐  ┌──────────────────┐              │
│    │  Streaming Mode  │  │ Non-Streaming    │              │
│    │  (Token-by-Token)│  │ (Eino Agent)     │              │
│    └──────────────────┘  └──────────────────┘              │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 9. Save Assistant Message to DB                             │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 10. Post-Processing (Background)                            │
│     - Update topic title (if 2-4 messages)                  │
│     - Auto-summarize (if threshold reached)                 │
│     - Incremental summary (if summary exists)               │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 11. Return Response to User                                  │
│     - MessageID, TopicID, ThreadID                          │
│     - Content, ToolCalls, Sources                           │
└─────────────────────────────────────────────────────────────┘
```

### 2. Session Management Flow

```
┌─────────────────────────────────────────────────────────────┐
│ getOrCreateSession(sessionID, userID)                        │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
              ┌──────────────┐
              │ Cache Hit?   │
              └──────┬───────┘
                     │
         ┌───────────┴───────────┐
         │ YES                   │ NO
         ▼                       ▼
┌─────────────────┐    ┌─────────────────────┐
│ Return Cached   │    │ Load from Database  │
│ Session         │    └──────────┬──────────┘
└─────────────────┘               │
                          ┌───────┴────────┐
                          │ Found in DB?   │
                          └───────┬────────┘
                                  │
                      ┌───────────┴───────────┐
                      │ YES                   │ NO
                      ▼                       ▼
            ┌──────────────────┐    ┌──────────────────┐
            │ Reconstruct from │    │ Create New       │
            │ DB History       │    │ Session          │
            └────────┬─────────┘    └────────┬─────────┘
                     │                       │
                     │                       ▼
                     │              ┌──────────────────┐
                     │              │ Save to DB       │
                     │              └────────┬─────────┘
                     │                       │
                     └───────────┬───────────┘
                                 │
                                 ▼
                    ┌──────────────────────────┐
                    │ Create Agent with Tools  │
                    └──────────┬───────────────┘
                               │
                               ▼
                    ┌──────────────────────────┐
                    │ Cache in Memory          │
                    └──────────┬───────────────┘
                               │
                               ▼
                    ┌──────────────────────────┐
                    │ Return Session           │
                    └──────────────────────────┘
```

### 3. History Summarization Flow

```
┌─────────────────────────────────────────────────────────────┐
│ After Each Chat Response (Background Goroutine)             │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ autoSummarizeIfNeeded()                                      │
│ - Check if summary threshold reached                         │
│ - Check if summary doesn't exist yet                        │
└────────────────────┬────────────────────────────────────────┘
                     │
         ┌───────────┴───────────┐
         │ Threshold Reached?    │
         │ AND No Summary?       │
         └───────────┬───────────┘
                     │
         ┌───────────┴───────────┐
         │ YES                   │ NO → Skip
         ▼                       
┌─────────────────────────────────────────────────────────────┐
│ Generate Initial Summary                                     │
│ 1. Get old messages (exclude recent N messages)             │
│ 2. Use summary model (Llama 3.2 1B/3B)                      │
│ 3. Generate compressed summary                               │
│ 4. Save to topic.history_summary                            │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ incrementalSummarizeIfNeeded()                              │
│ - Check if summary exists                                    │
│ - Check if enough new messages since last summary           │
└────────────────────┬────────────────────────────────────────┘
                     │
         ┌───────────┴───────────┐
         │ Enough New Messages?  │
         └───────────┬───────────┘
                     │
         ┌───────────┴───────────┐
         │ YES                   │ NO → Skip
         ▼                       
┌─────────────────────────────────────────────────────────────┐
│ Generate Incremental Summary                                 │
│ 1. Load existing summary                                     │
│ 2. Get new messages since last summary                      │
│ 3. Merge old summary + new messages                         │
│ 4. Update topic.history_summary                             │
│ 5. Update metadata (version, message count)                 │
└─────────────────────────────────────────────────────────────┘
```

---

## API Reference

### Constructor

#### NewAgentChatService

```go
func NewAgentChatService(
    app *application.App,
    db *database.Service,
    libService *llama.LibraryService,
    kbService *KnowledgeBaseService,
    toolsBridge *ToolsEngineBridge,
    contextBridge *ContextEngineBridge,
    threadService *ThreadManagementService,
) *AgentChatService
```

**Parameters:**
- `app`: Wails application instance untuk event emission
- `db`: Database service untuk persistence
- `libService`: Llama library service untuk model management
- `kbService`: Knowledge base service untuk RAG
- `toolsBridge`: Bridge ke tools engine (optional, dapat nil)
- `contextBridge`: Bridge ke context engine (optional, dapat nil)
- `threadService`: Thread management service (optional, dapat nil)

**Returns:**
- Initialized `*AgentChatService`

**Auto-Detection:**
- Deteksi title generation model (Llama 3.2 1B/3B)
- Deteksi summary generation model (Llama 3.2 1B/3B)
- Set default reasoning config (ReasoningDisabled)

**Example:**
```go
agentService := NewAgentChatService(
    app,
    dbService,
    llamaLibService,
    kbService,
    toolsBridge,
    contextBridge,
    threadService,
)
```

---

### Core Methods

#### Chat

```go
func (s *AgentChatService) Chat(
    ctx context.Context,
    req ChatRequest,
) (*ChatResponse, error)
```

**Description:**
Main method untuk memproses chat request dan menghasilkan response.

**Parameters:**
- `ctx`: Context untuk cancellation dan timeout
- `req`: ChatRequest dengan message dan konfigurasi

**Returns:**
- `*ChatResponse`: Response dengan message, IDs, dan metadata
- `error`: Error jika terjadi kegagalan

**Flow:**
1. Get/create session
2. Load context (summary, thread)
3. Validate model
4. Create topic (if first message)
5. Save user message
6. Process through context engine
7. Generate response (streaming/non-streaming)
8. Save assistant message
9. Background: Update title, auto-summarize
10. Return response

**Example:**
```go
response, err := agentService.Chat(ctx, ChatRequest{
    SessionID: "session-123",
    UserID:    "user-456",
    Message:   "What is the capital of France?",
    Stream:    true,
})
if err != nil {
    log.Fatal(err)
}
fmt.Println(response.Message)
```

---

#### ClearSession

```go
func (s *AgentChatService) ClearSession(sessionID string)
```

**Description:**
Menghapus session dari in-memory cache.

**Parameters:**
- `sessionID`: ID session yang akan dihapus

**Note:**
- Tidak menghapus dari database
- Hanya menghapus dari cache
- Session akan dimuat ulang dari DB jika diakses lagi

**Example:**
```go
agentService.ClearSession("session-123")
```

---

#### GetSessionHistory

```go
func (s *AgentChatService) GetSessionHistory(
    sessionID string,
) ([]*schema.Message, error)
```

**Description:**
Mengambil history messages dari session.

**Parameters:**
- `sessionID`: ID session

**Returns:**
- `[]*schema.Message`: Array of messages
- `error`: Error jika session tidak ditemukan

**Example:**
```go
messages, err := agentService.GetSessionHistory("session-123")
if err != nil {
    log.Fatal(err)
}
for _, msg := range messages {
    fmt.Printf("%s: %s\n", msg.Role, msg.Content)
}
```

---

### Reasoning Mode Methods

#### SetReasoningMode

```go
func (s *AgentChatService) SetReasoningMode(mode ReasoningMode) error
```

**Description:**
Mengatur reasoning mode untuk service.

**Parameters:**
- `mode`: ReasoningMode (Disabled, Enabled, Verbose)

**Returns:**
- `error`: Error jika mode invalid atau hardware tidak memadai

**Validation:**
- Validasi hardware specs
- Auto-switch ke mode yang sesuai jika hardware tidak cukup
- Log hardware requirements dan expected performance

**Example:**
```go
err := agentService.SetReasoningMode(ReasoningEnabled)
if err != nil {
    log.Fatal(err)
}
```

---

#### GetReasoningMode

```go
func (s *AgentChatService) GetReasoningMode() ReasoningMode
```

**Description:**
Mengambil current reasoning mode.

**Returns:**
- `ReasoningMode`: Current mode

**Example:**
```go
mode := agentService.GetReasoningMode()
fmt.Printf("Current mode: %s\n", mode)
```

---

#### ValidateModelForReasoningMode

```go
func (s *AgentChatService) ValidateModelForReasoningMode() error
```

**Description:**
Validasi apakah model yang dimuat sesuai dengan reasoning mode.

**Returns:**
- `error`: Error jika model tidak sesuai

**Example:**
```go
if err := agentService.ValidateModelForReasoningMode(); err != nil {
    log.Printf("Model mismatch: %v", err)
    agentService.SwitchToRecommendedModel()
}
```

---

#### SwitchToRecommendedModel

```go
func (s *AgentChatService) SwitchToRecommendedModel() error
```

**Description:**
Load model yang direkomendasikan untuk current reasoning mode.

**Returns:**
- `error`: Error jika gagal load model

**Example:**
```go
err := agentService.SwitchToRecommendedModel()
if err != nil {
    log.Fatal(err)
}
```

---

### Utility Methods

#### SetTitleModel

```go
func (s *AgentChatService) SetTitleModel(modelPath string)
```

**Description:**
Set model khusus untuk title generation.

**Parameters:**
- `modelPath`: Path ke model file

**Recommendation:**
- Gunakan model kecil (Llama 3.2 1B/3B)
- Non-reasoning model (no <think> tags)

**Example:**
```go
agentService.SetTitleModel("/path/to/llama-3.2-1b-instruct.gguf")
```

---

### Internal Methods

#### generateTopicTitle

```go
func (s *AgentChatService) generateTopicTitle(
    ctx context.Context,
    messages []*schema.Message,
    locale string,
) (string, error)
```

**Description:**
Generate title untuk conversation menggunakan LLM.

**Strategy:**
1. Gunakan title model (jika ada)
2. Fallback ke main model
3. Strip <think> tags jika ada
4. Truncate ke 50 karakter

**Parameters:**
- `ctx`: Context
- `messages`: Messages untuk di-summarize
- `locale`: Locale code (e.g., "en-US", "id-ID")

**Returns:**
- `string`: Generated title
- `error`: Error jika gagal

---

#### generateHistorySummary

```go
func (s *AgentChatService) generateHistorySummary(
    ctx context.Context,
    messages []*schema.Message,
) (string, error)
```

**Description:**
Generate summary dari history messages.

**Strategy (3-tier fallback):**
1. Summary model (BEST - Llama 3.2 1B/3B)
2. Title model (GOOD - small & fast)
3. Main model (FALLBACK - may be slow)

**Parameters:**
- `ctx`: Context
- `messages`: Messages untuk di-summarize

**Returns:**
- `string`: Generated summary (max 400 tokens)
- `error`: Error jika gagal

**Features:**
- Auto-detect optimal model
- Strip <think> tags
- Preserve key information
- Maintain chronological flow

---

#### generateIncrementalSummary

```go
func (s *AgentChatService) generateIncrementalSummary(
    ctx context.Context,
    existingSummary string,
    newMessages []*schema.Message,
) (string, error)
```

**Description:**
Update existing summary dengan new messages.

**Strategy:**
- Merge old summary + new messages
- Preserve important information
- Add new topics and developments
- Maintain conciseness

**Parameters:**
- `ctx`: Context
- `existingSummary`: Existing summary text
- `newMessages`: New messages to incorporate

**Returns:**
- `string`: Updated summary
- `error`: Error jika gagal

---

## Fitur-Fitur

### 1. Reasoning Modes

#### ReasoningDisabled (Default)

**Characteristics:**
- Model: Llama 3.2 1B/3B (non-reasoning)
- Speed: Fast (2-5 tokens/sec)
- Context: 16K tokens
- Max Turns: 40-50 turns
- Token Efficiency: High (no thinking overhead)

**Use Cases:**
- General conversation
- Quick Q&A
- Long conversations
- Low-resource systems

**Hardware Requirements:**
- RAM: 4GB minimum
- CPU: 4 cores minimum
- GPU: Optional

**Example:**
```go
agentService.SetReasoningMode(ReasoningDisabled)
```

---

#### ReasoningEnabled

**Characteristics:**
- Model: Qwen2.5-3B-Instruct with /no_think
- Speed: Medium (1-3 tokens/sec)
- Context: 16K tokens
- Max Turns: 15-20 turns
- Token Efficiency: Medium (reasoning compressed)

**Use Cases:**
- Complex problem solving
- Code generation
- Analysis tasks
- Medium-length conversations

**Hardware Requirements:**
- RAM: 8GB minimum
- CPU: 6 cores recommended
- GPU: Recommended

**Example:**
```go
agentService.SetReasoningMode(ReasoningEnabled)
```

---

#### ReasoningVerbose

**Characteristics:**
- Model: Qwen2.5-3B-Instruct (full reasoning)
- Speed: Slow (0.5-2 tokens/sec)
- Context: 16K tokens
- Max Turns: 3-5 turns only
- Token Efficiency: Low (full thinking process)

**Use Cases:**
- Deep analysis
- Debugging complex issues
- Educational purposes (show thinking)
- Short, focused conversations

**Hardware Requirements:**
- RAM: 16GB minimum
- CPU: 8 cores recommended
- GPU: Highly recommended

**Example:**
```go
agentService.SetReasoningMode(ReasoningVerbose)
```

---

### 2. Auto Summarization

#### Initial Summary

**Trigger Conditions:**
- Turn count ≥ threshold (mode-dependent)
- No existing summary
- Enough messages to summarize

**Thresholds:**
- ReasoningDisabled: 20 turns
- ReasoningEnabled: 10 turns
- ReasoningVerbose: No auto-summary (too short)

**Process:**
1. Get old messages (exclude recent N)
2. Generate summary using summary model
3. Save to `topic.history_summary`
4. Update metadata

**Example Metadata:**
```json
{
  "summarized_at": 1703001234567,
  "message_count": 40,
  "reasoning_mode": "disabled"
}
```

---

#### Incremental Summary

**Trigger Conditions:**
- Existing summary present
- New messages ≥ threshold since last summary

**Thresholds:**
- ReasoningDisabled: 10 new turns
- ReasoningEnabled: 5 new turns
- ReasoningVerbose: No incremental (too short)

**Process:**
1. Load existing summary
2. Get new messages since last summary
3. Merge old summary + new messages
4. Update `topic.history_summary`
5. Increment version in metadata

**Example Metadata:**
```json
{
  "summary_version": 3,
  "last_summarized_at": 1703001234567,
  "summarized_message_count": 80,
  "initial_summary_at": 1703001000000,
  "reasoning_mode": "disabled"
}
```

---

### 3. Streaming Support

#### Token-by-Token Streaming

**Flow:**
1. Frontend sends `ChatRequest` with `Stream: true`
2. Backend uses `generateWithTokenStreaming()`
3. Events emitted via Wails:
   - `chat:stream` (type: "start")
   - `chat:stream` (type: "chunk") - throttled to 1/sec
   - `chat:stream` (type: "complete")

**Event Format:**
```typescript
{
  type: "start" | "chunk" | "complete",
  session_id: string,
  message_id: string,
  content?: string,        // For chunk
  full_content?: string,   // For chunk (accumulated)
}
```

**Frontend Example:**
```typescript
// Listen for streaming events
app.Event.On("chat:stream", (event) => {
  if (event.type === "chunk") {
    updateMessage(event.message_id, event.content);
  } else if (event.type === "complete") {
    finalizeMessage(event.message_id, event.full_content);
  }
});

// Send chat request with streaming
await ChatService.Chat({
  session_id: "session-123",
  message: "Tell me a story",
  stream: true,
});
```

---

### 4. Topic & Thread Management

#### Auto Topic Creation

**Trigger:**
- First message in session
- No topicID provided

**Process:**
1. Create topic with placeholder title ("New Conversation")
2. Return topicID immediately
3. Background: Generate LLM title after first response
4. Update topic with generated title
5. Emit `chat:topic:updated` event

**Example:**
```go
// First message - auto-create topic
response, _ := agentService.Chat(ctx, ChatRequest{
    SessionID: "session-123",
    UserID:    "user-456",
    Message:   "What is AI?",
})
// response.TopicID = "topic-789" (auto-created)
```

---

#### Thread Branching

**Use Case:**
- User wants to explore alternative conversation path
- Create new thread from specific message

**Process:**
1. Provide `ThreadID` in request
2. Load thread messages from DB
3. Continue conversation in thread context

**Example:**
```go
// Branch from message-456
response, _ := agentService.Chat(ctx, ChatRequest{
    SessionID: "session-123",
    UserID:    "user-456",
    TopicID:   "topic-789",
    ThreadID:  "thread-new",
    ParentID:  "message-456",
    Message:   "What if we try a different approach?",
})
```

---

### 5. RAG (Knowledge Base Integration)

#### KB Search Tool

**Automatic Tool Creation:**
- If `KnowledgeBaseID` provided in request
- Creates `kbSearchTool` automatically
- Added to agent's tool list

**Tool Info:**
```json
{
  "name": "search_knowledge_base_name",
  "description": "Search the [KB Name] knowledge base",
  "parameters": {
    "query": "string - search query",
    "top_k": "integer - number of results (default: 5)"
  }
}
```

**Usage:**
```go
response, _ := agentService.Chat(ctx, ChatRequest{
    SessionID:       "session-123",
    UserID:          "user-456",
    Message:         "What does the documentation say about X?",
    KnowledgeBaseID: "kb-789",
})
// Agent automatically uses KB search tool
```

---

### 6. Context Engine Integration

#### Message Processing

**Flow:**
1. Messages collected in session
2. Passed to `ContextEngineBridge`
3. Context engine applies transformations:
   - Inject history summary
   - Add relevant context
   - Compress old messages
4. Processed messages sent to agent

**Example:**
```go
// Context engine automatically processes messages
messagesToAgent := session.Messages
if s.contextBridge != nil {
    processedMessages, _ := s.contextBridge.ProcessMessagesForAgent(
        ctx,
        session.Messages,
    )
    messagesToAgent = processedMessages
}
```

---

### 7. Tools Engine Integration

#### Tool Bridge

**Supported Tools:**
- Web search
- Calculator
- File operations
- Custom tools

**Usage:**
```go
response, _ := agentService.Chat(ctx, ChatRequest{
    SessionID: "session-123",
    UserID:    "user-456",
    Message:   "Search the web for latest AI news",
    Tools:     []string{"web_search"},
})
```

---

## Konfigurasi

### Environment Variables

```bash
# Model paths (optional - auto-detected if not set)
VERIDIUM_TITLE_MODEL="/path/to/llama-3.2-1b-instruct.gguf"
VERIDIUM_SUMMARY_MODEL="/path/to/llama-3.2-3b-instruct.gguf"

# Reasoning mode (optional - default: disabled)
VERIDIUM_REASONING_MODE="disabled" # disabled | enabled | verbose

# Database path
VERIDIUM_DB_PATH="./data/veridium.db"

# Model directory
VERIDIUM_MODELS_DIR="~/.llama-cpp/models"
```

---

### Reasoning Config

```go
type ReasoningConfig struct {
    Mode                 ReasoningMode
    SummaryThreshold     int  // Turns before initial summary
    IncrementalThreshold int  // Turns before incremental update
    KeepRecentMessages   int  // Messages to keep (not summarize)
}
```

**Defaults:**
```go
// ReasoningDisabled
ReasoningConfig{
    Mode:                 ReasoningDisabled,
    SummaryThreshold:     20,  // Summarize after 20 turns
    IncrementalThreshold: 10,  // Update every 10 turns
    KeepRecentMessages:   20,  // Keep last 20 messages
}

// ReasoningEnabled
ReasoningConfig{
    Mode:                 ReasoningEnabled,
    SummaryThreshold:     10,  // Summarize after 10 turns
    IncrementalThreshold: 5,   // Update every 5 turns
    KeepRecentMessages:   12,  // Keep last 12 messages
}

// ReasoningVerbose
ReasoningConfig{
    Mode:                 ReasoningVerbose,
    SummaryThreshold:     0,   // No auto-summary
    IncrementalThreshold: 0,   // No incremental
    KeepRecentMessages:   6,   // Keep last 6 messages
}
```

---

### Model Selection

#### Auto-Detection Logic

**Title Model Selection:**
```go
Priority:
1. Llama 3.2 1B (500MB-1GB) - BEST
2. Llama 3.2 3B (1GB-2GB) - Good
3. Other Llama models - OK
4. Mistral 7B - Acceptable
5. AVOID: Qwen (generates <think> tags)
6. AVOID: DeepSeek (reasoning model)
```

**Summary Model Selection:**
```go
Priority:
1. Llama 3.2 1B (500MB-1GB) - BEST
2. Llama 3.2 3B (1GB-2GB) - Good
3. Other Llama models - OK
4. Mistral 7B - Acceptable
5. AVOID: Qwen (generates <think> tags)
6. AVOID: DeepSeek (reasoning model)
```

**Chat Model Selection (by Reasoning Mode):**
```go
ReasoningDisabled:
  - Llama 3.2 1B/3B (preferred)
  - Mistral 7B
  - Gemma 2B/7B

ReasoningEnabled:
  - Qwen2.5-3B-Instruct (recommended)
  - DeepSeek-R1-Distill-Qwen-1.5B

ReasoningVerbose:
  - Qwen2.5-3B-Instruct (recommended)
  - DeepSeek-R1-Distill-Qwen-1.5B
```

---

## Best Practices

### 1. Session Management

#### ✅ DO:
```go
// Reuse sessions for same conversation
response1, _ := agentService.Chat(ctx, ChatRequest{
    SessionID: "session-123",
    Message:   "First question",
})

response2, _ := agentService.Chat(ctx, ChatRequest{
    SessionID: "session-123", // Same session
    Message:   "Follow-up question",
})
```

#### ❌ DON'T:
```go
// Don't create new session for each message
response1, _ := agentService.Chat(ctx, ChatRequest{
    SessionID: uuid.New().String(), // New session
    Message:   "First question",
})

response2, _ := agentService.Chat(ctx, ChatRequest{
    SessionID: uuid.New().String(), // Another new session (loses context!)
    Message:   "Follow-up question",
})
```

---

### 2. Reasoning Mode Selection

#### ✅ DO:
```go
// Choose mode based on use case
if isComplexTask {
    agentService.SetReasoningMode(ReasoningEnabled)
} else {
    agentService.SetReasoningMode(ReasoningDisabled)
}

// Validate model after mode change
if err := agentService.ValidateModelForReasoningMode(); err != nil {
    agentService.SwitchToRecommendedModel()
}
```

#### ❌ DON'T:
```go
// Don't use verbose mode for long conversations
agentService.SetReasoningMode(ReasoningVerbose)
// This will exhaust context in 3-5 turns!

// Don't ignore validation errors
agentService.SetReasoningMode(ReasoningEnabled)
// Model mismatch may cause poor performance
```

---

### 3. Streaming

#### ✅ DO:
```go
// Use streaming for better UX
response, _ := agentService.Chat(ctx, ChatRequest{
    SessionID: "session-123",
    Message:   "Tell me a long story",
    Stream:    true, // Enable streaming
})

// Listen for events in frontend
app.Event.On("chat:stream", handleStreamEvent)
```

#### ❌ DON'T:
```go
// Don't use streaming for short responses
response, _ := agentService.Chat(ctx, ChatRequest{
    SessionID: "session-123",
    Message:   "Hi", // Short response
    Stream:    true, // Unnecessary overhead
})
```

---

### 4. Knowledge Base Usage

#### ✅ DO:
```go
// Provide KB ID when relevant
response, _ := agentService.Chat(ctx, ChatRequest{
    SessionID:       "session-123",
    Message:         "What does the manual say about X?",
    KnowledgeBaseID: "kb-manual", // Relevant KB
})
```

#### ❌ DON'T:
```go
// Don't use KB for general knowledge
response, _ := agentService.Chat(ctx, ChatRequest{
    SessionID:       "session-123",
    Message:         "What is the capital of France?",
    KnowledgeBaseID: "kb-manual", // Irrelevant KB
})
```

---

### 5. Error Handling

#### ✅ DO:
```go
response, err := agentService.Chat(ctx, ChatRequest{
    SessionID: "session-123",
    Message:   "Hello",
})
if err != nil {
    log.Printf("Chat error: %v", err)
    // Handle gracefully
    return ErrorResponse{Message: "Sorry, something went wrong"}
}
```

#### ❌ DON'T:
```go
response, _ := agentService.Chat(ctx, ChatRequest{
    SessionID: "session-123",
    Message:   "Hello",
})
// Ignoring errors can lead to nil pointer panics
fmt.Println(response.Message) // May panic if error occurred
```

---

### 6. Context Management

#### ✅ DO:
```go
// Let auto-summarization handle context
// No manual intervention needed
response, _ := agentService.Chat(ctx, ChatRequest{
    SessionID: "session-123",
    Message:   "Continue conversation",
})
// Summary automatically created when threshold reached
```

#### ❌ DON'T:
```go
// Don't manually truncate messages
session.Messages = session.Messages[len(session.Messages)-10:]
// This breaks context and summary system
```

---

## Troubleshooting

### Problem 1: Model Generates `<think>` Tags

**Symptoms:**
```
⚠️  WARNING: Summary contains <think> tags (model: main), stripping...
```

**Cause:**
- Using reasoning model (Qwen, DeepSeek) for utility tasks
- Title/summary model not properly detected

**Solution:**
```go
// 1. Check current models
log.Printf("Title model: %s", agentService.titleModelPath)
log.Printf("Summary model: %s", agentService.summaryModelPath)

// 2. Manually set non-reasoning models
agentService.SetTitleModel("/path/to/llama-3.2-1b-instruct.gguf")

// 3. Or download recommended models
// Download Llama 3.2 1B/3B from Hugging Face
```

---

### Problem 2: Context Exhausted Too Quickly

**Symptoms:**
```
🚨 Context usage > 80%% Conversation may become unstable.
```

**Cause:**
- Using ReasoningVerbose mode (high token consumption)
- Summary threshold too high
- Model context size too small

**Solution:**
```go
// 1. Switch to more efficient mode
agentService.SetReasoningMode(ReasoningDisabled)

// 2. Lower summary threshold
config := agentService.GetReasoningConfig()
config.SummaryThreshold = 10 // Summarize earlier
agentService.reasoningConfig = config

// 3. Use model with larger context
// Load model with 32K context instead of 16K
```

---

### Problem 3: Slow Response Time

**Symptoms:**
- Response takes > 30 seconds
- High CPU usage

**Cause:**
- Using large model for utility tasks
- ReasoningVerbose mode
- No GPU acceleration

**Solution:**
```go
// 1. Check current mode
mode := agentService.GetReasoningMode()
if mode == ReasoningVerbose {
    agentService.SetReasoningMode(ReasoningDisabled)
}

// 2. Use smaller models for utility tasks
agentService.SetTitleModel("/path/to/llama-3.2-1b.gguf")

// 3. Enable GPU acceleration
// Set CUDA_VISIBLE_DEVICES or use Metal (macOS)

// 4. Use quantized models (Q4_K_M instead of Q8_0)
```

---

### Problem 4: Session Not Found

**Symptoms:**
```
Error: session not found: session-123
```

**Cause:**
- Session cleared from cache
- Database connection issue
- Session ID mismatch

**Solution:**
```go
// 1. Check if session exists in DB
dbSession, err := db.Queries().GetSession(ctx, db.GetSessionParams{
    ID:     sessionID,
    UserID: userID,
})

// 2. Session will auto-restore from DB on next request
// Just retry the request

// 3. If still fails, create new session
response, _ := agentService.Chat(ctx, ChatRequest{
    SessionID: uuid.New().String(), // New session
    UserID:    userID,
    Message:   message,
})
```

---

### Problem 5: Summary Not Generated

**Symptoms:**
- No summary after many turns
- `history_summary` field is NULL

**Cause:**
- Turn count below threshold
- ReasoningVerbose mode (no auto-summary)
- Background goroutine failed

**Solution:**
```go
// 1. Check current turn count
turnCount := len(session.Messages) / 2
threshold := agentService.reasoningConfig.SummaryThreshold
log.Printf("Turn count: %d, Threshold: %d", turnCount, threshold)

// 2. Check reasoning mode
mode := agentService.GetReasoningMode()
if mode == ReasoningVerbose {
    // No auto-summary in verbose mode
    agentService.SetReasoningMode(ReasoningEnabled)
}

// 3. Manually trigger summary (if needed)
// Note: This is internal method, not exposed
// Summary will be generated on next message
```

---

### Problem 6: Race Condition on Session Creation

**Symptoms:**
```
⚠️  Race condition detected, retrying load from DB...
```

**Cause:**
- Multiple requests with same sessionID arrive simultaneously
- UNIQUE constraint on session ID

**Solution:**
- **No action needed** - automatically handled by retry logic
- Service will retry loading from DB with exponential backoff
- If you see this frequently, consider:
  ```go
  // Add small delay before parallel requests
  time.Sleep(100 * time.Millisecond)
  ```

---

### Problem 7: Memory Leak

**Symptoms:**
- Memory usage grows over time
- Many cached sessions

**Cause:**
- Sessions never cleared from cache
- No session cleanup mechanism

**Solution:**
```go
// 1. Implement periodic cleanup
go func() {
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()
    
    for range ticker.C {
        // Clear inactive sessions
        agentService.sessionsMutex.Lock()
        for sessionID, session := range agentService.sessions {
            if time.Now().Unix() - session.UpdatedAt > 3600 {
                delete(agentService.sessions, sessionID)
            }
        }
        agentService.sessionsMutex.Unlock()
    }
}()

// 2. Or manually clear after conversation ends
agentService.ClearSession(sessionID)
```

---

## Performance Optimization

### 1. Model Selection

**Optimal Configuration:**
```
Chat Model (ReasoningDisabled):
  - Llama 3.2 3B Q4_K_M (1.9GB)
  - Speed: 2-5 tokens/sec
  - Quality: Good for general chat

Title Model:
  - Llama 3.2 1B Q4_K_M (0.9GB)
  - Speed: 5-10 tokens/sec
  - Quality: Sufficient for titles

Summary Model:
  - Llama 3.2 1B Q4_K_M (0.9GB)
  - Speed: 5-10 tokens/sec
  - Quality: Good for summarization
```

---

### 2. Context Management

**Thresholds by Mode:**
```
ReasoningDisabled (Long Conversations):
  - Summary threshold: 20 turns
  - Incremental: Every 10 turns
  - Keep recent: 20 messages
  - Expected: 40-50 turns total

ReasoningEnabled (Medium Conversations):
  - Summary threshold: 10 turns
  - Incremental: Every 5 turns
  - Keep recent: 12 messages
  - Expected: 15-20 turns total

ReasoningVerbose (Short Conversations):
  - Summary threshold: 0 (disabled)
  - Incremental: 0 (disabled)
  - Keep recent: 6 messages
  - Expected: 3-5 turns total
```

---

### 3. Streaming Optimization

**Throttling:**
```go
// Chunk events throttled to 1/second
const emitInterval = 1 * time.Second

// Prevents overwhelming frontend with events
// Reduces CPU usage on UI rendering
```

**Best Practices:**
- Use streaming for responses > 100 tokens
- Disable streaming for short responses
- Implement debouncing on frontend

---

### 4. Database Optimization

**Indexes:**
```sql
-- Session lookup
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_updated_at ON sessions(updated_at);

-- Message lookup
CREATE INDEX idx_messages_session_id ON messages(session_id);
CREATE INDEX idx_messages_topic_id ON messages(topic_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);

-- Topic lookup
CREATE INDEX idx_topics_user_id ON topics(user_id);
CREATE INDEX idx_topics_session_id ON topics(session_id);
```

---

### 5. Caching Strategy

**In-Memory Cache:**
- Store active sessions (last 1 hour)
- LRU eviction policy
- Max 100 sessions per user

**Database:**
- Full persistence
- No expiration
- Indexed for fast lookup

---

## Security Considerations

### 1. User Isolation

```go
// Always validate userID
response, err := agentService.Chat(ctx, ChatRequest{
    SessionID: sessionID,
    UserID:    userID, // MUST match authenticated user
    Message:   message,
})

// Database queries always filter by userID
dbSession, err := db.Queries().GetSession(ctx, db.GetSessionParams{
    ID:     sessionID,
    UserID: userID, // Prevents cross-user access
})
```

---

### 2. Input Validation

```go
// Validate request parameters
if req.SessionID == "" {
    return nil, fmt.Errorf("session_id required")
}
if req.UserID == "" {
    return nil, fmt.Errorf("user_id required")
}
if req.Message == "" {
    return nil, fmt.Errorf("message required")
}

// Sanitize message content
req.Message = strings.TrimSpace(req.Message)
if len(req.Message) > 10000 {
    return nil, fmt.Errorf("message too long")
}
```

---

### 3. Rate Limiting

```go
// Implement rate limiting per user
type RateLimiter struct {
    requests map[string][]time.Time
    mutex    sync.RWMutex
}

func (rl *RateLimiter) Allow(userID string) bool {
    rl.mutex.Lock()
    defer rl.mutex.Unlock()
    
    now := time.Now()
    requests := rl.requests[userID]
    
    // Keep only requests in last minute
    var recent []time.Time
    for _, t := range requests {
        if now.Sub(t) < time.Minute {
            recent = append(recent, t)
        }
    }
    
    // Max 60 requests per minute
    if len(recent) >= 60 {
        return false
    }
    
    recent = append(recent, now)
    rl.requests[userID] = recent
    return true
}
```

---

## Testing

### Unit Tests

```go
func TestChat(t *testing.T) {
    // Setup
    agentService := setupTestService(t)
    
    // Test
    response, err := agentService.Chat(context.Background(), ChatRequest{
        SessionID: "test-session",
        UserID:    "test-user",
        Message:   "Hello",
    })
    
    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, response.Message)
    assert.NotEmpty(t, response.MessageID)
}

func TestAutoSummarization(t *testing.T) {
    agentService := setupTestService(t)
    agentService.SetReasoningMode(ReasoningDisabled)
    
    // Send 20 turns (40 messages)
    for i := 0; i < 20; i++ {
        _, err := agentService.Chat(context.Background(), ChatRequest{
            SessionID: "test-session",
            UserID:    "test-user",
            Message:   fmt.Sprintf("Message %d", i),
        })
        assert.NoError(t, err)
    }
    
    // Wait for background summary
    time.Sleep(5 * time.Second)
    
    // Check summary exists
    topic, err := db.Queries().GetTopic(context.Background(), db.GetTopicParams{
        ID:     topicID,
        UserID: "test-user",
    })
    assert.NoError(t, err)
    assert.True(t, topic.HistorySummary.Valid)
    assert.NotEmpty(t, topic.HistorySummary.String)
}
```

---

### Integration Tests

```go
func TestEndToEndChat(t *testing.T) {
    // Setup full stack
    app := setupTestApp(t)
    dbService := setupTestDB(t)
    llamaService := setupTestLlama(t)
    kbService := setupTestKB(t)
    
    agentService := NewAgentChatService(
        app, dbService, llamaService, kbService,
        nil, nil, nil,
    )
    
    // Test conversation flow
    sessionID := uuid.New().String()
    userID := "test-user"
    
    // Message 1
    resp1, err := agentService.Chat(context.Background(), ChatRequest{
        SessionID: sessionID,
        UserID:    userID,
        Message:   "What is AI?",
    })
    assert.NoError(t, err)
    assert.NotEmpty(t, resp1.TopicID)
    
    // Message 2 (same topic)
    resp2, err := agentService.Chat(context.Background(), ChatRequest{
        SessionID: sessionID,
        UserID:    userID,
        TopicID:   resp1.TopicID,
        Message:   "Tell me more",
    })
    assert.NoError(t, err)
    assert.Equal(t, resp1.TopicID, resp2.TopicID)
}
```

---

## Migration Guide

### From LibraryChatService to AgentChatService

**Old Code:**
```go
libraryChatService := llama.NewLibraryChatService(libService, app)
response, err := libraryChatService.Chat(ctx, llama.ChatRequest{
    Messages: []llama.ChatMessage{
        {Role: "user", Content: "Hello"},
    },
})
```

**New Code:**
```go
agentChatService := NewAgentChatService(
    app, dbService, libService, kbService,
    toolsBridge, contextBridge, threadService,
)
response, err := agentChatService.Chat(ctx, ChatRequest{
    SessionID: sessionID,
    UserID:    userID,
    Message:   "Hello",
})
```

**Benefits:**
- ✅ Automatic session management
- ✅ Database persistence
- ✅ Context compression
- ✅ Tool integration
- ✅ RAG support
- ✅ Thread management

---

## Changelog

### Version 1.0.0 (Current)
- ✅ Initial implementation
- ✅ Three reasoning modes
- ✅ Auto summarization
- ✅ Streaming support
- ✅ RAG integration
- ✅ Thread management
- ✅ Utility model auto-detection

### Planned Features (v1.1.0)
- 🔄 Multi-user support
- 🔄 Session sharing
- 🔄 Export conversations
- 🔄 Custom tool creation
- 🔄 Advanced RAG strategies
- 🔄 Voice input/output

---

## References

### External Documentation
- [Eino Framework](https://github.com/cloudwego/eino)
- [llama.cpp](https://github.com/ggerganov/llama.cpp)
- [Wails v3](https://wails.io/)

### Internal Documentation
- [EINO_ARCHITECTURE.md](./EINO_ARCHITECTURE.md)
- [HISTORY_SUMMARY.md](./HISTORY_SUMMARY.md)
- [LLM_OPTIMIZATION_GUIDE.md](./LLM_OPTIMIZATION_GUIDE.md)

---

## Support

### Issues & Questions
- GitHub Issues: [veridium/issues](https://github.com/kawai-network/veridium/issues)
- Discussions: [veridium/discussions](https://github.com/kawai-network/veridium/discussions)

### Contributing
See [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines.

---

**Last Updated:** 2025-11-26  
**Version:** 1.0.0  
**Maintainer:** Veridium Team

