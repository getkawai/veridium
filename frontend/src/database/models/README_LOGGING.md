# 🔍 Database Models - Logging Infrastructure

## Ringkasan Implementasi

Semua 19 database models di direktori ini telah dilengkapi dengan comprehensive logging infrastructure menggunakan `ModelLogger` utility. Logging ini memungkinkan detailed debugging dan performance monitoring untuk semua operasi database.

## ✅ Status: COMPLETE

**Total Models**: 22 models  
**Logger Implementation**: ✅ 100% Complete  
**Date Completed**: November 4, 2025

## 📋 Daftar Models dengan Logging

### Core Models (Session & Messages)
| File | Model Class | Logger Context | Status |
|------|-------------|----------------|--------|
| `session.ts` | SessionModel | `Model:Session` | ✅ Full logging |
| `message.ts` | MessageModel | `Model:Message` | ✅ Full logging |
| `thread.ts` | ThreadModel | `Model:Thread` | ✅ Full logging |
| `topic.ts` | TopicModel | `Model:Topic` | ✅ Full logging |

### User & Authentication
| File | Model Class | Logger Context | Status |
|------|-------------|----------------|--------|
| `user.ts` | UserModel | `Model:User` | ✅ Logged |
| `apiKey.ts` | ApiKeyModel | `Model:ApiKey` | ✅ Logged |
| `oauthHandoff.ts` | OAuthHandoffModel | `Model:OAuthHandoff` | ✅ Logged |

### File & Content Management
| File | Model Class | Logger Context | Status |
|------|-------------|----------------|--------|
| `file.ts` | FileModel | `Model:File` | ✅ Logged |
| `document.ts` | DocumentModel | `Model:Document` | ✅ Logged |
| `chunk.ts` | ChunkModel | `Model:Chunk` | ✅ Logged |
| `embedding.ts` | EmbeddingModel | `Model:Embedding` | ✅ Logged |

### AI & Agent Infrastructure
| File | Model Class | Logger Context | Status |
|------|-------------|----------------|--------|
| `agent.ts` | AgentModel | `Model:Agent` | ✅ Logged |
| `aiModel.ts` | AiModelModel | `Model:AiModel` | ✅ Logged |
| `aiProvider.ts` | AiProviderModel | `Model:AiProvider` | ✅ Logged |
| `plugin.ts` | PluginModel | `Model:Plugin` | ✅ Logged |

### Knowledge & Collaboration
| File | Model Class | Logger Context | Status |
|------|-------------|----------------|--------|
| `knowledgeBase.ts` | KnowledgeBaseModel | `Model:KnowledgeBase` | ✅ Logged |
| `chatGroup.ts` | ChatGroupModel | `Model:ChatGroup` | ✅ Logged |
| `sessionGroup.ts` | SessionGroupModel | `Model:SessionGroup` | ✅ Logged |

### Image Generation
| File | Model Class | Logger Context | Status |
|------|-------------|----------------|--------|
| `generation.ts` | GenerationModel | `Model:Generation` | ✅ Logged |
| `generationBatch.ts` | GenerationBatchModel | `Model:GenerationBatch` | ✅ Logged |
| `generationTopic.ts` | GenerationTopicModel | `Model:GenerationTopic` | ✅ Logged |

### System Infrastructure
| File | Model Class | Logger Context | Status |
|------|-------------|----------------|--------|
| `asyncTask.ts` | AsyncTaskModel | `Model:AsyncTask` | ✅ Logged |

## 📖 Cara Menggunakan Logging

### 1. Logger Sudah Tersedia di Setiap Model

```typescript
// Setiap model sudah punya logger yang siap digunakan
export class SessionModel {
  private logger = createModelLogger('Session', 'SessionModel', 'database/models/session');
  
  // Logger bisa digunakan di method manapun
}
```

### 2. Method Entry/Exit Logging

```typescript
async query(params) {
  // Log method entry dengan params
  await this.logger.methodEntry('query', { userId: this.userId, ...params });
  
  // ... execute query ...
  const result = await DB.Query(...);
  
  // Log method exit dengan result dan duration (automatic)
  await this.logger.methodExit('query', result);
  return result;
}
```

### 3. Debug, Warning, dan Error Logging

```typescript
// Debug logging untuk informasi detail
await this.logger.debug('create', 'Creating new session', { sessionId, userId });

// Warning untuk kondisi yang tidak normal tapi bukan error
await this.logger.warn('findById', 'Session not found', { sessionId });

// Error logging dengan full context
await this.logger.error('delete', 'Failed to delete session', { 
  sessionId, 
  error: e.message 
});
```

## 🔎 Melihat dan Filtering Logs

### Filter by Model
```bash
# Lihat semua logs dari SessionModel
grep "[Model:Session]" logs.txt

# Lihat semua logs dari MessageModel  
grep "[Model:Message]" logs.txt
```

### Filter by Method
```bash
# Lihat semua calls ke SessionModel.query
grep "fullMethod=SessionModel.query" logs.txt

# Lihat semua calls ke MessageModel.create
grep "fullMethod=MessageModel.create" logs.txt
```

### Filter by User/Session/Topic
```bash
# Logs untuk specific user
grep "userId=user123" logs.txt

# Logs untuk specific session
grep "sessionId=sess456" logs.txt

# Logs untuk specific topic
grep "topicId=topic789" logs.txt
```

### Filter by Duration (Performance Analysis)
```bash
# Find slow operations (>1000ms)
grep "duration=" logs.txt | awk -F'duration=' '{print $2}' | awk -F'ms' '{if($1>1000)print}'
```

## 📊 Log Format Example

```
[Model:Session] → ENTER: SessionModel.query 
  timestamp=2025-11-04T10:30:00.123Z 
  class=SessionModel 
  method=query 
  fullMethod=SessionModel.query 
  path=database/models/session 
  params={"userId":"user123","sessionType":1} 
  userId=user123

[Model:Session] ← EXIT: SessionModel.query 
  timestamp=2025-11-04T10:30:00.138Z 
  class=SessionModel 
  method=query 
  fullMethod=SessionModel.query 
  path=database/models/session 
  duration=15.23ms 
  result={"count":5,"sessions":[...]} 
  count=5
```

## 🎯 Best Practices

### 1. Log Critical Methods
Tambahkan `methodEntry`/`methodExit` untuk:
- CRUD operations (create, update, delete)
- Complex queries
- Methods dengan external dependencies
- Transaction-heavy operations

### 2. Use Appropriate Log Levels
```typescript
// ✅ Good
await this.logger.debug('method', 'Normal operation details');
await this.logger.warn('method', 'Unusual but handled condition');  
await this.logger.error('method', 'Critical error', { error });

// ❌ Bad
await this.logger.error('method', 'Info message'); // Wrong level
```

### 3. Include Relevant Context
```typescript
// ✅ Good - includes userId, sessionId, and error details
await this.logger.error('delete', 'Delete failed', {
  userId: this.userId,
  sessionId,
  error: e.message,
  stack: e.stack
});

// ❌ Bad - missing context
await this.logger.error('delete', 'Error');
```

### 4. Smart Serialization
Logger automatically:
- Truncates large objects
- Limits array sizes
- Extracts common fields (userId, sessionId, etc.)
- Handles circular references

## 📚 Related Documentation

- **Logger Utility**: `../utils/logger.md` - Detailed API documentation
- **Usage Examples**: `../utils/logger.example.md` - Practical examples
- **Output Examples**: `../utils/logger-output-example.md` - Log format examples

## 🔧 Configuration

Logger dapat dikonfigurasi dengan options:

```typescript
const logger = createModelLogger('Session', 'SessionModel', 'database/models/session', {
  includeStack: true,  // Include call stack (default: false)
  maxDepth: 5          // Max object serialization depth (default: 3)
});
```

## 🚀 Performance Impact

- **Minimal overhead**: Async logging tidak block execution
- **Smart serialization**: Automatic truncation untuk large objects
- **Conditional logging**: Bisa di-disable via log level configuration
- **Production-ready**: Designed untuk production use

## 📈 Next Steps

1. **Monitoring Dashboard**: Aggregate logs untuk visualisasi
2. **Alert System**: Setup alerts untuk error patterns
3. **Performance Analytics**: Analyze duration metrics untuk optimization
4. **Custom Metrics**: Extract business metrics dari logs

---

**Last Updated**: November 4, 2025  
**Maintainer**: Development Team  
**Status**: ✅ Production Ready



