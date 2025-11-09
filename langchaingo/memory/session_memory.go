package memory

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kawai-network/veridium/langchaingo/llms"
	"github.com/kawai-network/veridium/langchaingo/schema"
)

// SessionType represents the type of a session
type SessionType string

const (
	// SessionTypeAgent represents an agent/assistant session
	SessionTypeAgent SessionType = "agent"
	// SessionTypeGroup represents a group session (future use)
	SessionTypeGroup SessionType = "group"
)

// Session represents a chat session.
type Session struct {
	ID        string      `json:"id"`
	Title     string      `json:"title"`
	Type      SessionType `json:"type"` // agent or group
	Config    AgentConfig `json:"config"`
	Meta      MetaData    `json:"meta"`
	GroupID   *string     `json:"groupId,omitempty"`
	Pinned    bool        `json:"pinned"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

// SessionMemory manages both sessions and their associated message history.
type SessionMemory struct {
	db        *sql.DB
	tableName string
}

// NewSessionMemory creates a new session memory manager.
func NewSessionMemory(db *sql.DB, tableName string) (*SessionMemory, error) {
	sm := &SessionMemory{
		db:        db,
		tableName: tableName,
	}

	// NOTE: Tables are created by database migrations in main.go (001_core_tables.sql)
	// This initSchema() serves as a fallback safety check for standalone usage.
	// Uses CREATE TABLE IF NOT EXISTS, so safe to call even if tables exist.
	if err := sm.initSchema(); err != nil {
		return nil, err
	}

	return sm, nil
}

// initSchema creates the necessary tables.
// NOTE: This is a fallback. Primary table creation happens via migrations (001_core_tables.sql).
// Safe to call multiple times due to CREATE TABLE IF NOT EXISTS.
func (sm *SessionMemory) initSchema() error {
	schema := fmt.Sprintf(`
		-- Session groups table
		CREATE TABLE IF NOT EXISTS session_groups (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			sort INTEGER,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		-- Sessions table
		CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			type TEXT NOT NULL DEFAULT 'agent',
			config TEXT,
			meta TEXT,
			group_id TEXT,
			pinned BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (group_id) REFERENCES session_groups(id) ON DELETE SET NULL
		);

		-- Enhance langchaingo_messages with session reference
		CREATE TABLE IF NOT EXISTS %s (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session_id TEXT NOT NULL,
			content TEXT NOT NULL,
			type TEXT NOT NULL,
			metadata TEXT,
			created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
		);

		-- Indexes
		CREATE INDEX IF NOT EXISTS idx_session_groups_sort ON session_groups(sort ASC, created_at ASC);
		CREATE INDEX IF NOT EXISTS idx_sessions_updated ON sessions(updated_at DESC);
		CREATE INDEX IF NOT EXISTS idx_sessions_pinned ON sessions(pinned DESC, updated_at DESC);
		CREATE INDEX IF NOT EXISTS idx_sessions_group ON sessions(group_id);
		CREATE INDEX IF NOT EXISTS idx_messages_session ON %s(session_id, created);
	`, sm.tableName, sm.tableName)

	_, err := sm.db.Exec(schema)
	return err
}

// CreateSession creates a new session.
func (sm *SessionMemory) CreateSession(ctx context.Context, title string, config *AgentConfig, meta *MetaData) (*Session, error) {
	session := &Session{
		ID:        uuid.New().String(),
		Title:     title,
		Type:      SessionTypeAgent, // Default to agent type
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Use provided config or default
	if config != nil {
		session.Config = *config
	} else {
		session.Config = DefaultAgentConfig()
	}

	// Use provided meta or default
	if meta != nil {
		session.Meta = *meta
	} else {
		session.Meta = DefaultMetaData()
	}

	// Marshal config and meta for database storage
	configJSON, err := json.Marshal(session.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal agent config: %w", err)
	}

	metaJSON, err := json.Marshal(session.Meta)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal meta: %w", err)
	}

	_, err = sm.db.ExecContext(ctx, `
		INSERT INTO sessions (id, title, type, config, meta, pinned, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, session.ID, session.Title, session.Type, string(configJSON), string(metaJSON),
		session.Pinned, session.CreatedAt, session.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return session, nil
}

// GetSessions returns all sessions.
func (sm *SessionMemory) GetSessions(ctx context.Context) ([]Session, error) {
	rows, err := sm.db.QueryContext(ctx, `
		SELECT id, title, type, config, meta, group_id, pinned, created_at, updated_at
		FROM sessions
		ORDER BY pinned DESC, updated_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		var s Session
		var configJSON, metaJSON string
		err := rows.Scan(&s.ID, &s.Title, &s.Type, &configJSON, &metaJSON,
			&s.GroupID, &s.Pinned, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// Unmarshal config
		if err := json.Unmarshal([]byte(configJSON), &s.Config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config for session %s: %w", s.ID, err)
		}

		// Unmarshal meta
		if err := json.Unmarshal([]byte(metaJSON), &s.Meta); err != nil {
			return nil, fmt.Errorf("failed to unmarshal meta for session %s: %w", s.ID, err)
		}

		sessions = append(sessions, s)
	}

	return sessions, nil
}

// GetSession returns a specific session.
func (sm *SessionMemory) GetSession(ctx context.Context, id string) (*Session, error) {
	var s Session
	var configJSON, metaJSON string
	err := sm.db.QueryRowContext(ctx, `
		SELECT id, title, type, config, meta, group_id, pinned, created_at, updated_at
		FROM sessions WHERE id = ?
	`, id).Scan(&s.ID, &s.Title, &s.Type, &configJSON, &metaJSON,
		&s.GroupID, &s.Pinned, &s.CreatedAt, &s.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Unmarshal config
	if err := json.Unmarshal([]byte(configJSON), &s.Config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config for session %s: %w", s.ID, err)
	}

	// Unmarshal meta
	if err := json.Unmarshal([]byte(metaJSON), &s.Meta); err != nil {
		return nil, fmt.Errorf("failed to unmarshal meta for session %s: %w", s.ID, err)
	}

	return &s, nil
}

// UpdateSessionParams represents parameters for updating a session
type UpdateSessionParams struct {
	Title   *string      `json:"title,omitempty"`
	Config  *AgentConfig `json:"config,omitempty"`
	Meta    *MetaData    `json:"meta,omitempty"`
	GroupID *string      `json:"groupId,omitempty"`
}

// UpdateSession updates a session (supports partial updates).
func (sm *SessionMemory) UpdateSession(ctx context.Context, id string, params UpdateSessionParams) error {
	if id == "" {
		return fmt.Errorf("session ID cannot be empty")
	}

	updates := []string{"updated_at = ?"}
	args := []interface{}{time.Now()}

	if params.Title != nil {
		if *params.Title == "" {
			return fmt.Errorf("session title cannot be empty")
		}
		if len(*params.Title) > 500 {
			return fmt.Errorf("session title too long (max 500 characters)")
		}
		updates = append(updates, "title = ?")
		args = append(args, *params.Title)
	}

	if params.Config != nil {
		configJSON, err := json.Marshal(*params.Config)
		if err != nil {
			return fmt.Errorf("failed to marshal config: %w", err)
		}
		updates = append(updates, "config = ?")
		args = append(args, string(configJSON))
	}

	if params.Meta != nil {
		metaJSON, err := json.Marshal(*params.Meta)
		if err != nil {
			return fmt.Errorf("failed to marshal meta: %w", err)
		}
		updates = append(updates, "meta = ?")
		args = append(args, string(metaJSON))
	}

	if params.GroupID != nil {
		updates = append(updates, "group_id = ?")
		args = append(args, *params.GroupID)
	}

	args = append(args, id)

	query := "UPDATE sessions SET "
	for i, update := range updates {
		if i > 0 {
			query += ", "
		}
		query += update
	}
	query += " WHERE id = ?"

	result, err := sm.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("session not found: %s", id)
	}

	return nil
}

// DeleteSession deletes a session and its messages.
func (sm *SessionMemory) DeleteSession(ctx context.Context, id string) error {
	// Messages will be cascade deleted due to foreign key
	_, err := sm.db.ExecContext(ctx, "DELETE FROM sessions WHERE id = ?", id)
	return err
}

// PinSession pins/unpins a session.
func (sm *SessionMemory) PinSession(ctx context.Context, id string, pinned bool) error {
	_, err := sm.db.ExecContext(ctx, `
		UPDATE sessions SET pinned = ?, updated_at = ? WHERE id = ?
	`, pinned, time.Now(), id)
	return err
}

// GetChatHistory returns a ChatMessageHistory for a specific session.
func (sm *SessionMemory) GetChatHistory(sessionID string) schema.ChatMessageHistory {
	return NewSqlite3ChatMessageHistory(
		sm.db,
		sm.tableName,
		sessionID,
	)
}

// Sqlite3ChatMessageHistory is a lightweight wrapper around SqliteChatMessageHistory.
type Sqlite3ChatMessageHistory struct {
	db        *sql.DB
	tableName string
	sessionID string
}

// NewSqlite3ChatMessageHistory creates a new chat message history for a session.
func NewSqlite3ChatMessageHistory(db *sql.DB, tableName, sessionID string) *Sqlite3ChatMessageHistory {
	return &Sqlite3ChatMessageHistory{
		db:        db,
		tableName: tableName,
		sessionID: sessionID,
	}
}

// Messages returns all messages for this session.
func (h *Sqlite3ChatMessageHistory) Messages(ctx context.Context) ([]llms.ChatMessage, error) {
	rows, err := h.db.QueryContext(ctx, fmt.Sprintf(`
		SELECT content, type FROM %s 
		WHERE session_id = ? 
		ORDER BY created ASC
	`, h.tableName), h.sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []llms.ChatMessage
	for rows.Next() {
		var content, msgType string
		if err := rows.Scan(&content, &msgType); err != nil {
			return nil, err
		}

		var msg llms.ChatMessage
		switch msgType {
		case "human":
			msg = llms.HumanChatMessage{Content: content}
		case "ai":
			msg = llms.AIChatMessage{Content: content}
		case "system":
			msg = llms.SystemChatMessage{Content: content}
		default:
			msg = llms.GenericChatMessage{Content: content, Role: msgType}
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// AddMessage adds a message to the history.
func (h *Sqlite3ChatMessageHistory) AddMessage(ctx context.Context, message llms.ChatMessage) error {
	_, err := h.db.ExecContext(ctx, fmt.Sprintf(`
		INSERT INTO %s (session_id, content, type, created)
		VALUES (?, ?, ?, ?)
	`, h.tableName), h.sessionID, message.GetContent(), message.GetType(), time.Now())
	return err
}

// AddUserMessage adds a user message.
func (h *Sqlite3ChatMessageHistory) AddUserMessage(ctx context.Context, text string) error {
	return h.AddMessage(ctx, llms.HumanChatMessage{Content: text})
}

// AddAIMessage adds an AI message.
func (h *Sqlite3ChatMessageHistory) AddAIMessage(ctx context.Context, text string) error {
	return h.AddMessage(ctx, llms.AIChatMessage{Content: text})
}

// Clear removes all messages for this session.
func (h *Sqlite3ChatMessageHistory) Clear(ctx context.Context) error {
	_, err := h.db.ExecContext(ctx, fmt.Sprintf(`
		DELETE FROM %s WHERE session_id = ?
	`, h.tableName), h.sessionID)
	return err
}

// SetMessages replaces all messages for this session.
func (h *Sqlite3ChatMessageHistory) SetMessages(ctx context.Context, messages []llms.ChatMessage) error {
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Clear existing messages
	if _, err := tx.ExecContext(ctx, fmt.Sprintf(`
		DELETE FROM %s WHERE session_id = ?
	`, h.tableName), h.sessionID); err != nil {
		return err
	}

	// Insert new messages
	for _, msg := range messages {
		if _, err := tx.ExecContext(ctx, fmt.Sprintf(`
			INSERT INTO %s (session_id, content, type, created)
			VALUES (?, ?, ?, ?)
		`, h.tableName), h.sessionID, msg.GetContent(), msg.GetType(), time.Now()); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// ============================================================================
// SESSION GROUPS
// ============================================================================

// SessionGroup represents a group of sessions for organization
type SessionGroup struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Sort      *int      `json:"sort,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateSessionGroup creates a new session group
func (sm *SessionMemory) CreateSessionGroup(ctx context.Context, name string) (*SessionGroup, error) {
	if name == "" {
		return nil, fmt.Errorf("session group name cannot be empty")
	}
	if len(name) > 200 {
		return nil, fmt.Errorf("session group name too long (max 200 characters)")
	}

	id := uuid.New().String()
	now := time.Now()

	_, err := sm.db.ExecContext(ctx, `
		INSERT INTO session_groups (id, name, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`, id, name, now, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create session group: %w", err)
	}

	return &SessionGroup{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// GetSessionGroups returns all session groups
func (sm *SessionMemory) GetSessionGroups(ctx context.Context) ([]SessionGroup, error) {
	rows, err := sm.db.QueryContext(ctx, `
		SELECT id, name, sort, created_at, updated_at
		FROM session_groups
		ORDER BY sort ASC, created_at ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query session groups: %w", err)
	}
	defer rows.Close()

	var groups []SessionGroup
	for rows.Next() {
		var g SessionGroup
		var sort sql.NullInt64
		if err := rows.Scan(&g.ID, &g.Name, &sort, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan session group: %w", err)
		}
		if sort.Valid {
			sortInt := int(sort.Int64)
			g.Sort = &sortInt
		}
		groups = append(groups, g)
	}

	return groups, nil
}

// GetSessionGroup returns a specific session group
func (sm *SessionMemory) GetSessionGroup(ctx context.Context, id string) (*SessionGroup, error) {
	if id == "" {
		return nil, fmt.Errorf("session group ID cannot be empty")
	}

	var g SessionGroup
	var sort sql.NullInt64
	err := sm.db.QueryRowContext(ctx, `
		SELECT id, name, sort, created_at, updated_at
		FROM session_groups
		WHERE id = ?
	`, id).Scan(&g.ID, &g.Name, &sort, &g.CreatedAt, &g.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session group not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session group: %w", err)
	}

	if sort.Valid {
		sortInt := int(sort.Int64)
		g.Sort = &sortInt
	}

	return &g, nil
}

// UpdateSessionGroupParams represents parameters for updating a session group
type UpdateSessionGroupParams struct {
	Name *string `json:"name,omitempty"`
	Sort *int    `json:"sort,omitempty"`
}

// UpdateSessionGroup updates a session group (supports partial updates)
func (sm *SessionMemory) UpdateSessionGroup(ctx context.Context, id string, params UpdateSessionGroupParams) error {
	if id == "" {
		return fmt.Errorf("session group ID cannot be empty")
	}

	updates := []string{"updated_at = ?"}
	args := []interface{}{time.Now()}

	if params.Name != nil {
		if *params.Name == "" {
			return fmt.Errorf("session group name cannot be empty")
		}
		if len(*params.Name) > 200 {
			return fmt.Errorf("session group name too long (max 200 characters)")
		}
		updates = append(updates, "name = ?")
		args = append(args, *params.Name)
	}

	if params.Sort != nil {
		updates = append(updates, "sort = ?")
		args = append(args, *params.Sort)
	}

	args = append(args, id)

	query := "UPDATE session_groups SET "
	for i, update := range updates {
		if i > 0 {
			query += ", "
		}
		query += update
	}
	query += " WHERE id = ?"

	result, err := sm.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update session group: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("session group not found: %s", id)
	}

	return nil
}

// DeleteSessionGroup deletes a session group
func (sm *SessionMemory) DeleteSessionGroup(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("session group ID cannot be empty")
	}

	// First, unset group_id for all sessions in this group
	_, err := sm.db.ExecContext(ctx, `
		UPDATE sessions SET group_id = NULL WHERE group_id = ?
	`, id)
	if err != nil {
		return fmt.Errorf("failed to unset group_id for sessions: %w", err)
	}

	// Then delete the group
	_, err = sm.db.ExecContext(ctx, "DELETE FROM session_groups WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete session group: %w", err)
	}

	return nil
}

// DeleteAllSessionGroups deletes all session groups
func (sm *SessionMemory) DeleteAllSessionGroups(ctx context.Context) error {
	// First, unset group_id for all sessions
	_, err := sm.db.ExecContext(ctx, `
		UPDATE sessions SET group_id = NULL
	`)
	if err != nil {
		return fmt.Errorf("failed to unset group_id for sessions: %w", err)
	}

	// Then delete all groups
	_, err = sm.db.ExecContext(ctx, "DELETE FROM session_groups")
	if err != nil {
		return fmt.Errorf("failed to delete all session groups: %w", err)
	}

	return nil
}

// UpdateSessionGroupOrder updates the sort order of multiple session groups
func (sm *SessionMemory) UpdateSessionGroupOrder(ctx context.Context, sortMap []struct {
	ID   string `json:"id"`
	Sort int    `json:"sort"`
}) error {
	if len(sortMap) == 0 {
		return nil
	}

	tx, err := sm.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, item := range sortMap {
		if item.ID == "" {
			continue
		}
		_, err := tx.ExecContext(ctx, `
			UPDATE session_groups SET sort = ?, updated_at = ? WHERE id = ?
		`, item.Sort, time.Now(), item.ID)
		if err != nil {
			return fmt.Errorf("failed to update sort for group %s: %w", item.ID, err)
		}
	}

	return tx.Commit()
}

// CountSessions returns the total number of sessions
func (sm *SessionMemory) CountSessions(ctx context.Context) (int, error) {
	var count int
	err := sm.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sessions`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count sessions: %w", err)
	}
	return count, nil
}

// CountSessionsByDate returns the number of sessions created before or on the specified date
// endDate format: "YYYY-MM-DD"
func (sm *SessionMemory) CountSessionsByDate(ctx context.Context, endDate string) (int, error) {
	if endDate == "" {
		return sm.CountSessions(ctx)
	}

	var count int
	err := sm.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM sessions 
		WHERE DATE(created_at) <= DATE(?)
	`, endDate).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count sessions by date: %w", err)
	}
	return count, nil
}

// SessionRankItem represents a session with its message count for ranking
type SessionRankItem struct {
	ID              string  `json:"id"`
	Title           string  `json:"title"`
	Avatar          *string `json:"avatar,omitempty"`
	BackgroundColor *string `json:"backgroundColor,omitempty"`
	Count           int     `json:"count"`
}

// RankSessions returns sessions ranked by message count
func (sm *SessionMemory) RankSessions(ctx context.Context) ([]SessionRankItem, error) {
	rows, err := sm.db.QueryContext(ctx, `
		SELECT 
			s.id,
			s.title,
			s.meta,
			COUNT(m.id) as message_count
		FROM sessions s
		LEFT JOIN langchaingo_messages m ON s.id = m.session_id
		GROUP BY s.id, s.title, s.meta
		HAVING message_count > 0
		ORDER BY message_count DESC
		LIMIT 20
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query session ranks: %w", err)
	}
	defer rows.Close()

	var items []SessionRankItem
	for rows.Next() {
		var (
			id           string
			title        string
			metaJSON     string
			messageCount int
		)
		if err := rows.Scan(&id, &title, &metaJSON, &messageCount); err != nil {
			return nil, fmt.Errorf("failed to scan session rank: %w", err)
		}

		// Parse meta to extract avatar and backgroundColor
		var meta MetaData
		if metaJSON != "" {
			if err := json.Unmarshal([]byte(metaJSON), &meta); err != nil {
				// If parsing fails, continue with empty meta
				meta = MetaData{}
			}
		}

		item := SessionRankItem{
			ID:    id,
			Title: title,
			Count: messageCount,
		}

		// Convert string to *string if not empty
		if meta.Avatar != "" {
			item.Avatar = &meta.Avatar
		}
		if meta.BackgroundColor != "" {
			item.BackgroundColor = &meta.BackgroundColor
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating session ranks: %w", err)
	}

	return items, nil
}
