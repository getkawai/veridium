# Contoh Penggunaan Logger

## Contoh 1: Basic Logging

```typescript
import { createModelLogger } from '@/utils/logger';

export class UserModel {
  private logger = createModelLogger('User');
  
  async findById(id: string) {
    await this.logger.methodEntry('findById', { id });
    
    const user = await DB.GetUser({ id });
    await this.logger.debug(`User found: ${user.name}`);
    
    await this.logger.methodExit('findById', { userId: user.id });
    return user;
  }
}
```

**Output:**
```
[Model:User] → ENTER: findById
  timestamp: 2025-11-04T10:30:45.123Z
  method: findById
  params: {
    "id": "user123"
  }
  id: user123

[Model:User] User found: John Doe

[Model:User] ← EXIT: findById
  timestamp: 2025-11-04T10:30:45.156Z
  method: findById
  duration: 33.45ms
  result: {
    "userId": "user123"
  }
  id: user123
```

## Contoh 2: Logging dengan Error Handling

```typescript
export class SessionModel {
  private logger = createModelLogger('Session');
  
  async delete(id: string) {
    await this.logger.methodEntry('delete', { id, userId: this.userId });
    
    try {
      const agents = await DB.GetSessionAgents({
        sessionId: id,
        userId: this.userId,
      });
      
      await this.logger.debug(`Found ${agents.length} agents to delete`);
      
      for (const agent of agents) {
        await DB.DeleteAgent({ id: agent.id, userId: this.userId });
        await this.logger.debug(`Deleted agent: ${agent.id}`);
      }
      
      await DB.DeleteSession({ id, userId: this.userId });
      
      await this.logger.methodExit('delete', { 
        deletedAgents: agents.length,
        success: true 
      });
      
    } catch (error) {
      await this.logger.methodError('delete', error, { id, userId: this.userId });
      throw error;
    }
  }
}
```

**Output (Success):**
```
[Model:Session] → ENTER: delete
  timestamp: 2025-11-04T10:30:45.123Z
  method: delete
  params: {
    "id": "sess123",
    "userId": "user456"
  }
  userId: user456
  id: sess123

[Model:Session] Found 2 agents to delete

[Model:Session] Deleted agent: agent789

[Model:Session] Deleted agent: agent012

[Model:Session] ← EXIT: delete
  timestamp: 2025-11-04T10:30:45.456Z
  method: delete
  duration: 333.00ms
  result: {
    "deletedAgents": 2,
    "success": true
  }
  count: 2
```

**Output (Error):**
```
[Model:Session] → ENTER: delete
  ...

[Model:Session] Found 2 agents to delete

[Model:Session] ✗ delete failed
  error: Database constraint violation
  stack: Error: Database constraint violation
    at DB.DeleteAgent (...)
    ...
  data: {
    "id": "sess123",
    "userId": "user456"
  }
```

## Contoh 3: Advanced Logging dengan Stack Trace

```typescript
export class MessageModel {
  // Enable stack trace untuk debugging mendalam
  private logger = createModelLogger('Message', {
    includeStack: true,
    maxDepth: 5,
  });
  
  async create(params: CreateMessageParams) {
    await this.logger.methodEntry('create', {
      role: params.role,
      hasFiles: !!params.files,
      userId: this.userId,
    });
    
    const message = await DB.CreateMessage({
      id: nanoid(),
      ...params,
      userId: this.userId,
    });
    
    await this.logger.info('Message created successfully', {
      messageId: message.id,
      role: message.role,
    });
    
    await this.logger.methodExit('create', { messageId: message.id });
    return message;
  }
}
```

**Output:**
```
[Model:Message] → ENTER: create
  timestamp: 2025-11-04T10:30:45.123Z
  method: create
  params: {
    "role": "user",
    "hasFiles": true,
    "userId": "user456"
  }
  userId: user456
  caller: at MessageModel.create (message.ts:567:25)

[Model:Message] Message created successfully
  data: {
    "messageId": "msg789",
    "role": "user"
  }

[Model:Message] ← EXIT: create
  timestamp: 2025-11-04T10:30:45.234Z
  method: create
  duration: 111.00ms
  result: {
    "messageId": "msg789"
  }
  id: msg789
```

## Contoh 4: Logging Complex Operations

```typescript
export class SessionModel {
  private logger = createModelLogger('Session');
  
  async query({ current = 0, pageSize = 9999 } = {}) {
    await this.logger.methodEntry('query', { 
      current, 
      pageSize, 
      userId: this.userId 
    });
    
    const offset = current * pageSize;
    
    const sessions = await DB.ListSessions({
      userId: this.userId,
      limit: pageSize,
      offset,
    });
    
    await this.logger.debug(`Retrieved ${sessions.length} sessions from DB`);
    
    const filtered = sessions.filter(
      (s) => s.slug !== INBOX_SESSION_ID,
    );
    
    await this.logger.debug(`Filtered to ${filtered.length} sessions (excluding inbox)`);
    
    // Enrich with agents
    const enriched = await Promise.all(
      filtered.map(async (session) => {
        const agents = await DB.GetSessionAgents({
          sessionId: session.id,
          userId: this.userId,
        });
        return { ...session, agents };
      }),
    );
    
    await this.logger.debug('Enrichment completed', {
      totalAgentsLoaded: enriched.reduce((sum, s) => sum + s.agents.length, 0),
    });
    
    await this.logger.methodExit('query', { count: enriched.length });
    return enriched;
  }
}
```

**Output:**
```
[Model:Session] → ENTER: query
  timestamp: 2025-11-04T10:30:45.123Z
  method: query
  params: {
    "current": 0,
    "pageSize": 9999,
    "userId": "user456"
  }
  userId: user456

[Model:Session] Retrieved 15 sessions from DB

[Model:Session] Filtered to 14 sessions (excluding inbox)

[Model:Session] Enrichment completed
  data: {
    "totalAgentsLoaded": 28
  }

[Model:Session] ← EXIT: query
  timestamp: 2025-11-04T10:30:45.678Z
  method: query
  duration: 555.00ms
  result: {
    "count": 14
  }
  count: 14
```

## Contoh 5: Logging Array dengan Smart Serialization

```typescript
export class FileModel {
  private logger = createModelLogger('File');
  
  async batchUpload(files: File[]) {
    await this.logger.methodEntry('batchUpload', {
      filesCount: files.length,
      totalSize: files.reduce((sum, f) => sum + f.size, 0),
      userId: this.userId,
    });
    
    const results = await Promise.all(
      files.map(async (file) => {
        const uploaded = await this.uploadFile(file);
        await this.logger.debug(`Uploaded: ${file.name}`);
        return uploaded;
      })
    );
    
    await this.logger.methodExit('batchUpload', {
      uploadedCount: results.length,
      successRate: `${(results.filter(r => r.success).length / results.length * 100).toFixed(1)}%`,
    });
    
    return results;
  }
}
```

**Output:**
```
[Model:File] → ENTER: batchUpload
  timestamp: 2025-11-04T10:30:45.123Z
  method: batchUpload
  params: {
    "filesCount": 5,
    "totalSize": 15728640,
    "userId": "user456"
  }
  userId: user456

[Model:File] Uploaded: document1.pdf
[Model:File] Uploaded: image1.png
[Model:File] Uploaded: data.csv
[Model:File] Uploaded: report.xlsx
[Model:File] Uploaded: presentation.pptx

[Model:File] ← EXIT: batchUpload
  timestamp: 2025-11-04T10:30:47.890Z
  method: batchUpload
  duration: 2767.00ms
  result: {
    "uploadedCount": 5,
    "successRate": "100.0%"
  }
  count: 5
```

## Filtering Logs

### Filter by User
```bash
grep "userId: user456" app.log
```

### Filter by Method
```bash
grep "method: query" app.log
```

### Filter by Duration (slow queries)
```bash
grep "duration:" app.log | awk -F': ' '$2 > 1000'
```

### Filter Errors Only
```bash
grep "✗" app.log
```

### Combined Filters
```bash
grep "Model:Session" app.log | grep "userId: user456" | grep "→ ENTER"
```

## Performance Analysis

Track method performance:

```bash
# Extract all durations
grep "duration:" app.log | awk '{print $NF}' | sort -n

# Average duration for specific method
grep "method: query" app.log | grep "duration:" | awk '{sum += $NF} END {print sum/NR}'
```

