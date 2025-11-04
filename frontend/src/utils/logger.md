# Database Model Logger

Utility untuk debugging database models menggunakan Wails log service.

## Penggunaan Dasar

### 1. Import Logger

```typescript
import { createModelLogger, timeOperation } from '@/utils/logger';
```

### 2. Inisialisasi Logger di Model

#### Basic (Recommended)

```typescript
export class SessionModel {
  private userId: string;
  // Include class name and file path untuk detail logging
  private logger = createModelLogger(
    'Session',              // Model name
    'SessionModel',         // Class name
    'database/models/session'  // File path
  );

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }
}
```

#### Dengan Opsi Advanced

```typescript
export class SessionModel {
  private userId: string;
  // Enable stack trace dan custom max depth
  private logger = createModelLogger(
    'Session',
    'SessionModel',
    'database/models/session',
    {
      includeStack: true,  // Include call stack in logs
      maxDepth: 5,         // Max object nesting depth
    }
  );

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }
}
```

#### Shorthand (tanpa class name & path)

```typescript
export class UserModel {
  private logger = createModelLogger('User');  // Auto: UserModel class
}
```

### 3. Log Method Entry dan Exit

```typescript
async findByIdOrSlug(idOrSlug: string) {
  await this.logger.methodEntry('findByIdOrSlug', { idOrSlug, userId: this.userId });
  
  // Your code here...
  
  await this.logger.methodExit('findByIdOrSlug', { found: !!result });
  return result;
}
```

### 4. Log Errors

```typescript
try {
  const session = await DB.GetSessionByIdOrSlug({
    id: idOrSlug,
    slug: idOrSlug,
    userId: this.userId,
  });
} catch (error) {
  await this.logger.methodError('findByIdOrSlug', error, { idOrSlug });
  return undefined;
}
```

### 5. Log Debug Information

```typescript
await this.logger.debug(`Retrieved ${sessions.length} sessions from DB`);
await this.logger.debug(`Filtered to ${filtered.length} sessions (excluding inbox)`);
```

### 6. Method Structure

Struktur umum method dengan logging:

```typescript
query = async ({ current = 0, pageSize = 9999 } = {}) => {
  await this.logger.methodEntry('query', { current, pageSize, userId: this.userId });
  
  // Your code here...
  const result = await DB.ListSessions({ userId: this.userId });
  await this.logger.debug(`Retrieved ${result.length} sessions`);
  
  await this.logger.methodExit('query', { count: result.length });
  return result;
};
```

## Log Levels

Logger mendukung 4 level logging:

### Debug
Untuk informasi detail debugging:
```typescript
await this.logger.debug('Processing batch', { batchSize: 100 });
```

### Info
Untuk informasi umum:
```typescript
await this.logger.info('Session created successfully', { sessionId });
```

### Warning
Untuk peringatan:
```typescript
await this.logger.warn('Session not found', { id });
```

### Error
Untuk errors:
```typescript
await this.logger.error('Failed to create session', error, { params });
```

## Performance Monitoring

Untuk monitoring performance, gunakan console.time dan console.timeEnd:

```typescript
query = async () => {
  console.time('query');
  await this.logger.methodEntry('query');
  
  // Your code here...
  const result = await DB.ListSessions({ userId: this.userId });
  
  await this.logger.methodExit('query', { count: result.length });
  console.timeEnd('query');
  return result;
};
```

## Format Output Log

Logs akan muncul dengan format detail:

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
    "userId": "user123"
  }
  userId: user123

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

### Detail Log Fields

**Method Entry (`→ ENTER`):**
- `timestamp` - ISO timestamp saat method dipanggil
- `class` - Class name (e.g., `SessionModel`)
- `method` - Method name (e.g., `query`)
- `fullMethod` - Full method signature (e.g., `SessionModel.query`)
- `path` - File path relative to src (e.g., `database/models/session`)
- `params` - Parameter lengkap yang diterima
- `userId`, `id`, `sessionId`, `topicId` - Individual fields untuk filtering
- `caller` - Call stack location (jika `includeStack: true`)

**Method Exit (`← EXIT`):**
- `timestamp` - ISO timestamp saat method selesai
- `class` - Class name
- `method` - Method name
- `fullMethod` - Full method signature
- `path` - File path
- `duration` - Waktu eksekusi dalam milliseconds
- `result` - Return value summary
- `count`, `id`, `arrayLength` - Individual fields untuk filtering

## Advanced Features

### Logger Options

```typescript
interface LoggerOptions {
  includeStack?: boolean;  // Include call stack (default: false)
  maxDepth?: number;       // Max object depth (default: 3)
}
```

### Automatic Duration Tracking

`methodEntry` dan `methodExit` otomatis menghitung durasi eksekusi:

```typescript
async query() {
  await this.logger.methodEntry('query', { params });
  // ... code execution ...
  await this.logger.methodExit('query', { result });
}
// Output includes: duration: 123.45ms
```

### Smart Value Serialization

Logger otomatis menyesuaikan output berdasarkan tipe data:

- **Array panjang**: `[Array(1000)]` - hanya menampilkan panjang
- **Object besar**: `{Object with 50 keys}` - hanya menampilkan jumlah keys
- **Nested objects**: Dibatasi sampai `maxDepth`
- **Circular references**: Handled dengan graceful error

### Important Fields Auto-Extraction

Logger otomatis extract field penting untuk filtering:

```typescript
await this.logger.methodEntry('query', { 
  userId: 'user123',
  sessionId: 'sess456', 
  topicId: 'topic789',
  other: 'data'
});

// Logs:
// params: { full object }
// userId: user123      <- extracted
// sessionId: sess456   <- extracted
// topicId: topic789    <- extracted
```

Ini memudahkan filtering log dengan grep:

```bash
# Filter by userId
grep "userId: user123" app.log

# Filter by specific session
grep "sessionId: sess456" app.log

# Filter by class
grep "class: SessionModel" app.log

# Filter by file path
grep "path: database/models/session" app.log

# Filter by full method
grep "fullMethod: SessionModel.query" app.log

# Combined: SessionModel.query with specific user
grep "SessionModel.query" app.log | grep "userId: user123"
```

## Tips dan Best Practices

1. **Selalu log method entry dan exit** - Automatic duration tracking & flow visualization
2. **Log informasi penting** - userId, sessionId, topicId di-extract otomatis
3. **Gunakan level logging yang tepat** - Debug untuk detail, Info untuk operasi normal
4. **Log errors dengan context** - Include parameters yang menyebabkan error
5. **Jangan log sensitive data** - Passwords, tokens, personal information, dll
6. **Log di step-step penting** - Membantu debugging flow logic
7. **Enable includeStack untuk deep debugging** - Tapi disable di production (performance)
8. **Gunakan maxDepth yang sesuai** - Balance antara detail dan readability

## Mengaktifkan/Menonaktifkan Log

Untuk mengontrol log level dari aplikasi:

```typescript
import { LogService, Level } from '@@/github.com/wailsapp/wails/v3/pkg/services/log';

// Set log level
await LogService.SetLogLevel(Level.Debug); // Show all logs
await LogService.SetLogLevel(Level.Info);  // Show info, warning, error
await LogService.SetLogLevel(Level.Warning); // Show warning, error only
await LogService.SetLogLevel(Level.Error); // Show error only

// Get current log level
const currentLevel = await LogService.LogLevel();
```

## Contoh Implementasi Lengkap

Lihat file-file berikut untuk contoh implementasi:
- `frontend/src/database/models/session.ts` - Query, Create, Update, Delete with logging
- `frontend/src/database/models/message.ts` - Complex query with multiple DB calls
- `frontend/src/database/models/thread.ts` - Simple CRUD with logging

## Debugging di Production

Log dapat dilihat di:
1. **Console output** - Saat development
2. **Wails debug window** - Saat testing
3. **Log files** - Dikonfigurasi di main.go

Untuk melihat logs di runtime:
```bash
# macOS/Linux
tail -f /path/to/app/logs/app.log

# Windows
Get-Content /path/to/app/logs/app.log -Wait
```

