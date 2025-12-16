# Add Session to Group Flow

This document explains the complete flow that occurs when a user adds an agent session to a group (chat group) in LobeChat.

## Overview

In LobeChat, a **chat group** (also called **group session**) is a multi-agent conversation environment where multiple AI agents can participate together. This is different from a **session group**, which is simply an organizational folder for managing sessions.

When a user adds an agent session to a group, a series of operations occur across the UI layer, state management, service layer, and database to establish the relationship between the agent and the group.

## Architecture Layers

The flow involves these architectural layers:

1. **UI Layer** - React components that handle user interactions
2. **State Management** - Zustand stores for managing application state
3. **Service Layer** - Abstraction layer for client/server operations
4. **Database Layer** - Drizzle ORM models and database operations

## Complete Flow

### 1. User Interaction (UI Layer)

**File**: `src/features/GroupChatSettings/AgentTeamMembersSettings.tsx`

The user clicks an "Add" button on an agent card in the group settings interface.

```typescript
const handleAgentAction = async (agentId: string, action: 'add' | 'remove') => {
  if (!activeGroupId) {
    console.error('No active group to perform action on');
    return;
  }
  
  setLoadingAgentId(agentId);
  
  try {
    if (action === 'add') {
      await addAgentsToGroup(activeGroupId, [agentId]);
    }
    
    // Refresh session data to reflect the changes
    await refreshSessions();
  } catch (error) {
    console.error(`Failed to add agent to group:`, error);
  } finally {
    setLoadingAgentId(null);
  }
};
```

**Key Points**:
- Sets loading state for the specific agent being added
- Calls the Zustand store action `addAgentsToGroup`
- Refreshes sessions to update the UI
- Handles errors and clears loading state

### 2. State Management (Zustand Store)

**File**: `src/store/chatGroup/action.ts`

The `addAgentsToGroup` action in the chat group store is invoked:

```typescript
addAgentsToGroup: async (groupId, agentIds) => {
  await chatGroupService.addAgentsToGroup(groupId, agentIds);
  await get().internal_refreshGroups();
}
```

**Key Points**:
- Delegates to the service layer for the actual operation
- Calls `internal_refreshGroups()` to update the local state with fresh data from the database

### 3. Service Layer Routing

The service layer has two implementations depending on the deployment mode:

#### Client-Side (PGLite)

**File**: `src/services/chatGroup/client.ts`

```typescript
async addAgentsToGroup(groupId: string, agentIds: string[]): Promise<ChatGroupAgentItem[]> {
  return this.chatGroupModel.addAgentsToGroup(groupId, agentIds);
}
```

#### Server-Side (PostgreSQL via TRPC)

**File**: `src/services/chatGroup/server.ts`

```typescript
addAgentsToGroup(groupId: string, agentIds: string[]): Promise<ChatGroupAgentItem[]> {
  return lambdaClient.group.addAgentsToGroup.mutate({ agentIds, groupId });
}
```

**Key Points**:
- Client service directly calls the database model
- Server service makes a TRPC mutation to the backend

### 4. TRPC Router (Server-Side Only)

**File**: `src/server/routers/lambda/group.ts`

```typescript
addAgentsToGroup: groupProcedure
  .input(
    z.object({
      agentIds: z.array(z.string()),
      groupId: z.string(),
    }),
  )
  .mutation(async ({ input, ctx }) => {
    return ctx.chatGroupModel.addAgentsToGroup(input.groupId, input.agentIds);
  })
```

**Key Points**:
- Validates input using Zod schema
- Uses authenticated procedure with database middleware
- Calls the database model with validated input

### 5. Database Model Layer

**File**: `packages/database/src/models/chatGroup.ts`

```typescript
async addAgentsToGroup(groupId: string, agentIds: string[]): Promise<ChatGroupAgentItem[]> {
  // 1. Verify group exists
  const group = await this.findById(groupId);
  if (!group) throw new Error('Group not found');

  // 2. Get existing agents to avoid duplicates
  const existingAgents = await this.getGroupAgents(groupId);
  const existingAgentIds = new Set(existingAgents.map((a) => a.id));

  // 3. Filter out agents already in the group
  const newAgentIds = agentIds.filter((id) => !existingAgentIds.has(id));

  if (newAgentIds.length === 0) {
    return [];
  }

  // 4. Create new agent-group relationships
  const newAgents: NewChatGroupAgent[] = newAgentIds.map((agentId) => ({
    agentId,
    chatGroupId: groupId,
    enabled: true,
    userId: this.userId,
  }));

  // 5. Insert into database
  return this.db.insert(chatGroupsAgents).values(newAgents).returning();
}
```

**Key Points**:
- Validates that the group exists
- Prevents duplicate agent additions
- Creates junction table records with default values
- Returns the newly created relationships

### 6. Database Schema

**File**: `packages/database/src/schemas/chatGroup.ts`

Two tables are involved:

#### `chat_groups` Table

Stores group metadata:

```typescript
export const chatGroups = pgTable('chat_groups', {
  id: text('id').primaryKey(),
  title: text('title'),
  description: text('description'),
  config: jsonb('config').$type<ChatGroupConfig>(),
  userId: text('user_id').references(() => users.id, { onDelete: 'cascade' }),
  groupId: text('group_id').references(() => sessionGroups.id, { onDelete: 'set null' }),
  pinned: boolean('pinned').default(false),
  ...timestamps,
});
```

#### `chat_groups_agents` Junction Table

Stores the many-to-many relationship between groups and agents:

```typescript
export const chatGroupsAgents = pgTable('chat_groups_agents', {
  chatGroupId: text('chat_group_id')
    .references(() => chatGroups.id, { onDelete: 'cascade' }),
  agentId: text('agent_id')
    .references(() => agents.id, { onDelete: 'cascade' }),
  userId: text('user_id')
    .references(() => users.id, { onDelete: 'cascade' }),
  enabled: boolean('enabled').default(true),
  order: integer('order').default(0),
  role: text('role').default('participant'),
  ...timestamps,
}, (t) => ({
  pk: primaryKey({ columns: [t.chatGroupId, t.agentId] }),
}));
```

**Key Points**:
- Composite primary key prevents duplicate agent-group pairs
- Cascade delete ensures cleanup when group or agent is deleted
- `enabled` flag allows toggling agent participation without removal
- `order` field supports custom agent ordering in the group
- `role` field can differentiate agent roles (e.g., 'moderator', 'participant')

### 7. State Refresh

After the database operation completes, the application refreshes its state:

#### Group State Refresh

**File**: `src/store/chatGroup/action.ts`

```typescript
internal_refreshGroups: async () => {
  await mutate([FETCH_GROUPS_KEY, true]);
  
  const groups = await chatGroupService.getGroups();
  const groupMap = groups.reduce((acc, group) => {
    acc[group.id] = group;
    return acc;
  }, {} as Record<string, ChatGroupItem>);
  
  set({ groups, groupMap }, false, n('refreshGroups'));
  syncChatStoreGroupMap(groupMap);
}
```

#### Session State Refresh

**File**: `src/store/session/slices/session/action.ts`

```typescript
refreshSessions: async () => {
  await mutate([FETCH_SESSIONS_KEY, true]);
}
```

This triggers SWR to refetch session data, which includes:

**File**: `src/services/session/client.ts`

```typescript
getGroupedSessions: async () => {
  const { sessions, sessionGroups } = await this.sessionModel.queryWithGroups();
  const chatGroups = await this.chatGroupModel.queryWithMemberDetails();

  const groupSessions = chatGroups.map((group) => {
    const { title, description, avatar, backgroundColor, groupId, ...rest } = group;
    return {
      ...rest,
      group: groupId,
      meta: { avatar, backgroundColor, description, title },
      type: 'group' as const,
    };
  });

  const allSessions = [...sessions, ...groupSessions].sort(
    (a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime(),
  );

  return { sessionGroups, sessions: allSessions };
}
```

**Key Points**:
- `queryWithMemberDetails()` fetches groups with their member agents
- Group sessions are transformed to match the session interface
- All sessions (regular + group) are combined and sorted by update time

### 8. UI Update

The React components re-render with the updated data:

**File**: `src/features/GroupChatSettings/AgentTeamMembersSettings.tsx`

```typescript
// Get member IDs from current session
const memberIds = useMemo(() => {
  return currentSession?.members?.map((member: any) => member.id) || [];
}, [currentSession?.members]);

// Separate agents into two groups: in group and not in group
const { agentsInGroup, agentsNotInGroup } = useMemo(() => {
  const inGroup: LobeAgentSession[] = [];
  const notInGroup: LobeAgentSession[] = [];

  agentSessions.forEach((agent) => {
    const agentId = agent.config?.id;
    if (!agentId || agent.id === currentSessionId) return;

    if (memberIds.includes(agentId)) {
      inGroup.push(agent);
    } else {
      notInGroup.push(agent);
    }
  });

  return { agentsInGroup: inGroup, agentsNotInGroup: notInGroup };
}, [agentSessions, memberIds, currentSessionId]);
```

**Key Points**:
- The agent moves from "Available Agents" to "Group Members" section
- Loading state is cleared
- Member count badges update automatically

## Data Flow Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│ 1. User clicks "Add" button on agent card                      │
│    (AgentTeamMembersSettings.tsx)                              │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│ 2. handleAgentAction() calls addAgentsToGroup()                │
│    - Sets loading state                                         │
│    - Calls Zustand store action                                 │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│ 3. Zustand Store: addAgentsToGroup()                           │
│    (src/store/chatGroup/action.ts)                             │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│ 4. Service Layer Routes Request                                 │
│    ├─ Client: chatGroupModel.addAgentsToGroup()                │
│    └─ Server: TRPC mutation → groupRouter.addAgentsToGroup     │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│ 5. Database Model: ChatGroupModel.addAgentsToGroup()           │
│    - Verify group exists                                        │
│    - Check for existing agents                                  │
│    - Filter duplicates                                          │
│    - Insert new records into chat_groups_agents table           │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│ 6. Database: Insert into chat_groups_agents                    │
│    {                                                            │
│      chatGroupId: "group-123",                                  │
│      agentId: "agent-456",                                      │
│      userId: "user-789",                                        │
│      enabled: true,                                             │
│      order: 0,                                                  │
│      role: "participant"                                        │
│    }                                                            │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│ 7. State Refresh                                                │
│    - internal_refreshGroups() → Updates group state            │
│    - refreshSessions() → Triggers SWR revalidation             │
│    - getGroupedSessions() → Fetches updated data               │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│ 8. UI Re-render                                                 │
│    - Agent moves to "Group Members" section                     │
│    - Loading state cleared                                      │
│    - Member count updated                                       │
└─────────────────────────────────────────────────────────────────┘
```

## Key Concepts

### Chat Group vs Session Group

- **Chat Group**: A multi-agent conversation environment where multiple AI agents participate together. This is the focus of this document.
- **Session Group**: An organizational folder for managing sessions (similar to folders in a file system).

### Junction Table Pattern

The `chat_groups_agents` table implements a many-to-many relationship:
- One group can have many agents
- One agent can belong to many groups
- Additional metadata (enabled, order, role) enriches the relationship

### Virtual Agents

The system supports "virtual agents" - agents created specifically for a group:

```typescript
const sessionId = await createSession({
  config: { virtual: true },
  meta: {
    avatar: DEFAULT_AVATAR,
    description: '',
    title: t('settingGroupMembers.defaultAgent'),
  },
}, false);
```

Virtual agents are automatically deleted when the group is deleted.

### Supervisor Mode

Groups can have a "supervisor" (also called "host member"):
- The supervisor orchestrates conversation flow
- Determines which agent speaks next
- Can be enabled/disabled via group config
- Represented by a special `HOST_MEMBER_ID = 'supervisor'`

## Error Handling

The flow includes error handling at multiple levels:

1. **UI Level**: Try-catch blocks with user feedback
2. **Database Level**: Validation checks (group exists, no duplicates)
3. **Schema Level**: Foreign key constraints and cascade rules
4. **Type Safety**: TypeScript and Zod validation

## Performance Considerations

### Race Condition Prevention

When creating a group with initial agents:

```typescript
if (agentIds && agentIds.length > 0) {
  await chatGroupService.addAgentsToGroup(group.id, agentIds);
  
  // Wait to ensure database transactions are committed
  await new Promise<void>((resolve) => {
    setTimeout(resolve, 100);
  });
}
```

This prevents race conditions where `loadGroups()` executes before member addition is fully persisted.

### Duplicate Prevention

The model checks for existing agents before insertion:

```typescript
const existingAgents = await this.getGroupAgents(groupId);
const existingAgentIds = new Set(existingAgents.map((a) => a.id));
const newAgentIds = agentIds.filter((id) => !existingAgentIds.has(id));
```

This avoids unnecessary database operations and constraint violations.

### Efficient State Updates

- Uses SWR for automatic cache invalidation
- Batches state updates to minimize re-renders
- Leverages React useMemo for expensive computations

## Related Operations

### Remove Agent from Group

Similar flow but calls `removeAgentFromGroup()`:

```typescript
async removeAgentFromGroup(groupId: string, agentId: string): Promise<void> {
  await this.db
    .delete(chatGroupsAgents)
    .where(and(
      eq(chatGroupsAgents.chatGroupId, groupId),
      eq(chatGroupsAgents.agentId, agentId)
    ));
}
```

### Update Agent in Group

Modify agent properties within a group:

```typescript
async updateAgentInGroup(
  groupId: string,
  agentId: string,
  updates: Partial<Pick<NewChatGroupAgent, 'order' | 'role'>>
): Promise<NewChatGroupAgent> {
  const [result] = await this.db
    .update(chatGroupsAgents)
    .set({ ...updates, updatedAt: new Date() })
    .where(and(
      eq(chatGroupsAgents.chatGroupId, groupId),
      eq(chatGroupsAgents.agentId, agentId)
    ))
    .returning();

  return result;
}
```

## Testing

The flow can be tested at multiple levels:

### Unit Tests

Test individual model methods:

```typescript
describe('ChatGroupModel.addAgentsToGroup', () => {
  it('should add agents to group', async () => {
    const result = await model.addAgentsToGroup(groupId, [agentId]);
    expect(result).toHaveLength(1);
    expect(result[0].agentId).toBe(agentId);
  });

  it('should prevent duplicate additions', async () => {
    await model.addAgentsToGroup(groupId, [agentId]);
    const result = await model.addAgentsToGroup(groupId, [agentId]);
    expect(result).toHaveLength(0);
  });
});
```

### Integration Tests

Test the complete flow from store to database:

```typescript
it('should add agent to group and refresh state', async () => {
  const { addAgentsToGroup } = useChatGroupStore.getState();
  await addAgentsToGroup(groupId, [agentId]);
  
  const groups = useChatGroupStore.getState().groups;
  const group = groups.find(g => g.id === groupId);
  expect(group.members).toContainEqual(expect.objectContaining({ id: agentId }));
});
```

## Conclusion

Adding a session to a group involves a well-orchestrated flow across multiple architectural layers. The system ensures data consistency through validation, prevents race conditions with timing controls, and maintains UI responsiveness through efficient state management. The use of junction tables provides flexibility for future enhancements like agent roles, ordering, and enable/disable toggles.
