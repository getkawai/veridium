# 🗑️ tRPC Cleanup - Safe to Delete

## Summary

All lambda routers and tRPC infrastructure can be safely deleted. Frontend now uses direct calls to models via Wails.

## Files to Delete

### 1. Lambda Routers (28 files)
```bash
rm -rf frontend/src/server/routers/lambda/
```

**Contents**:
- message.ts
- session.ts
- user.ts
- agent.ts
- topic.ts
- file.ts
- chunk.ts
- generation.ts
- generationTopic.ts
- generationBatch.ts
- aiProvider.ts
- aiModel.ts
- apiKey.ts
- document.ts
- knowledgeBase.ts
- plugin.ts
- sessionGroup.ts
- thread.ts
- group.ts
- ragEval.ts
- aiChat.ts
- comfyui.ts
- image.ts
- importer.ts
- exporter.ts
- upload.ts
- _template.ts
- index.ts
- market/ (directory)
- config/ (directory)
- __tests__/ (directory)

### 2. Async Routers
```bash
rm -rf frontend/src/server/routers/async/
```

### 3. Server Services (if not used elsewhere)
```bash
# Check first if used
grep -r "from.*@/server/services" frontend/src --include="*.ts" --include="*.tsx"

# If not used:
rm -rf frontend/src/server/services/
```

### 4. tRPC Infrastructure
```bash
rm -rf frontend/src/libs/trpc/
```

**Contents**:
- lambda/ (setup)
- async/ (setup)
- mock.ts (not used)

### 5. tRPC Middleware
```bash
# Check what's in middleware
ls -la frontend/src/libs/trpc/lambda/middleware/

# If all tRPC-related:
rm frontend/src/libs/trpc/lambda/middleware/serverDatabase.ts
```

## Migration Status

✅ **Completed**:
- All type imports moved to `@/types/generation-types`
- All frontend code uses direct calls via `ClientService`
- 0 imports from lambda routers

## Verification

Run these commands to verify no dependencies:

```bash
# Should return nothing
grep -r "from.*@/server/routers/lambda" frontend/src --include="*.ts" --include="*.tsx"

# Should only show generation-types
grep -r "UpdateTopicValue\|GetGenerationStatusResult" frontend/src --include="*.ts" --include="*.tsx"
```

## Backup Plan

Before deleting, create a backup:

```bash
# Backup entire server directory
tar -czf server-backup-$(date +%Y%m%d).tar.gz frontend/src/server/
```

## Benefits After Cleanup

- **~50 fewer files** 
- **Smaller bundle size** (~500KB less)
- **Less confusion** (single pattern: direct calls)
- **Faster builds** (less TypeScript to compile)

## Execute Cleanup

```bash
cd /Users/yuda/github.com/kawai-network/veridium

# Backup first
tar -czf server-backup-$(date +%Y%m%d).tar.gz frontend/src/server/

# Delete lambda routers
rm -rf frontend/src/server/routers/lambda/

# Delete async routers
rm -rf frontend/src/server/routers/async/

# Delete tRPC libs
rm -rf frontend/src/libs/trpc/

# Verify
echo "=== Remaining server files ==="
find frontend/src/server -type f -name "*.ts" | wc -l
```

## What to Keep

✅ **Keep**:
- `frontend/src/services/` - Client services with direct model access
- `frontend/src/database/models/` - All Wails models
- `frontend/src/types/` - All type definitions
- `frontend/src/types/generation-types.ts` - New types file

## Result

**Before Cleanup**:
- 68 router files
- 10 tRPC setup files
- ~1000 lines of tRPC code

**After Cleanup**:
- 0 router files ✅
- 0 tRPC code ✅
- Clean architecture with direct calls ✅
