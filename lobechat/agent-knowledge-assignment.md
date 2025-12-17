# Agent Knowledge Base and File Assignment

## Overview

This document explains the complete process of how agents in LobeChat can be assigned knowledge bases and files. The system allows agents to access specific knowledge sources (knowledge bases and files) during conversations, enabling them to provide contextually relevant responses based on the assigned knowledge.

## Purpose

The knowledge assignment system enables:

- **Contextual AI Responses**: Agents can reference specific documents and knowledge bases
- **Selective Knowledge Access**: Control which knowledge sources each agent can access
- **Enable/Disable Toggle**: Temporarily enable or disable knowledge sources without removing them
- **Unified Management**: Manage both knowledge bases and individual files in one interface

## Key Concepts

### Knowledge Types

1. **Knowledge Base** (`KnowledgeType.KnowledgeBase`)
   - Collection of documents organized as a knowledge base
   - Can contain multiple files and chunks
   - Has metadata (name, description, avatar)
   - Managed as a single unit

2. **File** (`KnowledgeType.File`)
   - Individual files uploaded by users
   - Can be assigned directly to agents
   - Excludes image files (filtered out)
   - Can exist independently or within knowledge bases

### Assignment States

- **Assigned & Enabled**: Knowledge is assigned to agent and actively used
- **Assigned & Disabled**: Knowledge is assigned but temporarily disabled
- **Not Assigned**: Knowledge is available but not assigned to agent

## Architecture

### Component Hierarchy

```
Knowledge Assignment System
├── UI Layer
│   ├── KnowledgeBaseModal/AssignKnowledgeBase/
│   │   ├── List.tsx - Main list component
│   │   └── Item/
│   │       ├── index.tsx - Individual item
│   │       └── Action.tsx - Assign/remove actions
│   │
├── State Management (Zustand)
│   └── store/agent/slices/chat/action.ts
│       ├── useFetchFilesAndKnowledgeBases()
│       ├── addKnowledgeBaseToAgent()
│       ├── addFilesToAgent()
│       ├── removeKnowledgeBaseFromAgent()
│       ├── removeFileFromAgent()
│       ├── toggleKnowledgeBase()
│       └── toggleFile()
│
├── Service Layer
│   └── services/agent.ts
│       └── getFilesAndKnowledgeBases()
│
├── TRPC Router
│   └── server/routers/lambda/agent.ts
│       └── getKnowledgeBasesAndFiles
│
├── Database Model
│   └── packages/database/src/models/agent.ts
│       └── getAgentAssignedKnowledge()
│
└── Database Schema
    ├── agents_knowledge_bases (junction table)
    └── agents_files (junction table)
```

## Complete Data Flow: `getKnowledgeBasesAndFiles`

### Overview

The `getKnowledgeBasesAndFiles` endpoint retrieves all available knowledge bases and files, marking which ones are assigned to a specific agent.

### Step-by-Step Process

#### 1. UI Component Initiates Request

**File**: `src/features/KnowledgeBaseModal/AssignKnowledgeBase/List.tsx`

```typescript
export const List = memo(() => {
  const useFetchFilesAndKnowledgeBases = useAgentStore(
    (s) => s.useFetchFilesAndKnowledgeBases
  );

  const { isLoading, error, data } = useFetchFilesAndKnowledgeBases();

  return (
    <Virtuoso
      itemContent={(index) => {
        const item = data![index];
        return <Item key={item.id} {...item} />;
      }}
      totalCount={data!.length}
    />
  );
});
```

**Key Points**:
- Uses SWR hook from Zustand store
- Displays loading state while fetching
- Renders virtualized list for performance
- Each item shows assignment status

#### 2. Zustand Store Hook (SWR)

**File**: `src/store/agent/slices/chat/action.ts`

```typescript
useFetchFilesAndKnowledgeBases: () => {
  return useClientDataSWR<KnowledgeItem[]>(
    [FETCH_AGENT_KNOWLEDGE_KEY, get().activeAgentId],
    ([, id]: string[]) => agentService.getFilesAndKnowledgeBases(id),
    {
      fallbackData: [],
      suspense: true,
    },
  );
}
```

**Key Points**:
- SWR key includes agent ID for cache isolation
- Suspense mode for React Suspense integration
- Fallback to empty array prevents undefined errors
- Automatic revalidation on focus/reconnect

**Cache Key**: `['FETCH_AGENT_KNOWLEDGE', agentId]`

#### 3. Service Layer Call

**File**: `src/services/agent.ts`

```typescript
class AgentService {
  getFilesAndKnowledgeBases = async (agentId: string) => {
    return lambdaClient.agent.getKnowledgeBasesAndFiles.query({ agentId });
  };
}
```

**Key Points**:
- Thin wrapper around TRPC client
- Type-safe with TRPC inference
- Single parameter: agent ID

#### 4. TRPC Router Endpoint

**File**: `src/server/routers/lambda/agent.ts`

```typescript
getKnowledgeBasesAndFiles: agentProcedure
  .input(
    z.object({
      agentId: z.string(),
    }),
  )
  .query(async ({ ctx, input }): Promise<KnowledgeItem[]> => {
    // Step 4.1: Fetch all knowledge bases
    const knowledgeBases = await ctx.knowledgeBaseModel.query();

    // Step 4.2: Fetch all files (excluding those in knowledge bases)
    const files = await ctx.fileModel.query({
      showFilesInKnowledgeBase: false,
    });

    // Step 4.3: Get agent's assigned knowledge
    const knowledge = await ctx.agentModel.getAgentAssignedKnowledge(input.agentId);

    // Step 4.4: Combine and format results
    return [
      // Files (excluding images)
      ...files
        .filter((file) => !file.fileType.startsWith('image'))
        .map((file) => ({
          enabled: knowledge.files.some((item) => item.id === file.id),
          fileType: file.fileType,
          id: file.id,
          name: file.name,
          type: KnowledgeType.File,
        })),
      
      // Knowledge bases
      ...knowledgeBases.map((knowledgeBase) => ({
        avatar: knowledgeBase.avatar,
        description: knowledgeBase.description,
        enabled: knowledge.knowledgeBases.some((item) => item.id === knowledgeBase.id),
        id: knowledgeBase.id,
        name: knowledgeBase.name,
        type: KnowledgeType.KnowledgeBase,
      })),
    ];
  })
```

**Process Breakdown**:

##### Step 4.1: Fetch All Knowledge Bases

```typescript
const knowledgeBases = await ctx.knowledgeBaseModel.query();
```

Retrieves all knowledge bases for the current user from the `knowledge_bases` table.

##### Step 4.2: Fetch All Files

```typescript
const files = await ctx.fileModel.query({
  showFilesInKnowledgeBase: false,
});
```

Retrieves all files for the current user, excluding files that are already part of knowledge bases.

**Why exclude files in knowledge bases?**
- Prevents duplicate entries
- Files in knowledge bases are accessed through the knowledge base
- Simplifies UI by showing knowledge bases as single units

##### Step 4.3: Get Agent's Assigned Knowledge

```typescript
const knowledge = await ctx.agentModel.getAgentAssignedKnowledge(input.agentId);
```

Returns:
```typescript
{
  files: Array<{ id, enabled, ...fileData }>,
  knowledgeBases: Array<{ id, enabled, ...kbData }>
}
```

##### Step 4.4: Combine and Format Results

**File Processing**:
```typescript
...files
  .filter((file) => !file.fileType.startsWith('image'))
  .map((file) => ({
    enabled: knowledge.files.some((item) => item.id === file.id),
    fileType: file.fileType,
    id: file.id,
    name: file.name,
    type: KnowledgeType.File,
  }))
```

**Filters**:
- Excludes image files (not useful for text-based knowledge)

**Mapping**:
- `enabled`: Checks if file is in agent's assigned files
- `type`: Set to `KnowledgeType.File`

**Knowledge Base Processing**:
```typescript
...knowledgeBases.map((knowledgeBase) => ({
  avatar: knowledgeBase.avatar,
  description: knowledgeBase.description,
  enabled: knowledge.knowledgeBases.some((item) => item.id === knowledgeBase.id),
  id: knowledgeBase.id,
  name: knowledgeBase.name,
  type: KnowledgeType.KnowledgeBase,
}))
```

**Mapping**:
- `enabled`: Checks if knowledge base is in agent's assigned knowledge bases
- `type`: Set to `KnowledgeType.KnowledgeBase`
- Includes metadata (avatar, description)

#### 5. Database Model Query

**File**: `packages/database/src/models/agent.ts`

```typescript
getAgentAssignedKnowledge = async (id: string) => {
  // Query assigned knowledge bases
  const knowledgeBaseResult = await this.db
    .select({ enabled: agentsKnowledgeBases.enabled, knowledgeBases })
    .from(agentsKnowledgeBases)
    .where(eq(agentsKnowledgeBases.agentId, id))
    .orderBy(desc(agentsKnowledgeBases.createdAt))
    .leftJoin(knowledgeBases, eq(knowledgeBases.id, agentsKnowledgeBases.knowledgeBaseId));

  // Query assigned files
  const fileResult = await this.db
    .select({ enabled: agentsFiles.enabled, files })
    .from(agentsFiles)
    .where(eq(agentsFiles.agentId, id))
    .orderBy(desc(agentsFiles.createdAt))
    .leftJoin(files, eq(files.id, agentsFiles.fileId));

  return {
    files: fileResult.map((item) => ({
      ...item.files,
      enabled: item.enabled,
    })),
    knowledgeBases: knowledgeBaseResult.map((item) => ({
      ...item.knowledgeBases,
      enabled: item.enabled,
    })),
  };
};
```

**SQL Queries**:

**Knowledge Bases Query**:
```sql
SELECT 
  agents_knowledge_bases.enabled,
  knowledge_bases.*
FROM agents_knowledge_bases
LEFT JOIN knowledge_bases 
  ON knowledge_bases.id = agents_knowledge_bases.knowledge_base_id
WHERE agents_knowledge_bases.agent_id = ?
ORDER BY agents_knowledge_bases.created_at DESC
```

**Files Query**:
```sql
SELECT 
  agents_files.enabled,
  files.*
FROM agents_files
LEFT JOIN files 
  ON files.id = agents_files.file_id
WHERE agents_files.agent_id = ?
ORDER BY agents_files.created_at DESC
```

**Key Points**:
- Uses LEFT JOIN to include full knowledge base/file data
- Ordered by creation date (newest first)
- Includes `enabled` flag from junction tables
- Returns empty arrays if no assignments

#### 6. Response Format

**Type**: `KnowledgeItem[]`

```typescript
interface KnowledgeItem {
  id: string;
  name: string;
  type: KnowledgeType.File | KnowledgeType.KnowledgeBase;
  enabled: boolean;
  
  // File-specific
  fileType?: string;
  
  // Knowledge Base-specific
  avatar?: string;
  description?: string;
}
```

**Example Response**:

```json
[
  {
    "id": "file_abc123",
    "name": "Product Documentation.pdf",
    "type": "file",
    "enabled": true,
    "fileType": "application/pdf"
  },
  {
    "id": "file_def456",
    "name": "API Reference.md",
    "type": "file",
    "enabled": false,
    "fileType": "text/markdown"
  },
  {
    "id": "kb_xyz789",
    "name": "Company Knowledge Base",
    "type": "knowledgeBase",
    "enabled": true,
    "avatar": "https://example.com/avatar.png",
    "description": "Internal company documentation"
  }
]
```

#### 7. UI Rendering

**File**: `src/features/KnowledgeBaseModal/AssignKnowledgeBase/Item/index.tsx`

Each item in the list displays:
- Icon (based on type: file or knowledge base)
- Name
- Description (for knowledge bases)
- File type badge (for files)
- Action button (assign/remove)
- Enable/disable toggle (if assigned)

## Database Schema

### Junction Tables

#### `agents_knowledge_bases`

Stores many-to-many relationships between agents and knowledge bases.

```sql
CREATE TABLE agents_knowledge_bases (
  agent_id TEXT NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
  knowledge_base_id TEXT NOT NULL REFERENCES knowledge_bases(id) ON DELETE CASCADE,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  enabled BOOLEAN DEFAULT true,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  accessed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  PRIMARY KEY (agent_id, knowledge_base_id)
);
```

**Key Points**:
- Composite primary key prevents duplicate assignments
- `enabled` flag for temporary disable without removal
- Cascade delete ensures cleanup
- Timestamps for audit trail

#### `agents_files`

Stores many-to-many relationships between agents and files.

```sql
CREATE TABLE agents_files (
  file_id TEXT NOT NULL REFERENCES files(id) ON DELETE CASCADE,
  agent_id TEXT NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  enabled BOOLEAN DEFAULT true,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  accessed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  PRIMARY KEY (file_id, agent_id, user_id)
);
```

**Key Points**:
- Composite primary key includes user_id
- Same structure as knowledge bases table
- Supports enable/disable toggle

## Assignment Operations

### 1. Assign Knowledge Base to Agent

**User Action**: User clicks "Assign" button on a knowledge base

**Flow**:

```
User clicks Assign
        ↓
Action.tsx: assignKnowledge()
        ↓
Store: addKnowledgeBaseToAgent(knowledgeBaseId)
        ↓
Service: agentService.createAgentKnowledgeBase()
        ↓
TRPC: agent.createAgentKnowledgeBase.mutate()
        ↓
Model: agentModel.createAgentKnowledgeBase()
        ↓
Database: INSERT INTO agents_knowledge_bases
        ↓
Refresh: internal_refreshAgentConfig()
        ↓
Refresh: internal_refreshAgentKnowledge()
        ↓
UI Updates
```

**Code**:

**UI** (`src/features/KnowledgeBaseModal/AssignKnowledgeBase/Item/Action.tsx`):
```typescript
const assignKnowledge = async () => {
  setLoading(true);
  if (type === KnowledgeType.KnowledgeBase) {
    await addKnowledgeBasesToAgent(id);
  } else {
    await addFilesToAgent([id], true);
  }
  setLoading(false);
};
```

**Store** (`src/store/agent/slices/chat/action.ts`):
```typescript
addKnowledgeBaseToAgent: async (knowledgeBaseId) => {
  const { activeAgentId, internal_refreshAgentConfig, internal_refreshAgentKnowledge } = get();
  if (!activeAgentId) return;

  await agentService.createAgentKnowledgeBase(activeAgentId, knowledgeBaseId, true);
  await internal_refreshAgentConfig(get().activeId);
  await internal_refreshAgentKnowledge();
}
```

**Service** (`src/services/agent.ts`):
```typescript
createAgentKnowledgeBase = async (
  agentId: string,
  knowledgeBaseId: string,
  enabled?: boolean,
) => {
  return lambdaClient.agent.createAgentKnowledgeBase.mutate({
    agentId,
    enabled,
    knowledgeBaseId,
  });
};
```

**TRPC Router** (`src/server/routers/lambda/agent.ts`):
```typescript
createAgentKnowledgeBase: agentProcedure
  .input(
    z.object({
      agentId: z.string(),
      knowledgeBaseId: z.string(),
      enabled: z.boolean().optional(),
    }),
  )
  .mutation(async ({ input, ctx }) => {
    return ctx.agentModel.createAgentKnowledgeBase(
      input.agentId,
      input.knowledgeBaseId,
      input.enabled,
    );
  })
```

**Model** (`packages/database/src/models/agent.ts`):
```typescript
createAgentKnowledgeBase = async (
  agentId: string,
  knowledgeBaseId: string,
  enabled: boolean = true,
) => {
  return this.db.insert(agentsKnowledgeBases).values({
    agentId,
    enabled,
    knowledgeBaseId,
    userId: this.userId,
  });
};
```

**SQL**:
```sql
INSERT INTO agents_knowledge_bases (agent_id, knowledge_base_id, user_id, enabled)
VALUES (?, ?, ?, ?)
```

### 2. Remove Knowledge Base from Agent

**User Action**: User clicks "Remove" button on an assigned knowledge base

**Flow**:

```
User clicks Remove
        ↓
Action.tsx: removeKnowledge()
        ↓
Store: removeKnowledgeBaseFromAgent(knowledgeBaseId)
        ↓
Service: agentService.deleteAgentKnowledgeBase()
        ↓
TRPC: agent.deleteAgentKnowledgeBase.mutate()
        ↓
Model: agentModel.deleteAgentKnowledgeBase()
        ↓
Database: DELETE FROM agents_knowledge_bases
        ↓
Refresh: internal_refreshAgentConfig()
        ↓
Refresh: internal_refreshAgentKnowledge()
        ↓
UI Updates
```

**Model** (`packages/database/src/models/agent.ts`):
```typescript
deleteAgentKnowledgeBase = async (agentId: string, knowledgeBaseId: string) => {
  return this.db
    .delete(agentsKnowledgeBases)
    .where(
      and(
        eq(agentsKnowledgeBases.agentId, agentId),
        eq(agentsKnowledgeBases.knowledgeBaseId, knowledgeBaseId),
        eq(agentsKnowledgeBases.userId, this.userId),
      ),
    );
};
```

**SQL**:
```sql
DELETE FROM agents_knowledge_bases
WHERE agent_id = ? 
  AND knowledge_base_id = ? 
  AND user_id = ?
```

### 3. Toggle Knowledge Base (Enable/Disable)

**User Action**: User toggles the enable/disable switch

**Flow**:

```
User toggles switch
        ↓
Store: toggleKnowledgeBase(id, enabled)
        ↓
Service: agentService.toggleKnowledgeBase()
        ↓
TRPC: agent.toggleKnowledgeBase.mutate()
        ↓
Model: agentModel.toggleKnowledgeBase()
        ↓
Database: UPDATE agents_knowledge_bases SET enabled = ?
        ↓
Refresh: internal_refreshAgentConfig()
        ↓
UI Updates
```

**Model** (`packages/database/src/models/agent.ts`):
```typescript
toggleKnowledgeBase = async (agentId: string, knowledgeBaseId: string, enabled?: boolean) => {
  return this.db
    .update(agentsKnowledgeBases)
    .set({ enabled })
    .where(
      and(
        eq(agentsKnowledgeBases.agentId, agentId),
        eq(agentsKnowledgeBases.knowledgeBaseId, knowledgeBaseId),
        eq(agentsKnowledgeBases.userId, this.userId),
      ),
    );
};
```

**SQL**:
```sql
UPDATE agents_knowledge_bases
SET enabled = ?
WHERE agent_id = ? 
  AND knowledge_base_id = ? 
  AND user_id = ?
```

### 4. Assign Files to Agent

**User Action**: User clicks "Assign" button on a file

**Model** (`packages/database/src/models/agent.ts`):
```typescript
createAgentFiles = async (agentId: string, fileIds: string[], enabled: boolean = true) => {
  // Step 1: Check for existing assignments
  const existingFiles = await this.db
    .select({ id: agentsFiles.fileId })
    .from(agentsFiles)
    .where(
      and(
        eq(agentsFiles.agentId, agentId),
        eq(agentsFiles.userId, this.userId),
        inArray(agentsFiles.fileId, fileIds),
      ),
    );

  // Step 2: Filter out duplicates
  const existingFilesIds = new Set(existingFiles.map((item) => item.id));
  const needToInsertFileIds = fileIds.filter((fileId) => !existingFilesIds.has(fileId));

  if (needToInsertFileIds.length === 0) return;

  // Step 3: Insert new assignments
  return this.db
    .insert(agentsFiles)
    .values(
      needToInsertFileIds.map((fileId) => ({ 
        agentId, 
        enabled, 
        fileId, 
        userId: this.userId 
      })),
    );
};
```

**Key Points**:
- Prevents duplicate assignments
- Supports batch assignment of multiple files
- Returns early if all files already assigned

## State Management

### SWR Cache Keys

```typescript
const FETCH_AGENT_CONFIG_KEY = 'FETCH_AGENT_CONFIG';
const FETCH_AGENT_KNOWLEDGE_KEY = 'FETCH_AGENT_KNOWLEDGE';
```

**Cache Structure**:
```typescript
[FETCH_AGENT_KNOWLEDGE_KEY, agentId] → KnowledgeItem[]
```

### Refresh Strategy

After any mutation (assign/remove/toggle), the system refreshes:

1. **Agent Config**: `internal_refreshAgentConfig(agentId)`
   - Revalidates agent configuration
   - Updates agent metadata

2. **Agent Knowledge**: `internal_refreshAgentKnowledge()`
   - Revalidates knowledge list
   - Updates UI immediately

**Implementation**:
```typescript
await agentService.createAgentKnowledgeBase(activeAgentId, knowledgeBaseId, true);
await internal_refreshAgentConfig(get().activeId);
await internal_refreshAgentKnowledge();
```

### Optimistic Updates

The system does NOT use optimistic updates for knowledge assignments because:
- Database validation is required (foreign key constraints)
- Enable/disable state must be accurate
- Error handling is important for user feedback

## Performance Optimizations

### 1. Virtualized List

**File**: `src/features/KnowledgeBaseModal/AssignKnowledgeBase/List.tsx`

```typescript
<Virtuoso
  itemContent={(index) => {
    const item = data![index];
    return <Item key={item.id} {...item} />;
  }}
  overscan={400}
  style={{ height: 500 }}
  totalCount={data!.length}
/>
```

**Benefits**:
- Renders only visible items
- Handles thousands of items smoothly
- Minimal memory footprint

### 2. SWR Caching

- Automatic cache invalidation
- Revalidation on focus/reconnect
- Prevents unnecessary refetches
- Shared cache across components

### 3. Batch Operations

Files can be assigned in batches:

```typescript
await addFilesToAgent([fileId1, fileId2, fileId3], true);
```

Single database transaction for multiple files.

### 4. Filtered Queries

**Exclude files in knowledge bases**:
```typescript
const files = await ctx.fileModel.query({
  showFilesInKnowledgeBase: false,
});
```

Reduces data transfer and processing.

**Exclude image files**:
```typescript
.filter((file) => !file.fileType.startsWith('image'))
```

Only shows relevant file types.

## Error Handling

### Database Errors

**Foreign Key Violations**:
- Agent doesn't exist
- Knowledge base doesn't exist
- File doesn't exist

**Duplicate Key Violations**:
- Knowledge base already assigned
- File already assigned

**Solution**: Model methods check for existing assignments before inserting.

### UI Error States

```typescript
const { isLoading, error, data } = useFetchFilesAndKnowledgeBases();

if (error) {
  return (
    <Center>
      <Icon icon={ServerCrash} />
      {t('networkError')}
    </Center>
  );
}
```

### Network Errors

- SWR automatically retries failed requests
- Error boundary catches unhandled errors
- User-friendly error messages

## Security Considerations

### User Isolation

All queries include `userId` filter:

```typescript
eq(agentsKnowledgeBases.userId, this.userId)
```

Prevents cross-user data access.

### Authorization

- TRPC `authedProcedure` ensures user is authenticated
- Database models initialized with `userId`
- Cascade delete ensures cleanup

### Input Validation

Zod schemas validate all inputs:

```typescript
.input(
  z.object({
    agentId: z.string(),
    knowledgeBaseId: z.string(),
    enabled: z.boolean().optional(),
  }),
)
```

## Use Cases

### 1. Customer Support Agent

**Scenario**: Create an agent with access to product documentation

**Steps**:
1. Create knowledge base with product docs
2. Assign knowledge base to agent
3. Agent can now reference docs in responses

### 2. Code Assistant Agent

**Scenario**: Agent with access to specific code files

**Steps**:
1. Upload code files
2. Assign files to agent
3. Agent can analyze and reference code

### 3. Temporary Knowledge Access

**Scenario**: Temporarily disable knowledge base without removing

**Steps**:
1. Toggle knowledge base to disabled
2. Agent stops using that knowledge
3. Re-enable when needed

## Testing

### Unit Tests

Test database model methods:

```typescript
describe('AgentModel.createAgentKnowledgeBase', () => {
  it('should assign knowledge base to agent', async () => {
    await model.createAgentKnowledgeBase(agentId, kbId, true);
    const knowledge = await model.getAgentAssignedKnowledge(agentId);
    expect(knowledge.knowledgeBases).toContainEqual(
      expect.objectContaining({ id: kbId, enabled: true })
    );
  });

  it('should prevent duplicate assignments', async () => {
    await model.createAgentKnowledgeBase(agentId, kbId, true);
    await expect(
      model.createAgentKnowledgeBase(agentId, kbId, true)
    ).rejects.toThrow();
  });
});
```

### Integration Tests

Test complete flow:

```typescript
describe('Knowledge Assignment Flow', () => {
  it('should assign and retrieve knowledge', async () => {
    const { addKnowledgeBaseToAgent, useFetchFilesAndKnowledgeBases } = 
      useAgentStore.getState();
    
    await addKnowledgeBaseToAgent(kbId);
    
    const { data } = useFetchFilesAndKnowledgeBases();
    const kb = data.find(item => item.id === kbId);
    
    expect(kb.enabled).toBe(true);
  });
});
```

## Troubleshooting

### Knowledge Not Appearing

**Issue**: Assigned knowledge doesn't show in list

**Solutions**:
1. Check if cache is stale: Force refresh
2. Verify database record exists
3. Check user permissions
4. Ensure agent ID is correct

### Enable/Disable Not Working

**Issue**: Toggle doesn't update state

**Solutions**:
1. Check network requests in DevTools
2. Verify TRPC endpoint is called
3. Check database update query
4. Force cache revalidation

### Files Missing from List

**Issue**: Some files don't appear

**Solutions**:
1. Check if files are in knowledge bases (excluded)
2. Verify file type (images excluded)
3. Check user ownership
4. Verify file upload completed

## Best Practices

### For Developers

1. **Always refresh after mutations**: Ensure UI stays in sync
2. **Use SWR hooks**: Leverage automatic caching and revalidation
3. **Validate inputs**: Use Zod schemas for type safety
4. **Handle errors gracefully**: Show user-friendly messages
5. **Test edge cases**: Duplicate assignments, missing data, etc.

### For Users

1. **Organize knowledge bases**: Group related files
2. **Use descriptive names**: Easy to identify knowledge sources
3. **Disable instead of remove**: Preserve assignments for later
4. **Regular cleanup**: Remove unused assignments
5. **Test agent responses**: Verify knowledge is being used

## Conclusion

The agent knowledge assignment system provides a flexible and powerful way to give agents access to specific knowledge sources. The architecture is designed for performance, security, and maintainability, with clear separation of concerns across UI, state management, services, and database layers. The use of junction tables allows for many-to-many relationships with additional metadata (enabled flag), and the SWR-based caching ensures a responsive user experience.

The `getKnowledgeBasesAndFiles` endpoint is the cornerstone of this system, efficiently combining all available knowledge sources with their assignment status for a specific agent, enabling users to easily manage which knowledge each agent can access.
