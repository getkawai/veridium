# Adapter Pattern for Database Migration

## Overview

We're using an **Adapter Pattern** to migrate from Drizzle ORM to Wails bindings while keeping the existing API intact. This allows gradual migration without breaking existing code.

## Pattern

```
Existing Code → Model Class (Adapter) → Wails Bindings → Go Backend → SQLite
```

The Model class acts as an adapter that:
1. **Keeps the same public API** - All methods remain unchanged
2. **Translates calls** - Converts Drizzle-style calls to Wails bindings
3. **Handles data transformation** - Converts between nullable types
4. **Maintains compatibility** - Existing code works without changes

## Example: SessionModel

### Before (Drizzle)
```typescript
export class SessionModel {
  private db: LobeChatDatabase; // Drizzle database
  
  query = async () => {
    return this.db.query.sessions.findMany({
      where: eq(sessions.userId, this.userId),
      with: { agentsToSessions: { with: { agent: true } } },
    });
  };
}
```

### After (Wails Adapter)
```typescript
export class SessionModel {
  // No db needed - uses Wails bindings directly
  
  query = async () => {
    // Call Wails bindings
    const sessions = await DB.ListSessions({
      userID: this.userId,
      limit: 9999,
      offset: 0,
    });
    
    // Enrich with related data
    return await Promise.all(
      sessions.map(async (session) => {
        const agents = await DB.GetSessionAgents({
          sessionID: session.id,
          userID: this.userId,
        });
        return { ...session, agentsToSessions: agents.map(a => ({ agent: a })) };
      })
    );
  };
}
```

### Usage (Unchanged!)
```typescript
// Existing code works exactly the same
const sessionModel = new SessionModel(db, userId);
const sessions = await sessionModel.query();
const session = await sessionModel.create({ type: 'agent', config: {...} });
```

## Benefits

1. **Zero Breaking Changes** - Existing code continues to work
2. **Gradual Migration** - Migrate one model at a time
3. **Type Safety** - Full TypeScript support maintained
4. **Easy Rollback** - Can revert if needed
5. **Testing** - Can test new implementation alongside old

## Migration Steps

### 1. Update Imports
```typescript
// Old imports
import { eq, and, desc } from 'drizzle-orm';
import { sessions, agents } from '../schemas';

// New imports
import {
  DB,
  type Session,
  type Agent,
  toNullString,
  parseNullableJSON,
  currentTimestampMs,
} from '@/types/database';
```

### 2. Remove Database Dependency
```typescript
// Old constructor
constructor(db: LobeChatDatabase, userId: string) {
  this.db = db;
  this.userId = userId;
}

// New constructor (db parameter ignored)
constructor(_db: any, userId: string) {
  this.userId = userId;
}
```

### 3. Replace Query Logic
```typescript
// Old (Drizzle)
query = async () => {
  return this.db.query.sessions.findMany({
    where: eq(sessions.userId, this.userId),
  });
};

// New (Wails)
query = async () => {
  return DB.ListSessions({
    userID: this.userId,
    limit: 9999,
    offset: 0,
  });
};
```

### 4. Handle Nullable Fields
```typescript
// Old (Drizzle auto-handles)
const title = session.title;

// New (explicit handling)
const title = getNullableString(session.title);
```

### 5. Handle Relationships
```typescript
// Old (Drizzle with clause)
return this.db.query.sessions.findMany({
  with: { agentsToSessions: { with: { agent: true } } },
});

// New (explicit joins)
const sessions = await DB.ListSessions({...});
const enriched = await Promise.all(
  sessions.map(async (session) => {
    const agents = await DB.GetSessionAgents({
      sessionID: session.id,
      userID: this.userId,
    });
    return { ...session, agentsToSessions: agents.map(a => ({ agent: a })) };
  })
);
```

## Common Patterns

### Pattern 1: Simple Query
```typescript
// Drizzle
findById = async (id: string) => {
  return this.db.query.sessions.findFirst({
    where: eq(sessions.id, id),
  });
};

// Wails
findById = async (id: string) => {
  return DB.GetSession({
    id,
    userID: this.userId,
  });
};
```

### Pattern 2: Create with Transaction
```typescript
// Drizzle
create = async (data) => {
  return this.db.transaction(async (trx) => {
    const agent = await trx.insert(agents).values({...}).returning();
    const session = await trx.insert(sessions).values({...}).returning();
    await trx.insert(agentsToSessions).values({...});
    return session[0];
  });
};

// Wails (backend handles transaction)
create = async (data) => {
  const now = currentTimestampMs();
  
  const session = await DB.CreateSession({
    id: nanoid(),
    userId: this.userId,
    ...data,
    createdAt: now,
    updatedAt: now,
  });
  
  const agent = await DB.CreateAgent({
    id: nanoid(),
    userId: this.userId,
    ...data.config,
    createdAt: now,
    updatedAt: now,
  });
  
  await DB.LinkAgentToSession({
    agentId: agent.id,
    sessionId: session.id,
    userID: this.userId,
  });
  
  return session;
};
```

### Pattern 3: Update
```typescript
// Drizzle
update = async (id: string, data) => {
  return this.db
    .update(sessions)
    .set(data)
    .where(eq(sessions.id, id))
    .returning();
};

// Wails
update = async (id: string, data) => {
  const updated = await DB.UpdateSession({
    id,
    userID: this.userId,
    title: toNullString(data.title),
    description: toNullString(data.description),
    updatedAt: currentTimestampMs(),
  });
  return [updated];
};
```

### Pattern 4: Delete with Cleanup
```typescript
// Drizzle
delete = async (id: string) => {
  return this.db.transaction(async (trx) => {
    await trx.delete(agentsToSessions).where(eq(agentsToSessions.sessionId, id));
    return trx.delete(sessions).where(eq(sessions.id, id));
  });
};

// Wails
delete = async (id: string) => {
  const agents = await DB.GetSessionAgents({
    sessionID: id,
    userID: this.userId,
  });
  
  for (const agent of agents) {
    await DB.UnlinkAgentFromSession({
      agentId: agent.id,
      sessionId: id,
      userID: this.userId,
    });
  }
  
  await DB.DeleteSession({
    id,
    userID: this.userId,
  });
};
```

## Type Compatibility

### Nullable Fields
```typescript
// Drizzle type
type Session = {
  title: string | null;
};

// Wails type
type Session = {
  title: NullString; // { String: string, Valid: boolean }
};

// Adapter handles conversion
const title = getNullableString(session.title); // string | undefined
```

### JSON Fields
```typescript
// Drizzle
type Agent = {
  chatConfig: Record<string, any>;
};

// Wails
type Agent = {
  chatConfig: NullString; // JSON as string
};

// Adapter handles conversion
const config = parseNullableJSON(agent.chatConfig); // Record<string, any> | undefined
```

### Timestamps
```typescript
// Drizzle
type Session = {
  createdAt: Date;
};

// Wails
type Session = {
  createdAt: number; // Unix timestamp in ms
};

// Adapter handles conversion
const date = new Date(session.createdAt);
```

## Testing Strategy

1. **Keep both implementations** during migration
2. **Run parallel tests** to verify compatibility
3. **Compare results** between Drizzle and Wails
4. **Gradual rollout** - one model at a time
5. **Monitor errors** in production

## Rollback Plan

If issues arise:

1. **Revert the model file** to Drizzle version
2. **Keep Wails bindings** for future use
3. **Fix issues** in the adapter
4. **Try again** when ready

## Next Steps

1. ✅ SessionModel migrated (example)
2. ⏳ AgentModel - migrate next
3. ⏳ MessageModel - migrate next
4. ⏳ TopicModel - migrate next
5. ⏳ FileModel - migrate next
6. ⏳ UserModel - migrate next

Once all models are migrated:
- Remove Drizzle dependencies
- Remove old schema files
- Clean up database client code

