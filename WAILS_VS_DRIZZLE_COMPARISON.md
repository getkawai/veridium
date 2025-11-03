# Wails vs Drizzle Comparison - SessionModel

Comparison antara implementasi SessionModel menggunakan Wails bindings vs Drizzle ORM.

## File Comparison

- **Drizzle**: `frontend/src/database/models/session.ts`
- **Wails**: `frontend/src/database/models/session.wails.ts`

## Key Differences

### 1. **Imports & Dependencies**

#### Drizzle
```typescript
import { Column, and, asc, count, desc, eq, gt, inArray, isNull, like, not, or, sql } from 'drizzle-orm';
import { agents, agentsToSessions, sessionGroups, sessions, topics } from '../schemas';
import { LobeChatDatabase } from '../type';
import { idGenerator } from '../utils/idGenerator';
```

#### Wails
```typescript
import { nanoid } from 'nanoid';
import { DB, toNullString, toNullJSON, parseNullableJSON, getNullableString, currentTimestampMs, boolToInt } from '@/types/database';
```

**Analysis**: 
- Drizzle requires schema definitions and ORM utilities
- Wails uses generated bindings and helper functions for type conversion

---

### 2. **Constructor**

#### Drizzle
```typescript
constructor(db: LobeChatDatabase, userId: string) {
  this.userId = userId;
  this.db = db;
}
```

#### Wails
```typescript
constructor(_db: any, userId: string) {
  this.userId = userId;
}
```

**Analysis**: 
- Drizzle needs database instance for queries
- Wails doesn't need db instance (uses global `DB` bindings)

---

### 3. **Query Method**

#### Drizzle (Elegant, Single Query)
```typescript
query = async ({ current = 0, pageSize = 9999 } = {}) => {
  const offset = current * pageSize;

  return this.db.query.sessions.findMany({
    limit: pageSize,
    offset,
    orderBy: [desc(sessions.updatedAt)],
    where: and(eq(sessions.userId, this.userId), not(eq(sessions.slug, INBOX_SESSION_ID))),
    with: { agentsToSessions: { columns: {}, with: { agent: true } }, group: true },
  });
};
```

#### Wails (Multiple Queries + Manual Enrichment)
```typescript
query = async ({ current = 0, pageSize = 9999 } = {}) => {
  const offset = current * pageSize;

  // Get sessions with agents
  const sessions = await DB.ListSessions({
    userId: this.userId,
    limit: pageSize,
    offset,
  });

  // Filter out inbox session
  const filtered = sessions.filter((s) => s.slug !== INBOX_SESSION_ID);

  // Enrich with agents and groups (N+1 queries!)
  const enriched = await Promise.all(
    filtered.map(async (session) => {
      const agents = await DB.GetSessionAgents({
        sessionId: session.id,
        userId: this.userId,
      });

      let group: any = undefined;
      if (session.groupId.Valid && session.groupId.String) {
        try {
          group = await DB.GetSessionGroup({
            id: session.groupId.String,
            userId: this.userId,
          });
        } catch {
          // Group not found
        }
      }

      return {
        ...session,
        agentsToSessions: agents.map((agent) => ({ agent })),
        group,
      };
    }),
  );

  return enriched;
};
```

**Analysis**: 
- ✅ **Drizzle**: Single query with JOIN, efficient
- ❌ **Wails**: N+1 query problem (1 query for sessions + N queries for agents + N queries for groups)

---

### 4. **findByIdOrSlug Method**

#### Drizzle (Single Query with OR)
```typescript
findByIdOrSlug = async (idOrSlug: string) => {
  const result = await this.db.query.sessions.findFirst({
    where: and(
      or(eq(sessions.id, idOrSlug), eq(sessions.slug, idOrSlug)),
      eq(sessions.userId, this.userId),
    ),
    with: { agentsToSessions: { columns: {}, with: { agent: true } }, group: true },
  });

  if (!result) return;

  return { ...result, agent: (result?.agentsToSessions?.[0] as any)?.agent } as any;
};
```

#### Wails (Optimized with New Query)
```typescript
findByIdOrSlug = async (idOrSlug: string) => {
  let session: Session | undefined;
  
  try {
    session = await DB.GetSessionByIdOrSlug({
      id: idOrSlug,
      slug: idOrSlug,
      userId: this.userId,
    });
  } catch {
    return undefined;
  }

  if (!session) return undefined;

  // Still needs separate queries for agents and groups
  const agents = await DB.GetSessionAgents({
    sessionId: session.id,
    userId: this.userId,
  });

  let group: any = undefined;
  if (session.groupId.Valid && session.groupId.String) {
    try {
      group = await DB.GetSessionGroup({
        id: session.groupId.String,
        userId: this.userId,
      });
    } catch {
      // Group not found
    }
  }

  return {
    ...session,
    agent: agents[0],
    agentsToSessions: agents.map((agent) => ({ agent })),
    group,
  } as any;
};
```

**Analysis**: 
- ✅ **Drizzle**: Single query with all relations
- ⚠️ **Wails**: Improved with `GetSessionByIdOrSlug`, but still needs 2-3 additional queries

---

### 5. **Count Method**

#### Drizzle (Flexible with SQL Builder)
```typescript
count = async (params?: { endDate?: string; range?: [string, string]; startDate?: string; }) => {
  const result = await this.db
    .select({ count: count(sessions.id) })
    .from(sessions)
    .where(
      genWhere([
        eq(sessions.userId, this.userId),
        params?.range ? genRangeWhere(params.range, sessions.createdAt, (date) => date.toDate()) : undefined,
        params?.endDate ? genEndDateWhere(params.endDate, sessions.createdAt, (date) => date.toDate()) : undefined,
        params?.startDate ? genStartDateWhere(params.startDate, sessions.createdAt, (date) => date.toDate()) : undefined,
      ]),
    );

  return result[0].count;
};
```

#### Wails (Two Separate Queries)
```typescript
count = async (params?: { endDate?: string; range?: [string, string]; startDate?: string; }) => {
  if (!params) {
    return await DB.CountSessions(this.userId);
  }

  let startTime: number;
  let endTime: number;

  if (params.range) {
    const [start, end] = params.range;
    startTime = new Date(start).getTime();
    endTime = new Date(end).getTime();
  } else {
    startTime = params.startDate ? new Date(params.startDate).getTime() : 0;
    endTime = params.endDate ? new Date(params.endDate).getTime() : Date.now();
  }

  return await DB.CountSessionsByDateRange({
    userId: this.userId,
    createdAt: startTime,
    createdAt2: endTime,
  });
};
```

**Analysis**: 
- ✅ **Drizzle**: Flexible SQL builder with helper functions
- ✅ **Wails**: Optimized with dedicated queries, but less flexible

---

### 6. **Rank Method**

#### Drizzle (Complex JOIN with Subquery)
```typescript
rank = async (limit: number = 10) => {
  const inboxResult = await this.db
    .select({ count: count(topics.id).as('count') })
    .from(topics)
    .where(and(eq(topics.userId, this.userId), isNull(topics.sessionId)));

  const inboxCount = inboxResult[0].count;

  if (!inboxCount || inboxCount === 0) return this._rank(limit);

  const result = await this._rank(limit ? limit - 1 : undefined);

  return [
    {
      avatar: DEFAULT_INBOX_AVATAR,
      backgroundColor: null,
      count: inboxCount,
      id: INBOX_SESSION_ID,
      title: 'inbox.title',
    },
    ...result,
  ].sort((a, b) => b.count - a.count);
};

_rank = async (limit: number = 10) => {
  return this.db
    .select({
      avatar: agents.avatar,
      backgroundColor: agents.backgroundColor,
      count: count(topics.id).as('count'),
      id: sessions.id,
      title: agents.title,
    })
    .from(sessions)
    .where(and(eq(sessions.userId, this.userId)))
    .leftJoin(topics, eq(sessions.id, topics.sessionId))
    .leftJoin(agentsToSessions, eq(sessions.id, agentsToSessions.sessionId))
    .leftJoin(agents, eq(agentsToSessions.agentId, agents.id))
    .groupBy(sessions.id, agentsToSessions.agentId, agents.id)
    .having(({ count }) => gt(count, 0))
    .orderBy(desc(sql`count`))
    .limit(limit);
};
```

#### Wails (Dedicated SQL Query)
```typescript
rank = async (limit: number = 10) => {
  // Get inbox count separately
  const inboxCount = await DB.CountTopicsBySession({
    sessionId: toNullString(''),
    userId: this.userId,
  });

  // Get ranked sessions
  const ranked = await DB.GetSessionRank({
    userId: this.userId,
    limit: inboxCount > 0 ? limit - 1 : limit,
  });

  const result = ranked.map((item) => ({
    id: item.id,
    title: getNullableString(item.title as any) || null,
    avatar: getNullableString(item.avatar as any) || null,
    backgroundColor: getNullableString(item.backgroundColor as any) || null,
    count: Number(item.topicCount) || 0,
  }));

  if (inboxCount > 0) {
    return [
      {
        id: INBOX_SESSION_ID,
        title: 'inbox.title',
        avatar: DEFAULT_INBOX_AVATAR,
        backgroundColor: null,
        count: inboxCount,
      },
      ...result,
    ].sort((a, b) => b.count - a.count);
  }

  return result;
};
```

**Analysis**: 
- ✅ **Drizzle**: Flexible, can build complex queries inline
- ✅ **Wails**: Optimized with pre-written SQL query (`GetSessionRank`)

---

### 7. **Create Method**

#### Drizzle (Transaction with Returning)
```typescript
create = async ({ id = idGenerator('sessions'), type = 'agent', session = {}, config = {}, slug }) => {
  return this.db.transaction(async (trx) => {
    if (slug) {
      const existResult = await trx.query.sessions.findFirst({
        where: and(eq(sessions.slug, slug), eq(sessions.userId, this.userId)),
      });

      if (existResult) return existResult;
    }

    if (type === 'group') {
      const result = await trx
        .insert(sessions)
        .values({
          ...session,
          createdAt: new Date(),
          id,
          slug,
          type,
          updatedAt: new Date(),
          userId: this.userId,
        })
        .returning();

      return result[0];
    }

    const newAgents = await trx
      .insert(agents)
      .values({
        ...config,
        createdAt: new Date(),
        id: idGenerator('agents'),
        updatedAt: new Date(),
        userId: this.userId,
      })
      .returning();

    const result = await trx
      .insert(sessions)
      .values({
        ...session,
        createdAt: new Date(),
        id,
        slug,
        type,
        updatedAt: new Date(),
        userId: this.userId,
      })
      .returning();

    await trx.insert(agentsToSessions).values({
      agentId: newAgents[0].id,
      sessionId: id,
      userId: this.userId,
    });

    return result[0];
  });
};
```

#### Wails (Manual Transaction Simulation)
```typescript
create = async ({ id = nanoid(), type = 'agent', session = {}, config = {}, slug }) => {
  if (slug) {
    try {
      const existing = await DB.GetSessionBySlug({
        slug,
        userId: this.userId,
      });
      if (existing) return existing;
    } catch {
      // Doesn't exist, continue
    }
  }

  const now = currentTimestampMs();

  // Create session
  const newSession = await DB.CreateSession({
    id,
    userId: this.userId,
    slug: slug || "",
    title: toNullString(session.title as any),
    description: toNullString(session.description as any),
    avatar: toNullString(session.avatar as any),
    backgroundColor: toNullString(session.backgroundColor as any),
    type: toNullString(type),
    groupId: toNullString(session.groupId as any),
    clientId: toNullString(session.clientId as any),
    pinned: boolToInt(false),
    createdAt: now,
    updatedAt: now,
  });

  // If agent type, create agent and link
  if (type === 'agent') {
    const agentId = nanoid();
    
    await DB.CreateAgent({
      id: agentId,
      userId: this.userId,
      slug: toNullString(undefined),
      title: toNullString(config.title as any),
      description: toNullString(config.description as any),
      tags: toNullJSON(config.tags || []),
      avatar: toNullString(config.avatar as any),
      backgroundColor: toNullString(config.backgroundColor as any),
      plugins: toNullJSON(config.plugins || []),
      clientId: toNullString(config.clientId as any),
      chatConfig: toNullJSON(config.chatConfig),
      fewShots: toNullJSON(config.fewShots),
      model: toNullString(config.model as any),
      params: toNullJSON(config.params),
      provider: toNullString(config.provider as any),
      systemRole: toNullString(config.systemRole as any),
      tts: toNullJSON(config.tts),
      virtual: boolToInt(false),
      openingMessage: toNullString(config.openingMessage as any),
      openingQuestions: toNullJSON(config.openingQuestions || []),
      createdAt: now,
      updatedAt: now,
    });

    await DB.LinkAgentToSession({
      agentId,
      sessionId: id,
      userId: this.userId,
    });
  }

  return newSession;
};
```

**Analysis**: 
- ✅ **Drizzle**: True transaction support, automatic rollback on error
- ❌ **Wails**: No transaction support (each query is separate), manual type conversion required

---

### 8. **Delete Method**

#### Drizzle (Transaction with Cleanup)
```typescript
delete = async (id: string) => {
  return this.db.transaction(async (trx) => {
    const links = await trx
      .select({ agentId: agentsToSessions.agentId })
      .from(agentsToSessions)
      .where(and(eq(agentsToSessions.sessionId, id), eq(agentsToSessions.userId, this.userId)));

    const agentIds = links.map((link) => link.agentId);

    await trx
      .delete(agentsToSessions)
      .where(and(eq(agentsToSessions.sessionId, id), eq(agentsToSessions.userId, this.userId)));

    const result = await trx
      .delete(sessions)
      .where(and(eq(sessions.id, id), eq(sessions.userId, this.userId)));

    await this.clearOrphanAgent(agentIds, trx);

    return result;
  });
};
```

#### Wails (Manual Cleanup)
```typescript
delete = async (id: string) => {
  const agents = await DB.GetSessionAgents({
    sessionId: id,
    userId: this.userId,
  });

  for (const agent of agents) {
    await DB.UnlinkAgentFromSession({
      agentId: agent.id,
      sessionId: id,
      userId: this.userId,
    });
  }

  await DB.DeleteSession({
    id,
    userId: this.userId,
  });

  for (const agent of agents) {
    const agentSessions = await DB.GetAgentSessions({
      agentId: agent.id,
      userId: this.userId,
    });

    if (agentSessions.length === 0) {
      await DB.DeleteAgent({
        id: agent.id,
        userId: this.userId,
      });
    }
  }
};
```

**Analysis**: 
- ✅ **Drizzle**: Atomic transaction, guaranteed consistency
- ❌ **Wails**: No transaction, potential data inconsistency if any query fails

---

### 9. **Batch Delete Method**

#### Drizzle (Efficient with inArray)
```typescript
batchDelete = async (ids: string[]) => {
  if (ids.length === 0) return { count: 0 };

  return this.db.transaction(async (trx) => {
    const links = await trx
      .select({ agentId: agentsToSessions.agentId })
      .from(agentsToSessions)
      .where(
        and(inArray(agentsToSessions.sessionId, ids), eq(agentsToSessions.userId, this.userId)),
      );

    const agentIds = [...new Set(links.map((link) => link.agentId))];

    await trx
      .delete(agentsToSessions)
      .where(
        and(inArray(agentsToSessions.sessionId, ids), eq(agentsToSessions.userId, this.userId)),
      );

    const result = await trx
      .delete(sessions)
      .where(and(inArray(sessions.id, ids), eq(sessions.userId, this.userId)));

    await this.clearOrphanAgent(agentIds, trx);

    return result;
  });
};
```

#### Wails (Optimized with BatchDeleteSessions)
```typescript
batchDelete = async (ids: string[]) => {
  if (ids.length === 0) return { count: 0 };

  const allAgents = await Promise.all(
    ids.map(async (id) => {
      try {
        return await DB.GetSessionAgents({
          sessionId: id,
          userId: this.userId,
        });
      } catch {
        return [];
      }
    })
  );

  const agentIds = [...new Set(allAgents.flat().map((a) => a.id))];

  await Promise.all(
    ids.flatMap((sessionId) =>
      agentIds.map((agentId) =>
        DB.UnlinkAgentFromSession({
          agentId,
          sessionId,
          userId: this.userId,
        }).catch(() => {})
      )
    )
  );

  await DB.BatchDeleteSessions({
    userId: this.userId,
    ids,
  });

  const orphanedAgents = await DB.GetOrphanedAgents(this.userId);
  await Promise.all(
    orphanedAgents.map((agent) =>
      DB.DeleteAgent({
        id: agent.id,
        userId: this.userId,
      })
    )
  );

  return { count: ids.length };
};
```

**Analysis**: 
- ✅ **Drizzle**: Single transaction, efficient `inArray`
- ⚠️ **Wails**: Optimized with `BatchDeleteSessions` and `GetOrphanedAgents`, but still multiple queries without transaction

---

### 10. **Update Method**

#### Drizzle (Simple and Clean)
```typescript
update = async (id: string, data: Partial<SessionItem>) => {
  return this.db
    .update(sessions)
    .set(data)
    .where(and(eq(sessions.id, id), eq(sessions.userId, this.userId)))
    .returning();
};
```

#### Wails (Manual Field Mapping)
```typescript
update = async (id: string, data: Partial<SessionItem>) => {
  const updated = await DB.UpdateSession({
    id,
    userId: this.userId,
    title: data.title !== undefined ? toNullString(getNullableString(data.title as any)) : toNullString(""),
    description: data.description !== undefined ? toNullString(getNullableString(data.description as any)) : toNullString(""),
    avatar: data.avatar !== undefined ? toNullString(getNullableString(avatar as any)) : toNullString(""),
    backgroundColor: data.backgroundColor !== undefined ? toNullString(getNullableString(data.backgroundColor as any)) : toNullString(""),
    groupId: data.groupId !== undefined ? toNullString(getNullableString(data.groupId as any)) : toNullString(""),
    pinned: data.pinned !== undefined ? data.pinned : 0,
    updatedAt: currentTimestampMs(),
  });

  return [updated];
};
```

**Analysis**: 
- ✅ **Drizzle**: Clean, uses spread operator
- ❌ **Wails**: Verbose, manual field mapping, type conversion overhead

---

## Summary

| Feature | Drizzle | Wails |
|---------|---------|-------|
| **Query Efficiency** | ✅ Single queries with JOINs | ❌ N+1 queries |
| **Transaction Support** | ✅ Full ACID transactions | ❌ No transactions |
| **Type Safety** | ✅ Full TypeScript inference | ⚠️ Manual type conversion |
| **Code Readability** | ✅ Clean, declarative | ❌ Verbose, imperative |
| **Flexibility** | ✅ SQL builder for complex queries | ⚠️ Requires pre-written queries |
| **Performance** | ✅ Optimized queries | ❌ Multiple round trips |
| **Developer Experience** | ✅ Excellent | ⚠️ Requires helper functions |

## Recommendations

### For Drizzle (Current)
- ✅ **Keep using Drizzle** for complex queries and transactions
- ✅ Great for rapid development
- ✅ Excellent type safety

### For Wails (Future)
To make Wails competitive with Drizzle, you need:

1. **Add Transaction Support** in Go backend
2. **Create More Efficient Queries** with JOINs
3. **Reduce Type Conversion Overhead** with better type mapping
4. **Add Query Builder** or more flexible query methods

### Hybrid Approach (Recommended)
- Use **Drizzle** for complex operations (create, delete with cleanup)
- Use **Wails** for simple CRUD and when you need backend validation
- Gradually migrate as Wails queries improve

## User Model Comparison

### getUserState Method

#### Drizzle (Single Query with JOIN)
```typescript
getUserState = async (decryptor: DecryptUserKeyVaults) => {
  const result = await this.db
    .select({
      avatar: users.avatar,
      email: users.email,
      firstName: users.firstName,
      // ... all fields
      settingsDefaultAgent: userSettings.defaultAgent,
      settingsGeneral: userSettings.general,
      // ... all settings fields
    })
    .from(users)
    .where(eq(users.id, this.userId))
    .leftJoin(userSettings, eq(users.id, userSettings.id));

  // Single query returns everything
  return { /* mapped data */ };
};
```

#### Wails (Optimized with JOIN Query)
```typescript
getUserState = async (decryptor: DecryptUserKeyVaults) => {
  // Single query with JOIN - efficient!
  const result = await DB.GetUserWithSettings(this.userId);
  
  if (!result) throw new UserNotFoundError();
  
  // Map the result
  return { /* mapped data */ };
};
```

**Analysis**: 
- ✅ **Both are efficient** - single query with JOIN
- ✅ **Wails improved** with `GetUserWithSettings` query

---

### updatePreference Method

#### Drizzle (Simple)
```typescript
updatePreference = async (value: Partial<UserPreference>) => {
  const user = await this.db.query.users.findFirst({ where: eq(users.id, this.userId) });
  if (!user) return;

  return this.db
    .update(users)
    .set({ preference: merge(user.preference, value) })
    .where(eq(users.id, this.userId));
};
```

#### Wails (More Verbose)
```typescript
updatePreference = async (value: Partial<UserPreference>) => {
  let user;
  
  try {
    user = await DB.GetUser(this.userId);
  } catch {
    return;
  }

  if (!user) return;

  const currentPreference = parseNullableJSON(user.preference as any) || {};
  const mergedPreference = merge(currentPreference, value);

  return await DB.UpdateUserPreference({
    id: this.userId,
    preference: toNullJSON(mergedPreference),
    updatedAt: currentTimestampMs(),
  });
};
```

**Analysis**: 
- ✅ **Drizzle**: Cleaner, less type conversion
- ⚠️ **Wails**: More verbose, manual JSON parsing/stringifying

---

## Conclusion

**Drizzle is currently superior** for this use case due to:
- Better query efficiency (no N+1 problem in complex queries)
- Transaction support
- Cleaner, more maintainable code
- Better developer experience
- Less type conversion overhead

**Wails has been improved** with:
- ✅ Efficient JOIN queries (`GetUserWithSettings`, `GetSessionByIdOrSlug`, etc.)
- ✅ Batch operations (`BatchDeleteSessions`)
- ✅ Optimized ranking queries (`GetSessionRank`)
- ✅ Orphaned resource cleanup (`GetOrphanedAgents`)

**Wails still needs**:
- ❌ Transaction support in Go
- ❌ Better type mapping (reduce `toNullString`, `toNullJSON` boilerplate)
- ❌ More flexible query builder

