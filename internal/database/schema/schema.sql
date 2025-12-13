-- Initial SQLite migration generated from schemas
-- This file contains the initial database setup for Wails SQLite

-- User settings
CREATE TABLE IF NOT EXISTS user_settings (
  id TEXT PRIMARY KEY,
  tts TEXT, -- JSON as text
  hotkey TEXT, -- JSON as text
  key_vaults TEXT,
  general TEXT, -- JSON as text
  language_model TEXT, -- JSON as text
  system_agent TEXT, -- JSON as text
  default_agent TEXT, -- JSON as text
  tool TEXT, -- JSON as text
  image TEXT -- JSON as text
);

-- User installed plugins
CREATE TABLE IF NOT EXISTS user_installed_plugins (
  identifier TEXT NOT NULL,
  type TEXT NOT NULL CHECK (type IN ('plugin', 'customPlugin')),
  manifest TEXT, -- JSON as text
  settings TEXT, -- JSON as text
  custom_params TEXT, -- JSON as text
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  PRIMARY KEY (identifier)
);

-- Global files table
CREATE TABLE IF NOT EXISTS global_files (
  hash_id TEXT PRIMARY KEY,
  file_type TEXT NOT NULL,
  size INTEGER NOT NULL,
  url TEXT NOT NULL,
  metadata TEXT, -- JSON as text
  creator TEXT,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- Files table
CREATE TABLE IF NOT EXISTS files (
  id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  file_type TEXT NOT NULL,
  file_hash TEXT REFERENCES global_files(hash_id),
  name TEXT NOT NULL,
  size INTEGER NOT NULL,
  url TEXT NOT NULL,
  source TEXT, -- JSON as text
  metadata TEXT, -- JSON as text
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- Knowledge bases
CREATE TABLE IF NOT EXISTS knowledge_bases (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  avatar TEXT,
  type TEXT,
  is_public INTEGER DEFAULT 0 NOT NULL, -- boolean as integer
  settings TEXT, -- JSON as text
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- Knowledge base files junction table
CREATE TABLE IF NOT EXISTS knowledge_base_files (
  knowledge_base_id TEXT NOT NULL REFERENCES knowledge_bases(id) ON DELETE CASCADE,
  file_id TEXT NOT NULL REFERENCES files(id) ON DELETE CASCADE,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  PRIMARY KEY (knowledge_base_id, file_id)
);

-- Agents table
CREATE TABLE IF NOT EXISTS agents (
	id TEXT PRIMARY KEY,
  title TEXT,
  description TEXT,
  tags TEXT DEFAULT '[]', -- JSON as text
  avatar TEXT,
  background_color TEXT,
  plugins TEXT DEFAULT '[]', -- JSON as text
  chat_config TEXT, -- JSON as text
  few_shots TEXT, -- JSON as text
  model TEXT,
  params TEXT DEFAULT '{}', -- JSON as text
  provider TEXT,
  system_role TEXT,
  tts TEXT, -- JSON as text
  virtual INTEGER DEFAULT 0 NOT NULL, -- boolean as integer
  opening_message TEXT,
  opening_questions TEXT DEFAULT '[]', -- JSON as text
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- Session groups
CREATE TABLE IF NOT EXISTS session_groups (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  sort INTEGER,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- Sessions table
CREATE TABLE IF NOT EXISTS sessions (
	id TEXT PRIMARY KEY,
  title TEXT,
  description TEXT,
  avatar TEXT,
  background_color TEXT,
  type TEXT DEFAULT 'agent',
  group_id TEXT REFERENCES session_groups(id) ON DELETE SET NULL,
  pinned INTEGER DEFAULT 0 NOT NULL, -- boolean as integer
  model TEXT, -- added for direct access
  tags TEXT DEFAULT '[]', -- JSON as text
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- Topics table
CREATE TABLE IF NOT EXISTS topics (
  id TEXT PRIMARY KEY,
  title TEXT,
  favorite INTEGER DEFAULT 0 NOT NULL, -- boolean as integer
  session_id TEXT REFERENCES sessions(id) ON DELETE CASCADE,
  group_id TEXT REFERENCES chat_groups(id) ON DELETE CASCADE,
  history_summary TEXT, -- JSON as text
  metadata TEXT, -- JSON as text
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- Threads table
CREATE TABLE IF NOT EXISTS threads (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  type TEXT NOT NULL CHECK (type IN ('continuation', 'standalone')),
  status TEXT DEFAULT 'active' CHECK (status IN ('active', 'deprecated', 'archived')),
  topic_id TEXT NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
  source_message_id TEXT NOT NULL,
  parent_thread_id TEXT REFERENCES threads(id) ON DELETE SET NULL,
  last_active_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- Messages table
CREATE TABLE IF NOT EXISTS messages (
  id TEXT PRIMARY KEY,
  role TEXT NOT NULL,
  content TEXT, -- JSON as text
  reasoning TEXT, -- JSON as text
  search TEXT, -- JSON as text
  metadata TEXT, -- JSON as text
  model TEXT, -- JSON as text
  provider TEXT, -- JSON as text
  favorite INTEGER DEFAULT 0 NOT NULL, -- boolean as integer
  error TEXT, -- JSON as text
  tools TEXT, -- JSON as text
  trace_id TEXT,
  observation_id TEXT,
  session_id TEXT REFERENCES sessions(id) ON DELETE CASCADE,
  topic_id TEXT REFERENCES topics(id) ON DELETE CASCADE,
  thread_id TEXT REFERENCES threads(id) ON DELETE CASCADE,
  parent_id TEXT REFERENCES messages(id) ON DELETE SET NULL,
  quota_id TEXT REFERENCES messages(id) ON DELETE SET NULL,
  agent_id TEXT REFERENCES agents(id) ON DELETE SET NULL,
  group_id TEXT REFERENCES chat_groups(id) ON DELETE SET NULL,
  target_id TEXT,
  message_group_id TEXT REFERENCES message_groups(id) ON DELETE CASCADE,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- Message plugins
CREATE TABLE IF NOT EXISTS message_plugins (
  id TEXT PRIMARY KEY REFERENCES messages(id) ON DELETE CASCADE,
  tool_call_id TEXT,
  type TEXT DEFAULT 'default' CHECK (type IN ('default', 'markdown', 'standalone', 'builtin')),
  api_name TEXT,
  arguments TEXT, -- JSON as text
  identifier TEXT,
  state TEXT, -- JSON as text
  error TEXT, -- JSON as text
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- Message TTS
CREATE TABLE IF NOT EXISTS message_tts (
  id TEXT PRIMARY KEY REFERENCES messages(id) ON DELETE CASCADE,
  content_md5 TEXT,
  file_id TEXT REFERENCES files(id) ON DELETE CASCADE,
  voice TEXT
);

-- Message translates
CREATE TABLE IF NOT EXISTS message_translates (
  id TEXT PRIMARY KEY REFERENCES messages(id) ON DELETE CASCADE,
  content TEXT, -- JSON as text
  "from" TEXT,
  "to" TEXT
);

-- Message queries
CREATE TABLE IF NOT EXISTS message_queries (
  id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  message_id TEXT NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
  rewrite_query TEXT,
  user_query TEXT
);

-- Message query chunks
CREATE TABLE IF NOT EXISTS message_query_chunks (
  message_id TEXT REFERENCES messages(id) ON DELETE CASCADE,
  query_id TEXT REFERENCES message_queries(id) ON DELETE CASCADE,
  chunk_id TEXT REFERENCES chunks(id) ON DELETE CASCADE,
  similarity INTEGER, -- We'll store similarity as integer for now
  PRIMARY KEY (chunk_id, message_id, query_id)
);

-- Message chunks
CREATE TABLE IF NOT EXISTS message_chunks (
  message_id TEXT REFERENCES messages(id) ON DELETE CASCADE,
  chunk_id TEXT REFERENCES chunks(id) ON DELETE CASCADE,
  PRIMARY KEY (chunk_id, message_id)
);

-- Messages files junction table
CREATE TABLE IF NOT EXISTS messages_files (
  file_id TEXT NOT NULL REFERENCES files(id) ON DELETE CASCADE,
  message_id TEXT NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
  PRIMARY KEY (file_id, message_id)
);

-- Chunks table (RAG)
CREATE TABLE IF NOT EXISTS chunks (
  id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  document_id TEXT REFERENCES documents(id) ON DELETE CASCADE,
  text TEXT,
  abstract TEXT,
  metadata TEXT, -- JSON as text
  chunk_index INTEGER,
  type TEXT,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- Unstructured chunks
CREATE TABLE IF NOT EXISTS unstructured_chunks (
  id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  text TEXT,
  metadata TEXT, -- JSON as text
  chunk_index INTEGER,
  type TEXT,
  parent_id TEXT,
  composite_id TEXT REFERENCES chunks(id) ON DELETE CASCADE,
  file_id TEXT REFERENCES files(id) ON DELETE CASCADE,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- File chunks junction table
CREATE TABLE IF NOT EXISTS file_chunks (
  file_id TEXT REFERENCES files(id) ON DELETE CASCADE,
  chunk_id TEXT REFERENCES chunks(id) ON DELETE CASCADE,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  PRIMARY KEY (file_id, chunk_id)
);

-- Agents to sessions junction table
CREATE TABLE IF NOT EXISTS agents_to_sessions (
  agent_id TEXT NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
  session_id TEXT NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
  PRIMARY KEY (agent_id, session_id)
);

-- Files to sessions junction table
CREATE TABLE IF NOT EXISTS files_to_sessions (
  file_id TEXT NOT NULL REFERENCES files(id) ON DELETE CASCADE,
  session_id TEXT NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
  PRIMARY KEY (file_id, session_id)
);

-- Agents knowledge bases junction table
CREATE TABLE IF NOT EXISTS agents_knowledge_bases (
  agent_id TEXT NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
  knowledge_base_id TEXT NOT NULL REFERENCES knowledge_bases(id) ON DELETE CASCADE,
  enabled INTEGER DEFAULT 1 NOT NULL, -- boolean as integer
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  PRIMARY KEY (agent_id, knowledge_base_id)
);

-- Agents files junction table
CREATE TABLE IF NOT EXISTS agents_files (
  file_id TEXT NOT NULL REFERENCES files(id) ON DELETE CASCADE,
  agent_id TEXT NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
  enabled INTEGER DEFAULT 1 NOT NULL, -- boolean as integer
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  PRIMARY KEY (file_id, agent_id)
);

-- Topic documents junction table
CREATE TABLE IF NOT EXISTS topic_documents (
  document_id TEXT NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
  topic_id TEXT NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  PRIMARY KEY (document_id, topic_id)
);

-- Document chunks junction table
CREATE TABLE IF NOT EXISTS document_chunks (
  document_id TEXT NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
  chunk_id TEXT NOT NULL REFERENCES chunks(id) ON DELETE CASCADE,
  page_index INTEGER,
  PRIMARY KEY (document_id, chunk_id)
);

-- Chat groups table
CREATE TABLE IF NOT EXISTS chat_groups (
  id TEXT PRIMARY KEY,
  title TEXT,
  description TEXT,
  config TEXT, -- JSON as text
  group_id TEXT,
  pinned INTEGER DEFAULT 0 NOT NULL, -- boolean as integer
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- Chat groups agents junction table
CREATE TABLE IF NOT EXISTS chat_groups_agents (
  chat_group_id TEXT NOT NULL REFERENCES chat_groups(id) ON DELETE CASCADE,
  agent_id TEXT NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
  enabled INTEGER DEFAULT 1 NOT NULL, -- boolean as integer
  sort_order INTEGER DEFAULT 0,
  role TEXT DEFAULT 'participant',
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  PRIMARY KEY (chat_group_id, agent_id)
);

-- Message groups table
CREATE TABLE IF NOT EXISTS message_groups (
  id TEXT PRIMARY KEY,
  title TEXT,
  description TEXT,
  topic_id TEXT REFERENCES topics(id) ON DELETE CASCADE,
  parent_group_id TEXT REFERENCES message_groups(id) ON DELETE CASCADE,
  parent_message_id TEXT,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- AI providers table
CREATE TABLE IF NOT EXISTS ai_providers (
  id TEXT NOT NULL,
  name TEXT,
  sort INTEGER, -- boolean as integer (0/1)
  enabled INTEGER, -- boolean as integer (0/1)
  fetch_on_client INTEGER, -- boolean as integer (0/1)
  check_model TEXT,
  logo TEXT,
  description TEXT,
  key_vaults TEXT,
  source TEXT,
  settings TEXT, -- JSON as text
  config TEXT, -- JSON as text
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  PRIMARY KEY (id)
);

-- AI models table
CREATE TABLE IF NOT EXISTS ai_models (
  id TEXT NOT NULL,
  display_name TEXT,
  description TEXT,
  organization TEXT,
  enabled INTEGER, -- boolean as integer (0/1)
  provider_id TEXT NOT NULL,
  type TEXT DEFAULT 'chat' NOT NULL,
  sort INTEGER, -- boolean as integer (0/1)
  pricing TEXT, -- JSON as text
  parameters TEXT DEFAULT '{}', -- JSON as text
  config TEXT, -- JSON as text
  abilities TEXT DEFAULT '{}', -- JSON as text
  context_window_tokens INTEGER,
  source TEXT,
  released_at TEXT,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  PRIMARY KEY (id, provider_id)
);

-- Async tasks table
CREATE TABLE IF NOT EXISTS async_tasks (
  id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  type TEXT,
  status TEXT,
  error TEXT,
  duration INTEGER,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- Generation topics table
CREATE TABLE IF NOT EXISTS generation_topics (
  id TEXT PRIMARY KEY,
  title TEXT,
  cover_url TEXT,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- Generation batches table
CREATE TABLE IF NOT EXISTS generation_batches (
  id TEXT PRIMARY KEY,
  generation_topic_id TEXT NOT NULL REFERENCES generation_topics(id) ON DELETE CASCADE,
  provider TEXT NOT NULL,
  model TEXT NOT NULL,
  prompt TEXT NOT NULL,
  width INTEGER,
  height INTEGER,
  ratio TEXT,
  config TEXT, -- JSON as text
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- Generations table
CREATE TABLE IF NOT EXISTS generations (
  id TEXT PRIMARY KEY,
  generation_batch_id TEXT NOT NULL REFERENCES generation_batches(id) ON DELETE CASCADE,
  async_task_id TEXT REFERENCES async_tasks(id) ON DELETE SET NULL,
  file_id TEXT REFERENCES files(id) ON DELETE CASCADE,
  seed INTEGER DEFAULT 0,
  asset TEXT, -- JSON as text
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- Documents table
CREATE TABLE IF NOT EXISTS documents (
  id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  title TEXT,
  content TEXT,
  file_type TEXT NOT NULL,
  filename TEXT,
  total_char_count INTEGER NOT NULL,
  total_line_count INTEGER NOT NULL,
  metadata TEXT, -- JSON as text
  pages TEXT, -- JSON as text
  source_type TEXT NOT NULL CHECK (source_type IN ('file', 'web', 'api')),
  source TEXT NOT NULL,
  file_id TEXT REFERENCES files(id) ON DELETE SET NULL,
  editor_data TEXT, -- JSON as text
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- API keys table
CREATE TABLE IF NOT EXISTS api_keys (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  key TEXT NOT NULL UNIQUE,
  enabled INTEGER DEFAULT 1 NOT NULL, -- boolean as integer
  expires_at INTEGER, -- timestamp_ms
  last_used_at INTEGER, -- timestamp_ms
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);
CREATE INDEX IF NOT EXISTS idx_messages_session_id ON messages(session_id);
CREATE INDEX IF NOT EXISTS idx_messages_topic_id ON messages(topic_id);
CREATE INDEX IF NOT EXISTS idx_messages_parent_id ON messages(parent_id);
CREATE INDEX IF NOT EXISTS idx_messages_quota_id ON messages(quota_id);
CREATE INDEX IF NOT EXISTS idx_messages_thread_id ON messages(thread_id);

CREATE INDEX IF NOT EXISTS idx_chunks_document_id ON chunks(document_id);

CREATE INDEX IF NOT EXISTS idx_files_file_hash ON files(file_hash);

CREATE INDEX IF NOT EXISTS idx_agents_title ON agents(title);
CREATE INDEX IF NOT EXISTS idx_agents_description ON agents(description);

CREATE INDEX IF NOT EXISTS idx_documents_source ON documents(source);
CREATE INDEX IF NOT EXISTS idx_documents_file_type ON documents(file_type);
CREATE INDEX IF NOT EXISTS idx_documents_file_id ON documents(file_id);

-- Drizzle migration tracking table
CREATE TABLE IF NOT EXISTS __drizzle_migrations (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  hash TEXT NOT NULL,
  created_at INTEGER
);

-- Insert initial migration record
INSERT OR IGNORE INTO __drizzle_migrations (hash, created_at) VALUES ('initial_sqlite_setup', strftime('%s', 'now') * 1000);-- Additional tables from migrations 0008, 0037, 0040

-- RAG Evaluation tables (from 0008_add_rag_evals.sql)
CREATE TABLE IF NOT EXISTS rag_eval_datasets (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

CREATE TABLE IF NOT EXISTS rag_eval_dataset_records (
  id TEXT PRIMARY KEY,
  dataset_id TEXT NOT NULL REFERENCES rag_eval_datasets(id) ON DELETE CASCADE,
  query TEXT NOT NULL,
  reference_answer TEXT,
  reference_contexts TEXT,  -- JSON as text
  metadata TEXT,  -- JSON as text
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

CREATE TABLE IF NOT EXISTS rag_eval_evaluations (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  dataset_id TEXT NOT NULL REFERENCES rag_eval_datasets(id) ON DELETE CASCADE,
  config TEXT,  -- JSON as text
  status TEXT NOT NULL CHECK (status IN ('pending', 'running', 'completed', 'failed')),
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

CREATE TABLE IF NOT EXISTS rag_eval_evaluation_records (
  id TEXT PRIMARY KEY,
  evaluation_id TEXT NOT NULL REFERENCES rag_eval_evaluations(id) ON DELETE CASCADE,
  dataset_record_id TEXT NOT NULL REFERENCES rag_eval_dataset_records(id) ON DELETE CASCADE,
  retrieved_contexts TEXT,  -- JSON as text
  generated_answer TEXT,
  metrics TEXT,  -- JSON as text
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- User Memory tables (from 0037_add_user_memory.sql, 0040_improve_user_memory_field.sql)
CREATE TABLE IF NOT EXISTS user_memories (
  id TEXT PRIMARY KEY,
  memory_category TEXT,
  memory_layer TEXT,
  memory_type TEXT,
  title TEXT,
  summary TEXT,
  summary_vector_1024 BLOB,  -- Store as BLOB in SQLite
  details TEXT,
  details_vector_1024 BLOB,  -- Store as BLOB in SQLite
  status TEXT,
  accessed_count INTEGER DEFAULT 0,
  last_accessed_at INTEGER NOT NULL,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

CREATE TABLE IF NOT EXISTS user_memories_contexts (
  id TEXT PRIMARY KEY,
  user_memory_ids TEXT,  -- JSON as text
  labels TEXT,  -- JSON as text
  extracted_labels TEXT,  -- JSON as text
  associated_objects TEXT,  -- JSON as text
  associated_subjects TEXT,  -- JSON as text
  title TEXT,
  title_vector BLOB,  -- Store as BLOB in SQLite
  description TEXT,
  description_vector BLOB,  -- Store as BLOB in SQLite
  type TEXT,
  current_status TEXT,
  score_impact REAL DEFAULT 0,
  score_urgency REAL DEFAULT 0,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

CREATE TABLE IF NOT EXISTS user_memories_experiences (
  id TEXT PRIMARY KEY,
  user_memory_id TEXT REFERENCES user_memories(id) ON DELETE CASCADE,
  labels TEXT,  -- JSON as text
  extracted_labels TEXT,  -- JSON as text
  type TEXT,
  situation TEXT,
  situation_vector BLOB,  -- Store as BLOB in SQLite
  reasoning TEXT,
  possible_outcome TEXT,
  action TEXT,
  action_vector BLOB,  -- Store as BLOB in SQLite
  key_learning TEXT,
  key_learning_vector BLOB,  -- Store as BLOB in SQLite
  metadata TEXT,  -- JSON as text
  score_confidence REAL DEFAULT 0,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

CREATE TABLE IF NOT EXISTS user_memories_identities (
  id TEXT PRIMARY KEY,
  user_memory_id TEXT REFERENCES user_memories(id) ON DELETE CASCADE,
  current_focuses TEXT,
  description TEXT,
  description_vector BLOB,  -- Store as BLOB in SQLite
  experience TEXT,
  extracted_labels TEXT,  -- JSON as text
  labels TEXT,  -- JSON as text
  relationship TEXT,
  role TEXT,
  type TEXT,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

CREATE TABLE IF NOT EXISTS user_memories_preferences (
  id TEXT PRIMARY KEY,
  context_id TEXT REFERENCES user_memories_contexts(id) ON DELETE CASCADE,
  user_memory_id TEXT REFERENCES user_memories(id) ON DELETE CASCADE,
  labels TEXT,  -- JSON as text
  extracted_labels TEXT,  -- JSON as text
  extracted_scopes TEXT,  -- JSON as text
  conclusion_directives TEXT,
  conclusion_directives_vector BLOB,  -- Store as BLOB in SQLite
  type TEXT,
  suggestions TEXT,
  score_priority REAL DEFAULT 0,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- User budgets and subscriptions (from various migrations)
CREATE TABLE IF NOT EXISTS user_budgets (
  id TEXT PRIMARY KEY,
  budget_type TEXT NOT NULL,
  amount REAL NOT NULL,
  period TEXT NOT NULL,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

CREATE TABLE IF NOT EXISTS user_subscriptions (
  id TEXT PRIMARY KEY,
  plan_id TEXT NOT NULL,
  status TEXT NOT NULL,
  started_at INTEGER NOT NULL,
  expires_at INTEGER,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_rag_eval_dataset_records_dataset_id ON rag_eval_dataset_records(dataset_id);
CREATE INDEX IF NOT EXISTS idx_rag_eval_evaluations_dataset_id ON rag_eval_evaluations(dataset_id);
CREATE INDEX IF NOT EXISTS idx_rag_eval_evaluation_records_evaluation_id ON rag_eval_evaluation_records(evaluation_id);

CREATE INDEX IF NOT EXISTS idx_user_memories_experiences_user_memory_id ON user_memories_experiences(user_memory_id);
CREATE INDEX IF NOT EXISTS idx_user_memories_identities_user_memory_id ON user_memories_identities(user_memory_id);
CREATE INDEX IF NOT EXISTS idx_user_memories_preferences_user_memory_id ON user_memories_preferences(user_memory_id);
