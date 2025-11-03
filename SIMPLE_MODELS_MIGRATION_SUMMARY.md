# Simple Models Migration Summary

## Migrated Models

Successfully migrated 3 simple models from Drizzle to Wails bindings:

### 1. Document Model (`document.wails.ts`)
**Original**: `frontend/src/database/models/document.ts`  
**Migrated**: `frontend/src/database/models/document.wails.ts`

**Methods**:
- ✅ `create` - Creates a new document
- ✅ `delete` - Deletes a document by ID
- ✅ `deleteAll` - Deletes all user documents
- ✅ `query` - Lists all user documents
- ✅ `findById` - Finds document by ID
- ✅ `update` - Updates document fields

**SQL Queries**:
- `GetDocument` - Get single document
- `ListDocuments` - List documents with pagination
- `CreateDocument` - Create new document
- `UpdateDocument` - Update document
- `DeleteDocument` - Delete single document
- `DeleteAllDocuments` - **NEW** - Delete all user documents

**Notes**:
- Simple CRUD operations
- No complex relationships
- Uses pagination for listing (limit/offset)
- Includes document chunks and topic documents relations

---

### 2. Embedding Model (`embedding.wails.ts`)
**Original**: `frontend/src/database/models/embedding.ts`  
**Migrated**: `frontend/src/database/models/embedding.wails.ts`

**Methods**:
- ✅ `create` - Creates a new embedding
- ✅ `bulkCreate` - Creates multiple embeddings with conflict handling
- ✅ `delete` - Deletes an embedding by ID
- ✅ `query` - Lists all user embeddings
- ✅ `findById` - Finds embedding by ID
- ✅ `countUsage` - Counts user's embeddings

**SQL Queries**:
- `GetEmbeddingsItem` - Get single embedding
- `ListEmbeddingsItems` - List all embeddings
- `CreateEmbeddingsItem` - Create new embedding
- `BulkCreateEmbeddingsItems` - **NEW** - Bulk insert with conflict ignore
- `DeleteEmbeddingsItem` - Delete single embedding
- `CountEmbeddingsItems` - **NEW** - Count embeddings

**Notes**:
- Query names changed to avoid conflicts with RAG queries (`GetEmbedding` → `GetEmbeddingsItem`)
- `bulkCreate` uses `Promise.all` - **NOT transactional** ⚠️
- No `tokens` field in schema (removed from SQL)
- `embeddings` stored as BLOB

---

### 3. OAuth Handoff Model (`oauthHandoff.wails.ts`)
**Original**: `frontend/src/database/models/oauthHandoff.ts`  
**Migrated**: `frontend/src/database/models/oauthHandoff.wails.ts`

**Methods**:
- ✅ `create` - Creates a new OAuth handoff (one-time use token)
- ✅ `fetchAndConsume` - Gets and immediately deletes the handoff (5 min TTL)
- ✅ `cleanupExpired` - Removes expired handoffs
- ✅ `exists` - Checks if handoff exists without consuming

**SQL Queries**:
- `GetOAuthHandoff` - Get handoff by ID
- `CreateOAuthHandoff` - Create new handoff
- `DeleteOAuthHandoff` - Delete handoff
- `GetOAuthHandoffByClient` - **NEW** - Get with client and timestamp filter
- `CleanupExpiredOAuthHandoffs` - **NEW** - Bulk delete expired

**Notes**:
- Used for secure credential passing between clients
- 5-minute TTL enforcement
- `create` simulates `onConflictDoNothing` with try-catch
- `fetchAndConsume` is **NOT atomic** (query then delete) ⚠️
- `cleanupExpired` returns 0 (SQLite exec doesn't return row count easily)

---

### 4. Knowledge Base Model (`knowledgeBase.wails.ts`)
**Original**: `frontend/src/database/models/knowledgeBase.ts`  
**Migrated**: `frontend/src/database/models/knowledgeBase.wails.ts`

**Methods**:
- ✅ `create` - Creates a new knowledge base
- ✅ `addFilesToKnowledgeBase` - Links files to knowledge base
- ✅ `delete` - Deletes a knowledge base by ID
- ✅ `deleteAll` - Deletes all user knowledge bases
- ✅ `removeFilesFromKnowledgeBase` - Unlinks files from knowledge base
- ✅ `query` - Lists all user knowledge bases
- ✅ `findById` - Finds knowledge base by ID
- ✅ `update` - Updates knowledge base fields
- ✅ `static findById` - Static method to find by ID without user context

**SQL Queries**:
- `GetKnowledgeBase` - Get single knowledge base
- `ListKnowledgeBases` - List knowledge bases
- `CreateKnowledgeBase` - Create new knowledge base
- `UpdateKnowledgeBase` - Update knowledge base
- `DeleteKnowledgeBase` - Delete single knowledge base
- `DeleteAllKnowledgeBases` - **NEW** - Delete all user knowledge bases
- `BatchLinkKnowledgeBaseToFiles` - **NEW** - Link files in batch
- `BatchUnlinkKnowledgeBaseFromFiles` - **NEW** - Unlink files in batch
- `LinkKnowledgeBaseToFile` - Link single file
- `UnlinkKnowledgeBaseFromFile` - Unlink single file
- `GetKnowledgeBaseFiles` - Get files in knowledge base

**Notes**:
- Simple CRUD with file relationships
- Batch operations use `Promise.all` - **NOT transactional** ⚠️
- Static method for finding by ID without user context
- Settings stored as JSON

---

## Schema Changes

### Documents
- Added `DeleteAllDocuments` query

### Embeddings
- Created new `embeddings.sql` query file
- Renamed queries to avoid conflicts with RAG module
- Schema uses `model`, `client_id` instead of `tokens`

### OAuth Handoffs
- Extended OIDC queries with client filtering and cleanup

### Knowledge Bases
- Added `DeleteAllKnowledgeBases` query
- Added `BatchLinkKnowledgeBaseToFiles` query
- Added `BatchUnlinkKnowledgeBaseFromFiles` query

---

## Common Patterns

### Type Conversions
```typescript
// Nullable strings
toNullString(value as any)
getNullableString(result.field as any)

// JSON fields (NOT nullable in schema)
JSON.stringify(object)  // To store
parseNullableJSON(result.field as any)  // To retrieve

// Integers
toNullInt(value as any)

// Timestamps
currentTimestampMs()  // Current time
new Date(result.createdAt)  // Parse from DB
```

### Error Handling
```typescript
try {
  const result = await DB.Query(...);
  return mapResult(result);
} catch {
  return undefined;  // or null
}
```

---

## Known Limitations

### 1. No Transaction Support
- `bulkCreate` in Embedding model uses `Promise.all`, not atomic
- `fetchAndConsume` in OAuth model has race condition (query then delete)

### 2. No Conflict Handling
- `onConflictDoNothing` simulated with try-catch
- May throw errors instead of silently ignoring

### 3. Row Count Not Available
- SQLite `:exec` queries don't return affected rows easily
- `cleanupExpired` returns 0 instead of actual count

### 4. Pagination
- Document listing uses simple limit/offset
- No cursor-based pagination

---

## Migration Status

✅ **Completed Models** (15/20+):
1. ✅ Session
2. ✅ User  
3. ✅ Message
4. ✅ Topic
5. ✅ Plugin
6. ✅ Thread
7. ✅ Generation
8. ✅ SessionGroup
9. ✅ Agent
10. ✅ APIKey
11. ✅ AsyncTask
12. ✅ Document
13. ✅ Embedding
14. ✅ OAuthHandoff
15. ✅ KnowledgeBase

🔄 **Remaining Models** (Complex - require more work):
- ⏳ ChatGroup (has complex agent relationships)
- ⏳ File (has complex relationships with chunks, knowledge bases, sessions)
- ⏳ Chunk (has semantic search with vector operations)
- ⏳ AIModel (configuration model)
- ⏳ AIProvider (configuration model)
- ⏳ GenerationBatch (batch operations)
- ⏳ GenerationTopic (topic-specific generation)

---

## Next Steps

1. **Test the migrations**:
   - Run the application
   - Test each model's CRUD operations
   - Verify data integrity

2. **Consider optimizations**:
   - Add proper transaction support for critical operations
   - Implement cursor pagination for large datasets
   - Add row count queries for bulk operations

3. **Remove old Drizzle code** (when ready):
   - Delete `*.ts` versions of models
   - Remove Drizzle dependencies
   - Clean up unused database utilities

4. **Documentation**:
   - Update API documentation
   - Create developer guide for new bindings
   - Document known limitations and workarounds

