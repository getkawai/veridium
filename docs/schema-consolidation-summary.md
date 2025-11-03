# Schema Consolidation Summary

## Overview

Successfully consolidated PostgreSQL migrations (0000-0040) into a complete SQLite-compatible schema for the Veridium application.

## What Was Done

### 1. Schema Audit
- Analyzed 41 PostgreSQL migrations from `frontend/src/database/migrations/`
- Identified PostgreSQL-specific features that needed conversion:
  - `vector` data types → `BLOB`
  - `jsonb` → `TEXT`
  - `timestamp with time zone` → `INTEGER` (Unix milliseconds)
  - HNSW indexes → Standard indexes
  - `DO $$ BEGIN` blocks → Removed
  - `ALTER TABLE` statements → Consolidated into CREATE TABLE
  - Reserved keywords (`index`, `order`) → Renamed (`chunk_index`, `sort_order`)

### 2. Schema Consolidation
- Started with `0000_initial_sqlite_setup.sql` as the base (61 tables)
- Added missing tables from later migrations:
  - RAG evaluation tables (4 tables)
  - User memory tables (5 tables)
  - User budgets and subscriptions (2 tables)
- Fixed SQLite reserved keyword conflicts
- Removed PostgreSQL schema qualifiers (`drizzle.`)

### 3. Final Schema
- **Total tables: 72**
- **Location:** `internal/database/schema/schema.sql`
- **Embedded in:** `database_service.go` using `//go:embed`

## Tables Included

### Core Tables (61 from base schema)
- Users & Authentication: `users`, `user_settings`, `user_installed_plugins`
- Sessions: `sessions`, `session_groups`
- Agents: `agents`, `agents_files`, `agents_knowledge_bases`, `agents_to_sessions`
- Messages: `messages`, `message_plugins`, `message_tts`, `message_translates`, `message_queries`, `message_query_chunks`, `message_chunks`, `message_groups`
- Topics & Threads: `topics`, `threads`, `topic_documents`
- Files: `files`, `global_files`, `files_to_sessions`, `messages_files`
- Knowledge Base: `knowledge_bases`, `knowledge_base_files`
- RAG: `chunks`, `unstructured_chunks`, `embeddings`, `file_chunks`, `document_chunks`
- Documents: `documents`
- Chat Groups: `chat_groups`, `chat_groups_agents`
- AI Infrastructure: `ai_providers`, `ai_models`
- Tasks: `async_tasks`
- Generations: `generation_topics`, `generation_batches`, `generations`
- API Keys: `api_keys`
- RBAC: `rbac_roles`, `rbac_permissions`, `rbac_role_permissions`, `rbac_user_roles`
- NextAuth: `nextauth_accounts`, `nextauth_sessions`, `nextauth_verificationtokens`, `nextauth_authenticators`
- OIDC: `oidc_authorization_codes`, `oidc_access_tokens`, `oidc_refresh_tokens`, `oidc_device_codes`, `oidc_interactions`, `oidc_grants`, `oidc_clients`, `oidc_sessions`, `oidc_consents`
- OAuth: `oauth_handoffs`
- Migrations: `__drizzle_migrations`

### Additional Tables (11 added)
- RAG Evaluation: `rag_eval_datasets`, `rag_eval_dataset_records`, `rag_eval_evaluations`, `rag_eval_evaluation_records`
- User Memory: `user_memories`, `user_memories_contexts`, `user_memories_experiences`, `user_memories_identities`, `user_memories_preferences`
- User Management: `user_budgets`, `user_subscriptions`

## Key Changes

### Data Type Conversions
- `vector(1024)` → `BLOB` (for embeddings and vector storage)
- `jsonb` → `TEXT` (JSON stored as text)
- `varchar(255)` → `TEXT`
- `bigint` → `INTEGER`
- `numeric` → `REAL`
- `timestamp with time zone` → `INTEGER` (Unix milliseconds)

### Reserved Keyword Fixes
- `index` → `chunk_index` (in chunks and unstructured_chunks tables)
- `order` → `sort_order` (in chat_groups_agents table)

### PostgreSQL-Specific Removals
- `CREATE EXTENSION IF NOT EXISTS vector` (pgvector)
- `USING hnsw` indexes
- `USING btree` indexes
- `DO $$ BEGIN ... END $$` blocks
- `ALTER TABLE ... SET (autovacuum_...)` statements
- Schema qualifiers (`public.`, `drizzle.`)

## Testing

Successfully tested schema initialization:
```bash
go run test_schema_temp.go
✓ Schema initialized successfully!
  Total tables created: 72
```

## Next Steps

Phase 1 of the migration plan is complete. The consolidated schema is ready for:
1. sqlc code generation
2. Go service implementation
3. Wails bindings creation
4. Frontend integration

## Files Created/Modified

### Created
- `internal/database/schema/schema.sql` - Complete SQLite schema (72 tables)
- `scripts/consolidate_schema.py` - Schema consolidation script
- `docs/schema-consolidation-summary.md` - This document

### Modified
- `database_service.go` - Updated to use embedded schema via `//go:embed`

## Schema Statistics

- **Total SQL lines:** ~1,100
- **Total tables:** 72
- **Total indexes:** 40+
- **Foreign key relationships:** 100+
- **Migrations consolidated:** 41 (0000-0040)

## Notes

- Vector columns are stored as BLOB - will need special handling for vector operations
- All timestamps use Unix milliseconds (INTEGER)
- JSON data stored as TEXT - will need JSON parsing in Go
- Schema is compatible with SQLite 3.x
- All foreign key constraints are properly defined
- Indexes are created for performance-critical columns

