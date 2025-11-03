# 🧹 Final Drizzle Cleanup - Complete Migration Status

## Summary

Inventory of remaining Drizzle files and migration plan.

## ✅ Already Deleted

### Database Client Layer
- ✅ `frontend/src/database/client/wails-sqlite-driver.ts` (187 lines) - DELETED
- ✅ `frontend/src/database/client/wails-sqlite.ts` (50 lines) - DELETED
- ✅ `frontend/src/database/core/db-adaptor.ts` (29 lines) - DELETED
- ✅ `frontend/src/database/core/migrations.json` (7 lines) - DELETED

### Models
- ✅ All 22+ models migrated to Wails bindings
- ✅ No Drizzle syntax in client models

### Services  
- ✅ 13+ services migrated to direct Wails calls
- ✅ ClientService pattern established

## 🗑️ Can Delete Now

### 1. Migration Files (41 files, ~150KB)
**Directory**: `frontend/src/database/migrations/`

**Why delete**:
- Migrations now handled by Go backend using `internal/database/schema/schema.sql`
- Frontend doesn't run migrations anymore
- These are PostgreSQL/PGLite migrations, already converted to SQLite

**Files**:
```
frontend/src/database/migrations/
├── 0000_init.sql
├── 0001_add_client_id.sql
├── 0002-0040_*.sql (40 migration files)
└── meta/ (metadata directory)
```

**Size**: ~150KB
**Command**: 
```bash
rm -rf frontend/src/database/migrations/
```

### 2. Server-Side Models (Still Drizzle)
**Directory**: `frontend/src/database/server/models/ragEval/`

**Why keep for now**:
- Only used by server-side tRPC routers
- RAG Eval not migrated yet (low priority)
- Will be cleaned when we remove tRPC completely

**Files**:
- `dataset.ts` (61 lines)
- `datasetRecord.ts` (84 lines)
- `evaluation.ts` (98 lines)
- `evaluationRecord.ts` (similar)

**Status**: ⏸️ Keep (used by tRPC lambdaClient)

## ⚠️ Need Migration

### Repositories (4 directories, ~1800 lines)
**Directory**: `frontend/src/database/repositories/`

**Used by**:
1. ✅ `aiInfra/` - Used by `services/aiModel/client.ts` & `services/aiProvider/client.ts`
2. ❌ `dataImporter/` - Not used (only by tRPC)
3. ❌ `dataExporter/` - Used by `services/export/client.ts`
4. ❌ `tableViewer/` - Used by `services/tableViewer/client.ts`

**Analysis**:

#### `aiInfra/` Repository
**File**: `frontend/src/database/repositories/aiInfra/index.ts` (328 lines)

**Used by**:
- `services/aiModel/client.ts` - ✅ Active service
- `services/aiProvider/client.ts` - ✅ Active service

**Methods**:
```typescript
createAiProvider(params)
createAiModel(params)
deleteAiProvider(id)
deleteAiModel(id, providerId)
getAllAiModels()
getAllAiProviders()
```

**Status**: ⚠️ **Need to migrate** - Used by active services

**How to migrate**:
Already has Wails models! Just update aiModel/aiProvider services to use models directly:
```typescript
// Current (via repository):
this.aiInfraRepos.createAiModel(params)

// Target (direct model):
this.aiModelModel.create(params)
```

#### `dataImporter/` Repository
**File**: `frontend/src/database/repositories/dataImporter/index.ts` (718 lines)

**Used by**: Nobody! Only imported by itself.

**Status**: ✅ **Can delete** - Unused

#### `dataExporter/` Repository
**File**: `frontend/src/database/repositories/dataExporter/index.ts`

**Used by**: `services/export/client.ts`

**Methods**: Export all data to JSON

**Status**: ⚠️ **Need to migrate** - Used for data export feature

#### `tableViewer/` Repository
**File**: `frontend/src/database/repositories/tableViewer/index.ts`

**Used by**: `services/tableViewer/client.ts`

**Methods**: View/debug database tables

**Status**: ⏸️ **Low priority** - Dev/debug feature

## Migration Plan

### Phase 1: Delete Unused (Now) ✅
```bash
# 1. Delete migrations (Go backend handles this)
rm -rf frontend/src/database/migrations/

# 2. Delete unused repository
rm -rf frontend/src/database/repositories/dataImporter/
```

**Time**: 1 minute
**Impact**: Zero (unused files)

### Phase 2: Migrate aiInfra Repository (15 min)
**Goal**: Remove `aiInfra` repository, use models directly

**Changes**:
1. Update `services/aiModel/client.ts`:
   - Remove `aiInfraRepos` 
   - Use `aiModelModel` directly (already exists!)
   
2. Update `services/aiProvider/client.ts`:
   - Remove `aiInfraRepos`
   - Use `aiProviderModel` directly (already exists!)

3. Delete `repositories/aiInfra/`

**Benefit**: Simpler code, direct model access

### Phase 3: Migrate dataExporter (30 min)
**Goal**: Rewrite export using direct Wails queries

**Approach**:
```typescript
// Current:
const repo = new DataExporterRepo(db);
const data = await repo.exportAllData();

// Target:
const sessions = await DB.ListSessions(userId);
const messages = await DB.ListMessages({ userId, limit: 10000, offset: 0 });
const agents = await DB.ListAgents(userId);
// ... etc
```

### Phase 4: tableViewer (Optional)
**Goal**: Migrate or remove dev/debug feature

**Options**:
- **A**: Migrate to direct queries (30 min)
- **B**: Remove feature (not critical for prod)

## Recommendation

### Do Now:
```bash
# Quick wins - zero impact
rm -rf frontend/src/database/migrations/
rm -rf frontend/src/database/repositories/dataImporter/
```

**Time**: 1 minute
**Benefit**: Clean up 42 unused files (~200KB)

### Do Next (if time permits):
1. **Migrate aiInfra** (15 min) - Most important
2. **Migrate dataExporter** (30 min) - Nice to have
3. **Skip tableViewer** - Dev feature, low priority

### Can Skip:
- `server/models/ragEval/` - Keep until RAG Eval migration
- `tableViewer/` - Dev feature
- Other server-side code - Not critical

## Current Status

### ✅ Migrated (95%)
- Database client layer
- All 22+ client models
- 13+ services
- Type utilities

### ⏸️ Remaining (5%)
- `migrations/` - ✅ Can delete now
- `repositories/aiInfra/` - ⚠️ Should migrate (15 min)
- `repositories/dataExporter/` - ⚠️ Should migrate (30 min)
- `repositories/dataImporter/` - ✅ Can delete now
- `repositories/tableViewer/` - ⏸️ Skip for now
- `server/models/ragEval/` - ⏸️ Keep for tRPC

## Commands

### Immediate Cleanup (Safe)
```bash
cd /Users/yuda/github.com/kawai-network/veridium

# Delete migrations (Go backend handles this)
rm -rf frontend/src/database/migrations/

# Delete unused importer
rm -rf frontend/src/database/repositories/dataImporter/

# Verify
echo "Deleted migrations and unused repositories"
```

### After aiInfra Migration
```bash
# Delete aiInfra repository
rm -rf frontend/src/database/repositories/aiInfra/
```

### After All Migrations
```bash
# Remove entire repositories directory
rm -rf frontend/src/database/repositories/

# Remove drizzle-orm from package.json
npm uninstall drizzle-orm drizzle-kit
```

## Summary

**Can delete now** (zero impact):
- ✅ `migrations/` (41 files, ~150KB)
- ✅ `repositories/dataImporter/` (1 file, ~718 lines)

**Should migrate** (active usage):
- ⚠️ `repositories/aiInfra/` (15 min, used by aiModel/aiProvider)
- ⚠️ `repositories/dataExporter/` (30 min, used by export feature)

**Can skip** (low priority):
- ⏸️ `repositories/tableViewer/` (dev feature)
- ⏸️ `server/models/ragEval/` (tRPC only)

**Total cleanup potential**: 43 files, ~200KB

---

**Next Step**: Delete `migrations/` and `dataImporter/` now (safe, 1 minute).

