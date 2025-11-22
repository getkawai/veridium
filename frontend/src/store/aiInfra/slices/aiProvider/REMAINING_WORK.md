# Remaining Work - Operations Not Yet Migrated

## Overview

This document tracks operations that still use the service layer and have not been migrated to direct DB calls.

**Status**: Phase 1-3 completed (13 operations migrated), but some operations remain unmigrated.

---

## Unmigrated Operations

### AI Provider Operations (action.ts)

#### Read Operations:
- [ ] `refreshAiProviderDetail` - Fetches single provider detail
  - Current: `aiProviderService.getAiProviderById()`
  - Target: `DB.GetAIProvider()`
  - Complexity: Low
  - Impact: Medium (called frequently)

- [ ] `refreshAiProviderList` - Refreshes provider list
  - Current: `aiProviderService.getAiProviderList()`
  - Target: Already migrated in `internal_fetchAiProviderList`
  - Note: This is a wrapper that calls the migrated function
  - Complexity: Low (just needs refactoring)

- [ ] `refreshAiProviderRuntimeState` - Refreshes runtime state
  - Current: `aiProviderService.getAiProviderRuntimeState()`
  - Target: Already migrated in `internal_fetchAiProviderRuntimeState`
  - Note: This is a wrapper that calls the migrated function
  - Complexity: Low (just needs refactoring)

### AI Model Operations (aiModel/action.ts)

#### Write Operations:
- [ ] `batchToggleAiModels` - Toggle multiple models at once
  - Current: `aiModelService.batchToggleAiModels()`
  - Target: `DB.ToggleAIModelEnabled()` in parallel
  - Complexity: Low (similar to existing batch operations)
  - Impact: Low (infrequent operation)

- [ ] `clearModelsByProvider` - Clear all models for a provider
  - Current: `aiModelService.clearModelsByProvider()`
  - Target: `DB.DeleteAIModelsByProvider()`
  - Complexity: Medium (needs backend support)
  - Impact: Low (rare operation)

- [ ] `clearRemoteModels` - Clear remote models
  - Current: `aiModelService.clearRemoteModels()`
  - Target: `DB.DeleteAIModelsBySource()` with source='remote'
  - Complexity: Medium (needs backend support)
  - Impact: Low (rare operation)

- [ ] `updateAiModel` - Update single model
  - Current: `aiModelService.updateAiModel()`
  - Target: `DB.UpdateAIModel()`
  - Complexity: Low (similar to existing update operations)
  - Impact: Medium (called occasionally)

#### Read Operations:
- [ ] `refreshAiModelList` - Refreshes model list
  - Current: `aiModelService.getAiProviderModelList()`
  - Target: `DB.ListAIModels()`
  - Complexity: Low
  - Impact: High (called frequently)

- [ ] `internal_fetchAiProviderModels` - Fetches models for a provider
  - Current: `aiModelService.getAiProviderModelList()`
  - Target: `DB.ListAIModelsByProvider()`
  - Complexity: Low
  - Impact: High (called frequently)

---

## Migration Priority

### High Priority (Frequent Operations):
1. `refreshAiModelList` - Called after every model operation
2. `internal_fetchAiProviderModels` - Called when switching providers
3. `refreshAiProviderDetail` - Called when viewing provider details

### Medium Priority (Occasional Operations):
4. `updateAiModel` - Called when editing models
5. `batchToggleAiModels` - Called when bulk toggling

### Low Priority (Rare Operations):
6. `clearModelsByProvider` - Rarely used
7. `clearRemoteModels` - Rarely used
8. Refactor wrapper functions (`refreshAiProviderList`, `refreshAiProviderRuntimeState`)

---

## Estimated Effort

| Operation | Complexity | Time | Backend Changes |
|-----------|------------|------|-----------------|
| `refreshAiProviderDetail` | Low | 30min | None |
| `refreshAiModelList` | Low | 30min | None |
| `internal_fetchAiProviderModels` | Low | 30min | None |
| `updateAiModel` | Low | 30min | None |
| `batchToggleAiModels` | Low | 30min | None |
| `clearModelsByProvider` | Medium | 1hr | May need new query |
| `clearRemoteModels` | Medium | 1hr | May need new query |
| Wrapper refactoring | Low | 30min | None |

**Total Estimated Time: ~5 hours**

---

## Why These Were Not Migrated in Phase 1-3

1. **Focus on Core CRUD**: Phase 1-3 focused on the most common CRUD operations
2. **Backend Support**: Some operations may need new backend queries
3. **Wrapper Functions**: Some are just wrappers around already-migrated functions
4. **Diminishing Returns**: These operations are less frequently called

---

## Next Steps

### Option A: Complete Migration (Phase 5)
- Migrate all remaining operations
- Delete service layer completely
- Full architecture simplification

### Option B: Hybrid Approach (Current)
- Keep service layer for rare operations
- Migrated operations use direct DB calls
- Best of both worlds: performance + completeness

### Option C: Stop Here
- 13 operations migrated (most important ones)
- Service layer kept for remaining operations
- Good enough for production use

---

## Current Status

**Migrated**: 13 operations (Phase 1-3)
**Remaining**: 8 operations (documented above)
**Service Layer**: Still required (cannot be deleted yet)

**Recommendation**: Option B (Hybrid Approach) - Keep current state, migrate more operations as needed.

---

## Migration Pattern Reference

For future migrations, use these patterns:

### Read Pattern:
```typescript
const USE_DIRECT_DB_CALLS = true;

if (USE_DIRECT_DB_CALLS) {
  const userId = getUserId();
  const data = await DB.GetAIProvider({ id, userId });
  // ... process data
} else {
  const data = await aiProviderService.getAiProviderById(id);
}
```

### Batch Write Pattern:
```typescript
await Promise.all(
  items.map(item => DB.UpdateAIModel({ ...item }))
);
```

### Validation Pattern:
```typescript
try {
  const existing = await DB.GetAIProvider({ id, userId });
  if (existing) throw new Error('Already exists');
} catch (e) {
  if (!e.message?.includes('not found')) throw e;
}
```

---

**Last Updated**: Phase 3 completion
**Next Review**: When performance issues arise or when migrating more operations

