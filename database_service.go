package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type DatabaseService struct{}

const initialSchemaSQL = `-- Initial SQLite migration generated from schemas
-- This file contains the initial database setup for Wails SQLite

-- Users table
CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY NOT NULL,
  username TEXT UNIQUE,
  email TEXT,
  avatar TEXT,
  phone TEXT,
  first_name TEXT,
  last_name TEXT,
  is_onboarded INTEGER DEFAULT 0 NOT NULL, -- boolean as integer
  clerk_created_at INTEGER, -- timestamp_ms
  email_verified_at INTEGER, -- timestamp_ms
  preference TEXT DEFAULT '{}', -- JSON as text
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000), -- timestamp_ms
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000) -- timestamp_ms
);

-- User settings
CREATE TABLE IF NOT EXISTS user_settings (
  id TEXT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
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
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  identifier TEXT NOT NULL,
  type TEXT NOT NULL CHECK (type IN ('plugin', 'customPlugin')),
  manifest TEXT, -- JSON as text
  settings TEXT, -- JSON as text
  custom_params TEXT, -- JSON as text
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  PRIMARY KEY (user_id, identifier)
);

-- Global files table
CREATE TABLE IF NOT EXISTS global_files (
  hash_id TEXT PRIMARY KEY,
  file_type TEXT NOT NULL,
  size INTEGER NOT NULL,
  url TEXT NOT NULL,
  metadata TEXT, -- JSON as text
  creator TEXT NOT NULL REFERENCES users(id) ON DELETE SET NULL,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  accessed_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- Files table
CREATE TABLE IF NOT EXISTS files (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  file_type TEXT NOT NULL,
  file_hash TEXT REFERENCES global_files(hash_id),
  name TEXT NOT NULL,
  size INTEGER NOT NULL,
  url TEXT NOT NULL,
  source TEXT, -- JSON as text
  client_id TEXT,
  metadata TEXT, -- JSON as text
  chunk_task_id TEXT REFERENCES async_tasks(id) ON DELETE SET NULL,
  embedding_task_id TEXT REFERENCES async_tasks(id) ON DELETE SET NULL,
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
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  client_id TEXT,
  is_public INTEGER DEFAULT 0 NOT NULL, -- boolean as integer
  settings TEXT, -- JSON as text
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- Knowledge base files junction table
CREATE TABLE IF NOT EXISTS knowledge_base_files (
  knowledge_base_id TEXT NOT NULL REFERENCES knowledge_bases(id) ON DELETE CASCADE,
  file_id TEXT NOT NULL REFERENCES files(id) ON DELETE CASCADE,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  PRIMARY KEY (knowledge_base_id, file_id)
);

-- Agents table
CREATE TABLE IF NOT EXISTS agents (
  id TEXT PRIMARY KEY,
  slug TEXT UNIQUE,
  title TEXT,
  description TEXT,
  tags TEXT DEFAULT '[]', -- JSON as text
  avatar TEXT,
  background_color TEXT,
  plugins TEXT DEFAULT '[]', -- JSON as text
  client_id TEXT,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
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
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  client_id TEXT,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- Sessions table
CREATE TABLE IF NOT EXISTS sessions (
  id TEXT PRIMARY KEY,
  slug TEXT NOT NULL,
  title TEXT,
  description TEXT,
  avatar TEXT,
  background_color TEXT,
  type TEXT DEFAULT 'agent',
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  group_id TEXT REFERENCES session_groups(id) ON DELETE SET NULL,
  client_id TEXT,
  pinned INTEGER DEFAULT 0 NOT NULL, -- boolean as integer
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  UNIQUE(slug, user_id)
);

-- Topics table
CREATE TABLE IF NOT EXISTS topics (
  id TEXT PRIMARY KEY,
  title TEXT,
  favorite INTEGER DEFAULT 0 NOT NULL, -- boolean as integer
  session_id TEXT REFERENCES sessions(id) ON DELETE CASCADE,
  group_id TEXT REFERENCES chat_groups(id) ON DELETE CASCADE,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  client_id TEXT,
  history_summary TEXT, -- JSON as text
  metadata TEXT, -- JSON as text
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- Threads table
CREATE TABLE IF NOT EXISTS threads (
  id TEXT PRIMARY KEY,
  title TEXT,
  type TEXT NOT NULL CHECK (type IN ('continuation', 'standalone')),
  status TEXT DEFAULT 'active' CHECK (status IN ('active', 'deprecated', 'archived')),
  topic_id TEXT NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
  source_message_id TEXT NOT NULL,
  parent_thread_id TEXT REFERENCES threads(id) ON DELETE SET NULL,
  client_id TEXT,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
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
  client_id TEXT,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
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
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  UNIQUE(client_id, user_id)
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
  client_id TEXT,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  UNIQUE(client_id, user_id)
);

-- Message TTS
CREATE TABLE IF NOT EXISTS message_tts (
  id TEXT PRIMARY KEY REFERENCES messages(id) ON DELETE CASCADE,
  content_md5 TEXT,
  file_id TEXT REFERENCES files(id) ON DELETE CASCADE,
  voice TEXT,
  client_id TEXT,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  UNIQUE(client_id, user_id)
);

-- Message translates
CREATE TABLE IF NOT EXISTS message_translates (
  id TEXT PRIMARY KEY REFERENCES messages(id) ON DELETE CASCADE,
  content TEXT, -- JSON as text
  "from" TEXT,
  "to" TEXT,
  client_id TEXT,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  UNIQUE(client_id, user_id)
);

-- Message queries
CREATE TABLE IF NOT EXISTS message_queries (
  id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  message_id TEXT NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
  rewrite_query TEXT,
  user_query TEXT,
  client_id TEXT,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  embeddings_id TEXT REFERENCES embeddings(id) ON DELETE SET NULL,
  UNIQUE(client_id, user_id)
);

-- Message query chunks
CREATE TABLE IF NOT EXISTS message_query_chunks (
  message_id TEXT REFERENCES messages(id) ON DELETE CASCADE,
  query_id TEXT REFERENCES message_queries(id) ON DELETE CASCADE,
  chunk_id TEXT REFERENCES chunks(id) ON DELETE CASCADE,
  similarity INTEGER, -- We'll store similarity as integer for now
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  PRIMARY KEY (chunk_id, message_id, query_id)
);

-- Message chunks
CREATE TABLE IF NOT EXISTS message_chunks (
  message_id TEXT REFERENCES messages(id) ON DELETE CASCADE,
  chunk_id TEXT REFERENCES chunks(id) ON DELETE CASCADE,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  PRIMARY KEY (chunk_id, message_id)
);

-- Messages files junction table
CREATE TABLE IF NOT EXISTS messages_files (
  file_id TEXT NOT NULL REFERENCES files(id) ON DELETE CASCADE,
  message_id TEXT NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  PRIMARY KEY (file_id, message_id)
);

-- Chunks table (RAG)
CREATE TABLE IF NOT EXISTS chunks (
  id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  text TEXT,
  abstract TEXT,
  metadata TEXT, -- JSON as text
  chunk_index INTEGER,
  type TEXT,
  client_id TEXT,
  user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  UNIQUE(client_id, user_id)
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
  client_id TEXT,
  user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
  file_id TEXT REFERENCES files(id) ON DELETE CASCADE,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- Embeddings table
CREATE TABLE IF NOT EXISTS embeddings (
  id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  chunk_id TEXT UNIQUE REFERENCES chunks(id) ON DELETE CASCADE,
  embeddings BLOB, -- Store embeddings as binary data
  model TEXT,
  client_id TEXT,
  user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
  UNIQUE(client_id, user_id)
);

-- File chunks junction table
CREATE TABLE IF NOT EXISTS file_chunks (
  file_id TEXT REFERENCES files(id) ON DELETE CASCADE,
  chunk_id TEXT REFERENCES chunks(id) ON DELETE CASCADE,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  PRIMARY KEY (file_id, chunk_id)
);

-- Agents to sessions junction table
CREATE TABLE IF NOT EXISTS agents_to_sessions (
  agent_id TEXT NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
  session_id TEXT NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  PRIMARY KEY (agent_id, session_id)
);

-- Files to sessions junction table
CREATE TABLE IF NOT EXISTS files_to_sessions (
  file_id TEXT NOT NULL REFERENCES files(id) ON DELETE CASCADE,
  session_id TEXT NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  PRIMARY KEY (file_id, session_id)
);

-- Agents knowledge bases junction table
CREATE TABLE IF NOT EXISTS agents_knowledge_bases (
  agent_id TEXT NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
  knowledge_base_id TEXT NOT NULL REFERENCES knowledge_bases(id) ON DELETE CASCADE,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
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
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  PRIMARY KEY (file_id, agent_id, user_id)
);

-- Topic documents junction table
CREATE TABLE IF NOT EXISTS topic_documents (
  document_id TEXT NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
  topic_id TEXT NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  PRIMARY KEY (document_id, topic_id)
);

-- Document chunks junction table
CREATE TABLE IF NOT EXISTS document_chunks (
  document_id TEXT NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
  chunk_id TEXT NOT NULL REFERENCES chunks(id) ON DELETE CASCADE,
  page_index INTEGER,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  PRIMARY KEY (document_id, chunk_id)
);

-- Chat groups table
CREATE TABLE IF NOT EXISTS chat_groups (
  id TEXT PRIMARY KEY,
  title TEXT,
  description TEXT,
  config TEXT, -- JSON as text
  client_id TEXT,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  pinned INTEGER DEFAULT 0 NOT NULL, -- boolean as integer
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- Chat groups agents junction table
CREATE TABLE IF NOT EXISTS chat_groups_agents (
  chat_group_id TEXT NOT NULL REFERENCES chat_groups(id) ON DELETE CASCADE,
  agent_id TEXT NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  enabled INTEGER DEFAULT 1 NOT NULL, -- boolean as integer
  agent_order INTEGER DEFAULT 0,
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
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  parent_group_id TEXT REFERENCES message_groups(id) ON DELETE CASCADE,
  client_id TEXT,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- AI providers table
CREATE TABLE IF NOT EXISTS ai_providers (
  id TEXT NOT NULL,
  name TEXT,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
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
  PRIMARY KEY (id, user_id)
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
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  pricing TEXT, -- JSON as text
  parameters TEXT DEFAULT '{}', -- JSON as text
  config TEXT, -- JSON as text
  abilities TEXT DEFAULT '{}', -- JSON as text
  context_window_tokens INTEGER,
  source TEXT,
  released_at TEXT,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  PRIMARY KEY (id, provider_id, user_id)
);

-- Async tasks table
CREATE TABLE IF NOT EXISTS async_tasks (
  id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  type TEXT,
  status TEXT,
  error TEXT,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  duration INTEGER,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- Generation topics table
CREATE TABLE IF NOT EXISTS generation_topics (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title TEXT,
  cover_url TEXT,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- Generation batches table
CREATE TABLE IF NOT EXISTS generation_batches (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
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
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
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
  id TEXT PRIMARY KEY,
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
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  client_id TEXT,
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
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- RBAC tables
CREATE TABLE IF NOT EXISTS rbac_roles (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE,
  display_name TEXT NOT NULL,
  description TEXT,
  is_system INTEGER DEFAULT 0 NOT NULL, -- boolean as integer
  is_active INTEGER DEFAULT 1 NOT NULL, -- boolean as integer
  metadata TEXT DEFAULT '{}', -- JSON as text
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

CREATE TABLE IF NOT EXISTS rbac_permissions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  code TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  description TEXT,
  category TEXT NOT NULL,
  is_active INTEGER DEFAULT 1 NOT NULL, -- boolean as integer
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

CREATE TABLE IF NOT EXISTS rbac_role_permissions (
  role_id INTEGER NOT NULL REFERENCES rbac_roles(id) ON DELETE CASCADE,
  permission_id INTEGER NOT NULL REFERENCES rbac_permissions(id) ON DELETE CASCADE,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE IF NOT EXISTS rbac_user_roles (
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role_id INTEGER NOT NULL REFERENCES rbac_roles(id) ON DELETE CASCADE,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  expires_at INTEGER, -- timestamp_ms
  PRIMARY KEY (user_id, role_id)
);

-- NextAuth tables
CREATE TABLE IF NOT EXISTS nextauth_accounts (
  access_token TEXT,
  expires_at INTEGER,
  id_token TEXT,
  provider TEXT NOT NULL,
  provider_account_id TEXT NOT NULL,
  refresh_token TEXT,
  scope TEXT,
  session_state TEXT,
  token_type TEXT,
  type TEXT NOT NULL,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  PRIMARY KEY (provider, provider_account_id)
);

CREATE TABLE IF NOT EXISTS nextauth_sessions (
  expires INTEGER NOT NULL,
  session_token TEXT PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS nextauth_verificationtokens (
  expires INTEGER NOT NULL,
  identifier TEXT NOT NULL,
  token TEXT NOT NULL,
  PRIMARY KEY (identifier, token)
);

CREATE TABLE IF NOT EXISTS nextauth_authenticators (
  counter INTEGER NOT NULL,
  credential_backed_up INTEGER NOT NULL,
  credential_device_type TEXT NOT NULL,
  credential_id TEXT NOT NULL UNIQUE,
  credential_public_key TEXT NOT NULL,
  provider_account_id TEXT NOT NULL,
  transports TEXT,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  PRIMARY KEY (user_id, credential_id)
);

-- OIDC tables
CREATE TABLE IF NOT EXISTS oidc_authorization_codes (
  id TEXT PRIMARY KEY,
  data TEXT NOT NULL,
  expires_at INTEGER NOT NULL,
  consumed_at INTEGER,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  client_id TEXT NOT NULL,
  grant_id TEXT,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

CREATE TABLE IF NOT EXISTS oidc_access_tokens (
  id TEXT PRIMARY KEY,
  data TEXT NOT NULL,
  expires_at INTEGER NOT NULL,
  consumed_at INTEGER,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  client_id TEXT NOT NULL,
  grant_id TEXT,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

CREATE TABLE IF NOT EXISTS oidc_refresh_tokens (
  id TEXT PRIMARY KEY,
  data TEXT NOT NULL,
  expires_at INTEGER NOT NULL,
  consumed_at INTEGER,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  client_id TEXT NOT NULL,
  grant_id TEXT,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

CREATE TABLE IF NOT EXISTS oidc_device_codes (
  id TEXT PRIMARY KEY,
  data TEXT NOT NULL,
  expires_at INTEGER NOT NULL,
  consumed_at INTEGER,
  user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
  client_id TEXT NOT NULL,
  grant_id TEXT,
  user_code TEXT,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

CREATE TABLE IF NOT EXISTS oidc_interactions (
  id TEXT PRIMARY KEY,
  data TEXT NOT NULL,
  expires_at INTEGER NOT NULL,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

CREATE TABLE IF NOT EXISTS oidc_grants (
  id TEXT PRIMARY KEY,
  data TEXT NOT NULL,
  expires_at INTEGER NOT NULL,
  consumed_at INTEGER,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  client_id TEXT NOT NULL,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

CREATE TABLE IF NOT EXISTS oidc_clients (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  client_secret TEXT,
  redirect_uris TEXT NOT NULL,
  grants TEXT NOT NULL,
  response_types TEXT NOT NULL,
  scopes TEXT NOT NULL,
  token_endpoint_auth_method TEXT,
  application_type TEXT,
  client_uri TEXT,
  logo_uri TEXT,
  policy_uri TEXT,
  tos_uri TEXT,
  is_first_party INTEGER DEFAULT 0 NOT NULL, -- boolean as integer
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

CREATE TABLE IF NOT EXISTS oidc_sessions (
  id TEXT PRIMARY KEY,
  data TEXT NOT NULL,
  expires_at INTEGER NOT NULL,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

CREATE TABLE IF NOT EXISTS oidc_consents (
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  client_id TEXT NOT NULL REFERENCES oidc_clients(id) ON DELETE CASCADE,
  scopes TEXT NOT NULL,
  expires_at INTEGER,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  PRIMARY KEY (user_id, client_id)
);

-- OAuth handoffs
CREATE TABLE IF NOT EXISTS oauth_handoffs (
  id TEXT PRIMARY KEY,
  client TEXT NOT NULL,
  payload TEXT NOT NULL, -- JSON as text
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now') * 1000)
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);
CREATE INDEX IF NOT EXISTS idx_messages_user_id ON messages(user_id);
CREATE INDEX IF NOT EXISTS idx_messages_session_id ON messages(session_id);
CREATE INDEX IF NOT EXISTS idx_messages_topic_id ON messages(topic_id);
CREATE INDEX IF NOT EXISTS idx_messages_parent_id ON messages(parent_id);
CREATE INDEX IF NOT EXISTS idx_messages_quota_id ON messages(quota_id);
CREATE INDEX IF NOT EXISTS idx_messages_thread_id ON messages(thread_id);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_id_user_id ON sessions(id, user_id);

CREATE INDEX IF NOT EXISTS idx_topics_user_id ON topics(user_id);
CREATE INDEX IF NOT EXISTS idx_topics_id_user_id ON topics(id, user_id);

CREATE INDEX IF NOT EXISTS idx_chunks_user_id ON chunks(user_id);
CREATE INDEX IF NOT EXISTS idx_chunks_client_id_user_id ON chunks(client_id, user_id);

CREATE INDEX IF NOT EXISTS idx_unstructured_chunks_client_id_user_id ON unstructured_chunks(client_id, user_id);

CREATE INDEX IF NOT EXISTS idx_embeddings_client_id_user_id ON embeddings(client_id, user_id);
CREATE INDEX IF NOT EXISTS idx_embeddings_chunk_id ON embeddings(chunk_id);

CREATE INDEX IF NOT EXISTS idx_files_file_hash ON files(file_hash);
CREATE INDEX IF NOT EXISTS idx_files_client_id_user_id ON files(client_id, user_id);

CREATE INDEX IF NOT EXISTS idx_agents_client_id_user_id ON agents(client_id, user_id);
CREATE INDEX IF NOT EXISTS idx_agents_title ON agents(title);
CREATE INDEX IF NOT EXISTS idx_agents_description ON agents(description);

CREATE INDEX IF NOT EXISTS idx_knowledge_bases_client_id_user_id ON knowledge_bases(client_id, user_id);

CREATE INDEX IF NOT EXISTS idx_documents_source ON documents(source);
CREATE INDEX IF NOT EXISTS idx_documents_file_type ON documents(file_type);
CREATE INDEX IF NOT EXISTS idx_documents_file_id ON documents(file_id);

CREATE INDEX IF NOT EXISTS idx_generation_topics_user_id ON generation_topics(user_id);
CREATE INDEX IF NOT EXISTS idx_generation_batches_user_id ON generation_batches(user_id);
CREATE INDEX IF NOT EXISTS idx_generations_user_id ON generations(user_id);

CREATE INDEX IF NOT EXISTS idx_rbac_role_permissions_role_id ON rbac_role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_rbac_role_permissions_permission_id ON rbac_role_permissions(permission_id);
CREATE INDEX IF NOT EXISTS idx_rbac_user_roles_user_id ON rbac_user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_rbac_user_roles_role_id ON rbac_user_roles(role_id);

-- Drizzle migration tracking table
CREATE TABLE IF NOT EXISTS __drizzle_migrations (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  hash TEXT NOT NULL,
  created_at INTEGER
);

-- Insert initial migration record
INSERT OR IGNORE INTO __drizzle_migrations (hash, created_at) VALUES ('initial_sqlite_setup', strftime('%s', 'now') * 1000);`

func (d *DatabaseService) InitializeDatabase() error {
	// Get user config directory for database location
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Printf("Warning: Could not get user config dir, using current directory: %v", err)
		userConfigDir = "."
	}

	// Create app data directory
	appDataDir := filepath.Join(userConfigDir, "veridium")
	if err := os.MkdirAll(appDataDir, 0755); err != nil {
		return err
	}

	// Database file path
	dbPath := filepath.Join(appDataDir, "veridium.db")

	// Open SQLite database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Execute schema migration
	_, err = db.Exec(initialSchemaSQL)
	if err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

func (d *DatabaseService) GetDatabasePath() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Printf("Warning: Could not get user config dir, using current directory: %v", err)
		userConfigDir = "."
	}

	appDataDir := filepath.Join(userConfigDir, "veridium")
	if err := os.MkdirAll(appDataDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(appDataDir, "veridium.db"), nil
}
