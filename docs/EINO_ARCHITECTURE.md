# Arsitektur Eino Integration - Veridium

> **Comprehensive Architecture Documentation**  
> CloudWeGo Eino integration with Veridium - Complete implementation guide, status tracking, and database integration strategy.
>
> **Version**: 1.4  
> **Status**: Phase 3 Complete (95%) - Thread Management Added!  
> **Last Updated**: November 2025

## 📑 Table of Contents

1. [Overview](#overview)
2. [Quick Status Overview](#-quick-status-overview)
3. [Tujuan Integrasi](#-tujuan-integrasi)
4. [Arsitektur Diagram](#-arsitektur-diagram)
5. [Komponen Detail & Status](#️-komponen-detail--status-implementasi)
   - [Llama Eino Model Adapter](#1-llama-eino-model-adapter)
   - [Knowledge Base Service](#2-knowledge-base-service)
   - [RAG Workflow](#3-rag-workflow)
   - [RAG Agent](#4-rag-agent)
   - [Agent Chat Service](#5-agent-chat-service-)
   - [Eino Adapters (Chromem)](#6-eino-adapters-chromem)
   - [Database Integration](#7-database-integration)
   - [Main Application Wiring](#8-main-application-wiring)
   - [Session Management Architecture](#9-session-management-architecture-)
   - [Thread Management Service](#10-thread-management-service-)
6. [Komponen Belum Diimplementasi](#-komponen-belum-diimplementasikan)
7. [Data Flow - Complete Path](#-data-flow---complete-path)
8. [Test Coverage](#-test-coverage)
9. [How to Use](#-how-to-use)
10. [Next Steps](#-next-steps)
11. [Status Implementasi Detail](#-status-implementasi-detail)
12. [Session Storage Strategy](#️-penting-session-storage-strategy)
13. [Quick Start](#-quick-start)
14. [Roadmap](#-roadmap)
15. [Progress Breakdown](#-progress-breakdown)
16. [Implementation Checklist](#-implementation-checklist)
17. [Key Takeaways](#-key-takeaways)
18. [Migration Strategy](#-migration-strategy-frontend--backend)
19. [References](#-references)

---

## Overview

Dokumen ini menjelaskan arsitektur lengkap integrasi CloudWeGo Eino dengan Veridium, termasuk status implementasi setiap komponen.

### 📊 Quick Status Overview

```
Phase 1: Core Backend           ✅ 100% COMPLETE
├─ Llama Eino Model Adapter     ✅ DONE
├─ Knowledge Base Service       ✅ DONE
├─ RAG Workflow                 ✅ DONE
├─ RAG Agent (ADK)              ✅ DONE
├─ Agent Chat Service           ✅ DONE
├─ Database Integration         ✅ DONE (KB)
└─ E2E Test                     ✅ DONE

Phase 2: DB Persistence         ✅ 100% COMPLETE
├─ Session persistence          ✅ DONE (hybrid: DB + cache)
├─ Message history              ✅ DONE (SQLite storage)
├─ Agent reconstruction         ✅ DONE (auto-load on restart)
├─ DB conversion helpers        ✅ DONE (Eino ↔ DB)
├─ Auto-save messages           ✅ DONE (user + assistant)
└─ Session timestamp updates    ✅ DONE

Phase 3: Integration            ✅ 100% COMPLETE
├─ Context Engine Bridge        ✅ DONE (message processing)
├─ Tools Engine Bridge          ✅ DONE (tool integration)
├─ Agent integration            ✅ DONE (bridges in AgentChatService)
└─ Main.go wiring               ✅ DONE

Phase 4: Production             ⏳ NEXT PRIORITY
├─ Frontend UI                  ⏳ TODO (highest priority)
├─ TypeScript types             ⏳ TODO
├─ Streaming                    ⏳ TODO (placeholder exists)
├─ Advanced RAG                 ⏳ TODO (optional)
└─ Monitoring                   ⏳ TODO

Overall Progress: ██████████████████████░░ 95%
```

### 🎉 Phase 3 Complete + Thread Management!

**Tools, Context Engine, Auto Topics & Thread Branching terintegrasi!**
```
✅ Tools Engine Bridge:    Existing tools accessible by agent
✅ Context Engine Bridge:  Message preprocessing before agent
✅ Agent integration:      Bridges wired into AgentChatService
✅ Dynamic tool loading:   Tools specified per request
✅ Context processing:     History, templates, placeholders
✅ Auto Topic Generation:  LLM generates title after first response
✅ Thread Management:      Conversation branching & multi-threading
🎯 Next priority:          Frontend UI untuk full system usage
```

---

## 🎯 Tujuan Integrasi

Mengintegrasikan CloudWeGo Eino sebagai framework orchestration untuk:
1. **Knowledge Base & RAG** - Sistem manajemen pengetahuan dengan Retrieval-Augmented Generation
2. **Agent Development** - AI agents dengan tool calling dan reasoning capabilities
3. **Local-First** - Semua berjalan lokal tanpa dependensi external API
4. **Hybrid Storage** - SQLite untuk metadata, Chromem untuk vectors
5. **Existing Integration** - Memanfaatkan komponen yang sudah ada (context engine, tools engine)

---

## 📐 Arsitektur Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                        FRONTEND (React)                          │
│  - TypeScript types                                             │
│  - Wails bindings                                               │
│  - Session ID management (localStorage/state)                   │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                    WAILS APPLICATION                             │
│  Services exposed via application.NewService()                  │
└────────────────────────┬────────────────────────────────────────┘
                         │
        ┌────────────────┼────────────────┐
        │                │                │
        ▼                ▼                ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│   KB Service │ │  Agent Chat  │ │ RAG Workflow │
│              │ │   Service    │ │              │
│              │ │  ┌─────────┐ │ │              │
│              │ │  │ SESSION │ │ │              │
│              │ │  │ MANAGER │ │ │              │
│              │ │  │ (Memory)│ │ │              │
│              │ │  └─────────┘ │ │              │
└──────┬───────┘ └──────┬───────┘ └──────┬───────┘
       │                │                │
       │                └────────┬───────┘
       │                         │
       ▼                         ▼
┌──────────────────────────────────────────┐
│         EINO COMPONENTS                  │
│  ┌────────────────────────────────────┐ │
│  │  Llama Eino Model                  │ │
│  │  - ToolCallingChatModel            │ │
│  │  - Generate() / Stream()           │ │
│  └────────────────────────────────────┘ │
│                                          │
│  ┌────────────────────────────────────┐ │
│  │  ADK Agent (per session)           │ │
│  │  - ChatModelAgent                  │ │
│  │  - Tool binding                    │ │
│  │  - Multi-turn conversations        │ │
│  │  - Message history                 │ │
│  └────────────────────────────────────┘ │
│                                          │
│  ┌────────────────────────────────────┐ │
│  │  Retriever / Indexer               │ │
│  │  - ChromemAdapter.Retriever        │ │
│  │  - ChromemAdapter.Indexer          │ │
│  │  - FileManager                     │ │
│  └────────────────────────────────────┘ │
└──────────────┬───────────────────────────┘
               │
        ┌──────┴──────┐
        │             │
        ▼             ▼
┌──────────────┐ ┌──────────────┐
│   Llama.cpp  │ │   Chromem    │
│  (Library)   │ │  (Vector DB) │
└──────────────┘ └──────────────┘
        │             │
        └──────┬──────┘
               │
               ▼
┌─────────────────────────────┐
│     STORAGE LAYER           │
│  ┌──────────┐ ┌──────────┐ │
│  │  SQLite  │ │  Disk    │ │
│  │ Metadata │ │ Vectors  │ │
│  └──────────┘ └──────────┘ │
└─────────────────────────────┘
```

---

## 🏗️ Komponen Detail & Status Implementasi

### 1. **Llama Eino Model Adapter** 
**File**: `internal/llama/eino_adapter.go`  
**Status**: ✅ **IMPLEMENTED**

#### Fungsi:
- Adapter untuk Llama.cpp agar kompatibel dengan Eino's `ToolCallingChatModel` interface
- Bridge antara local LLM dan Eino framework

#### API:
```go
type LlamaEinoModel struct {
    libService *LibraryService
    tools      []*schema.ToolInfo
}

// BaseChatModel interface
func (m *LlamaEinoModel) Generate(ctx, input, opts) (*Message, error)
func (m *LlamaEinoModel) Stream(ctx, input, opts) (*StreamReader[*Message], error)

// ToolCallingChatModel interface
func (m *LlamaEinoModel) WithTools(tools) (ToolCallingChatModel, error)
```

#### Features:
- ✅ Synchronous generation
- ✅ Streaming generation
- ✅ Tool calling support
- ✅ Message conversion (Eino ↔ Llama)
- ✅ Sampler configuration (temperature, top-p, top-k, max tokens)
- ✅ Context management

---

### 2. **Knowledge Base Service**
**File**: `internal/services/knowledge_base.go`  
**Status**: ✅ **IMPLEMENTED**

#### Fungsi:
- Manajemen knowledge bases
- Upload & parsing dokumen
- Vector indexing & retrieval
- Hybrid storage (SQLite + Chromem)

#### API:
```go
type KnowledgeBaseService struct {
    db            *database.Service
    chromemDB     *chromem.DB
    collections   map[string]*chromem.Collection
    indexers      map[string]*chromemAdapter.Indexer
    retrievers    map[string]*chromemAdapter.Retriever
    fileManagers  map[string]*chromemAdapter.FileManager
}

// CRUD Operations
func (s *KnowledgeBaseService) CreateKnowledgeBase(ctx, name, desc, userID) (string, error)
func (s *KnowledgeBaseService) GetKnowledgeBase(ctx, kbID, userID) (KnowledgeBasis, error)
func (s *KnowledgeBaseService) ListKnowledgeBases(ctx, userID) ([]KnowledgeBasis, error)
func (s *KnowledgeBaseService) UpdateKnowledgeBase(ctx, kbID, name, desc, userID) error
func (s *KnowledgeBaseService) DeleteKnowledgeBase(ctx, kbID, userID) error

// File Operations
func (s *KnowledgeBaseService) AddFileToKnowledgeBase(ctx, kbID, filePath, metadata, userID) error
func (s *KnowledgeBaseService) RemoveFileFromKnowledgeBase(ctx, kbID, fileID, userID) error

// Query Operations
func (s *KnowledgeBaseService) QueryKnowledgeBase(ctx, kbID, query, topK, userID) ([]*Document, error)

// Eino Integration
func (s *KnowledgeBaseService) GetRetriever(ctx, kbID, userID) (*Retriever, error)
func (s *KnowledgeBaseService) GetIndexer(ctx, kbID, userID) (*Indexer, error)
func (s *KnowledgeBaseService) GetFileManager(ctx, kbID, userID) (*FileManager, error)
```

#### Features:
- ✅ CRUD knowledge bases
- ✅ Multi-user support
- ✅ File upload & automatic parsing (DOCX, PDF, XLSX, HTML, TXT, MD)
- ✅ Automatic chunking & embedding
- ✅ Vector storage dengan persistence
- ✅ On-demand loading
- ✅ Metadata storage di SQLite
- ✅ Vector storage di Chromem

#### Storage:
- **SQLite**: `knowledge_bases` table, `knowledge_base_files` junction
- **Chromem**: Collections per KB, persistent to disk
- **File System**: Document copies di `kb-assets/{kb_id}/`

---

### 3. **RAG Workflow**
**File**: `internal/services/rag_workflow.go`  
**Status**: ✅ **IMPLEMENTED**

#### Fungsi:
- Implementasi Retrieval-Augmented Generation
- Orchestration antara retrieval dan generation

#### API:
```go
type RAGWorkflow struct {
    kbService *KnowledgeBaseService
}

func (w *RAGWorkflow) BuildContext(ctx, req) (string, []*Document, error)
func (w *RAGWorkflow) FormatContextForLLM(docs) string
func (w *RAGWorkflow) ExecuteRAG(ctx, req) (*RAGResponse, error)
```

#### Features:
- ✅ Document retrieval dengan semantic search
- ✅ Context building untuk LLM
- ✅ Source tracking
- ✅ Configurable TopK

---

### 4. **RAG Agent**
**File**: `internal/services/rag_agent.go`  
**Status**: ✅ **IMPLEMENTED**

#### Fungsi:
- AI Agent dengan RAG capabilities
- Tool calling untuk KB search
- Built with Eino ADK

#### API:
```go
type RAGAgent struct {
    agent       adk.Agent
    kbService   *KnowledgeBaseService
    ragWorkflow *RAGWorkflow
}

func NewRAGAgent(ctx, config, kbService) (*RAGAgent, error)
func (a *RAGAgent) Run(ctx, userMessage) (string, error)
func (a *RAGAgent) Stream(ctx, userMessage) (*AsyncIterator[*AgentEvent], error)
```

#### Features:
- ✅ Eino ADK ChatModelAgent
- ✅ Tool calling support
- ✅ KB search tools (auto-generated per KB)
- ✅ Multi-turn conversations
- ✅ Configurable max iterations

---

### 5. **Agent Chat Service** 🆕
**File**: `internal/services/agent_chat_service.go`  
**Status**: ✅ **FULLY IMPLEMENTED** (with DB persistence)

#### Fungsi:
- **Service utama untuk chat dengan AI agent**
- Replaces/complements existing `LibraryChatService`
- Exposes agent capabilities ke frontend
- **Hybrid session storage**: In-memory cache + SQLite persistence

#### API:
```go
type AgentChatService struct {
    app         *application.App
    libService  *llama.LibraryService
    llamaModel  *llama.LlamaEinoModel
    kbService   *KnowledgeBaseService
    ragWorkflow *RAGWorkflow
    sessions    map[string]*AgentSession
}

// Main Chat API
func (s *AgentChatService) Chat(ctx, req) (*ChatResponse, error)
func (s *AgentChatService) ChatStream(ctx, req) (*ChatResponse, error) // TODO: streaming

// Session Management
func (s *AgentChatService) ClearSession(sessionID)
func (s *AgentChatService) GetSessionHistory(sessionID) ([]*Message, error)
```

#### Request/Response:
```go
type ChatRequest struct {
    SessionID       string         // Required for multi-turn
    UserID          string
    Message         string
    KnowledgeBaseID string         // Optional: KB untuk RAG
    Tools           []string       // Optional: additional tools
    Context         map[string]any // Optional: session context
    Temperature     float32
    MaxTokens       int
    Stream          bool
}

type ChatResponse struct {
    SessionID    string
    Message      string
    ToolCalls    []ToolCall     // Tool yang dipanggil agent
    Sources      []*Document    // Sources dari KB (jika ada)
    FinishReason string
    Usage        *TokenUsage
}
```

#### Features:
- ✅ **Session management** - Multi-turn conversations dengan state
- ✅ **Database persistence** - Sessions & messages saved to SQLite
- ✅ **Auto-load on restart** - Reconstruct agent with full history
- ✅ **Hybrid storage** - Fast in-memory cache + reliable DB
- ✅ **Auto-creates agents** - Lazy initialization per session
- ✅ **KB integration** - Auto-adds KB search tool jika KB specified
- ✅ **Tool execution** - Agent bisa call tools dan return results
- ✅ **Source tracking** - Return sources dari KB queries
- ✅ **Wails integration** - Exposed ke frontend via Wails bindings
- ✅ **Auto Topic Generation** - LLM generates conversation title after first response
- ⏳ **Streaming** - TODO: implement dengan Wails v3 events API

#### Session Management:
```go
type AgentSession struct {
    SessionID       string
    UserID          string
    Agent           adk.Agent       // Eino agent instance
    Messages        []*schema.Message // Conversation history
    KnowledgeBaseID string
    Tools           []tool.BaseTool
    Context         map[string]any
    CreatedAt       int64
    UpdatedAt       int64
}

// Stored in: AgentChatService.sessions (in-memory map)
// Key: SessionID (string, e.g., "session-123")
// Lifecycle: 
//   - Created: On first message with new SessionID
//   - Updated: On every message exchange
//   - Cleared: Via ClearSession() or app restart
```

---

### 9. **Session Management Architecture** 🔑
**Location**: `internal/services/agent_chat_service.go`  
**Status**: ✅ **IMPLEMENTED**

#### Posisi Session:
Session disimpan **in-memory** di `AgentChatService`:

```go
type AgentChatService struct {
    db            *database.Service              // Database service
    // ... other fields ...
    sessions      map[string]*AgentSession       // In-memory cache
    sessionsMutex sync.RWMutex                   // Thread-safe access
}
```

#### Lifecycle Session:

```
┌─────────────────────────────────────────────────────────────┐
│                    SESSION LIFECYCLE                         │
└─────────────────────────────────────────────────────────────┘

1. CREATION (Lazy)
   Frontend → ChatRequest{SessionID: "sess-123"} 
   ↓
   AgentChatService.getOrCreateSession()
   ├─ Check: sessions["sess-123"] exists?
   │  ├─ YES → Return existing session
   │  └─ NO  → Create new session
   │           ├─ Initialize ADK Agent
   │           ├─ Bind LlamaEinoModel
   │           ├─ Add KB search tool (if KB specified)
   │           └─ Store in sessions map
   ↓
   
2. USAGE (Every message)
   session.Messages.Append(userMessage)
   ↓
   agent.Run(session.Messages) // ADK uses full history
   ↓
   session.Messages.Append(assistantMessage)
   ↓
   session.UpdatedAt = now()
   
3. PERSISTENCE
   ✅ Sessions are HYBRID (Memory + DB)
   ├─ Saved to SQLite database
   ├─ Survives app restart
   └─ Auto-loads on next session access
   
4. CLEANUP
   Manual: AgentChatService.ClearSession(sessionID)
   ├─ Removes from sessions map
   └─ Releases agent resources
   
   Automatic: On app shutdown
   ├─ All sessions cleared
   └─ Memory freed
```

#### Session Scope:
- **Per User**: Each UserID can have multiple SessionIDs
- **Per KB**: Session can be bound to one KnowledgeBaseID
- **Per Agent**: Each session has its own ADK Agent instance
- **Multi-turn**: Session maintains full conversation history

#### Thread Safety:
```go
func (s *AgentChatService) getOrCreateSession(...) {
    s.sessionsMutex.Lock()
    defer s.sessionsMutex.Unlock()
    
    // Safe concurrent access
}
```

#### Session ID Management:

**Frontend Responsibility**:
```typescript
// Generate on client side
const sessionID = `session-${uuidv4()}`; 

// Or use conversation ID from UI
const sessionID = `conv-${conversationId}`;

// Store in localStorage or React state
localStorage.setItem('currentSessionID', sessionID);
```

**Backend** (read-only):
```go
// AgentChatService only reads SessionID
// Does NOT generate or persist SessionIDs
req.SessionID // Provided by frontend
```

#### Session vs Conversation:
```
Session (Backend):
  ├─ In-memory state
  ├─ Agent instance
  ├─ Message history
  └─ Bound to KB & tools

Conversation (Frontend):
  ├─ UI component
  ├─ Message display
  ├─ SessionID reference
  └─ Persisted in frontend state/localStorage
  
Relationship: 1 Conversation = 1 Session
```

---

### 6. **Eino Adapters** (Chromem)
**Path**: `pkg/eino-adapters/chromem/`  
**Status**: ✅ **ALREADY IMPLEMENTED** (by user)

#### Components:
- ✅ `Indexer` - ChromemAdapter.Indexer
- ✅ `Retriever` - ChromemAdapter.Retriever  
- ✅ `FileManager` - Document parsing & chunking

#### Supported Formats:
- DOCX, XLSX, PDF, HTML, TXT, MD

---

### 7. **Database Integration**
**Status**: ✅ **FULLY IMPLEMENTED** 

#### Tables Used (KB, Files, Sessions & Threads):
- ✅ `knowledge_bases` - KB metadata
- ✅ `knowledge_base_files` - KB-file junction
- ✅ `files` - File metadata
- ✅ `agents` - Agent metadata
- ✅ `agents_knowledge_bases` - Agent-KB junction
- ✅ **`sessions`** - Session metadata (Phase 2)
- ✅ **`messages`** - Full conversation history (Phase 2)
- ✅ **`topics`** - Conversation topics with auto-generated titles (Phase 3)
- ✅ **`threads`** - Conversation branches for multi-threading (Phase 3)

#### Tables Available for Future Use:
- ⏳ `agents_to_sessions` - Link agents to sessions (optional)
- ⏳ `files_to_sessions` - Attach files to sessions (optional)
- ⏳ `session_groups` - Organize sessions into groups

#### SQLC Generated Queries Available:

**Sessions** (`internal/database/queries/sessions.sql`):
```sql
-- Already generated by SQLC:
GetSession(id, user_id)
GetSessionBySlug(slug, user_id)
ListSessions(user_id, limit, offset)
CreateSession(id, slug, title, desc, user_id, ...)
UpdateSession(title, desc, avatar, ...)
DeleteSession(id, user_id)
SearchSessions(user_id, query, limit)
GetSessionWithGroup(id, user_id)
```

**Messages** (`internal/database/queries/messages.sql`):
```sql
-- Already generated by SQLC:
GetMessage(id)
ListMessages(user_id, role, limit, offset)
ListMessagesBySessionId(session_id, limit, offset)
ListMessagesByThread(user_id, thread_id)  -- NEW: Phase 3
CreateMessage(id, session_id, role, content, ...)
UpdateMessage(id, content, ...)
DeleteMessage(id)
```

**Topics** (`internal/database/queries/topics.sql`):
```sql
-- Already generated by SQLC:
CreateTopic(id, title, session_id, user_id, ...)
GetTopic(id, user_id)
ListTopics(user_id, limit, offset)
CountTopicsBySession(session_id, user_id)  -- NEW: Phase 3
UpdateTopic(id, title, ...)
DeleteTopic(id, user_id)
```

**Threads** (`internal/database/queries/threads.sql`):
```sql
-- Already generated by SQLC:
CreateThread(id, title, type, status, topic_id, source_message_id, ...)
ListThreadsByTopic(topic_id, user_id)
GetThread(id, user_id)
UpdateThread(id, title, status, last_active_at, ...)
DeleteThread(id, user_id)
```

**All available via**:
```go
import db "github.com/kawai-network/veridium/internal/database/generated"

// Example usage:
session, err := queries.GetSession(ctx, db.GetSessionParams{
    ID:     sessionID,
    UserID: userID,
})

messages, err := queries.ListMessagesBySessionId(ctx, db.ListMessagesBySessionIdParams{
    SessionID: sessionID,
    Limit:     100,
    Offset:    0,
})
```

#### Database Schema for Agent Sessions:

```sql
-- Sessions (conversation containers)
CREATE TABLE sessions (
  id TEXT PRIMARY KEY,
  slug TEXT NOT NULL,              -- URL-friendly identifier
  title TEXT,                      -- Session title
  description TEXT,
  avatar TEXT,                     -- Session avatar/icon
  background_color TEXT,
  type TEXT DEFAULT 'agent',       -- 'agent' | 'chat' | 'group'
  user_id TEXT NOT NULL,
  group_id TEXT,                   -- Optional grouping
  client_id TEXT,
  pinned INTEGER DEFAULT 0,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL,
  UNIQUE(slug, user_id)
);

-- Messages (conversation history)
CREATE TABLE messages (
  id TEXT PRIMARY KEY,
  role TEXT NOT NULL,              -- 'system' | 'user' | 'assistant' | 'tool'
  content TEXT,                    -- Message content
  name TEXT,                       -- For tool/function messages
  session_id TEXT NOT NULL,        -- Link to session
  parent_id TEXT,                  -- For threading
  user_id TEXT NOT NULL,
  topic_id TEXT,                   -- Optional topic/thread
  model TEXT,
  tool_call_id TEXT,               -- For tool responses
  tool_calls TEXT,                 -- JSON: tool calls made
  client_id TEXT,
  error TEXT,
  metadata TEXT,                   -- JSON: extra metadata
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);

-- Agent-Session Links
CREATE TABLE agents_to_sessions (
  agent_id TEXT NOT NULL,
  session_id TEXT NOT NULL,
  user_id TEXT NOT NULL,
  PRIMARY KEY (agent_id, session_id)
);

-- File-Session Links (for context)
CREATE TABLE files_to_sessions (
  file_id TEXT NOT NULL,
  session_id TEXT NOT NULL,
  user_id TEXT NOT NULL,
  PRIMARY KEY (file_id, session_id)
);

-- Topics (conversation titles, auto-generated)
CREATE TABLE topics (
  id TEXT PRIMARY KEY,
  title TEXT,
  favorite INTEGER DEFAULT 0,
  session_id TEXT,
  group_id TEXT,
  user_id TEXT NOT NULL,
  client_id TEXT,
  history_summary TEXT,
  metadata TEXT,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);

-- Threads (conversation branches for multi-threading)
CREATE TABLE threads (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  type TEXT NOT NULL CHECK (type IN ('continuation', 'standalone')),
  status TEXT DEFAULT 'active' CHECK (status IN ('active', 'deprecated', 'archived')),
  topic_id TEXT NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
  source_message_id TEXT NOT NULL,
  parent_thread_id TEXT REFERENCES threads(id) ON DELETE SET NULL,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  client_id TEXT,
  last_active_at INTEGER NOT NULL,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);
```

#### Integration Strategy:

**Phase 1** (Was: In-memory only):
```go
// In-memory only
AgentChatService.sessions = map[string]*AgentSession{}
```

**Phase 2** (Current - IMPLEMENTED ✅):
```go
// ✅ IMPLEMENTED - This is the actual code now running!
func (s *AgentChatService) getOrCreateSession(ctx, req) (*AgentSession, error) {
    s.sessionsMutex.Lock()
    defer s.sessionsMutex.Unlock()

    // 1. Check in-memory cache first (fast path)
    if session, exists := s.sessions[req.SessionID]; exists {
        log.Printf("♻️  Reusing cached session: %s", req.SessionID)
        return session, nil
    }

    // 2. Try to load from DB
    dbSession, dbMessages, err := s.loadSessionFromDB(ctx, req.SessionID, req.UserID)
    if err == nil {
        // Session exists in DB - reconstruct from history
        log.Printf("📂 Loading session from DB: %s (%d messages)", req.SessionID, len(dbMessages))
        
        // Convert DB messages to Eino messages
        einoMessages := make([]*schema.Message, 0, len(dbMessages))
        for _, dbMsg := range dbMessages {
            einoMsg, _ := convertDBMessageToEino(&dbMsg)
            einoMessages = append(einoMessages, einoMsg)
        }
        
        // Reconstruct agent with history
        session := &AgentSession{
            SessionID:   req.SessionID,
            UserID:      req.UserID,
            Messages:    einoMessages,
            KnowledgeBaseID: req.KnowledgeBaseID,
            Agent:       createAgentWithHistory(...),
            DBSession:   dbSession,
        }
        s.sessions[req.SessionID] = session
        return session, nil
    }
    
    // 3. Session doesn't exist - create new
    log.Printf("🆕 Creating new session: %s", req.SessionID)
    dbSession, _ = s.createSessionInDB(ctx, req.SessionID, req.UserID, req.KnowledgeBaseID)
    
    session := &AgentSession{
        SessionID: req.SessionID,
        Agent:     createNewAgent(...),
        DBSession: dbSession,
        ...
    }
    s.sessions[req.SessionID] = session
    return session, nil
}

func (s *AgentChatService) Chat(ctx, req) (*ChatResponse, error) {
    session, _ := s.getOrCreateSession(ctx, req)
    
    // Add & save user message
    userMsg := &schema.Message{Role: schema.User, Content: req.Message}
    session.Messages = append(session.Messages, userMsg)
    s.saveMessageToDB(ctx, userMsg, session.SessionID, session.UserID)
    
    // Run agent
    output, _ := session.Agent.Run(ctx, ...)
    
    // Add & save assistant message
    assistantMsg := &schema.Message{Role: schema.Assistant, Content: output.Message, ...}
    session.Messages = append(session.Messages, assistantMsg)
    s.saveMessageToDB(ctx, assistantMsg, session.SessionID, session.UserID)
    
    // Update session timestamp
    s.updateSessionTimestamp(ctx, session.SessionID, session.UserID)
    
    return &ChatResponse{...}, nil
}
```

**Helper Methods Added** (Phase 2):
```go
// DB ↔ Eino conversions
func convertDBMessageToEino(dbMsg *db.Message) (*schema.Message, error)
func convertEinoMessageToDB(einoMsg *schema.Message, sessionID, userID string) db.CreateMessageParams

// DB operations
func (s *AgentChatService) loadSessionFromDB(ctx, sessionID, userID string) (*db.Session, []db.Message, error)
func (s *AgentChatService) createSessionInDB(ctx, sessionID, userID, kbID string) (*db.Session, error)
func (s *AgentChatService) saveMessageToDB(ctx, msg *schema.Message, sessionID, userID string) error
func (s *AgentChatService) updateSessionTimestamp(ctx, sessionID, userID string) error
```

---

### 8. **Main Application Wiring**
**File**: `main.go`  
**Status**: ✅ **IMPLEMENTED**

#### Initialization Flow:
```go
1. LibraryService (Llama.cpp) ✅
2. KnowledgeBaseService ✅
   - ChromemDB initialization
   - Embedding function (Llama.cpp)
   - Asset directories
3. LibraryChatService (existing) ✅
4. ToolsEngineService & Bridge ✅
5. ContextEngine & Bridge ✅
6. AgentChatService (with bridges) ✅
7. ThreadManagementService ✅
8. Register all to Wails app ✅
```

#### Code Snippet:
```go
// Phase 3: Initialize integration bridges
toolsEngineService := NewToolsEngineService()
toolsBridge := services.NewToolsEngineBridge(toolsEngineService.engine)

contextEngine := contextengine.New(contextengine.Config{
    EnableHistoryCount: true,
    HistoryCount:       20,
})
contextBridge := services.NewContextEngineBridge(contextEngine)

// Initialize Agent Chat Service
agentChatService := services.NewAgentChatService(
    app,
    dbService,
    libService,
    kbService,
    toolsBridge,   // Phase 3
    contextBridge, // Phase 3
)
app.RegisterService(application.NewService(agentChatService))

// Initialize Thread Management Service
threadManagementService := services.NewThreadManagementService(app, dbService)
app.RegisterService(application.NewService(threadManagementService))
```

---

## ✅ Phase 3: Tools & Context Engine Integration (IMPLEMENTED)

### 1. **Context Engine Bridge** ✅
**File**: `internal/services/context_engine_bridge.go`  
**Status**: ✅ **FULLY IMPLEMENTED**

#### Fungsi:
- Bridge existing `contextengine.Engine` dengan Eino agent system
- Pre-process messages sebelum dikirim ke agent
- Apply all context processors (history truncate, templates, placeholders, etc.)

#### API:
```go
type ContextEngineBridge struct {
    contextEngine *contextengine.Engine
}

func NewContextEngineBridge(contextEngine *contextengine.Engine) *ContextEngineBridge

// Process messages through context engine
func (b *ContextEngineBridge) ProcessMessagesForAgent(
    ctx context.Context, 
    einoMessages []*schema.Message,
) ([]*schema.Message, error)

// Merge processed messages with RAG context
func (b *ContextEngineBridge) MergeWithRAGContext(
    processedMessages []*schema.Message,
    ragContext string,
) []*schema.Message
```

#### Features:
- ✅ Automatic message type conversion (Eino ↔ ContextEngine)
- ✅ History truncation & management
- ✅ Template & placeholder processing
- ✅ Optional RAG context merging
- ✅ Graceful fallback on errors

---

### 2. **Tools Engine Bridge** ✅
**File**: `internal/services/tools_engine_bridge.go`  
**Status**: ✅ **FULLY IMPLEMENTED**

#### Fungsi:
- Bridge existing `toolsengine.ToolsEngine` dengan Eino agent system
- Dynamic tool loading per request
- Expose existing tools to agents

#### API:
```go
type ToolsEngineBridge struct {
    toolsEngine *toolsengine.ToolsEngine
}

func NewToolsEngineBridge(toolsEngine *toolsengine.ToolsEngine) *ToolsEngineBridge

// Get tools for agent by IDs
func (b *ToolsEngineBridge) GetToolsForAgent(toolIDs []string) []tool.BaseTool

// Get all enabled tools
func (b *ToolsEngineBridge) getAllEnabledTools() []tool.BaseTool

// Get tool by ID
func (b *ToolsEngineBridge) GetToolByID(toolID string) (tool.BaseTool, error)

// Get available tool descriptions
func (b *ToolsEngineBridge) GetAvailableToolNames(ctx context.Context) ([]ToolDescription, error)
```

#### Features:
- ✅ Dynamic tool loading by ID
- ✅ Support for all existing tools (calculator, web-search, etc.)
- ✅ Direct Eino `tool.BaseTool` interface
- ✅ Tool enable/disable support
- ✅ Tool execution via bridge

---

### 3. **Agent Integration** ✅
**File**: `internal/services/agent_chat_service.go`  
**Status**: ✅ **FULLY IMPLEMENTED**

#### Changes:
```go
type AgentChatService struct {
    // ... existing fields ...
    toolsBridge   *ToolsEngineBridge   // Phase 3: NEW
    contextBridge *ContextEngineBridge // Phase 3: NEW
}

func NewAgentChatService(
    app *application.App,
    db *database.Service,
    libService *llama.LibraryService,
    kbService *KnowledgeBaseService,
    toolsBridge *ToolsEngineBridge,     // Phase 3: NEW
    contextBridge *ContextEngineBridge, // Phase 3: NEW
) *AgentChatService
```

#### Usage in `Chat()` method:
1. **Context Processing** (before agent):
```go
// Process messages through context engine (if available)
messagesToAgent := session.Messages
if s.contextBridge != nil {
    processedMessages, _ := s.contextBridge.ProcessMessagesForAgent(ctx, session.Messages)
    messagesToAgent = processedMessages
}
```

2. **Tools Loading** (during session creation):
```go
// Add tools from tools engine bridge
if s.toolsBridge != nil && len(req.Tools) > 0 {
    bridgeTools := s.toolsBridge.GetToolsForAgent(req.Tools)
    tools = append(tools, bridgeTools...)
}
```

#### Features:
- ✅ Optional bridges (can be nil)
- ✅ Graceful fallback if bridges fail
- ✅ Tools specified per ChatRequest
- ✅ Context processing for all messages
- ✅ Compatible with existing RAG & DB features

---

### 4. **Main Application Wiring** ✅
**File**: `main.go`  
**Status**: ✅ **FULLY IMPLEMENTED**

#### Initialization Flow:
```go
// Initialize tools engine bridge
toolsEngineService := NewToolsEngineService()
toolsBridge := services.NewToolsEngineBridge(toolsEngineService.engine)

// Initialize context engine bridge
contextEngine := contextengine.New(contextengine.Config{
    EnableHistoryCount: true,
    HistoryCount:       20,
})
contextBridge := services.NewContextEngineBridge(contextEngine)

// Create Agent Chat Service with bridges
agentChatService := services.NewAgentChatService(
    app,
    dbService,
    libService,
    kbService,
    toolsBridge,   // Phase 3
    contextBridge, // Phase 3
)
```

---

### 5. **Auto Topic Generation** ✅
**File**: `internal/services/agent_chat_service.go`  
**Status**: ✅ **FULLY IMPLEMENTED**

#### Fungsi:
- Automatically generate conversation title after first response
- Uses LLM to create concise, meaningful titles
- Saves topic to database for future reference
- Inspired by frontend `chainSummaryTitle.ts`

#### API:
```go
// Generate topic title using LLM
func (s *AgentChatService) generateTopicTitle(
    ctx context.Context,
    messages []*schema.Message,
    locale string,
) (string, error)

// Create topic for session after first response
func (s *AgentChatService) createTopicForSession(
    ctx context.Context,
    sessionID, userID string,
    messages []*schema.Message,
) error
```

#### How it works:
1. **Trigger**: After first assistant response (when `len(messages) == 2`)
2. **Background**: Runs in goroutine to not block response
3. **LLM Prompt**: Professional summarizer with strict rules:
   - Maximum 10 words
   - Maximum 50 characters
   - No punctuation or special characters
   - Language-aware (default: en-US)
4. **Database**: Saves to `topics` table with link to session
5. **De-duplication**: Checks if topic already exists for session

#### Example:
```
User: "What is CloudWeGo Eino?"
Assistant: "CloudWeGo Eino is a framework..."

→ LLM generates: "CloudWeGo Eino Framework Explanation"
→ Topic saved to database with session_id
```

#### Features:
- ✅ Automatic title generation
- ✅ Background processing (non-blocking)
- ✅ Database persistence
- ✅ De-duplication check
- ✅ Graceful fallback ("New Conversation")
- ✅ Locale support
- ✅ Same format as frontend `chainSummaryTitle`

---

### 10. **Thread Management Service** 🆕
**File**: `internal/services/thread_management_service.go`  
**Status**: ✅ **FULLY IMPLEMENTED**

#### Fungsi:
- Manage conversation branches (threads) within topics
- Allow users to create alternative conversation paths from any message
- Track thread relationships (parent-child, source message)
- Support multiple active threads per topic

#### API:
```go
type ThreadManagementService struct {
    app *application.App
    db  *database.Service
}

// Thread Creation
func (s *ThreadManagementService) CreateThread(
    ctx context.Context,
    req CreateThreadRequest,
) (*CreateThreadResponse, error)

// Thread Listing
func (s *ThreadManagementService) ListThreadsByTopic(
    ctx context.Context,
    req ListThreadsRequest,
) ([]ThreadInfo, error)

// Thread Retrieval
func (s *ThreadManagementService) GetThread(
    ctx context.Context,
    threadID, userID string,
) (*ThreadInfo, error)

// Thread Status Management
func (s *ThreadManagementService) UpdateThreadStatus(
    ctx context.Context,
    threadID, userID string,
    status ThreadStatus,
) error

// Thread Messages
func (s *ThreadManagementService) GetThreadMessages(
    ctx context.Context,
    threadID, userID string,
) ([]*db.Message, error)
```

#### Request/Response Types:
```go
type CreateThreadRequest struct {
    Title           string     // Thread title
    Type            ThreadType // "continuation" | "standalone"
    TopicID         string     // Parent topic
    SourceMessageID string     // Message to branch from
    ParentThreadID  string     // Optional parent thread
    UserID          string     // Owner
}

type ThreadInfo struct {
    ID              string
    Title           string
    Type            string     // "continuation" | "standalone"
    Status          string     // "active" | "deprecated" | "archived"
    TopicID         string
    SourceMessageID string     // Message that was branched from
    ParentThreadID  string
    LastActiveAt    int64
    CreatedAt       int64
    UpdatedAt       int64
}

// Thread Types
const (
    ThreadTypeContinuation = "continuation"  // Continue from message
    ThreadTypeStandalone   = "standalone"    // Fresh start
)

// Thread Status
const (
    ThreadStatusActive     = "active"        // Currently active
    ThreadStatusDeprecated = "deprecated"    // Old/unused
    ThreadStatusArchived   = "archived"      // Archived
)
```

#### Database Schema:
```sql
CREATE TABLE IF NOT EXISTS threads (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  type TEXT NOT NULL CHECK (type IN ('continuation', 'standalone')),
  status TEXT DEFAULT 'active' CHECK (status IN ('active', 'deprecated', 'archived')),
  topic_id TEXT NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
  source_message_id TEXT NOT NULL,
  parent_thread_id TEXT REFERENCES threads(id) ON DELETE SET NULL,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  last_active_at INTEGER NOT NULL,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);
```

#### Features:
- ✅ Create conversation branches from any message
- ✅ Multiple thread types (continuation, standalone)
- ✅ Thread status management (active, deprecated, archived)
- ✅ Parent-child thread relationships
- ✅ Track source message for each thread
- ✅ List all threads for a topic
- ✅ Get messages within a specific thread
- ✅ Update thread metadata
- ✅ Delete threads

#### Relationship Hierarchy:
```
Session (sess-001)
  │
  ├─ Topic (topic-main) [Auto-generated title]
  │    │
  │    ├─ Thread A (active) ← Default/main thread
  │    │    ├─ Message 1 (user)
  │    │    ├─ Message 2 (assistant)
  │    │    ├─ Message 3 (user) ──┐
  │    │    └─ Message 4 (assistant) │
  │    │                              │
  │    └─ Thread B (active) ←────────┘ Branch from Message 3
  │         ├─ Message 5 (assistant) <- Alternative response
  │         └─ Message 6 (user)
  │
  └─ Topic (topic-other)
       └─ ...
```

#### Use Cases:

**1. Alternative Response Path**:
```
User: "Explain Eino"
Assistant A: "Eino is a framework..." (Message 2)

→ User clicks "branch" on Message 2
→ Creates Thread B from Message 2
→ New alternative response: "Let me explain differently..."
```

**2. Explore Different Directions**:
```
Thread A: Technical deep-dive
Thread B: High-level overview
Thread C: Practical examples
→ All branched from same source message
```

**3. Conversation Versioning**:
```
Main Thread: Current conversation
Deprecated Threads: Previous attempts/versions
→ Keep history without cluttering main view
```

#### Frontend Integration:
```typescript
// User clicks "branching" button on a message
const thread = await ThreadManagementService.CreateThread({
    source_message_id: messageId,
    topic_id: currentTopicId,
    user_id: currentUserId,
    title: "Alternative discussion",
    type: "continuation"
});

// Switch to new thread
setActiveThreadId(thread.thread_id);

// Continue conversation in new thread
await AgentChatService.Chat({
    session_id: sessionId,
    user_id: userId,
    message: userMessage,
    thread_id: thread.thread_id  // Messages go to this thread
});
```

#### SQLC Queries Used:
```sql
-- threads.sql
CreateThread
ListThreadsByTopic
GetThread
UpdateThread
DeleteThread

-- messages.sql (extended)
ListMessagesByThread  -- NEW: Get messages for a thread
```

#### Integration Points:
- **AgentChatService**: Will be extended to handle `thread_id` in chat requests
- **Topics**: Each thread belongs to a topic
- **Messages**: Each message can have a `thread_id` field
- **Frontend**: Branching UI action button (from `Actions.tsx`)

---

## ❌ Komponen Belum Diimplementasikan

### 1. **Frontend Integration**
**Status**: ⏳ **NOT YET IMPLEMENTED**

#### Yang Perlu:
- Generate TypeScript types untuk KB & Agent services
- UI untuk KB management (create, upload files, query)
- Chat UI dengan KB selection
- Source citation display

#### File Target:
- `frontend/src/types/knowledge-base.ts` (new)
- `frontend/src/types/agent-chat.ts` (new)
- `frontend/src/services/kb-service.ts` (new)
- `frontend/src/components/KnowledgeBase/` (new)

#### TypeScript Types:
```typescript
interface ChatRequest {
    session_id: string;
    user_id: string;
    message: string;
    knowledge_base_id?: string;
    tools?: string[];
    context?: Record<string, any>;
    temperature?: number;
    max_tokens?: number;
    stream?: boolean;
}

interface ChatResponse {
    session_id: string;
    message: string;
    tool_calls?: ToolCall[];
    sources?: Document[];
    finish_reason: string;
    usage?: TokenUsage;
}
```

---

### 4. **Streaming Implementation**
**Status**: ⏳ **PARTIALLY IMPLEMENTED**

#### Current State:
- `LlamaEinoModel.Stream()` ✅ implemented
- `AgentChatService.ChatStream()` ⏳ placeholder (returns sync result)

#### Yang Perlu:
- Understand Wails v3 events API
- Implement proper SSE/WebSocket streaming
- Frontend listener untuk stream events

---

### 5. **Advanced Features**
**Status**: ⏳ **NOT YET IMPLEMENTED**

#### Multi-Agent Orchestration:
- Transfer between agents
- Sub-agent delegation
- Agent handoff

#### Advanced RAG:
- Hybrid search (vector + BM25)
- Re-ranking
- Query expansion
- Multi-step reasoning

#### Monitoring & Observability:
- Token usage tracking
- Tool call tracing
- Performance metrics
- Error logging

---

## 🔄 Data Flow - Complete Path

### Scenario: User Chat dengan KB (Multi-turn)

```
┌─────────────────────────────────────────────────────────────────┐
│                     FIRST MESSAGE (Session Creation)             │
└─────────────────────────────────────────────────────────────────┘

1. Frontend generates SessionID: "session-abc123"
   └─ Store in localStorage/React state
   ↓
2. Frontend sends ChatRequest
   {
     session_id: "session-abc123",
     user_id: "user-001",
     message: "What is CloudWeGo Eino?",
     knowledge_base_id: "kb-xyz",
   }
   ↓
3. AgentChatService.Chat() receives request
   ↓
4. getOrCreateSession("session-abc123")
   ├─ Lock: sessionsMutex
   ├─ Check: sessions["session-abc123"] exists?
   │  └─ NO → Create new session:
   │         ├─ Initialize ADK Agent with:
   │         │  ├─ LlamaEinoModel
   │         │  ├─ KB Search Tool (for kb-xyz)
   │         │  ├─ System instruction
   │         │  └─ Max iterations = 5
   │         ├─ Create AgentSession{
   │         │    SessionID: "session-abc123",
   │         │    Agent: agent,
   │         │    Messages: [],
   │         │    KnowledgeBaseID: "kb-xyz",
   │         │  }
   │         └─ sessions["session-abc123"] = session
   ├─ Unlock: sessionsMutex
   └─ Return: session
   ↓
5. Add user message to session.Messages
   session.Messages = [
     {Role: "user", Content: "What is CloudWeGo Eino?"}
   ]
   ↓
6. agent.Run(session.Messages)
   ├─ ADK Agent processes with LlamaEinoModel
   ├─ Agent decides to call KB search tool
   ├─ Tool: search_knowledge_base_kb_xyz
   │  ├─ kbSearchTool.InvokableRun(query="CloudWeGo Eino")
   │  ├─ KBService.QueryKnowledgeBase("kb-xyz", "CloudWeGo Eino", topK=5)
   │  ├─ Retriever.Retrieve() (Eino adapter)
   │  ├─ Chromem.Query() (vector search)
   │  └─ Return: [doc1, doc2, doc3] with relevant chunks
   ├─ Agent receives tool result
   ├─ Agent generates final response using context from docs
   └─ Return: AgentOutput{
        Message: "CloudWeGo Eino is...",
        ToolCalls: [{tool: "search_knowledge_base_kb_xyz", ...}]
      }
   ↓
7. Update session.Messages
   session.Messages = [
     {Role: "user", Content: "What is CloudWeGo Eino?"},
     {Role: "assistant", Content: "CloudWeGo Eino is...", 
      ToolCalls: [...]},
   ]
   session.UpdatedAt = now()
   ↓
8. AgentChatService collects response data:
   ├─ Final message
   ├─ Tool calls made
   ├─ Sources from KB (doc1, doc2, doc3)
   └─ Usage statistics
   ↓
9. Return ChatResponse to frontend
   {
     session_id: "session-abc123",
     message: "CloudWeGo Eino is...",
     sources: [doc1, doc2, doc3],
     tool_calls: [{...}],
   }
   ↓
10. Frontend displays response
    └─ Keep SessionID for next message


┌─────────────────────────────────────────────────────────────────┐
│                   FOLLOW-UP MESSAGE (Session Reuse)              │
└─────────────────────────────────────────────────────────────────┘

1. User asks follow-up: "Tell me more about ADK"
   ↓
2. Frontend sends ChatRequest
   {
     session_id: "session-abc123",  // SAME SessionID
     user_id: "user-001",
     message: "Tell me more about ADK",
     knowledge_base_id: "kb-xyz",   // SAME KB
   }
   ↓
3. AgentChatService.Chat() receives request
   ↓
4. getOrCreateSession("session-abc123")
   ├─ Lock: sessionsMutex
   ├─ Check: sessions["session-abc123"] exists?
   │  └─ YES → Return existing session
   │         ├─ session.Agent (same agent instance)
   │         └─ session.Messages (has history)
   └─ Unlock: sessionsMutex
   ↓
5. Add user message to EXISTING session.Messages
   session.Messages = [
     {Role: "user", Content: "What is CloudWeGo Eino?"},
     {Role: "assistant", Content: "CloudWeGo Eino is...", ToolCalls: [...]},
     {Role: "user", Content: "Tell me more about ADK"},  // NEW
   ]
   ↓
6. agent.Run(session.Messages)  // Agent has FULL context
   ├─ ADK Agent processes with history
   ├─ Understands "ADK" refers to Eino ADK from previous context
   ├─ May call KB tool again for more details
   └─ Generates contextual response
   ↓
7. Update session.Messages with assistant response
   session.Messages = [
     {Role: "user", Content: "What is CloudWeGo Eino?"},
     {Role: "assistant", Content: "CloudWeGo Eino is...", ToolCalls: [...]},
     {Role: "user", Content: "Tell me more about ADK"},
     {Role: "assistant", Content: "Eino ADK is...", ToolCalls: [...]},  // NEW
   ]
   ↓
8. Return ChatResponse to frontend
   {
     session_id: "session-abc123",  // SAME SessionID
     message: "Eino ADK is...",
     sources: [...],
     tool_calls: [...],
   }


┌─────────────────────────────────────────────────────────────────┐
│                    SESSION CLEANUP                               │
└─────────────────────────────────────────────────────────────────┘

Option 1: Manual (User closes conversation)
  Frontend calls: AgentChatService.ClearSession("session-abc123")
  ↓
  sessions["session-abc123"] = nil
  └─ Agent & history released

Option 2: Automatic (App restart)
  All sessions cleared from memory
  └─ Frontend must create new session on next launch
```

---

## 📊 Test Coverage

### Implemented Tests:
- ✅ End-to-end test: `cmd/test-kb-rag/main.go`
  - KB creation
  - Document upload
  - Query KB
  - RAG workflow
  - Agent execution

### Missing Tests:
- ⏳ Unit tests per service
- ⏳ Integration tests
- ⏳ Performance benchmarks

---

## 🚀 How to Use

### 1. Create Knowledge Base
```go
kbID, err := kbService.CreateKnowledgeBase(ctx, "My KB", "Description", userID)
```

### 2. Add Documents
```go
err := kbService.AddFileToKnowledgeBase(ctx, kbID, "/path/to/doc.pdf", metadata, userID)
```

### 3. Chat with Agent
```go
req := ChatRequest{
    SessionID:       "session-123",
    UserID:          "user-001",
    Message:         "What is CloudWeGo Eino?",
    KnowledgeBaseID: kbID,
    Temperature:     0.7,
    MaxTokens:       500,
}

resp, err := agentChatService.Chat(ctx, req)
// resp.Message = agent's response
// resp.Sources = documents used from KB
// resp.ToolCalls = tools executed
```

---

## 📝 Next Steps

### Priority 1 (High Impact):
1. **Streaming Implementation** - Real-time responses
2. **Frontend Integration** - UI untuk KB & chat
3. **Tools Engine Bridge** - Extend agent capabilities

### Priority 2 (Enhancement):
4. **Context Engine Bridge** - Better message processing
5. **Multi-agent** - Agent orchestration
6. **Advanced RAG** - Better retrieval quality

### Priority 3 (Nice to Have):
7. **Monitoring** - Usage tracking
8. **Testing** - Comprehensive test suite
9. **Documentation** - API docs, tutorials

---

---

## 📊 Status Implementasi Detail

### ✅ **SUDAH DIIMPLEMENTASI** (95% Complete)

| Komponen | Status | File | Fungsi Utama |
|----------|--------|------|--------------|
| **Llama Eino Model** | ✅ DONE | `internal/llama/eino_adapter.go` | Adapter Llama.cpp → Eino |
| **Knowledge Base Service** | ✅ DONE | `internal/services/knowledge_base.go` | CRUD KB, file upload, query |
| **RAG Workflow** | ✅ DONE | `internal/services/rag_workflow.go` | Retrieval + context building |
| **RAG Agent** | ✅ DONE | `internal/services/rag_agent.go` | ADK agent dengan KB tools |
| **Agent Chat Service** | ✅ DONE | `internal/services/agent_chat_service.go` | **Main service untuk chat** |
| **Tools Engine Bridge** | ✅ DONE | `internal/services/tools_engine_bridge.go` | Integrate existing tools |
| **Context Engine Bridge** | ✅ DONE | `internal/services/context_engine_bridge.go` | Message preprocessing |
| **Thread Management** | ✅ DONE | `internal/services/thread_management_service.go` | Conversation branching |
| **Database Queries** | ✅ DONE | `internal/database/queries/*.sql` | SQLC queries (KB, sessions, topics, threads) |
| **Main App Wiring** | ✅ DONE | `main.go` | Semua services terintegrasi |
| **E2E Test** | ✅ DONE | `cmd/test-kb-rag/main.go` | Full test suite |

### ⏳ **BELUM DIIMPLEMENTASI** (25% Remaining)

| Komponen | Status | Priority | Target File |
|----------|--------|----------|-------------|
| **Session Persistence** | ⚠️ CRITICAL | HIGH | `agent_chat_service.go` update |
| **Context Engine Bridge** | ⏳ TODO | MEDIUM | `internal/services/context_rag_bridge.go` |
| **Tools Engine Bridge** | ⏳ TODO | MEDIUM | `internal/services/agent_tools_bridge.go` |
| **Frontend Types** | ⏳ TODO | HIGH | `frontend/src/types/*.ts` |
| **Streaming** | ⏳ TODO | HIGH | `agent_chat_service.go` update |
| **UI Components** | ⏳ TODO | HIGH | `frontend/src/components/*` |

---

## ⚠️ PENTING: Session Storage Strategy

### Current Implementation (In-Memory):
```go
type AgentChatService struct {
    sessions map[string]*AgentSession  // ❌ In-memory only
}
```

**Masalah**:
- ❌ Lost on app restart
- ❌ Not shareable across devices
- ❌ No conversation history persistence

### Database sudah punya tabel `sessions`! 

```sql
CREATE TABLE IF NOT EXISTS sessions (
  id TEXT PRIMARY KEY,
  slug TEXT NOT NULL,
  title TEXT,
  description TEXT,
  avatar TEXT,
  background_color TEXT,
  type TEXT DEFAULT 'agent',
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  group_id TEXT REFERENCES session_groups(id) ON DELETE SET NULL,
  client_id TEXT,
  pinned INTEGER DEFAULT 0 NOT NULL,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL,
  UNIQUE(slug, user_id)
);

-- Plus junction tables:
CREATE TABLE agents_to_sessions (...)  -- Link agent to session
CREATE TABLE files_to_sessions (...)   -- Link files to session
```

### ✅ Recommended: Hybrid Storage Strategy

```
┌──────────────────────────────────────────────────────────┐
│          HYBRID SESSION STORAGE                          │
└──────────────────────────────────────────────────────────┘

DATABASE (SQLite) - Persistent
├─ sessions table
│  ├─ Session metadata (id, title, user_id, etc.)
│  ├─ Linked agents (agents_to_sessions)
│  └─ Linked KB (via agent → agents_knowledge_bases)
│
└─ messages table (already exists)
   └─ Full conversation history

MEMORY (AgentChatService) - Runtime only
├─ Active agent instances (ADK Agent)
├─ Eino model instances (LlamaEinoModel)
└─ Cached retrievers & tools

Flow:
1. Load session from DB (metadata + messages)
2. Reconstruct agent in-memory with history
3. Run agent with new message
4. Save new messages to DB
5. Keep agent in-memory for performance
6. Cleanup after inactivity timeout
```

### Implementation Priority:

**Phase 1 (Current)**: ✅ In-memory sessions
- Works for single-session usage
- Fast prototyping
- No persistence

**Phase 2 (Next)**: ⚠️ Hybrid persistence (RECOMMENDED)
- Use existing `sessions` table
- Use existing `messages` table
- Load/save on demand
- Lazy agent reconstruction

**Phase 3 (Future)**: 🚀 Advanced features
- Session groups
- Multi-device sync
- Conversation export/import

---

## 🚀 Quick Start

### Build & Run
```bash
# Build main app
go build -o bin/veridium .

# Run E2E test
go run cmd/test-kb-rag/main.go
```

### Example Usage (Backend)

#### 1. Create Knowledge Base
```go
kbID, err := kbService.CreateKnowledgeBase(ctx, "My KB", "desc", userID)
```

#### 2. Upload Documents
```go
// Support: DOCX, PDF, XLSX, HTML, TXT, MD
err := kbService.AddFileToKnowledgeBase(ctx, kbID, "/path/file.pdf", metadata, userID)
// Otomatis: parse → chunk → embed → store
```

#### 3. Chat dengan Agent (+ RAG)
```go
req := services.ChatRequest{
    SessionID:       "sess-123",
    UserID:          "user-001", 
    Message:         "Jelaskan tentang Eino",
    KnowledgeBaseID: kbID,  // KB akan digunakan untuk RAG
    Temperature:     0.7,
    MaxTokens:       500,
}

resp, err := agentChatService.Chat(ctx, req)
// resp.Message = jawaban agent
// resp.Sources = dokumen dari KB yang digunakan
// resp.ToolCalls = tools yang dipanggil agent
```

#### 4. Multi-turn Conversations
```go
// Session otomatis maintained
req1 := ChatRequest{SessionID: "sess-1", Message: "Apa itu Eino?"}
resp1 := agentChatService.Chat(ctx, req1)

req2 := ChatRequest{SessionID: "sess-1", Message: "Jelaskan lebih detail"}
resp2 := agentChatService.Chat(ctx, req2)
// Agent punya context dari req1
```

### Example Usage (Frontend)

```typescript
// Via Wails binding
const response = await AgentChatService.Chat({
    session_id: "my-session",
    user_id: "user-123",
    message: "What is CloudWeGo Eino?",
    knowledge_base_id: "kb-001",  // Optional
    temperature: 0.7,
    max_tokens: 500
});

console.log(response.message);      // Agent's answer
console.log(response.sources);      // Documents from KB
console.log(response.tool_calls);   // Tools used
```

---

## 🎯 Roadmap

### Phase 1: Core Backend ✅ **DONE** (Current)
- [x] Llama Eino Model Adapter
- [x] Knowledge Base Service (full CRUD)
- [x] RAG Workflow
- [x] RAG Agent (Eino ADK)
- [x] Agent Chat Service (in-memory sessions)
- [x] Database Integration (KB, files, queries)
- [x] Main App Wiring
- [x] E2E Test

### Phase 2: Session Persistence ✅ **COMPLETE**
- [x] **Update AgentChatService** to use existing `sessions` table
- [x] Load session metadata from DB
- [x] Load conversation history from `messages` table
- [x] Reconstruct agent with history
- [x] Save new messages to DB (user + assistant)
- [x] Update session timestamps
- [x] Hybrid storage (in-memory cache + DB)
- [x] DB ↔ Eino message conversion

### Phase 3: Integration ✅ **COMPLETE**
- [x] Context Engine Bridge
- [x] Tools Engine Bridge
- [x] Agent integration with bridges
- [x] Auto Topic Generation
- [x] Thread Management Service
- [x] Main.go wiring
- [ ] Advanced RAG features (hybrid search, re-ranking) - Optional for later

### Phase 4: Production Features
- [ ] Frontend UI & TypeScript types
- [ ] Streaming (Wails v3 events)
- [ ] Monitoring & analytics
- [ ] Comprehensive testing suite
- [ ] API documentation

---

## 📈 Progress Breakdown

**Overall: 95% Complete**

- ✅ Backend infrastructure: **100%**
- ✅ Core services: **100%**
- ✅ Eino integration: **100%**
- ✅ **Session persistence: 100%** (hybrid: DB + cache)
- ✅ **Tools & Context integration: 100%**
- ✅ **Auto Topic Generation: 100%**
- ✅ **Thread Management: 100%** ← NEW!
- ⏳ Advanced features: **60%**
- ⏳ Frontend: **0%**

---

## 🎉 Summary

### ✅ Siap Digunakan Sekarang:
1. **Knowledge Base** - Create, upload, query documents
2. **Agent Chat** - Multi-turn dengan RAG capabilities
3. **Tool Calling** - Agent bisa search KB secara otomatis
4. **Session Management** - Context-aware conversations (in-memory)
5. **Source Tracking** - Tahu informasi dari mana

### ✅ Baru Ditambahkan (Phase 3 + Enhancements):
1. **Tools Engine Bridge** - ✅ Existing tools accessible by agent
2. **Context Engine Bridge** - ✅ Message preprocessing (history, templates)
3. **Dynamic Tools** - ✅ Tools specified per request
4. **Context Processing** - ✅ All messages processed before agent
5. **Graceful Fallback** - ✅ Works with/without bridges
6. **Auto Topic Generation** - ✅ LLM creates title after first response
7. **Thread Management** - ✅ Conversation branching & multi-threading

### ⏳ Untuk Production:
1. Frontend UI & components
2. Streaming responses
3. Integration dengan existing engines
4. Advanced RAG features
5. Monitoring & analytics

### 🔧 Database Tables (Status):
- ✅ `sessions` - Session metadata *(USED)*
- ✅ `messages` - Conversation history *(USED)*
- ✅ `topics` - Conversation titles *(USED - Phase 3)*
- ✅ `threads` - Conversation branches *(USED - Phase 3)*
- ✅ `knowledge_bases` - KB metadata *(USED)*
- ✅ `knowledge_base_files` - KB-file links *(USED)*
- ⏳ `agents_to_sessions` - Link agents to sessions *(available)*
- ⏳ `files_to_sessions` - Link files to sessions *(available)*
- ⏳ `session_groups` - Organize sessions *(available)*

**Bottom line**: Core system **production-ready** untuk backend single-session usage. **Next critical step**: Integrate with existing database sessions for persistence.

---

## 📋 Implementation Checklist

### ✅ Phase 1: Core Backend (COMPLETED)
- [x] Llama.cpp → Eino adapter (`LlamaEinoModel`)
- [x] Knowledge Base CRUD operations
- [x] Document parsing & chunking (DOCX, PDF, XLSX, HTML, TXT, MD)
- [x] Vector storage with Chromem
- [x] Eino Retriever & Indexer adapters
- [x] RAG workflow (retrieval + context building)
- [x] ADK Agent with tool calling
- [x] Agent Chat Service (in-memory sessions)
- [x] KB search tool (auto-generated per KB)
- [x] Multi-turn conversations support
- [x] Source tracking
- [x] Main app integration
- [x] E2E test suite

### ✅ Phase 2: Database Persistence (COMPLETED)
- [x] Update `AgentChatService` to use `sessions` table
- [x] Load session metadata from DB
- [x] Load conversation history from `messages` table
- [x] Convert DB messages → Eino schema messages
- [x] Reconstruct agent with full history
- [x] Save user messages to DB
- [x] Save assistant responses to DB
- [x] Save tool calls to DB (JSON in `tools` column)
- [x] Update session timestamps
- [x] Handle session resume after app restart
- [x] Hybrid storage (cache + DB)
- [x] Thread-safe operations
- [ ] Link agents to sessions (`agents_to_sessions`) - Optional for later
- [ ] Session cleanup/timeout - Optional for later

### ✅ Phase 3: Integration (COMPLETED)
- [x] Context Engine Bridge
  - [x] Pre-process messages with context engine
  - [x] Automatic message type conversion
  - [x] History truncation & template processing
  - [x] Merge context with RAG results (optional)
- [x] Tools Engine Bridge
  - [x] Auto-discover tools from tools engine
  - [x] Direct Eino `tool.BaseTool` interface
  - [x] Dynamic tool loading per request
  - [x] Tool enable/disable support
- [x] Agent Integration
  - [x] Wire bridges into AgentChatService
  - [x] Optional bridge support (graceful fallback)
  - [x] Main.go initialization
- [x] Auto Topic Generation
  - [x] LLM-powered title generation
  - [x] Background processing
  - [x] Database persistence
  - [x] De-duplication check
- [x] Thread Management
  - [x] Thread creation with branching
  - [x] List threads by topic
  - [x] Get thread details
  - [x] Update thread status
  - [x] Get messages by thread
  - [x] Thread types (continuation, standalone)
  - [x] Thread status (active, deprecated, archived)
- [ ] Advanced RAG (Optional for Phase 4)
  - [ ] Hybrid search (vector + BM25)
  - [ ] Re-ranking
  - [ ] Query expansion
  - [ ] Multi-step reasoning

### 📱 Phase 4: Frontend & Production (FUTURE)
- [ ] Generate TypeScript types
  - [ ] `ChatRequest` interface
  - [ ] `ChatResponse` interface
  - [ ] `KnowledgeBase` types
  - [ ] `Session` types
- [ ] UI Components
  - [ ] Knowledge Base management
  - [ ] Session list & groups
  - [ ] Chat interface
  - [ ] Source citation display
  - [ ] Tool call visualization
- [ ] Streaming implementation
  - [ ] Wails v3 events API
  - [ ] Frontend event listeners
  - [ ] Streaming UI updates
- [ ] Monitoring & Analytics
  - [ ] Token usage tracking
  - [ ] Performance metrics
  - [ ] Error logging
- [ ] Testing
  - [ ] Unit tests per service
  - [ ] Integration tests
  - [ ] Performance benchmarks
- [ ] Documentation
  - [ ] API documentation
  - [ ] User guides
  - [ ] Architecture diagrams

---

## 🔑 Key Takeaways

### ✅ What Works Now:
1. **Complete RAG Pipeline**: Documents → Parse → Chunk → Embed → Store → Retrieve → Generate
2. **Agent-based Chat**: Eino ADK agents with tool calling and reasoning
3. **Knowledge Base**: Full CRUD with automatic indexing
4. **Multi-turn**: Conversations with context (in-memory)
5. **Local-First**: Everything runs locally with Llama.cpp

### ✅ Phase 2 Achievement:
**Session persistence is now FULLY implemented!**
- ✅ Tables used: `sessions`, `messages`
- ✅ SQLC queries integrated
- ✅ AgentChatService uses hybrid storage
- ✅ Sessions survive app restart
- ✅ Auto-reconstruction with full history

### ✅ Phase 3 Achievement:
**Complete ecosystem integration + conversation management!**
- ✅ Tools Engine Bridge - Dynamic tool loading
- ✅ Context Engine Bridge - Message preprocessing
- ✅ Auto Topic Generation - LLM-powered titles
- ✅ Thread Management - Conversation branching
- ✅ Tables used: `topics`, `threads`
- 🎯 **Next step**: Frontend UI for full system usage

### 🎯 Priority Order:
1. ~~**Session DB Integration**~~ ✅ **COMPLETE (Phase 2)**
2. ~~**Tools/Context Engine Integration**~~ ✅ **COMPLETE (Phase 3)**
3. ~~**Auto Topic Generation**~~ ✅ **COMPLETE (Phase 3)**
4. ~~**Thread Management**~~ ✅ **COMPLETE (Phase 3)**
5. **Frontend UI** (HIGH) - Make it usable for end-users
6. **Streaming** (HIGH) - Better UX
7. **Advanced Features** (LOW) - Nice to have

### 💡 Architecture Strengths:
- **Modular**: Each component is independent
- **Extensible**: Easy to add new tools, engines, models
- **Type-safe**: SQLC generates type-safe DB queries
- **Standards-based**: Uses Eino interfaces throughout
- **Local-first**: No external API dependencies

### 🚧 Technical Debt:
- ~~Session persistence~~ ✅ **FIXED in Phase 2**
- Streaming implementation (placeholder only)
- Frontend types generation
- Comprehensive test coverage
- Performance optimization (batching message saves)
- Session cleanup/timeout (optional enhancement)

---

## 🔄 Migration Strategy: Frontend → Backend

> **NEXT PHASE: Frontend Migration**  
> Replace frontend business logic dengan backend AgentChatService  
> **Timeline**: 3 weeks | **Status**: Ready to implement

### 📊 Current Problem

**Frontend** (`generateAIChat.ts`): **1147 lines!**
- ❌ Message creation & validation
- ❌ Topic auto-generation (with LLM)
- ❌ Thread management
- ❌ Message persistence (IndexedDB)
- ❌ Context engineering
- ❌ Tool orchestration
- ❌ RAG workflow

**Backend** (`AgentChatService`): **95% ready but unused!**
- ✅ Agent execution (Eino ADK)
- ✅ RAG retrieval
- ✅ Tool calling
- ✅ Session persistence (SQLite)
- ✅ Topic generation
- ✅ Thread management

**Problem**: Backend hanya jadi "LLM proxy", semua logic di frontend!

### 🎯 Solution

**Move ALL business logic to backend**, frontend hanya UI layer.

### 📋 3-Week Implementation Plan

#### **Week 1: Backend API Extension**

**Goal**: Make AgentChatService API frontend-ready

**Changes**:
```go
// File: internal/services/agent_chat_service.go

// 1. Extend ChatRequest
type ChatRequest struct {
    // ... existing ...
    TopicID  string `json:"topic_id"`  // NEW
    ThreadID string `json:"thread_id"` // NEW
    ParentID string `json:"parent_id"` // NEW
}

// 2. Extend ChatResponse
type ChatResponse struct {
    // ... existing ...
    MessageID string `json:"message_id"` // NEW - created msg ID
    TopicID   string `json:"topic_id"`   // NEW - may be auto-created
    ThreadID  string `json:"thread_id"`  // NEW
    CreatedAt int64  `json:"created_at"` // NEW
}

// 3. Add to AgentSession
type AgentSession struct {
    // ... existing ...
    TopicID  string // NEW
    ThreadID string // NEW
}

// 4. Update Chat method
func (s *AgentChatService) Chat(ctx, req) (*ChatResponse, error) {
    session, _ := s.getOrCreateSession(ctx, req)
    
    // Auto-create topic if needed
    topicID := req.TopicID
    if topicID == "" && len(session.Messages) >= 2 {
        topicID, _ = s.createTopicSync(ctx, session, req.UserID)
    }
    
    // Load thread context if specified
    if req.ThreadID != "" {
        threadMsgs, _ := s.threadService.GetThreadMessages(ctx, req.ThreadID, req.UserID)
        session.Messages = s.convertToEino(threadMsgs)
    }
    
    // Save with topic_id & thread_id
    userMsgID, _ := s.saveMessageToDB(ctx, userMsg, session.SessionID, req.UserID, req.ThreadID, topicID)
    
    // ... agent execution ...
    
    assistantMsgID, _ := s.saveMessageToDB(ctx, assistantMsg, session.SessionID, req.UserID, req.ThreadID, topicID)
    
    return &ChatResponse{
        MessageID: assistantMsgID,
        TopicID:   topicID,
        ThreadID:  req.ThreadID,
        // ...
    }, nil
}
```

**Tasks**:
- [ ] Day 1-2: Extend structs
- [ ] Day 2-3: Add topic auto-creation
- [ ] Day 3-4: Update message persistence
- [ ] Day 4: Wire ThreadManagementService
- [ ] Day 5: Test with curl

#### **Week 2: Frontend Refactoring**

**Goal**: Reduce 1147 lines → 200 lines

**New File**: `frontend/src/services/backendAgentChat.ts`
```typescript
import { AgentChatService } from '@@/github.com/kawai-network/veridium/internal/services';

class BackendAgentChatService {
  async sendMessage(params: SendMessageParams) {
    const response = await AgentChatService.Chat(params);
    if (response.error) throw new Error(response.error);
    return response;
  }
}

export const backendAgentChat = new BackendAgentChatService();
```

**Simplify**: `frontend/src/store/chat/slices/aiChat/actions/generateAIChat.ts`
```typescript
const sendMessage = async (params: SendMessageParams) => {
  try {
    set({ isCreatingMessage: true });
    
    // Backend handles EVERYTHING!
    const response = await backendAgentChat.sendMessage({
      session_id: get().activeId,
      user_id: getCurrentUserId(),
      message: params.message,
      topic_id: get().activeTopicId || '',
      thread_id: get().activeThreadId || '',
      tools: getEnabledTools(),
    });
    
    // Update UI only
    set(state => ({
      activeTopicId: response.topic_id || state.activeTopicId,
      messagesMap: {
        ...state.messagesMap,
        [messageMapKey(response.session_id, response.topic_id)]: [
          ...state.messagesMap[messageMapKey(response.session_id, response.topic_id)],
          {
            id: response.message_id,
            role: 'assistant',
            content: response.message,
            createdAt: response.created_at,
          },
        ],
      },
      isCreatingMessage: false,
    }));
    
    return response;
  } catch (error) {
    set({ isCreatingMessage: false });
    throw error;
  }
};
```

**Remove**:
- ❌ `internal_coreProcessMessage` (150+ lines)
- ❌ `internal_fetchAIChatMessage` (100+ lines)
- ❌ Topic generation (50+ lines)
- ❌ Context engineering (30+ lines)
- ❌ Tools orchestration (40+ lines)

**Tasks**:
- [ ] Day 1: Create backend wrapper
- [ ] Day 2-3: Simplify generateAIChat.ts
- [ ] Day 4-5: Test integration

#### **Week 3: Testing & Polish**

**Tasks**:
- [ ] Day 1-2: E2E testing (full chat flow, topic creation, thread branching)
- [ ] Day 3-4: Edge cases (errors, network failures)
- [ ] Day 5: Performance optimization

### ✅ Success Criteria

**Backend**:
- ✅ Handles full chat lifecycle
- ✅ Topic auto-creation works
- ✅ Thread management integrated
- ✅ All data in SQLite

**Frontend**:
- ✅ < 300 lines in generateAIChat.ts
- ✅ No business logic
- ✅ UI updates only

**Data**:
- ✅ SQLite single source of truth
- ✅ No IndexedDB for chat

### 🚀 Start Now (First 4 Hours)

**Step 1**: Extend `ChatRequest`/`ChatResponse` (1 hour)
```go
// Add TopicID, ThreadID, MessageID fields
```

**Step 2**: Add topic auto-creation (2 hours)
```go
// Implement createTopicSync() and integrate
```

**Step 3**: Test (1 hour)
```bash
go build && ./bin/veridium
curl -X POST http://localhost:8080/api/chat \
  -d '{"session_id":"test","user_id":"me","message":"hi"}'
```

### 📈 Expected Benefits

1. **Code Reduction**: 1147 lines → 200 lines (82% reduction)
2. **Better Architecture**: Clear separation of concerns
3. **Single Source of Truth**: Backend owns all logic
4. **Easier Maintenance**: Changes in one place
5. **Faster Feature Development**: Backend changes auto-benefit frontend

---

## 📚 References

### Code Locations:
- **Services**: `internal/services/`
- **Database**: `internal/database/`
- **Eino Adapters**: `pkg/eino-adapters/`
- **Llama Integration**: `internal/llama/`, `pkg/yzma/llama/`
- **Tests**: `cmd/test-kb-rag/`

### Key Files:
- Architecture: `docs/EINO_ARCHITECTURE.md` (this file)
- Main: `main.go`
- Agent Chat: `internal/services/agent_chat_service.go`
- KB Service: `internal/services/knowledge_base.go`
- Eino Model: `internal/llama/eino_adapter.go`
- DB Schema: `internal/database/schema/schema.sql`

### External Dependencies:
- **CloudWeGo Eino**: `cloudwego/eino/`
- **Chromem**: `pkg/chromem/`
- **Llama.cpp**: `pkg/yzma/llama/`
- **SQLC**: Generates `internal/database/generated/`
- **Wails v3**: Application framework

