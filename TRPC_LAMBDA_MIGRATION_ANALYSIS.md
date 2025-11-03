# tRPC Lambda Migration to Wails Analysis

## Current State

`@libs/trpc/lambda` is still used by these routers:

### 1. MCP Routers (2 files)
- `desktop/mcp.ts` - MCP for desktop
- `tools/mcp.ts` - MCP for tools/web

**Operations**:
- `getStdioMcpServerManifest` - Get MCP server manifest
- `getStreamableMcpServerManifest` - Get streamable manifest
- `listTools` - List available tools
- `listResources` - List available resources
- `listPrompts` - List available prompts
- `callTool` - Execute MCP tool
- `validMcpServerInstallable` - Check if MCP can be installed

**Complexity**: High
- External MCP server communication
- Stdio subprocess management
- HTTP streaming
- Complex service logic

### 2. Search Router (1 file)
- `tools/search.ts`

**Operations**:
- `crawlPages` - Web crawling
- `query` - Search query
- `webSearch` - Web search

**Complexity**: High
- External API calls (Jina, Browserless)
- Web crawling
- Search engine integration

### 3. PgTable Router (1 file)
- `desktop/pgTable.ts`

**Operations**:
- `getAllTables` - List database tables
- `getTableData` - Get table data with pagination
- `getTableDetails` - Get table schema

**Complexity**: Low
- Uses `TableViewerRepo` (already using Wails SQLite)
- Simple CRUD operations

## Migration Feasibility

### ❌ **NOT Feasible**: MCP & Search Routers

**Why?**
1. **External I/O Operations**:
   - MCP: Subprocess management, stdio streams
   - Search: HTTP requests to external APIs
   - Both require Node.js runtime capabilities

2. **Complex Business Logic**:
   - MCP service has complex client management
   - Search has multiple provider integrations
   - Would require rewriting entire services in Go

3. **No Database Operations**:
   - These don't use database at all
   - Migration won't provide performance benefits
   - Just adds complexity

4. **Wails Limitations**:
   - Wails is for database/data layer, not for:
     - External API calls
     - Subprocess management
     - Streaming responses
     - Complex service orchestration

**Recommendation**: **KEEP as tRPC** - These are proper use cases for tRPC server-side routers.

### ✅ **Feasible**: PgTable Router

**Why?**
1. **Pure Database Operations**:
   - Only queries SQLite
   - Uses `TableViewerRepo` which already uses Wails

2. **Simple Logic**:
   - No external I/O
   - No complex business logic
   - Just data retrieval

3. **Direct Binding Possible**:
   ```go
   // Go service
   type TableViewerService struct {
       db *database.Service
   }
   
   func (s *TableViewerService) GetAllTables() ([]TableInfo, error) {
       // Direct SQL queries
   }
   ```

**Recommendation**: **CAN migrate** but **LOW PRIORITY** (dev tool only).

## Architecture Comparison

### Current (tRPC)
```
Frontend → tRPC Router → tRPC Procedure → Service → External API/DB
         (HTTP)        (Auth/Middleware)   (Logic)   (I/O)
```

### If Migrated to Wails
```
Frontend → Wails Binding → Go Service → External API/DB
         (Direct Call)     (Logic)      (I/O)
```

**Problem**: Go services would need to:
- Manage MCP stdio processes (complex in Go)
- Make HTTP requests to search APIs (possible but unnecessary)
- Handle streaming responses (complex)

## Decision Matrix

| Router | Database Operations | External I/O | Complexity | Feasible? | Recommended? |
|--------|---------------------|--------------|------------|-----------|--------------|
| MCP (desktop) | ❌ None | ✅ Yes (stdio) | Very High | ❌ No | ❌ Keep tRPC |
| MCP (tools) | ❌ None | ✅ Yes (HTTP) | Very High | ❌ No | ❌ Keep tRPC |
| Search | ❌ None | ✅ Yes (HTTP) | High | ❌ No | ❌ Keep tRPC |
| PgTable | ✅ Yes | ❌ None | Low | ✅ Yes | ⏸️ Low Priority |

## Recommendation

### ✅ **KEEP tRPC Lambda Infrastructure**

**Reasons**:
1. **Proper Use Case**: MCP and Search are exactly what tRPC is designed for:
   - Server-side operations
   - External API integration
   - Complex service orchestration
   - Type-safe client-server communication

2. **No Performance Benefit**: 
   - These operations are I/O bound (external APIs)
   - Database is not the bottleneck
   - Migration won't improve performance

3. **High Migration Cost**:
   - Need to rewrite MCP service in Go (complex subprocess management)
   - Need to rewrite Search service in Go (HTTP clients)
   - Need to handle streaming in Go
   - 10-20 hours of work for minimal benefit

4. **Maintenance Burden**:
   - Go isn't better than Node.js for these operations
   - Lose ecosystem benefits (npm packages for MCP, search APIs)
   - Harder to debug and test

### ⏸️ **Optional**: Migrate PgTable (Low Priority)

If you want 100% Drizzle-free:
1. Create Go `TableViewerService`
2. Bind to Wails
3. Replace tRPC calls with direct Wails calls
4. **Effort**: 1-2 hours
5. **Benefit**: Remove one more tRPC dependency (dev tool only)

## Current Migration Summary

### ✅ Successfully Migrated (Wails)
- All database models (Session, User, Message, etc.)
- All client services (direct calls, no HTTP)
- All repositories (TableViewer, DataExporter already use Wails SQLite)

### ✅ Kept as tRPC (Proper Use Cases)
- MCP routers (external subprocess/HTTP)
- Search routers (external APIs)
- PgTable router (dev tool, low priority)

### ❌ Removed (No Longer Needed)
- Lambda message router (database → Wails)
- Lambda generation router (database → Wails)
- Other lambda routers (database operations → Wails)

## Final Architecture

```
┌─────────────────────────────────────────┐
│           Frontend React                │
├─────────────────────────────────────────┤
│                                         │
│  Database Operations (95%)              │
│  └─→ Direct Wails Bindings → Go → SQL  │
│                                         │
│  External Services (5%)                 │
│  └─→ tRPC → Node.js → External APIs    │
│                                         │
└─────────────────────────────────────────┘
```

**Result**: 
- ✅ 95% of operations use Wails (fast, direct)
- ✅ 5% use tRPC (proper use case: external I/O)
- ✅ Clean separation of concerns
- ✅ Best tool for each job

## Conclusion

**Answer**: **NO, tidak perlu dimigrate**.

MCP dan Search routers adalah **proper use cases for tRPC** karena mereka:
- Tidak ada database operations
- Fokus pada external I/O
- Butuh Node.js ecosystem
- Tidak dapat performance benefit dari Go

`@libs/trpc/lambda` infrastructure **TETAP DIPERLUKAN** dan ini adalah **correct architecture**.

Yang sudah dimigrate (database operations) sudah **optimal**. Yang tersisa (MCP, Search) **sebaiknya tetap di tRPC**.

