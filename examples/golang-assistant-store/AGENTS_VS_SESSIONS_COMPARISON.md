# 📊 Perbandingan Kolom: `agents` vs `sessions`

## 🎯 Overview

Dokumen ini membandingkan struktur kolom antara tabel `agents` (template/configuration) dan `sessions` (chat instances).

---

## 📋 Side-by-Side Comparison

### **Common Columns (Ada di Kedua Tabel)**

| Column | `agents` | `sessions` | Notes |
|--------|----------|-----------|-------|
| `id` | text (PK) | text (PK) | Primary key, auto-generated |
| `slug` | varchar(100), unique | varchar(100), unique | URL-friendly identifier |
| `title` | varchar(255) | text | Display name |
| `description` | varchar(1000) | text | Description text |
| `avatar` | text | text | Avatar URL/emoji |
| `background_color` | text | text | Background color |
| `user_id` | text (FK, NOT NULL) | text (FK, NOT NULL) | Owner reference |
| `client_id` | text | text | Client sync identifier |
| `created_at` | timestamp | timestamp | Creation time |
| `updated_at` | timestamp | timestamp | Last update time |
| `accessed_at` | timestamp | timestamp | Last access time |

### **Unique to `agents` (Configuration/Template Data)**

| Column | Type | Default | Description |
|--------|------|---------|-------------|
| `system_role` | text | - | **System prompt** - defines agent behavior |
| `model` | text | - | AI model (e.g., "gpt-4", "claude-3") |
| `provider` | text | - | AI provider (e.g., "openai", "anthropic") |
| `chat_config` | jsonb | - | Chat configuration (temperature, max_tokens, etc.) |
| `params` | jsonb | `{}` | Additional parameters |
| `tts` | jsonb | - | Text-to-speech configuration |
| `few_shots` | jsonb | - | Few-shot examples for prompting |
| `plugins` | jsonb | `[]` | Array of plugin IDs |
| `tags` | jsonb | `[]` | Array of tags/categories |
| `virtual` | boolean | false | Whether it's a virtual agent |
| `opening_message` | text | - | Initial greeting message |
| `opening_questions` | text[] | `[]` | Suggested opening questions |

### **Unique to `sessions` (Instance/Conversation Data)**

| Column | Type | Default | Description |
|--------|------|---------|-------------|
| `type` | text (enum) | 'agent' | Session type: 'agent' or 'group' |
| `group_id` | text (FK) | null | Reference to session_groups |
| `pinned` | boolean | false | Whether session is pinned to top |

---

## 🔍 Detailed Column Analysis

### 1. **Identity Columns**

#### `id` (Both tables)
```typescript
// agents
id: text('id')
  .primaryKey()
  .$defaultFn(() => idGenerator('agents'))
  .notNull()

// sessions  
id: text('id')
  .$defaultFn(() => idGenerator('sessions'))
  .primaryKey()
```

**Purpose:**
- Unique identifier for each record
- Auto-generated with prefix (`agents_xxx`, `sessions_xxx`)

---

#### `slug` (Both tables)
```typescript
// agents
slug: varchar('slug', { length: 100 })
  .$defaultFn(() => randomSlug(4))
  .unique()

// sessions
slug: varchar('slug', { length: 100 })
  .notNull()
  .$defaultFn(() => randomSlug())
```

**Purpose:**
- URL-friendly identifier
- Used in routes: `/chat/[slug]`, `/agent/[slug]`
- **Difference:** agents slug is optional, sessions slug is required

---

### 2. **Display Metadata**

#### `title` (Both tables)
```typescript
// agents
title: varchar('title', { length: 255 })

// sessions
title: text('title')
```

**Purpose:**
- Display name shown in UI
- **Difference:** agents has 255 char limit, sessions unlimited

**Example:**
- Agent: "GPT-4 Code Assistant"
- Session: "Help with React Hooks Implementation"

---

#### `description` (Both tables)
```typescript
// agents
description: varchar('description', { length: 1000 })

// sessions
description: text('description')
```

**Purpose:**
- Longer description text
- **Difference:** agents limited to 1000 chars, sessions unlimited

---

#### `avatar` (Both tables)
```typescript
// Both identical
avatar: text('avatar')
```

**Purpose:**
- Avatar image URL or emoji
- Can be overridden in session

**Example:**
- Agent: "🤖"
- Session: "🚀" (custom for this conversation)

---

#### `background_color` (Both tables)
```typescript
// Both identical
backgroundColor: text('background_color')
```

**Purpose:**
- Background color for UI
- Hex color code or CSS color name

---

### 3. **Ownership & Sync**

#### `user_id` (Both tables)
```typescript
// Both identical
userId: text('user_id')
  .references(() => users.id, { onDelete: 'cascade' })
  .notNull()
```

**Purpose:**
- Owner of the agent/session
- CASCADE delete: if user deleted, all their agents/sessions deleted

---

#### `client_id` (Both tables)
```typescript
// Both identical
clientId: text('client_id')
```

**Purpose:**
- Client-side identifier for sync
- Used in PGLite ↔ Server sync

---

### 4. **Agent-Specific: AI Configuration**

#### `system_role` (agents only) ⭐
```typescript
systemRole: text('system_role')
```

**Purpose:**
- **THE MOST IMPORTANT COLUMN**
- System prompt that defines agent behavior
- Tells the AI "who" it is and how to behave

**Example:**
```
You are an expert React developer with 10 years of experience.
You provide clear, concise code examples and follow best practices.
Always explain your reasoning and suggest optimizations.
```

---

#### `model` (agents only)
```typescript
model: text('model')
```

**Purpose:**
- AI model identifier

**Examples:**
- `"gpt-4-turbo"`
- `"claude-3-opus"`
- `"gemini-pro"`

---

#### `provider` (agents only)
```typescript
provider: text('provider')
```

**Purpose:**
- AI provider name

**Examples:**
- `"openai"`
- `"anthropic"`
- `"google"`

---

#### `chat_config` (agents only)
```typescript
chatConfig: jsonb('chat_config').$type<LobeAgentChatConfig>()
```

**Purpose:**
- Chat configuration parameters

**Example JSON:**
```json
{
  "temperature": 0.7,
  "max_tokens": 2000,
  "top_p": 1,
  "frequency_penalty": 0,
  "presence_penalty": 0,
  "historyCount": 10
}
```

---

#### `params` (agents only)
```typescript
params: jsonb('params').default({})
```

**Purpose:**
- Additional custom parameters
- Flexible storage for future extensions

---

#### `tts` (agents only)
```typescript
tts: jsonb('tts').$type<LobeAgentTTSConfig>()
```

**Purpose:**
- Text-to-speech configuration

**Example JSON:**
```json
{
  "voice": "alloy",
  "speed": 1.0,
  "showAllLocaleVoice": false
}
```

---

#### `few_shots` (agents only)
```typescript
fewShots: jsonb('few_shots')
```

**Purpose:**
- Few-shot learning examples
- Provides examples to guide AI behavior

**Example JSON:**
```json
[
  {
    "user": "How do I center a div?",
    "assistant": "Use flexbox: display: flex; justify-content: center; align-items: center;"
  }
]
```

---

#### `plugins` (agents only)
```typescript
plugins: jsonb('plugins').$type<string[]>().default([])
```

**Purpose:**
- Array of enabled plugin IDs

**Example:**
```json
["web-search", "code-interpreter", "dalle-3"]
```

---

#### `tags` (agents only)
```typescript
tags: jsonb('tags').$type<string[]>().default([])
```

**Purpose:**
- Categorization tags

**Example:**
```json
["coding", "react", "frontend", "tutorial"]
```

---

#### `virtual` (agents only)
```typescript
virtual: boolean('virtual').default(false)
```

**Purpose:**
- Whether agent is virtual (system-generated)
- Virtual agents might have special behavior

---

#### `opening_message` (agents only)
```typescript
openingMessage: text('opening_message')
```

**Purpose:**
- First message shown when starting new chat

**Example:**
```
"Hi! I'm your React expert. How can I help you today?"
```

---

#### `opening_questions` (agents only)
```typescript
openingQuestions: text('opening_questions').array().default([])
```

**Purpose:**
- Suggested questions to start conversation

**Example:**
```json
[
  "How do I optimize React performance?",
  "What are the best practices for state management?",
  "Can you review my component code?"
]
```

---

### 5. **Session-Specific: Conversation Data**

#### `type` (sessions only)
```typescript
type: text('type', { enum: ['agent', 'group'] }).default('agent')
```

**Purpose:**
- Type of session
- `'agent'`: Single agent conversation
- `'group'`: Multi-agent group chat

---

#### `group_id` (sessions only)
```typescript
groupId: text('group_id').references(() => sessionGroups.id, { onDelete: 'set null' })
```

**Purpose:**
- Reference to session group (folder/category)
- Allows organizing sessions into groups

**Example:**
```
Group: "Work Projects"
  ├─ Session: "API Design Discussion"
  ├─ Session: "Code Review"
  └─ Session: "Bug Fixing"
```

---

#### `pinned` (sessions only)
```typescript
pinned: boolean('pinned').default(false)
```

**Purpose:**
- Whether session is pinned to top of list
- Pinned sessions appear first regardless of update time

---

### 6. **Timestamps**

#### `created_at` (Both tables)
```typescript
createdAt: timestamp with time zone DEFAULT now() NOT NULL
```

**Purpose:**
- When record was created
- Never changes

---

#### `updated_at` (Both tables)
```typescript
updatedAt: timestamp with time zone DEFAULT now() NOT NULL
```

**Purpose:**
- When record was last modified
- **Critical for sessions:** Used to sort conversations by recency

---

#### `accessed_at` (Both tables)
```typescript
accessedAt: timestamp with time zone DEFAULT now() NOT NULL
```

**Purpose:**
- When record was last accessed/viewed
- Used for analytics and cleanup

---

## 📊 Visual Comparison

```
┌─────────────────────────────────────────────────────────────────┐
│                         AGENTS TABLE                             │
├─────────────────────────────────────────────────────────────────┤
│ Identity:                                                        │
│   ✓ id, slug, title, description, avatar, background_color      │
│                                                                  │
│ AI Configuration (UNIQUE):                                       │
│   ⭐ system_role        ← The "brain" of the agent              │
│   ⭐ model, provider    ← Which AI to use                       │
│   ⭐ chat_config        ← How AI behaves                        │
│   ⭐ plugins            ← What tools it has                     │
│   ⭐ few_shots          ← Learning examples                     │
│   ⭐ tts                ← Voice settings                        │
│   ⭐ opening_message    ← First greeting                        │
│   ⭐ opening_questions  ← Suggested starts                      │
│   ⭐ tags               ← Categories                            │
│                                                                  │
│ Ownership:                                                       │
│   ✓ user_id, client_id                                          │
│                                                                  │
│ Timestamps:                                                      │
│   ✓ created_at, updated_at, accessed_at                         │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                        SESSIONS TABLE                            │
├─────────────────────────────────────────────────────────────────┤
│ Identity:                                                        │
│   ✓ id, slug, title, description, avatar, background_color      │
│                                                                  │
│ Conversation Data (UNIQUE):                                      │
│   ⭐ type               ← 'agent' or 'group'                    │
│   ⭐ group_id           ← Folder/category                       │
│   ⭐ pinned             ← Pin to top                            │
│                                                                  │
│ Ownership:                                                       │
│   ✓ user_id, client_id                                          │
│                                                                  │
│ Timestamps:                                                      │
│   ✓ created_at, updated_at, accessed_at                         │
│                                                                  │
│ Related Data (via relations):                                    │
│   → messages (1:N)      ← Chat history                          │
│   → topics (1:N)        ← Conversation topics                   │
│   → files (N:M)         ← Attached files                        │
│   → agent (N:1)         ← Which agent is used                   │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🎯 Key Insights

### 1. **Data Separation Philosophy**

| Aspect | `agents` | `sessions` |
|--------|----------|-----------|
| **What** | Configuration | Instance |
| **Reusability** | Reusable template | One-time use |
| **Changes** | Rarely (by design) | Frequently (every message) |
| **Size** | Small (~10 columns config) | Large (with messages) |

### 2. **Column Count**

```
agents:   ~20 columns (mostly configuration)
sessions: ~13 columns (mostly metadata)
```

**Why agents has more columns?**
- Needs to store complete AI configuration
- Self-contained template
- No dependency on other tables for core functionality

**Why sessions has fewer columns?**
- Relies on relations (messages, agent)
- Metadata-focused
- Heavy data is in related tables

### 3. **JSONB Usage**

**agents (Heavy JSONB usage):**
- `chat_config` - Complex configuration object
- `params` - Flexible parameters
- `tts` - Voice configuration
- `few_shots` - Array of examples
- `plugins` - Array of plugin IDs
- `tags` - Array of tags

**sessions (Minimal JSONB usage):**
- None! All simple types

**Why?**
- Agents need flexible configuration
- Sessions are simpler, more structured

### 4. **Most Important Columns**

**For agents:**
1. `system_role` - Defines agent personality/behavior
2. `model` - Which AI model to use
3. `chat_config` - How AI responds

**For sessions:**
1. `updated_at` - For sorting conversations
2. `pinned` - For prioritizing important chats
3. `type` - For group chat vs single agent

---

## 💡 Practical Examples

### Example 1: Creating an Agent

```sql
INSERT INTO agents (
  id, slug, title, description,
  system_role, model, provider,
  chat_config, tags, plugins,
  opening_message, opening_questions,
  user_id
) VALUES (
  'agent_abc123',
  'react-expert',
  'React Expert',
  'Expert in React development and best practices',
  'You are an expert React developer...',
  'gpt-4-turbo',
  'openai',
  '{"temperature": 0.7, "max_tokens": 2000}',
  '["react", "frontend", "javascript"]',
  '["web-search"]',
  'Hi! I''m your React expert.',
  '["How to optimize performance?", "Best state management?"]',
  'user_xyz'
);
```

### Example 2: Creating a Session

```sql
INSERT INTO sessions (
  id, slug, title, description,
  type, group_id, pinned,
  user_id
) VALUES (
  'sess_def456',
  'react-hooks-help',
  'Help with React Hooks',
  'Discussion about useEffect and custom hooks',
  'agent',
  'group_work',
  false,
  'user_xyz'
);

-- Link session to agent
INSERT INTO agents_to_sessions (
  agent_id, session_id, user_id
) VALUES (
  'agent_abc123',
  'sess_def456',
  'user_xyz'
);
```

---

## 🔗 Relationship Summary

```
agents (1) ←→ (N) agents_to_sessions (N) ←→ (1) sessions
                                                    ↓
                                                messages (N)
                                                    ↓
                                                topics (N)
```

**Key Point:**
- Agent stores "HOW to chat" (configuration)
- Session stores "WHAT was chatted" (history)

---

## 📝 Conclusion

| Question | Answer |
|----------|--------|
| **Which has more columns?** | `agents` (~20) > `sessions` (~13) |
| **Which uses more JSONB?** | `agents` (6 JSONB columns) > `sessions` (0 JSONB) |
| **Which changes more often?** | `sessions` (every message) > `agents` (rarely) |
| **Which is more complex?** | `agents` (AI config) > `sessions` (metadata) |
| **Which is more important for chat list?** | `sessions` (shows conversations) |
| **Which is more important for AI behavior?** | `agents` (defines how AI responds) |

**Final Insight:**
- **`agents`** = Configuration-heavy, reusable template
- **`sessions`** = Metadata-light, conversation instance

Think of it like:
- **Agent** = Recipe (ingredients, instructions)
- **Session** = Meal (what you actually cooked and ate)

🎯 You query `sessions` in `/chat` because users want to see their **meals** (conversations), not **recipes** (agent templates)!

