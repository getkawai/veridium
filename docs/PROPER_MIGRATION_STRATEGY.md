# Proper Migration Strategy

> **Purpose**: Step-by-step migration from frontend to backend based on complete flow understanding  
> **Date**: November 2025  
> **Status**: 🎯 Ready to Execute

---

## 📋 Table of Contents

1. [Current State](#current-state)
2. [Migration Principles](#migration-principles)
3. [Phase-by-Phase Plan](#phase-by-phase-plan)
4. [Testing Strategy](#testing-strategy)
5. [Rollback Plan](#rollback-plan)

---

## 🎯 Current State

### ✅ **What's Working**

1. **Original Frontend Flow** - Fully restored
   - `generateAIChat.ts` - Complete original implementation
   - `contextEngineeringBackend.ts` - Message preprocessing
   - `chatService` - Model/provider routing
   - All state management intact

2. **Backend Infrastructure** - Partially working
   - ✅ Eino agent integration
   - ✅ Knowledge base (RAG)
   - ✅ Tools engine bridge
   - ✅ Context engine bridge
   - ✅ Session persistence (ID + slug support)
   - ✅ Topic auto-creation
   - ✅ Thread management service

### ⚠️ **What's Not Working**

1. **Backend-Frontend Integration**
   - Backend creates sessions/topics correctly
   - But frontend doesn't use backend yet
   - Two parallel systems running

2. **Missing Backend Features**
   - No streaming support (Wails v3 limitation)
   - Tool calling not recursive
   - No message map updates

---

## 🎓 Migration Principles

### **1. Incremental Migration**
- Migrate **one function at a time**
- Keep both systems running in parallel
- Use feature flags to switch between old/new

### **2. Test-Driven**
- Test each function before moving to next
- Compare old vs new output
- Ensure UI behavior identical

### **3. Backward Compatible**
- Don't break existing functionality
- Keep fallbacks to old system
- Allow easy rollback

### **4. Preserve State**
- Frontend state management stays
- Backend only handles business logic
- Clear separation of concerns

---

## 📅 Phase-by-Phase Plan

### **Phase 1: Backend Session Management** ✅ (DONE)

**Goal**: Backend handles session creation and loading

**Changes**:
- ✅ `AgentChatService.getOrCreateSession()` - Handle ID/slug
- ✅ `AgentChatService.loadSessionFromDB()` - Load by ID or slug
- ✅ Fix race condition in session creation

**Testing**:
- ✅ Create new session
- ✅ Load existing session by ID
- ✅ Load existing session by slug
- ✅ Handle concurrent requests

---

### **Phase 2: Message Creation** (NEXT)

**Goal**: Backend handles message persistence

**Current Frontend**:
```typescript
// generateAIChat.ts:273
const id = await get().internal_createMessage(newMessage, {
  tempMessageId,
  skipRefresh: !onlyAddUserMessage && newMessage.fileList?.length === 0,
});
```

**Backend API**:
```go
// Already exists in AgentChatService
func (s *AgentChatService) Chat(ctx context.Context, req models.ChatRequest) (models.ChatResponse, error)
```

**Migration Steps**:

1. **Add Feature Flag**
   ```typescript
   // frontend/src/config/features.ts
   export const USE_BACKEND_MESSAGE_CREATION = false; // Start with false
   ```

2. **Create Wrapper Function**
   ```typescript
   // frontend/src/services/backendAgentChat.ts
   export async function createMessage(params: CreateMessageParams): Promise<string> {
     if (!USE_BACKEND_MESSAGE_CREATION) {
       // Use old method
       return messageService.createMessage(params);
     }
     
     // Use backend
     const response = await AgentChatService.Chat({
       session_id: params.sessionId,
       user_id: FALLBACK_USER_ID,
       topic_id: params.topicId,
       thread_id: params.threadId,
       message: params.content,
       tools: [], // No tools for user message
       temperature: 0.7,
       max_tokens: 2000,
     });
     
     return response.user_message_id;
   }
   ```

3. **Update Frontend**
   ```typescript
   // generateAIChat.ts
   import { createMessage } from '@/services/backendAgentChat';
   
   // Replace:
   const id = await get().internal_createMessage(newMessage, options);
   
   // With:
   const id = await createMessage(newMessage);
   ```

4. **Test**
   - Create user message
   - Verify ID returned
   - Check DB persistence
   - Verify UI update

5. **Enable Feature Flag**
   ```typescript
   export const USE_BACKEND_MESSAGE_CREATION = true;
   ```

**Testing Checklist**:
- [ ] User message created in DB
- [ ] Message ID returned correctly
- [ ] UI shows message immediately
- [ ] Topic ID preserved
- [ ] Thread ID preserved
- [ ] Files attached correctly

---

### **Phase 3: Topic Auto-Creation** (AFTER PHASE 2)

**Goal**: Backend handles topic creation and title generation

**Current Frontend**:
```typescript
// generateAIChat.ts:228-269
if (!onlyAddUserMessage && !activeTopicId && agentConfig.enableAutoCreateTopic) {
  const featureLength = chats.length + 2;
  if (featureLength >= agentConfig.autoCreateTopicThreshold) {
    const topicId = await get().createTopic();
    newMessage.topicId = topicId;
    // ... copy messages, switch topic ...
  }
}
```

**Backend API**:
```go
// Already implemented in AgentChatService.Chat()
// Auto-creates topic if needed and returns topic_id
```

**Migration Steps**:

1. **Add Feature Flag**
   ```typescript
   export const USE_BACKEND_TOPIC_CREATION = false;
   ```

2. **Update Backend to Return More Info**
   ```go
   // models/chat_models.go
   type ChatResponse struct {
       // ... existing fields ...
       TopicCreated bool   `json:"topic_created"`
       TopicTitle   string `json:"topic_title"`
   }
   ```

3. **Update Frontend Logic**
   ```typescript
   // generateAIChat.ts
   if (USE_BACKEND_TOPIC_CREATION) {
     // Backend handles topic creation
     const response = await backendAgentChat.sendMessage({
       session_id: activeId,
       user_id: FALLBACK_USER_ID,
       message: message,
       // ... other params ...
     });
     
     if (response.topic_created && response.topic_id) {
       // Update frontend state
       set({ activeTopicId: response.topic_id });
       
       // Copy messages to new topic map
       const mapKey = chatSelectors.currentChatKey(get());
       const newMaps = {
         ...get().messagesMap,
         [messageMapKey(activeId, response.topic_id)]: get().messagesMap[mapKey],
       };
       set({ messagesMap: newMaps });
       
       // Switch to new topic
       await get().switchTopic(response.topic_id, true);
     }
   } else {
     // Use old frontend logic
     // ... existing code ...
   }
   ```

4. **Test**
   - First message in session
   - Topic auto-created
   - Title generated
   - Frontend switches to topic
   - Messages copied correctly

5. **Enable Feature Flag**

**Testing Checklist**:
- [ ] Topic created at threshold
- [ ] Topic title generated by LLM
- [ ] Frontend switches to topic
- [ ] Messages copied to new topic map
- [ ] Old messages deleted from inbox
- [ ] UI shows topic in sidebar

---

### **Phase 4: AI Response Generation** (COMPLEX)

**Goal**: Backend handles AI response generation

**Current Frontend**:
```typescript
// generateAIChat.ts:340-344
await internal_coreProcessMessage(messages, id, {
  isWelcomeQuestion,
  ragQuery: get().internal_shouldUseRAG() ? message : undefined,
  threadId: currentActiveThreadId,
});
```

**Challenges**:
1. **Context Engineering** - Currently in frontend
2. **Tool Calling Loop** - Recursive, complex
3. **Streaming** - Wails v3 limitation
4. **Message Updates** - Real-time UI updates

**Migration Steps**:

1. **Keep Context Engineering in Frontend** (for now)
   - Too complex to migrate immediately
   - Backend can add later
   - Frontend preprocessing still works

2. **Add Feature Flag**
   ```typescript
   export const USE_BACKEND_AI_GENERATION = false;
   ```

3. **Create Hybrid Approach**
   ```typescript
   // generateAIChat.ts
   async function generateAIResponse(messages, userMessageId, params) {
     if (USE_BACKEND_AI_GENERATION) {
       // Use backend
       const response = await AgentChatService.Chat({
         session_id: get().activeId,
         user_id: FALLBACK_USER_ID,
         topic_id: get().activeTopicId,
         thread_id: params?.threadId,
         message: messages[messages.length - 1].content,
         tools: get().enabledTools,
         temperature: agentConfig.temperature,
         max_tokens: agentConfig.maxTokens,
       });
       
       // Update UI with response
       await get().internal_updateMessageContent(
         response.assistant_message_id,
         response.content
       );
       
       return response;
     } else {
       // Use old frontend logic
       return await internal_coreProcessMessage(messages, userMessageId, params);
     }
   }
   ```

4. **Handle Tool Calling**
   ```typescript
   // Backend returns tool_calls in response
   if (response.tool_calls && response.tool_calls.length > 0) {
     // Trigger tool execution
     await get().triggerToolCalls(
       response.assistant_message_id,
       response.trace_id,
       params?.threadId
     );
   }
   ```

5. **Test**
   - Simple message (no tools)
   - Message with RAG
   - Message with tools
   - Tool calling loop
   - Error handling

**Testing Checklist**:
- [ ] AI response generated
- [ ] Content updated in UI
- [ ] Tool calls executed
- [ ] Recursive tool calling works
- [ ] RAG context included
- [ ] Error messages shown

---

### **Phase 5: Streaming Support** (FUTURE)

**Goal**: Implement real-time streaming

**Challenges**:
- Wails v3 event API unclear
- Need to research proper implementation
- May require custom protocol

**Options**:

1. **Server-Sent Events (SSE)**
   - Use HTTP endpoint for streaming
   - Bypass Wails event system
   - More standard approach

2. **WebSocket**
   - Full duplex communication
   - More complex setup
   - Better for real-time

3. **Polling**
   - Simplest approach
   - Higher latency
   - Fallback option

**Decision**: Defer to Phase 5 after core functionality working

---

## 🧪 Testing Strategy

### **Unit Tests**

1. **Backend Tests**
   ```bash
   cd internal/services
   go test -v ./...
   ```

2. **Frontend Tests**
   ```bash
   cd frontend
   npm test
   ```

### **Integration Tests**

1. **Session Management**
   - Create session by slug
   - Load session by ID
   - Load session by slug
   - Concurrent session creation

2. **Message Flow**
   - User message → DB
   - AI response → DB
   - Message with files
   - Message with tools

3. **Topic Management**
   - Auto-create topic
   - Generate title
   - Switch topic
   - List topics

4. **Thread Management**
   - Create thread
   - Branch conversation
   - List threads
   - Thread messages

### **E2E Tests**

1. **Basic Chat Flow**
   ```
   User: "Hello"
   → Session created
   → User message saved
   → AI response generated
   → UI updated
   ```

2. **Topic Creation Flow**
   ```
   User: "Message 1"
   User: "Message 2"
   User: "Message 3" (threshold reached)
   → Topic created
   → Title generated
   → Frontend switches
   → Messages copied
   ```

3. **Tool Calling Flow**
   ```
   User: "Search for X"
   → Tool call detected
   → Tool executed
   → Result returned
   → AI processes result
   → Final response
   ```

4. **RAG Flow**
   ```
   User: "Question about docs"
   → RAG enabled
   → Chunks retrieved
   → Context added
   → AI response with context
   ```

---

## 🔄 Rollback Plan

### **If Phase 2 Fails**

1. **Disable Feature Flag**
   ```typescript
   export const USE_BACKEND_MESSAGE_CREATION = false;
   ```

2. **Verify Old System Works**
   - Test message creation
   - Check DB persistence
   - Verify UI updates

### **If Phase 3 Fails**

1. **Disable Feature Flag**
   ```typescript
   export const USE_BACKEND_TOPIC_CREATION = false;
   ```

2. **Verify Old System Works**
   - Test topic creation
   - Check title generation
   - Verify topic switch

### **If Phase 4 Fails**

1. **Disable Feature Flag**
   ```typescript
   export const USE_BACKEND_AI_GENERATION = false;
   ```

2. **Verify Old System Works**
   - Test AI response
   - Check tool calling
   - Verify streaming

### **Complete Rollback**

If all else fails:

```bash
# Restore original files
cd frontend/src/store/chat/slices/aiChat/actions
mv generateAIChat.broken.ts generateAIChat.ts

# Restore context engineering
git restore frontend/src/services/chat/contextEngineeringBackend.ts
git restore frontend/src/services/contextEngineBackend.ts

# Restore chat service
git restore frontend/src/services/chat/index.ts
```

---

## 📊 Progress Tracking

### **Phase 1: Backend Session Management** ✅
- [x] Fix ID/slug handling
- [x] Fix race condition
- [x] Test session creation
- [x] Test session loading

### **Phase 2: Message Creation** 🔄 (IN PROGRESS)
- [ ] Add feature flag
- [ ] Create wrapper function
- [ ] Update frontend
- [ ] Test message creation
- [ ] Enable feature flag

### **Phase 3: Topic Auto-Creation** ⏳ (PENDING)
- [ ] Add feature flag
- [ ] Update backend response
- [ ] Update frontend logic
- [ ] Test topic creation
- [ ] Enable feature flag

### **Phase 4: AI Response Generation** ⏳ (PENDING)
- [ ] Add feature flag
- [ ] Create hybrid approach
- [ ] Handle tool calling
- [ ] Test AI generation
- [ ] Enable feature flag

### **Phase 5: Streaming Support** ⏳ (FUTURE)
- [ ] Research Wails v3 streaming
- [ ] Choose implementation approach
- [ ] Implement streaming
- [ ] Test streaming
- [ ] Enable streaming

---

## 🎯 Success Criteria

### **Phase 2 Success**
- ✅ User messages created via backend
- ✅ Message IDs returned correctly
- ✅ DB persistence verified
- ✅ UI updates immediately
- ✅ No regressions in existing features

### **Phase 3 Success**
- ✅ Topics auto-created at threshold
- ✅ Titles generated by LLM
- ✅ Frontend switches correctly
- ✅ Messages copied to new topic
- ✅ UI shows topic in sidebar

### **Phase 4 Success**
- ✅ AI responses generated via backend
- ✅ Tool calling works recursively
- ✅ RAG context included
- ✅ Error handling robust
- ✅ UI updates correctly

### **Phase 5 Success**
- ✅ Streaming works in real-time
- ✅ No lag or buffering
- ✅ Error handling during stream
- ✅ Cancellation works
- ✅ UI shows streaming text

---

## 📝 Next Steps

1. **Start Phase 2** - Message Creation
   - Add feature flag
   - Create wrapper function
   - Test thoroughly

2. **Document Progress**
   - Update this document
   - Track issues encountered
   - Document solutions

3. **Communicate**
   - Keep user informed
   - Ask for feedback
   - Adjust plan as needed

---

**Status**: ✅ Ready to Execute Phase 2

