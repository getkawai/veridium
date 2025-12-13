-- Memory Queries for Infinite Memory Architecture
-- Supports MemGPT-style conversation memory with semantic search

-- ============================================================================
-- USER MEMORIES (Core Memory Facts)
-- ============================================================================

-- name: CreateUserMemory :one
INSERT INTO user_memories (
    id, memory_category, memory_layer, memory_type,
    title, summary, summary_vector_1024, details, details_vector_1024,
    status, accessed_count, last_accessed_at, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetUserMemory :one
SELECT * FROM user_memories WHERE id = ?;

-- name: ListUserMemories :many
SELECT * FROM user_memories
ORDER BY last_accessed_at DESC
LIMIT ? OFFSET ?;

-- name: ListUserMemoriesByCategory :many
SELECT * FROM user_memories
WHERE memory_category = ?
ORDER BY last_accessed_at DESC
LIMIT ? OFFSET ?;

-- name: ListUserMemoriesByType :many
SELECT * FROM user_memories
WHERE memory_type = ?
ORDER BY last_accessed_at DESC
LIMIT ? OFFSET ?;

-- name: UpdateUserMemory :one
UPDATE user_memories
SET title = COALESCE(?, title),
    summary = COALESCE(?, summary),
    summary_vector_1024 = COALESCE(?, summary_vector_1024),
    details = COALESCE(?, details),
    details_vector_1024 = COALESCE(?, details_vector_1024),
    status = COALESCE(?, status),
    updated_at = ?
WHERE id = ?
RETURNING *;

-- name: UpdateMemoryAccessCount :exec
UPDATE user_memories
SET accessed_count = accessed_count + 1,
    last_accessed_at = ?,
    updated_at = ?
WHERE id = ?;

-- name: DeleteUserMemory :exec
DELETE FROM user_memories WHERE id = ?;

-- name: CountUserMemories :one
SELECT COUNT(*) as count FROM user_memories;

-- name: GetRecentMemories :many
SELECT * FROM user_memories
ORDER BY created_at DESC
LIMIT ?;

-- name: GetMostAccessedMemories :many
SELECT * FROM user_memories
ORDER BY accessed_count DESC, last_accessed_at DESC
LIMIT ?;

-- name: SearchMemoriesByTitle :many
SELECT * FROM user_memories
WHERE title LIKE '%' || ? || '%'
ORDER BY last_accessed_at DESC
LIMIT ?;

-- ============================================================================
-- USER MEMORIES EXPERIENCES (Situational Memory)
-- ============================================================================

-- name: CreateUserMemoryExperience :one
INSERT INTO user_memories_experiences (
    id, user_memory_id, labels, extracted_labels, type,
    situation, situation_vector, reasoning, possible_outcome,
    action, action_vector, key_learning, key_learning_vector,
    metadata, score_confidence, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetUserMemoryExperience :one
SELECT * FROM user_memories_experiences WHERE id = ?;

-- name: ListExperiencesByMemoryId :many
SELECT * FROM user_memories_experiences
WHERE user_memory_id = ?
ORDER BY created_at DESC;

-- name: DeleteUserMemoryExperience :exec
DELETE FROM user_memories_experiences WHERE id = ?;

-- ============================================================================
-- USER MEMORIES IDENTITIES (User Profile Memory)
-- ============================================================================

-- name: CreateUserMemoryIdentity :one
INSERT INTO user_memories_identities (
    id, user_memory_id, current_focuses, description, description_vector,
    experience, extracted_labels, labels, relationship, role, type,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetUserMemoryIdentity :one
SELECT * FROM user_memories_identities WHERE id = ?;

-- name: ListIdentitiesByMemoryId :many
SELECT * FROM user_memories_identities
WHERE user_memory_id = ?
ORDER BY created_at DESC;

-- name: DeleteUserMemoryIdentity :exec
DELETE FROM user_memories_identities WHERE id = ?;

-- ============================================================================
-- USER MEMORIES PREFERENCES (User Preference Memory)
-- ============================================================================

-- name: CreateUserMemoryPreference :one
INSERT INTO user_memories_preferences (
    id, context_id, user_memory_id, labels, extracted_labels, extracted_scopes,
    conclusion_directives, conclusion_directives_vector, type, suggestions,
    score_priority, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetUserMemoryPreference :one
SELECT * FROM user_memories_preferences WHERE id = ?;

-- name: ListPreferencesByMemoryId :many
SELECT * FROM user_memories_preferences
WHERE user_memory_id = ?
ORDER BY score_priority DESC, created_at DESC;

-- name: DeleteUserMemoryPreference :exec
DELETE FROM user_memories_preferences WHERE id = ?;

-- ============================================================================
-- USER MEMORIES CONTEXTS (Contextual Memory)
-- ============================================================================

-- name: CreateUserMemoryContext :one
INSERT INTO user_memories_contexts (
    id, user_memory_ids, labels, extracted_labels, associated_objects,
    associated_subjects, title, title_vector, description, description_vector,
    type, current_status, score_impact, score_urgency, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetUserMemoryContext :one
SELECT * FROM user_memories_contexts WHERE id = ?;

-- name: ListUserMemoryContexts :many
SELECT * FROM user_memories_contexts
ORDER BY score_impact DESC, created_at DESC
LIMIT ? OFFSET ?;

-- name: DeleteUserMemoryContext :exec
DELETE FROM user_memories_contexts WHERE id = ?;

-- ============================================================================
-- BATCH OPERATIONS
-- ============================================================================

-- name: BatchDeleteUserMemories :exec
DELETE FROM user_memories
WHERE id IN (sqlc.slice('ids'));

-- name: GetUserMemoriesByIds :many
SELECT * FROM user_memories
WHERE id IN (sqlc.slice('ids'));

-- ============================================================================
-- CONVERSATION MEMORY LINKING (Link memories to messages/sessions)
-- ============================================================================

-- name: GetMemoriesBySessionContext :many
SELECT um.* FROM user_memories um
WHERE um.memory_category IN ('conversation', 'context', 'fact')
ORDER BY um.last_accessed_at DESC
LIMIT ?;

-- name: ArchiveOldMemories :exec
UPDATE user_memories
SET status = 'archived',
    updated_at = ?
WHERE last_accessed_at < ?
  AND status != 'archived';
