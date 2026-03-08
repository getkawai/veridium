# Memory Architecture

This document defines the active memory architecture in Veridium.

## Current Runtime Split

1. `MuninnDB` is the active backend for conversational memory.
2. `DuckDB` remains the vector store for RAG/file chunk retrieval.
3. Legacy `search_memory` tool registration is intentionally disabled.

## Active Path (Conversation Memory)

The active path is initialized in `internal/app/context.go`:

1. `InitMemoryServices()` creates `MuninnMemoryBackend`.
2. `MemoryIntegration` is constructed in Muninn mode.
3. Runtime memory integration is enforced to use Muninn backend.

Main files:

- `internal/services/muninn_memory.go`
- `internal/services/memory_integration.go`
- `internal/app/context.go`

## DuckDB Role

DuckDB is still used for vector retrieval in the RAG pipeline.

Main files:

- `internal/services/duckdb_store.go`
- `internal/services/vector_search.go`
- `internal/services/rag_processor.go`

This is separate from conversational memory.

## Operational Guardrails

1. Do not re-enable legacy `search_memory` tool by default.
2. Keep startup logs explicit about backend roles.
3. Add tests when changing memory wiring:
   - integration must stay in Muninn mode
   - legacy tool registration remains disabled

## Next Migration Step

If full unification is desired later:

1. Introduce a `VectorStore` interface for retrieval adapters.
2. Add a Muninn-backed retrieval adapter.
3. Run A/B evaluation (quality + latency) before replacing DuckDB path.
