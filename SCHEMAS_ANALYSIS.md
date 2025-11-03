# Schemas Directory Analysis

## Current Status

The `frontend/src/database/schemas/` directory contains **Drizzle ORM schema definitions** used for:

### 1. Type Definitions (Main Usage) 📋
```typescript
// Drizzle-inferred types
export type NewChatGroup = typeof chatGroups.$inferInsert;
export type ChatGroupItem = typeof chatGroups.$inferSelect;
export type NewChatGroupAgent = typeof chatGroupsAgents.$inferInsert;
```

**Used by**:
- `frontend/src/services/chatGroup/client.ts` - Type imports only
- `frontend/src/services/chatGroup/type.ts` - Interface definitions
- `frontend/src/server/services/nextAuthUser/index.ts` - Server-side NextAuth

### 2. LobeChatDatabase Type 🔗
```typescript
// frontend/src/database/type.ts
import * as schema from './schemas';

export type LobeChatDatabaseSchema = typeof schema;
export type LobeChatDatabase = BaseSQLiteDatabase<'sync', any, LobeChatDatabaseSchema>;
```

**Used by**:
- Server-side services (file, generation, aiChat, document, agent, user, chunk)
- tRPC async context
- RAG Eval models

## Can We Delete Schemas? ❓

### Option A: ✅ YES - Replace with Go Types

**Pros**:
- 100% Drizzle-free
- Single source of truth (Go backend)
- Type safety from `sqlc`
- Smaller bundle (no Drizzle schemas)

**Cons**:
- Need to refactor 3 files
- Need to update server-side code
- 2-3 hours of work

**Migration Path**:
```typescript
// BEFORE (Drizzle)
import { NewChatGroup, ChatGroupItem } from '@/database/schemas/chatGroup';

// AFTER (Go-generated)
import { ChatGroup, CreateChatGroupParams } from '@/types/database';

// Type mapping:
ChatGroupItem → ChatGroup (Go model)
NewChatGroup  → Omit<CreateChatGroupParams, 'userId'> (Go params)
```

### Option B: ⏸️ KEEP - Minimal Impact

**Pros**:
- Zero work needed
- No breaking changes
- Only 3 files affected
- Schemas are just type definitions (no runtime code)

**Cons**:
- Still have Drizzle schemas in bundle
- Duplicate type definitions

## Detailed Usage Analysis

### 1. ChatGroup Service (2 files)

**frontend/src/services/chatGroup/client.ts**:
```typescript
import {
  ChatGroupAgentItem,    // typeof agents.$inferInsert
  ChatGroupItem,         // typeof chatGroups.$inferSelect
  NewChatGroup,          // typeof chatGroups.$inferInsert
  NewChatGroupAgent,     // typeof chatGroupsAgents.$inferInsert
} from '@/database/schemas/chatGroup';
```

**Migration**:
```typescript
// Use Go-generated types
import { 
  ChatGroup,                    // replaces ChatGroupItem
  CreateChatGroupParams,        // replaces NewChatGroup
  ChatGroupsAgent,              // replaces ChatGroupAgentItem
  LinkChatGroupToAgentParams    // replaces NewChatGroupAgent
} from '@/types/database';
```

### 2. Server-Side NextAuth (1 file)

**frontend/src/server/services/nextAuthUser/index.ts**:
```typescript
import {
  nextauthAccounts,
  nextauthAuthenticators,
  nextauthSessions,
  nextauthVerificationTokens,
  users,
} from '@/database/schemas';
```

**Status**: Server-side code, likely for tRPC routes. Can be kept or migrated.

### 3. Database Type (1 file)

**frontend/src/database/type.ts**:
```typescript
import * as schema from './schemas';

export type LobeChatDatabaseSchema = typeof schema;
export type LobeChatDatabase = BaseSQLiteDatabase<'sync', any, LobeChatDatabaseSchema>;
```

**Used by**: 15 server-side services and tRPC contexts.

**Migration**: Replace `LobeChatDatabase` with a simpler type or just use `any` since we're moving away from Drizzle.

## Impact Assessment

### Files Directly Using Schemas: **3**
1. `frontend/src/services/chatGroup/client.ts`
2. `frontend/src/services/chatGroup/type.ts`
3. `frontend/src/server/services/nextAuthUser/index.ts`

### Files Using LobeChatDatabase Type: **15**
- Server services (8 files)
- tRPC contexts (2 files)
- RAG Eval models (4 files)
- database/type.ts (1 file)

## Recommendation

### 🎯 Option A: Delete Schemas (Clean Slate)

**Effort**: 2-3 hours
**Risk**: Low (only type changes)
**Benefit**: 100% Drizzle-free

**Steps**:
1. ✅ Replace ChatGroup service types with Go types
2. ✅ Update/simplify `LobeChatDatabase` type
3. ✅ Fix server-side NextAuth imports (or keep as legacy)
4. ✅ Delete schemas directory
5. ✅ Test and verify

### Alternative: Keep Schemas (Low Priority)

If server-side code (NextAuth, tRPC) is still heavily using Drizzle, we can:
- Keep schemas for backward compatibility
- Mark as "legacy types only"
- Migrate gradually post-launch

## Decision Required

**Question**: Apakah mau 100% clean (delete schemas) atau keep untuk server-side compatibility?

**My Recommendation**: **Delete schemas** - karena:
1. Only 3 files using them directly
2. Server-side code bisa pakai `any` atau simplified types
3. Go types sudah complete
4. Konsisten dengan migration strategy

---

**Next Steps** (jika user pilih delete):
1. Migrate ChatGroup service to Go types
2. Simplify LobeChatDatabase type
3. Fix NextAuth imports
4. Delete schemas directory
5. Verify and test

