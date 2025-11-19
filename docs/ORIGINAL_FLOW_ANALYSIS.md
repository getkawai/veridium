# Original Frontend Flow Analysis

> **Purpose**: Understanding the complete flow before migration  
> **Date**: November 2025  
> **Status**: 🔍 Analysis Complete

---

## 📋 Table of Contents

1. [High-Level Flow](#high-level-flow)
2. [Detailed Step-by-Step](#detailed-step-by-step)
3. [Key Components](#key-components)
4. [State Management](#state-management)
5. [Critical Findings](#critical-findings)

---

## 🎯 High-Level Flow

```
User Input
    ↓
sendMessage()
    ↓
[Auto Topic Creation?] ──Yes──→ createTopic() + switchTopic()
    ↓ No
internal_createMessage() (user message)
    ↓
internal_coreProcessMessage()
    ↓
[RAG Enabled?] ──Yes──→ internal_retrieveChunks()
    ↓ No
internal_createMessage() (assistant placeholder)
    ↓
internal_fetchAIChatMessage()
    ↓
[Tool Calling?] ──Yes──→ triggerToolCalls() (loop)
    ↓ No
summaryTopicTitle() (if new topic)
    ↓
Done
```

---

## 🔍 Detailed Step-by-Step

### **1. sendMessage() - Entry Point**

**Location**: `generateAIChat.ts.backup:161-377`

**Key Variables**:
```typescript
const {
  internal_coreProcessMessage,
  activeTopicId,        // Current topic ID (null for inbox)
  activeId,             // Session ID (e.g., "inbox")
  activeThreadId,       // Thread ID (for branching)
  sendMessageInServer,
} = get();
```

**Flow**:

#### 1.1. **Validation**
```typescript
if (!activeId) return;
if (!message && !hasFile) return;
```

#### 1.2. **Server Mode Check**
```typescript
if (isServerMode)
  return sendMessageInServer({ message, files, onlyAddUserMessage, isWelcomeQuestion });
```

#### 1.3. **Create User Message Object**
```typescript
const newMessage: CreateMessageParams = {
  content: message,
  files: fileIdList,
  role: 'user',
  sessionId: activeId,      // ← "inbox" (slug, not ID!)
  topicId: activeTopicId,   // ← null initially
  threadId: activeThreadId, // ← undefined initially
};
```

#### 1.4. **Auto Topic Creation** (Lines 228-269)

**Condition**:
```typescript
if (!onlyAddUserMessage && !activeTopicId && agentConfig.enableAutoCreateTopic) {
  const featureLength = chats.length + 2; // +2 for user + assistant
  
  if (featureLength >= agentConfig.autoCreateTopicThreshold) {
    // CREATE TOPIC!
  }
}
```

**Process**:
1. Create **temp message** for optimistic UI update
2. Call `createTopic()` → returns `topicId`
3. Update `newMessage.topicId = topicId`
4. **Copy messages** to new topic in `messagesMap`
5. Set topic loading state

#### 1.5. **Create User Message in DB**
```typescript
const id = await get().internal_createMessage(newMessage, {
  tempMessageId,
  skipRefresh: !onlyAddUserMessage && newMessage.fileList?.length === 0,
});
```

#### 1.6. **Switch to New Topic** (Lines 287-304)

**If topic was created**:
```typescript
if (!!newTopicId) {
  await get().switchTopic(newTopicId, true);
  await get().internal_fetchMessages(activeId, newTopicId);
  
  // Delete previous messages from inbox
  const newMaps = { ...get().messagesMap, [messageMapKey(activeId, null)]: [] };
  set({ messagesMap: newMaps }, false, 'internal_copyMessages');
}
```

**⚠️ CRITICAL**: `switchTopic()` changes `activeTopicId` in store!

#### 1.7. **Early Return for User-Only Messages**
```typescript
if (onlyAddUserMessage) {
  set({ isCreatingMessage: false }, false, 'creatingMessage/start');
  return;
}
```

#### 1.8. **Prepare for AI Response**
```typescript
const messages = chatSelectors.activeBaseChats(get());
const currentActiveThreadId = get().activeThreadId; // Re-read after topic switch!

await internal_coreProcessMessage(messages, id, {
  isWelcomeQuestion,
  ragQuery: get().internal_shouldUseRAG() ? message : undefined,
  threadId: currentActiveThreadId,
});
```

#### 1.9. **Post-Processing**
```typescript
// Generate topic title
await summaryTitle();

// Add files to agent (server mode only)
await addFilesToAgent();
```

---

### **2. internal_coreProcessMessage() - Core Logic**

**Location**: `generateAIChat.ts.backup:389-850`

**Parameters**:
```typescript
async (
  originalMessages: UIChatMessage[],
  userMessageId: string,
  params?: {
    traceId?: string;
    isWelcomeQuestion?: boolean;
    inSearchWorkflow?: boolean;
    ragQuery?: string;
    threadId?: string;
    inPortalThread?: boolean;
    groupId?: string;
    agentId?: string;
    agentConfig?: any;
  }
)
```

**Flow**:

#### 2.1. **Get Agent Config**
```typescript
const { model, provider, chatConfig } = agentSelectors.currentAgentConfig(agentStoreState);
```

#### 2.2. **RAG Flow** (if `params.ragQuery` exists)

```typescript
if (params?.ragQuery) {
  // 1. Retrieve chunks
  const { chunks, queryId, rewriteQuery } = await get().internal_retrieveChunks(
    userMessageId,
    params.ragQuery,
    messages.map((m) => m.content).slice(0, messages.length - 1),
  );
  
  // 2. Build knowledge base context
  const knowledgeBaseQAContext = knowledgeBaseQAPrompts({
    chunks,
    userQuery: lastMsg.content,
    rewriteQuery,
    knowledge: agentSelectors.currentEnabledKnowledge(agentStoreState),
  });
  
  // 3. Append context to user message
  messages.push({
    ...lastMsg,
    content: (lastMsg.content + '\n\n' + knowledgeBaseQAContext).trim(),
  });
  
  fileChunks = chunks.map((c) => ({ id: c.id, similarity: c.similarity }));
}
```

#### 2.3. **Create Assistant Placeholder**
```typescript
const assistantMessage: CreateMessageParams = {
  role: 'assistant',
  content: LOADING_FLAT,
  fromModel: model,
  fromProvider: provider,
  parentId: userMessageId,
  sessionId: get().activeId,
  topicId: activeTopicId,
  threadId: params?.threadId, // ← Important for thread branching!
  fileChunks,
  ragQueryId,
};

const assistantId = await get().internal_createMessage(assistantMessage);
```

#### 2.4. **Check Tool/Search Capabilities**
```typescript
const isModelSupportToolUse = aiModelSelectors.isModelSupportToolUse(model, provider);
const isAgentEnableSearch = agentChatConfigSelectors.isAgentEnableSearch(agentStoreState);
const useModelSearch = ((isProviderHasBuiltinSearch || isModelHasBuiltinSearch) && useModelBuiltinSearch) || isModelBuiltinSearchInternal;
```

#### 2.5. **Fetch AI Response**
```typescript
const { content, isFunctionCall, traceId } = await internal_fetchAIChatMessage({
  messages,
  messageId: assistantId,
  params,
  model,
  provider,
});
```

**This calls**: `chatService.getChatCompletion()` or `chatService.getChatCompletionStream()`

#### 2.6. **Tool Calling Loop** (if applicable)
```typescript
if (isFunctionCall) {
  await triggerToolCalls(assistantId, traceId, params?.threadId, params?.inPortalThread);
}
```

---

### **3. internal_fetchAIChatMessage() - API Call**

**Location**: `generateAIChat.ts.backup:852-1057`

**Key Steps**:

#### 3.1. **Get Enabled Tools**
```typescript
const enabledToolIds = get().enabledTools;
```

#### 3.2. **Context Engineering** (Pre-processing)
```typescript
const oaiMessages = await contextEngineeringBackend({
  enableHistoryCount: agentChatConfigSelectors.enableHistoryCount(agentStoreState),
  historyCount: agentChatConfigSelectors.historyCount(agentStoreState) + 2,
  historySummary: options?.historySummary,
  inputTemplate: chatConfig.inputTemplate,
  isWelcomeQuestion: options?.isWelcomeQuestion,
  messages,
  model: payload.model,
  provider: payload.provider!,
  sessionId: options?.trace?.sessionId,
  systemRole: getNullableString(agentConfig.systemRole as any),
  tools: enabledToolIds,
});
```

**⚠️ CRITICAL**: This is where message history is processed!

#### 3.3. **Call Chat Service**
```typescript
if (enableStreaming) {
  await chatService.getChatCompletionStream(payload, {
    onMessageHandle: (text) => {
      get().internal_updateMessageContent(messageId, text);
    },
    onFinish: async (content, { traceId, observationId, toolCalls }) => {
      // Update message
      // Handle tool calls
    },
  });
} else {
  const data = await chatService.getChatCompletion(payload);
  // Update message
}
```

---

## 🔑 Key Components

### **State Variables**

| Variable | Type | Description | Example |
|----------|------|-------------|---------|
| `activeId` | `string` | Session ID (can be slug!) | `"inbox"` |
| `activeTopicId` | `string \| null` | Current topic ID | `null` or `"topic-123"` |
| `activeThreadId` | `string \| undefined` | Current thread ID | `undefined` or `"thread-456"` |
| `messagesMap` | `Record<string, UIChatMessage[]>` | Messages by key | `{ "inbox:null": [...], "inbox:topic-123": [...] }` |
| `enabledTools` | `string[]` | Active tool IDs | `["web-search", "calculator"]` |

### **Message Map Key**

```typescript
function messageMapKey(sessionId: string, topicId: string | null): string {
  return `${sessionId}:${topicId}`;
}
```

**Examples**:
- Inbox (no topic): `"inbox:null"`
- With topic: `"inbox:topic-123"`
- With thread: Messages still use topic key, but have `threadId` field

### **Services**

1. **chatService** (`@/services/chat`)
   - `getChatCompletion()` - Non-streaming
   - `getChatCompletionStream()` - Streaming
   - Handles model/provider routing

2. **messageService** (`@/services/message`)
   - `createMessage()` - Save to DB
   - `updateMessage()` - Update in DB
   - `removeMessage()` - Delete from DB

3. **contextEngineeringBackend** (`@/services/chat/contextEngineeringBackend`)
   - Pre-processes messages
   - Applies history count
   - Applies input templates
   - Handles placeholders

---

## ⚠️ Critical Findings

### **1. Session ID vs Slug Confusion**

**Problem**: Frontend uses `activeId = "inbox"` (slug), but DB has:
```sql
CREATE TABLE sessions (
  id TEXT PRIMARY KEY,     -- "RfXAMJTxkC1bEbCHr1e5X"
  slug TEXT,               -- "inbox"
  user_id TEXT,
  UNIQUE(slug, user_id)
);
```

**Impact**: Backend must handle **both ID and slug** in queries!

### **2. Topic Creation Timing**

**Flow**:
1. Check if topic needed (line 228)
2. Create **temp message** (line 242)
3. Call `createTopic()` (line 245)
4. Update message with `topicId` (line 254)
5. **Copy messages** to new topic map (line 258-263)
6. Save user message to DB (line 273)
7. **Switch topic** (line 293)
8. **Fetch messages** (line 298)
9. **Delete old messages** (line 302)

**⚠️ CRITICAL**: Topic switch happens **AFTER** user message is created!

### **3. Thread ID Handling**

**Thread ID is passed through**:
```typescript
sendMessage({ ... })
  → internal_coreProcessMessage(messages, id, { threadId: currentActiveThreadId })
    → internal_createMessage({ threadId: params?.threadId })
      → messageService.createMessage({ threadId })
```

**⚠️ CRITICAL**: `threadId` must be preserved through entire chain!

### **4. Context Engineering**

**Original flow**:
```
Messages → contextEngineeringBackend() → Processed Messages → chatService
```

**What it does**:
- Truncates history to `historyCount`
- Applies input templates
- Replaces placeholders
- Formats for OpenAI API

**⚠️ CRITICAL**: This is **essential** for proper message handling!

### **5. Tool Calling Loop**

**Flow**:
```
internal_fetchAIChatMessage()
  → Returns { isFunctionCall: true }
    → triggerToolCalls()
      → Execute tools
        → internal_resendMessage() (recursive!)
```

**⚠️ CRITICAL**: Tool calling is **recursive** - can loop multiple times!

### **6. Message Map Structure**

**Key insight**: Messages are stored in `messagesMap` by `sessionId:topicId`:

```typescript
{
  "inbox:null": [msg1, msg2, msg3],           // Inbox without topic
  "inbox:topic-123": [msg4, msg5, msg6],      // Topic 123
  "inbox:topic-456": [msg7, msg8, msg9],      // Topic 456
}
```

**Thread messages**: Same key, but with `threadId` field in message object!

---

## 🎯 Migration Implications

### **What Backend MUST Handle**

1. ✅ **Session by ID or Slug** - Already fixed
2. ❌ **Topic auto-creation** - Backend does this, but timing is wrong
3. ❌ **Thread ID preservation** - Not tested
4. ❌ **Context engineering** - Removed from frontend, must be in backend
5. ❌ **Tool calling loop** - Backend must handle recursion
6. ❌ **Message map updates** - Frontend still needs to update its map
7. ❌ **Topic title generation** - Backend does this, but needs to return topic_id

### **What Frontend MUST Still Do**

1. ✅ **State management** (`activeId`, `activeTopicId`, `activeThreadId`)
2. ✅ **Optimistic UI updates** (temp messages)
3. ✅ **Message map management** (`messagesMap`)
4. ✅ **Topic switching** (`switchTopic()`)
5. ✅ **Message refresh** (`refreshMessages()`)

### **Critical Missing Pieces**

1. **Context Engineering**: Removed from frontend, but backend doesn't do it!
2. **Tool Calling Loop**: Backend does one call, but doesn't loop
3. **Topic Switch Coordination**: Backend creates topic, but frontend must switch
4. **Message Map Updates**: Frontend needs to know when to update map

---

## 📝 Conclusion

**The migration was incomplete because**:

1. ❌ Didn't understand session ID vs slug distinction
2. ❌ Didn't understand topic creation timing
3. ❌ Didn't understand context engineering importance
4. ❌ Didn't understand tool calling recursion
5. ❌ Didn't understand message map structure
6. ❌ Removed too much frontend logic that was still needed

**Next Steps**:

1. **Rollback** to original `generateAIChat.ts`
2. **Fix backend** to properly handle ID/slug
3. **Keep context engineering** in frontend OR implement in backend
4. **Migrate incrementally** - one function at a time
5. **Test each step** before moving to next

---

**Status**: ✅ Analysis Complete - Ready for proper migration strategy

