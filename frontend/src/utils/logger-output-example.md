# Contoh Output Log yang Sebenarnya

## Contoh 1: SessionModel.query

**Code:**
```typescript
export class SessionModel {
  private logger = createModelLogger('Session', 'SessionModel', 'database/models/session');
  
  async query({ current = 0, pageSize = 9999 } = {}) {
    await this.logger.methodEntry('query', { current, pageSize, userId: this.userId });
    
    const sessions = await DB.ListSessions({ userId: this.userId, limit: pageSize });
    await this.logger.debug(`Retrieved ${sessions.length} sessions from DB`);
    
    const filtered = sessions.filter(s => s.slug !== INBOX_SESSION_ID);
    await this.logger.debug(`Filtered to ${filtered.length} sessions (excluding inbox)`);
    
    await this.logger.methodExit('query', { count: filtered.length });
    return filtered;
  }
}
```

**Output:**
```
[Model:Session] → ENTER: SessionModel.query
  timestamp: 2025-11-04T10:30:45.123Z
  class: SessionModel
  method: query
  fullMethod: SessionModel.query
  path: database/models/session
  params: {
    "current": 0,
    "pageSize": 9999,
    "userId": "user_abc123"
  }
  userId: user_abc123

[Model:Session] Retrieved 15 sessions from DB

[Model:Session] Filtered to 14 sessions (excluding inbox)

[Model:Session] ← EXIT: SessionModel.query
  timestamp: 2025-11-04T10:30:45.456Z
  class: SessionModel
  method: query
  fullMethod: SessionModel.query
  path: database/models/session
  duration: 333.45ms
  result: {
    "count": 14
  }
  count: 14
```

## Contoh 2: MessageModel.create dengan File Upload

**Code:**
```typescript
export class MessageModel {
  private logger = createModelLogger('Message', 'MessageModel', 'database/models/message');
  
  async create(params: CreateMessageParams, id: string = nanoid(14)) {
    await this.logger.methodEntry('create', { 
      id, 
      role: params.role, 
      hasPlugin: !!params.plugin,
      filesCount: params.files?.length || 0,
      userId: this.userId 
    });
    
    const message = await DB.CreateMessage({ id, ...params });
    await this.logger.methodExit('create', { messageId: message.id });
    return message;
  }
}
```

**Output:**
```
[Model:Message] → ENTER: MessageModel.create
  timestamp: 2025-11-04T10:35:12.789Z
  class: MessageModel
  method: create
  fullMethod: MessageModel.create
  path: database/models/message
  params: {
    "id": "msg_xyz789",
    "role": "user",
    "hasPlugin": false,
    "filesCount": 3,
    "userId": "user_abc123"
  }
  userId: user_abc123
  id: msg_xyz789

[Model:Message] ← EXIT: MessageModel.create
  timestamp: 2025-11-04T10:35:13.012Z
  class: MessageModel
  method: create
  fullMethod: MessageModel.create
  path: database/models/message
  duration: 223.45ms
  result: {
    "messageId": "msg_xyz789"
  }
  id: msg_xyz789
```

## Contoh 3: ThreadModel.delete dengan Error

**Code:**
```typescript
export class ThreadModel {
  private logger = createModelLogger('Thread', 'ThreadModel', 'database/models/thread');
  
  async delete(id: string) {
    await this.logger.methodEntry('delete', { id, userId: this.userId });
    
    try {
      await DB.DeleteThread({ id, userId: this.userId });
      await this.logger.methodExit('delete', { id });
    } catch (error) {
      await this.logger.methodError('delete', error, { id });
      throw error;
    }
  }
}
```

**Output (Success):**
```
[Model:Thread] → ENTER: ThreadModel.delete
  timestamp: 2025-11-04T10:40:25.123Z
  class: ThreadModel
  method: delete
  fullMethod: ThreadModel.delete
  path: database/models/thread
  params: {
    "id": "thread_def456",
    "userId": "user_abc123"
  }
  userId: user_abc123
  id: thread_def456

[Model:Thread] ← EXIT: ThreadModel.delete
  timestamp: 2025-11-04T10:40:25.234Z
  class: ThreadModel
  method: delete
  fullMethod: ThreadModel.delete
  path: database/models/thread
  duration: 111.00ms
  result: {
    "id": "thread_def456"
  }
  id: thread_def456
```

**Output (Error):**
```
[Model:Thread] → ENTER: ThreadModel.delete
  timestamp: 2025-11-04T10:40:25.123Z
  class: ThreadModel
  method: delete
  fullMethod: ThreadModel.delete
  path: database/models/thread
  params: {
    "id": "thread_def456",
    "userId": "user_abc123"
  }
  userId: user_abc123
  id: thread_def456

[Model:Thread] ✗ ThreadModel.delete failed
  error: Foreign key constraint violation
  stack: Error: Foreign key constraint violation
    at ThreadModel.delete (database/models/thread.ts:52)
    at async SessionService.deleteThread (services/session.ts:123)
    ...
  data: {
    "id": "thread_def456"
  }
```

## Contoh 4: Nested Operations (Session dengan Agent)

**Output:**
```
[Model:Session] → ENTER: SessionModel.delete
  timestamp: 2025-11-04T11:00:00.000Z
  class: SessionModel
  method: delete
  fullMethod: SessionModel.delete
  path: database/models/session
  params: {
    "id": "sess_ghi789",
    "userId": "user_abc123"
  }
  userId: user_abc123
  id: sess_ghi789

[Model:Session] Found 2 agents linked to session sess_ghi789

[Model:Session] Unlinked agent agent_001 from session sess_ghi789

[Model:Session] Unlinked agent agent_002 from session sess_ghi789

[Model:Session] Deleted session sess_ghi789

[Model:Session] Deleted orphaned agent agent_001

[Model:Session] Deleted orphaned agent agent_002

[Model:Session] ← EXIT: SessionModel.delete
  timestamp: 2025-11-04T11:00:00.567Z
  class: SessionModel
  method: delete
  fullMethod: SessionModel.delete
  path: database/models/session
  duration: 567.00ms
  result: {
    "id": "sess_ghi789",
    "deletedAgents": 2
  }
  id: sess_ghi789
  count: 2
```

## Filtering Examples

### 1. Semua operasi SessionModel
```bash
grep "class: SessionModel" app.log
```

**Output:**
```
[Model:Session] → ENTER: SessionModel.query
  class: SessionModel
[Model:Session] ← EXIT: SessionModel.query
  class: SessionModel
[Model:Session] → ENTER: SessionModel.create
  class: SessionModel
[Model:Session] ← EXIT: SessionModel.create
  class: SessionModel
[Model:Session] → ENTER: SessionModel.delete
  class: SessionModel
[Model:Session] ← EXIT: SessionModel.delete
  class: SessionModel
```

### 2. Semua operasi query di berbagai model
```bash
grep "method: query" app.log
```

**Output:**
```
[Model:Session] → ENTER: SessionModel.query
  method: query
[Model:Session] ← EXIT: SessionModel.query
  method: query
[Model:Message] → ENTER: MessageModel.query
  method: query
[Model:Message] ← EXIT: MessageModel.query
  method: query
```

### 3. Track specific user activity
```bash
grep "userId: user_abc123" app.log
```

**Output:**
```
[Model:Session] → ENTER: SessionModel.query
  userId: user_abc123
[Model:Message] → ENTER: MessageModel.create
  userId: user_abc123
[Model:Thread] → ENTER: ThreadModel.delete
  userId: user_abc123
```

### 4. Slow queries (>500ms)
```bash
grep "duration:" app.log | awk -F': ' '{
  duration=$2;
  gsub(/ms/, "", duration);
  if(duration > 500) print $0
}'
```

**Output:**
```
[Model:Session] ← EXIT: SessionModel.delete
  duration: 567.00ms
[Model:Message] ← EXIT: MessageModel.query
  duration: 1234.56ms
```

### 5. Specific file path operations
```bash
grep "path: database/models/session" app.log
```

**Output:**
```
[Model:Session] → ENTER: SessionModel.query
  path: database/models/session
[Model:Session] ← EXIT: SessionModel.query
  path: database/models/session
[Model:Session] → ENTER: SessionModel.create
  path: database/models/session
```

### 6. Full method signature tracking
```bash
grep "fullMethod: SessionModel.delete" app.log
```

**Output:**
```
[Model:Session] → ENTER: SessionModel.delete
  fullMethod: SessionModel.delete
[Model:Session] ← EXIT: SessionModel.delete
  fullMethod: SessionModel.delete
```

### 7. All errors
```bash
grep "✗" app.log
```

**Output:**
```
[Model:Thread] ✗ ThreadModel.delete failed
  error: Foreign key constraint violation
[Model:Session] ✗ SessionModel.update failed
  error: Session not found
```

### 8. Performance analysis - average duration per method
```bash
grep "fullMethod: SessionModel.query" app.log | \
grep "duration:" | \
awk '{print $NF}' | \
sed 's/ms//' | \
awk '{sum+=$1; count++} END {print "Average:", sum/count, "ms"}'
```

**Output:**
```
Average: 345.67 ms
```

## Tips Filtering

1. **Combine multiple filters:**
   ```bash
   grep "SessionModel.query" app.log | grep "userId: user_abc123"
   ```

2. **Time-based filtering:**
   ```bash
   grep "2025-11-04T10:30" app.log
   ```

3. **Count occurrences:**
   ```bash
   grep "SessionModel.query" app.log | wc -l
   ```

4. **Show context (5 lines before and after):**
   ```bash
   grep -C 5 "✗" app.log
   ```

5. **Export to file:**
   ```bash
   grep "userId: user_abc123" app.log > user_activity.log
   ```

