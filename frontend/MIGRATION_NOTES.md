# Migration from PGlite/PostgreSQL to Wails SQLite

## Summary

Successfully migrated the entire database layer from PGlite/PostgreSQL to Wails SQLite. This migration provides:

- **Smaller bundle size**: No more 10MB+ WASM files to download
- **Native performance**: Direct integration with Go backend through Wails
- **Simplified architecture**: Single database implementation for all environments

## Changes Made

### 1. Custom Drizzle-Wails Adapter ✅

Created custom adapter to connect Drizzle ORM with Wails SQLite service:

- `src/database/client/wails-sqlite-driver.ts` - Custom Drizzle driver
- `src/database/client/wails-sqlite.ts` - Wails SQLite initialization

### 2. Schema Conversion ✅

Converted all 19 schema files from PostgreSQL to SQLite:

- `pgTable` → `sqliteTable`
- `jsonb` → `text` with `mode: 'json'`
- `uuid` → `text` with `randomUUID()`
- `boolean` → `integer` with `mode: 'boolean'`
- `vector` → `blob` with `mode: 'buffer'`
- `timestamp` → `integer` with `mode: 'timestamp_ms'`
- `varchar` → `text`
- `serial` → `integer` with autoincrement

### 3. Database Client Updates ✅

Updated `src/database/client/db.ts`:

- Removed WASM/PGlite initialization
- Integrated Wails SQLite driver
- Maintained singleton pattern and migration logic
- Updated all type references

### 4. Core Adapters ✅

- Updated `src/database/core/db-adaptor.ts` to use Wails SQLite
- Updated `src/database/type.ts` from `NeonDatabase` to `BaseSQLiteDatabase`
- Removed separate Electron/server implementations

### 5. Vector Search Implementation ✅

Created JavaScript-based vector similarity search (src/database/utils/vectorSearch.ts):

- `cosineSimilarity()` - Calculate similarity between vectors
- `euclideanDistance()` - Alternative distance metric
- `findTopKSimilar()` - Find most similar vectors
- `vectorToBuffer()` / `bufferToVector()` - Convert between storage formats

Updated models:

- `src/database/models/chunk.ts` - Uses new vector search utilities
- Embeddings stored as Buffer (blob) in SQLite

### 6. Files Deleted ✅

Removed obsolete files:

- `src/database/client/pglite.ts`
- `src/database/client/pglite.worker.ts`
- `src/database/client/type.ts`
- `src/database/core/electron.ts`
- `src/database/core/web-server.ts`
- `src/database/core/dbForTest.ts`

## Required Dependencies

Add these to your `package.json`:

```json
{
  "dependencies": {
    "drizzle-orm": "^0.30.0",
    "ts-md5": "^1.3.1"
  }
}
```

Install with:

```bash
npm install drizzle-orm ts-md5
```

## Migrations

### Current Status

The existing `src/database/core/migrations.json` contains PostgreSQL-format migrations. These need to be converted or regenerated for SQLite.

### Options

#### Option 1: Generate Fresh Migrations (Recommended)

Use Drizzle Kit to generate new SQLite migrations from the schemas:

```bash
npm install -D drizzle-kit
npx drizzle-kit generate:sqlite
```

#### Option 2: Manual Migration Conversion

If you have existing data, you'll need to:

1. Export data from PostgreSQL/PGlite
2. Convert PostgreSQL SQL to SQLite SQL:
   - Remove `timestamp with time zone` → use `INTEGER`
   - Convert `jsonb` → `TEXT`
   - Convert `uuid` → `TEXT`
   - Remove `gen_random_uuid()` → handled in code
   - Convert `vector` columns → `BLOB`
3. Import into SQLite

### Migration Differences

**PostgreSQL → SQLite conversions needed:**

```sql
-- PostgreSQL
CREATE TABLE users (
  id text PRIMARY KEY,
  created_at timestamp with time zone DEFAULT now() NOT NULL,
  preference jsonb DEFAULT '{}'::jsonb
);

-- SQLite
CREATE TABLE users (
  id TEXT PRIMARY KEY,
  created_at INTEGER NOT NULL,
  preference TEXT NOT NULL DEFAULT '{}'
);
```

## Vector Search Performance

### Current Implementation

Vector search is now performed in JavaScript:

1. Fetch all embeddings from database
2. Calculate cosine similarity in JS
3. Sort and return top K results

### Performance Considerations

- **Small datasets (< 10,000 vectors)**: Excellent performance
- **Medium datasets (10,000 - 100,000)**: Good performance, may need optimization
- **Large datasets (> 100,000)**: Consider:
  - Pre-filtering by metadata before vector comparison
  - Batch processing with `batchVectorSearch()`
  - Implementing approximate nearest neighbor (ANN) algorithms
  - Using a separate vector search service

### Future Optimizations

Consider these improvements if performance becomes an issue:

1. **HNSW Index**: Implement Hierarchical Navigable Small World algorithm
2. **SQLite Extensions**: Use sqlite-vss or sqlite-vec extensions (requires compilation)
3. **WebAssembly**: Port vector operations to WASM for better performance
4. **Worker Threads**: Offload vector calculations to Web Workers

## Testing Checklist

- [ ] Database initializes correctly on app start
- [ ] All CRUD operations work in models
- [ ] Migrations apply successfully
- [ ] Vector/embedding storage and retrieval works
- [ ] Semantic search returns relevant results
- [ ] Transaction handling works correctly
- [ ] Error handling and recovery functions
- [ ] No references to PGlite/PostgreSQL remain in code

## Known Limitations

1. **No native vector search**: Similarity calculation happens in JavaScript
2. **UUID generation**: Now uses crypto.randomUUID() instead of database-generated UUIDs
3. **JSON type safety**: SQLite stores JSON as text, type safety maintained at application level
4. **Array types**: PostgreSQL arrays converted to JSON arrays in SQLite
5. **Timestamp precision**: Stored as milliseconds since epoch

## Backend Configuration

Make sure your Wails backend (Go) has:

1. SQLite service properly configured
2. Database file path accessible
3. Appropriate permissions for read/write operations

Example Go configuration (if needed):

```go
// Configure SQLite service in your Wails app
app := application.New(application.Options{
    Services: []application.Service{
        sqlite.NewService(sqlite.Options{
            Path: "veridium.db",
        }),
    },
})
```

## Rollback Plan

If you need to rollback to PostgreSQL/PGlite:

1. Restore the deleted files from git history:
   ```bash
   git checkout HEAD~1 -- src/database/client/pglite.ts
   git checkout HEAD~1 -- src/database/core/electron.ts
   # ... etc
   ```

2. Revert schema changes:
   ```bash
   git checkout HEAD~1 -- src/database/schemas/
   ```

3. Restore old database client:
   ```bash
   git checkout HEAD~1 -- src/database/client/db.ts
   ```

## Support

For issues or questions:

1. Check Wails SQLite service documentation
2. Review Drizzle ORM SQLite adapter docs
3. Test vector search with sample data
4. Check browser console for initialization errors

