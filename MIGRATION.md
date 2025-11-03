# Migration Guide: Drizzle ORM → Wails Bindings

## Overview

We've migrated from Drizzle ORM (frontend database) to Wails bindings (backend database with sqlc). This provides:

- ✅ Type-safe database operations
- ✅ Backend-managed transactions
- ✅ Single source of truth (Go backend)
- ✅ Better performance (no frontend SQLite)
- ✅ Automatic TypeScript bindings

## Import Changes

### Before (Drizzle)

```typescript
// Old imports
import { DB } from '@/database/client/db';
import { sessionService } from '@/services/session';
import type { SessionModel } from '@/database/models/session';
```

### After (Wails Bindings)

```typescript
// New imports
import { DB, Session, CreateSessionParams } from '@/database';

// Or specific imports
import { GetSession, CreateSession } from '@/database';
import type { Session } from '@/database';
```

## API Changes

### 1. Sessions

#### Before
```typescript
// Drizzle
const sessions = await sessionService.getSessions(userId);
const session = await sessionService.createSession({
  id: nanoid(),
  userId,
  title: 'New Session',
  // ...
});
```

#### After
```typescript
// Wails
const sessions = await DB.ListSessions({
  userID: userId,
  limit: 100,
  offset: 0,
});

const session = await DB.CreateSession({
  id: nanoid(),
  userId,
  title: toNullString('New Session'),
  slug: toNullString('new-session'),
  type: toNullString('agent'),
  pinned: 0,
  createdAt: currentTimestampMs(),
  updatedAt: currentTimestampMs(),
});
```

### 2. Agents

#### Before
```typescript
// Drizzle
const agents = await agentService.getAgents(userId);
const agent = await agentService.createAgent({
  userId,
  title: 'My Agent',
  // ...
});
```

#### After
```typescript
// Wails
const agents = await DB.ListAgents({
  userID: userId,
  limit: 100,
  offset: 0,
});

const agent = await DB.CreateAgent({
  id: nanoid(),
  userId,
  title: toNullString('My Agent'),
  description: toNullString('Agent description'),
  virtual: 0,
  createdAt: currentTimestampMs(),
  updatedAt: currentTimestampMs(),
});
```

### 3. Messages

#### Before
```typescript
// Drizzle
const messages = await messageService.getMessages(sessionId);
const message = await messageService.createMessage({
  sessionId,
  role: 'user',
  content: 'Hello',
  // ...
});
```

#### After
```typescript
// Wails
const messages = await DB.ListMessages({
  userID: userId,
  sessionID: toNullString(sessionId),
  limit: 100,
  offset: 0,
});

const message = await DB.CreateMessage({
  id: nanoid(),
  userId,
  sessionId: toNullString(sessionId),
  role: 'user',
  content: toNullJSON({ text: 'Hello' }),
  favorite: 0,
  createdAt: currentTimestampMs(),
  updatedAt: currentTimestampMs(),
});
```

## Nullable Fields

### Before (Drizzle)
```typescript
// Drizzle handled nulls automatically
const agent = {
  title: 'My Agent', // Can be null
  description: undefined, // Becomes NULL
};
```

### After (Wails)
```typescript
// Use helper functions for nullable fields
import { toNullString, toNullJSON, parseNullableJSON } from '@/database';

const agent = {
  title: toNullString('My Agent'),
  description: toNullString(undefined), // { String: '', Valid: false }
  chatConfig: toNullJSON({ temperature: 0.7 }),
};

// Extract values
const title = parseNullableJSON(agent.title); // 'My Agent' or undefined
const config = parseNullableJSON(agent.chatConfig); // { temperature: 0.7 } or undefined
```

## Timestamps

### Before (Drizzle)
```typescript
// Drizzle auto-handled timestamps
const session = await sessionService.createSession({
  title: 'New Session',
  // createdAt, updatedAt auto-generated
});
```

### After (Wails)
```typescript
// Explicitly provide timestamps
import { currentTimestampMs } from '@/database';

const session = await DB.CreateSession({
  id: nanoid(),
  title: toNullString('New Session'),
  createdAt: currentTimestampMs(),
  updatedAt: currentTimestampMs(),
});
```

## Boolean Fields

SQLite stores booleans as integers (0/1).

### Before (Drizzle)
```typescript
const session = {
  pinned: true, // Boolean
};
```

### After (Wails)
```typescript
import { boolToInt, intToBool } from '@/database';

const session = {
  pinned: boolToInt(true), // 1
};

// When reading
const isPinned = intToBool(session.pinned); // true
```

## JSON Fields

### Before (Drizzle)
```typescript
const agent = {
  chatConfig: { temperature: 0.7 }, // Auto-serialized
};
```

### After (Wails)
```typescript
import { toNullJSON, parseNullableJSON } from '@/database';

// Writing
const agent = {
  chatConfig: toNullJSON({ temperature: 0.7 }),
};

// Reading
const config = parseNullableJSON<{ temperature: number }>(agent.chatConfig);
```

## Service Layer Removal

All service files have been removed:
- ❌ `services/session/`
- ❌ `services/agent/`
- ❌ `services/message/`
- ❌ etc.

Use direct database calls instead:

```typescript
// Before
import { sessionService } from '@/services/session';
await sessionService.getSessions(userId);

// After
import { DB } from '@/database';
await DB.ListSessions({ userID: userId, limit: 100, offset: 0 });
```

## Repository Pattern (Optional)

If you want to keep business logic separate, create lightweight repositories:

```typescript
// repositories/session.ts
import { DB, Session, toNullString, currentTimestampMs } from '@/database';
import { nanoid } from 'nanoid';

export class SessionRepository {
  static async getAll(userId: string): Promise<Session[]> {
    return DB.ListSessions({
      userID: userId,
      limit: 1000,
      offset: 0,
    });
  }

  static async create(userId: string, title: string): Promise<Session> {
    return DB.CreateSession({
      id: nanoid(),
      userId,
      slug: toNullString(title.toLowerCase().replace(/\s+/g, '-')),
      title: toNullString(title),
      type: toNullString('agent'),
      pinned: 0,
      createdAt: currentTimestampMs(),
      updatedAt: currentTimestampMs(),
    });
  }
}
```

## Type Definitions

All types are now generated from Go:

```typescript
// Available types (auto-generated)
import type {
  User,
  Session,
  SessionGroup,
  Agent,
  Message,
  Topic,
  File,
  // ... 284 types total
} from '@/database';

// Param types for create/update operations
import type {
  CreateUserParams,
  UpdateUserParams,
  CreateSessionParams,
  UpdateSessionParams,
  // ... all CRUD params
} from '@/database';
```

## Migration Checklist

- [ ] Replace all `import` statements from old services
- [ ] Update all database calls to use `DB.*` methods
- [ ] Add nullable field helpers (`toNullString`, etc.)
- [ ] Add explicit timestamps (`currentTimestampMs()`)
- [ ] Convert booleans to integers (`boolToInt()`)
- [ ] Convert JSON fields (`toNullJSON()`)
- [ ] Remove old service files
- [ ] Test all database operations
- [ ] Update tests

## Common Patterns

### Pagination
```typescript
const page = 1;
const pageSize = 20;
const sessions = await DB.ListSessions({
  userID: userId,
  limit: pageSize,
  offset: (page - 1) * pageSize,
});
```

### Search
```typescript
const agents = await DB.SearchAgents({
  userID: userId,
  title: toNullString(`%${searchQuery}%`),
  title_2: toNullString(`%${searchQuery}%`), // For description search
  limit: 50,
});
```

### Relationships
```typescript
// Get session with agents
const session = await DB.GetSession({ id: sessionId, userID: userId });
const agents = await DB.GetSessionAgents({ sessionID: sessionId, userID: userId });
```

## Need Help?

- Check `DATABASE.md` for database architecture
- Check `frontend/src/types/database.ts` for type utilities
- Check generated bindings in `frontend/bindings/`
- All 416 database methods are available via `DB.*`

