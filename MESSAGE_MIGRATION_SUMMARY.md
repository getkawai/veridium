# Message Model Migration Summary

## ✅ Migration Complete

File: `frontend/src/database/models/message.wails.ts` (807 lines)

### 📊 Status Overview

| Category | Status | Notes |
|----------|--------|-------|
| **Compilation** | ✅ Pass | No TypeScript errors |
| **Basic CRUD** | ✅ Complete | All methods migrated |
| **Complex Queries** | ⚠️ N+1 Problem | Requires optimization |
| **Transactions** | ❌ Not Supported | Data consistency risk |
| **Type Safety** | ⚠️ Partial | Many `as any` casts needed |

---

## 🔧 New SQL Queries Added

### Count & Statistics
```sql
-- name: CountMessages :one
SELECT COUNT(*) FROM messages WHERE user_id = ?;

-- name: CountMessagesByDateRange :one
SELECT COUNT(*) FROM messages
WHERE user_id = ? AND created_at >= ? AND created_at <= ?;

-- name: CountMessageWords :one
SELECT SUM(LENGTH(content)) as total_length
FROM messages WHERE user_id = ?;

-- name: CountMessageWordsByDateRange :one
SELECT SUM(LENGTH(content)) as total_length
FROM messages
WHERE user_id = ? AND created_at >= ? AND created_at <= ?;
```

### Batch Operations
```sql
-- name: BatchDeleteMessages :exec
DELETE FROM messages
WHERE user_id = ? AND id IN (sqlc.slice('ids'));

-- name: DeleteAllMessages :exec
DELETE FROM messages WHERE user_id = ?;
```

### Ranking & Search
```sql
-- name: RankModels :many
SELECT model as id, COUNT(*) as count
FROM messages
WHERE user_id = ? AND model IS NOT NULL AND model != ''
GROUP BY model
ORDER BY count DESC, model ASC
LIMIT ?;

-- name: SearchMessagesByKeyword :many
SELECT * FROM messages
WHERE user_id = ? AND content LIKE ?
ORDER BY created_at DESC
LIMIT ?;
```

### Upsert Operations
```sql
-- name: UpsertMessageTTS :one
INSERT INTO message_tts (id, content_md5, file_id, voice, client_id, user_id)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    content_md5 = excluded.content_md5,
    file_id = excluded.file_id,
    voice = excluded.voice
RETURNING *;

-- name: UpsertMessageTranslate :one
INSERT INTO message_translates (id, content, "from", "to", client_id, user_id)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    content = excluded.content,
    "from" = excluded."from",
    "to" = excluded."to"
RETURNING *;
```

### Delete Operations
```sql
-- name: DeleteMessageTTS :exec
DELETE FROM message_tts WHERE id = ? AND user_id = ?;

-- name: DeleteMessageTranslate :exec
DELETE FROM message_translates WHERE id = ? AND user_id = ?;

-- name: DeleteMessageQuery :exec
DELETE FROM message_queries WHERE id = ? AND user_id = ?;
```

### Query Chunks
```sql
-- name: LinkMessageQueryToChunk :exec
INSERT INTO message_query_chunks (message_id, query_id, chunk_id, similarity, user_id)
VALUES (?, ?, ?, ?, ?);

-- name: GetMessageQueryChunks :many
SELECT
    mqc.message_id,
    mqc.similarity,
    c.id,
    c.text,
    f.id as file_id,
    f.name as filename,
    f.file_type,
    f.url as file_url
FROM message_query_chunks mqc
LEFT JOIN chunks c ON mqc.chunk_id = c.id
LEFT JOIN file_chunks fc ON c.id = fc.chunk_id
LEFT JOIN files f ON fc.file_id = f.id
WHERE mqc.message_id IN (sqlc.slice('messageIds')) AND mqc.user_id = ?;
```

---

## ⚠️ Critical Issues

### 1. N+1 Query Problem in `query()` Method

**Drizzle (Efficient - 4 queries total)**:
```typescript
// 1 query: Get messages with plugins, translates, TTS (JOINs)
const result = await this.db.select({ /* all fields */ })
  .from(messages)
  .leftJoin(messagePlugins, ...)
  .leftJoin(messageTranslates, ...)
  .leftJoin(messageTTS, ...);

// 1 query: Get all files for all messages (IN clause)
const files = await this.db.select()
  .from(messagesFiles)
  .where(inArray(messagesFiles.messageId, messageIds));

// 1 query: Get all chunks (IN clause)
// 1 query: Get all queries (IN clause)
```

**Wails (Inefficient - N*4 + 3 queries)**:
```typescript
// 1 query: Get messages (no JOINs)
const messages = await DB.ListMessages(...);

// N queries: Get files for each message
await Promise.all(messages.map(m => DB.GetMessageFiles(m.id)));

// N queries: Get plugins for each message
await Promise.all(messages.map(m => DB.GetMessagePlugin(m.id)));

// N queries: Get translates for each message
await Promise.all(messages.map(m => DB.GetMessageTranslate(m.id)));

// N queries: Get TTS for each message
await Promise.all(messages.map(m => DB.GetMessageTTS(m.id)));
```

**Performance Impact**:
- For 100 messages: **401 queries** vs **4 queries**
- For 1000 messages: **4001 queries** vs **4 queries**

**Recommended Solution**: Create comprehensive SQL query with all JOINs:
```sql
-- name: QueryMessagesWithRelations :many
SELECT 
    m.*,
    mp.tool_call_id as plugin_tool_call_id,
    mp.type as plugin_type,
    mp.api_name as plugin_api_name,
    mp.arguments as plugin_arguments,
    mp.identifier as plugin_identifier,
    mp.state as plugin_state,
    mp.error as plugin_error,
    mt.content as translate_content,
    mt."from" as translate_from,
    mt."to" as translate_to,
    mtt.content_md5 as tts_content_md5,
    mtt.file_id as tts_file_id,
    mtt.voice as tts_voice
FROM messages m
LEFT JOIN message_plugins mp ON m.id = mp.id
LEFT JOIN message_translates mt ON m.id = mt.id
LEFT JOIN message_tts mtt ON m.id = mtt.id
WHERE m.user_id = ? 
  AND (? = '' OR m.session_id = ?)
  AND (? = '' OR m.topic_id = ?)
  AND (? = '' OR m.group_id = ?)
ORDER BY m.created_at ASC
LIMIT ? OFFSET ?;
```

### 2. No Transaction Support

**Drizzle**:
```typescript
return this.db.transaction(async (trx) => {
  const [item] = await trx.insert(messages).values(...).returning();
  
  if (message.role === 'tool') {
    await trx.insert(messagePlugins).values(...);
  }
  
  if (files) {
    await trx.insert(messagesFiles).values(...);
  }
  
  return item; // All or nothing!
});
```

**Wails**:
```typescript
// No transaction - each operation is separate!
const item = await DB.CreateMessage(...);

if (message.role === 'tool') {
  await DB.CreateMessagePlugin(...); // If this fails, message is already created!
}

if (files) {
  await Promise.all(...); // Partial failure possible!
}

return item;
```

**Impact**: Data inconsistency if any operation fails mid-way.

**Recommended Solution**: Implement transaction support in Go backend:
```go
func (s *Service) CreateMessageWithTransaction(ctx context.Context, params CreateMessageParams) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := s.queries.WithTx(tx)
    
    // All operations in transaction
    message, err := qtx.CreateMessage(ctx, ...)
    if err != nil {
        return err
    }
    
    if params.Role == "tool" {
        err = qtx.CreateMessagePlugin(ctx, ...)
        if err != nil {
            return err
        }
    }
    
    return tx.Commit()
}
```

### 3. Heatmap Query Inefficiency

**Current Implementation**:
```typescript
// Fetch ALL messages and group in memory
const messages = await DB.ListMessages({
  userId: this.userId,
  sessionId: toNullString(''),
  limit: 100000, // ⚠️ Fetching 100k rows!
  offset: 0,
});

// Group in memory
const dateCountMap = new Map<string, number>();
for (const message of messages) {
  const date = dayjs(message.createdAt).format('YYYY-MM-DD');
  dateCountMap.set(date, (dateCountMap.get(date) || 0) + 1);
}
```

**Recommended Solution**: Use SQL GROUP BY:
```sql
-- name: GetMessageHeatmaps :many
SELECT 
    DATE(created_at / 1000, 'unixepoch') as date,
    COUNT(*) as count
FROM messages
WHERE user_id = ? 
  AND created_at >= ? 
  AND created_at <= ?
GROUP BY DATE(created_at / 1000, 'unixepoch')
ORDER BY date ASC;
```

---

## 📝 Type Conversion Challenges

### NullString Handling
```typescript
// Wails requires explicit NullString construction
sessionId: toNullString(sessionId as any)

// vs Drizzle's automatic handling
sessionId: sessionId
```

### NullInt64 for Similarity
```typescript
// Wails requires manual struct construction
similarity: { Int64: chunk.similarity || 0, Valid: !!chunk.similarity } as any

// vs Drizzle's automatic handling
similarity: chunk.similarity
```

### JSON Fields
```typescript
// Wails requires stringify/parse
metadata: toNullJSON(metadata)
const parsed = parseNullableJSON(item.metadata as any)

// vs Drizzle's automatic handling
metadata: metadata
const parsed = item.metadata
```

### Type Casting
```typescript
// Many `as any` and `as unknown as` casts needed
return messages as unknown as DBMessageItem[];
content: toNullString((normalizedMessage as any).reasoning as any)
```

---

## 🎯 Methods Migrated

### Query Methods ✅
- `query()` - ⚠️ N+1 problem
- `findById()` - ✅ Optimized
- `findMessageQueriesById()` - ✅ Optimized
- `queryAll()` - ✅ Single query
- `queryBySessionId()` - ✅ Single query
- `queryByKeyword()` - ✅ Single query (new `SearchMessagesByKeyword`)

### Count Methods ✅
- `count()` - ✅ Efficient (new `CountMessages`, `CountMessagesByDateRange`)
- `countWords()` - ✅ Efficient (new `CountMessageWords`, `CountMessageWordsByDateRange`)
- `rankModels()` - ✅ Single GROUP BY query (new `RankModels`)
- `getHeatmaps()` - ⚠️ In-memory grouping (should use SQL GROUP BY)
- `hasMoreThanN()` - ✅ Uses LIMIT

### Create Methods ✅
- `create()` - ❌ No transaction support
- `createNewMessage()` - ❌ No transaction support
- `batchCreate()` - ⚠️ Sequential creates (no batch insert)
- `createMessageQuery()` - ✅ Single query

### Update Methods ✅
- `update()` - ❌ No transaction support
- `updateMetadata()` - ✅ Single query
- `updatePluginState()` - ✅ Single query
- `updateMessagePlugin()` - ✅ Single query
- `updateTranslate()` - ✅ Upsert (new `UpsertMessageTranslate`)
- `updateTTS()` - ✅ Upsert (new `UpsertMessageTTS`)
- `updateMessageRAG()` - ✅ Batch link

### Delete Methods ✅
- `deleteMessage()` - ❌ No transaction support
- `deleteMessages()` - ✅ Batch delete (new `BatchDeleteMessages`)
- `deleteMessageTranslate()` - ✅ Single query (new `DeleteMessageTranslate`)
- `deleteMessageTTS()` - ✅ Single query (new `DeleteMessageTTS`)
- `deleteMessageQuery()` - ✅ Single query (new `DeleteMessageQuery`)
- `deleteMessagesBySession()` - ✅ Single query
- `deleteAllMessages()` - ✅ Single query (new `DeleteAllMessages`)

---

## 🚀 Recommendations

### For Production

1. **Use Drizzle for Message operations** - Complex queries with transactions are critical for data integrity
2. **If migrating to Wails**:
   - ✅ Create `QueryMessagesWithRelations` SQL query with all JOINs
   - ✅ Implement transaction support in Go backend
   - ✅ Add batch query for files/chunks (use IN clause)
   - ✅ Add `GetMessageHeatmaps` with GROUP BY
3. **Use Wails for simple operations** - count, search, delete

### Query Improvements Needed

```sql
-- 1. Comprehensive query method
-- name: QueryMessagesWithRelations :many
SELECT m.*, mp.*, mt.*, mtt.*
FROM messages m
LEFT JOIN message_plugins mp ON m.id = mp.id
LEFT JOIN message_translates mt ON m.id = mt.id
LEFT JOIN message_tts mtt ON m.id = mtt.id
WHERE m.user_id = ? AND (? = '' OR m.session_id = ?)
ORDER BY m.created_at ASC
LIMIT ? OFFSET ?;

-- 2. Batch file query
-- name: GetMessageFilesBatch :many
SELECT mf.message_id, f.*
FROM messages_files mf
LEFT JOIN files f ON mf.file_id = f.id
WHERE mf.message_id IN (sqlc.slice('messageIds')) AND mf.user_id = ?;

-- 3. Heatmap query
-- name: GetMessageHeatmaps :many
SELECT 
    DATE(created_at / 1000, 'unixepoch') as date,
    COUNT(*) as count
FROM messages
WHERE user_id = ? AND created_at BETWEEN ? AND ?
GROUP BY DATE(created_at / 1000, 'unixepoch')
ORDER BY date ASC;

-- 4. Get messages by tool_call_id
-- name: GetMessagesByToolCallId :many
SELECT m.*
FROM messages m
INNER JOIN message_plugins mp ON m.id = mp.id
WHERE mp.tool_call_id IN (sqlc.slice('toolCallIds')) AND m.user_id = ?;
```

---

## 📊 Performance Comparison

| Operation | Drizzle | Wails (Current) | Wails (Optimized) |
|-----------|---------|-----------------|-------------------|
| Query 100 messages | 4 queries | 401 queries | 4 queries |
| Create message | 1 transaction | 3+ queries | 1 transaction |
| Count messages | 1 query | 1 query | 1 query |
| Rank models | 1 query | 1 query | 1 query |
| Get heatmaps | 1 query | 1 query + in-memory | 1 query |
| Delete message | 1 transaction | 2+ queries | 1 transaction |

---

## ✅ Files Ready for Review

- ✅ `message.ts` (Drizzle - original, 790 lines)
- ✅ `message.wails.ts` (Wails - migrated, 807 lines)
- ✅ No TypeScript errors
- ⚠️ **Functional but not production-ready** due to:
  - N+1 query problem in `query()` method
  - No transaction support in create/delete operations
  - In-memory grouping in `getHeatmaps()`

---

## 🎯 Next Steps

1. ✅ **Message model migration complete**
2. ⏭️ **Continue with other models**:
   - `topic.ts` / `topic.wails.ts`
   - `agent.ts` / `agent.wails.ts`
   - `file.ts` / `file.wails.ts`
   - `rag.ts` / `rag.wails.ts`
3. 🔄 **Update comparison document** with message model findings
4. 🚀 **Implement optimizations** if moving to production with Wails

---

## 📌 Key Takeaways

| Aspect | Drizzle ✅ | Wails (Current) ❌ | Wails (Optimized) ⚠️ |
|--------|-----------|-------------------|---------------------|
| **Query Efficiency** | Single queries with JOINs | N+1 queries | Requires custom SQL |
| **Transaction Support** | Full ACID transactions | No transactions | Needs Go implementation |
| **Type Safety** | Full TypeScript inference | Many `as any` casts | Same |
| **Code Readability** | Clean, declarative | Verbose, imperative | Same |
| **Developer Experience** | Excellent | Requires helper functions | Same |
| **Performance** | Excellent | Poor (N+1) | Good (with optimizations) |

**Conclusion**: Drizzle is currently superior for complex message operations. Wails can be viable with significant optimizations (custom JOIN queries + transaction support in Go).

