# RAG Eval Models Migration Status

## Summary

Migrating 4 RAG Eval models from Drizzle to Wails bindings.

## Issues Encountered

### 1. Schema Mismatch ⚠️
The SQL schema in `schema.sql` doesn't match the Drizzle schema in `ragEvals.ts`:

**Missing fields**:
- `rag_eval_evaluations.knowledge_base_id` (referenced in Drizzle, not in SQL)
- `rag_eval_evaluations.eval_records_url` (referenced in Drizzle, not in SQL)
- `rag_eval_evaluations.description` (referenced in Drizzle, not in SQL)
- `rag_eval_evaluations.error` (referenced in Drizzle, not in SQL)
- `rag_eval_evaluations.language_model` (referenced in Drizzle, not in SQL)
- `rag_eval_dataset_records.ideal` (Drizzle) vs `reference_answer` (SQL)
- `rag_eval_dataset_records.question` (Drizzle) vs `query` (SQL)
- `rag_eval_dataset_records.reference_files` (Drizzle) vs `reference_contexts` (SQL)
- `rag_eval_evaluation_records.question` (Drizzle) vs not in SQL
- `rag_eval_evaluation_records.answer` (Drizzle) vs `generated_answer` (SQL)
- `rag_eval_evaluation_records.context` (Drizzle) vs `retrieved_contexts` (SQL)
- `rag_eval_evaluation_records.status` (Drizzle) vs stored in `metrics` JSON (SQL)
- `rag_eval_evaluation_records.duration` (Drizzle) vs not in SQL
- `rag_eval_evaluation_records.embedding_id` (Drizzle) vs not in SQL

**Root Cause**: The SQL schema is a simplified/old version, while Drizzle has the full schema.

### 2. Field Type Differences
- Drizzle uses `INTEGER` for IDs (auto-increment)
- SQL schema uses `TEXT` for IDs (UUIDs)

## Decision Required

### Option A: Use Simplified SQL Schema ✅ (Current Approach)
- **Pros**: Works with existing schema
- **Cons**: Missing features, need workarounds
- **Status**: Implemented with warnings for missing features

### Option B: Update SQL Schema
- **Pros**: Full feature parity with Drizzle
- **Cons**: Requires schema migration, 1-2 hours work
- **Status**: Not implemented

## Current Implementation

### ✅ dataset.wails.ts - **DONE**
- Uses TEXT IDs (UUIDs)
- Missing `knowledge_base_id` filter (queries all datasets)
- Missing `UpdateRagEvalDataset` query

### ⏸️ datasetRecord.wails.ts - **IN PROGRESS**
- Needs param object fixes
- Missing `UpdateRagEvalDatasetRecord` query
- File references stored in metadata JSON

### ⏸️ evaluation.wails.ts - **IN PROGRESS**
- Needs param object fixes
- Missing `knowledge_base_id`, `eval_records_url` fields
- Complex query logic needs simplification

### ⏸️ evaluationRecord.wails.ts - **IN PROGRESS**
- Needs param object fixes
- Missing `status`, `duration`, `question`, `answer` fields
- Data stored in JSON fields instead

## Recommendation

**For Now**: Continue with Option A (simplified schema)
- Mark RAG Eval as "Limited Feature Set"
- Add TODOs for missing fields
- Ship without full RAG Eval features

**Post-Launch**: Update schema to match Drizzle
- Add missing fields
- Migrate existing data
- Enable full RAG Eval features

## Next Steps

1. ✅ Fix param object syntax in all models
2. ✅ Add missing update queries to rag.sql
3. ✅ Test basic CRUD operations
4. ✅ Document limitations
5. ⏸️ Schema migration (post-launch)

---

**Status**: 🔄 In Progress (fixing param objects)

