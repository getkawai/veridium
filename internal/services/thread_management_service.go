/*
 * Copyright 2025 Veridium Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kawai-network/veridium/pkg/xlog"

	"github.com/google/uuid"
	"github.com/kawai-network/veridium/internal/database"
	db "github.com/kawai-network/veridium/internal/database/generated"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// ThreadManagementService handles thread creation and management for conversation branching
type ThreadManagementService struct {
	app *application.App
	db  *database.Service
}

// ThreadType represents the type of thread
type ThreadType string

const (
	ThreadTypeContinuation ThreadType = "continuation" // Continue from a message
	ThreadTypeStandalone   ThreadType = "standalone"   // Standalone thread
)

// ThreadStatus represents the status of a thread
type ThreadStatus string

const (
	ThreadStatusActive     ThreadStatus = "active"
	ThreadStatusDeprecated ThreadStatus = "deprecated"
	ThreadStatusArchived   ThreadStatus = "archived"
)

// CreateThreadRequest represents a request to create a new thread
type CreateThreadRequest struct {
	SourceMessageID string     `json:"source_message_id"` // Message to branch from
	TopicID         string     `json:"topic_id"`          // Topic this thread belongs to
	Title           string     `json:"title,omitempty"`
	Type            ThreadType `json:"type"` // continuation or standalone
	ParentThreadID  string     `json:"parent_thread_id,omitempty"`
}

// CreateThreadResponse represents the response from creating a thread
type CreateThreadResponse struct {
	ThreadID  string `json:"thread_id"`
	Title     string `json:"title"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at"`
	Error     string `json:"error,omitempty"`
}

// ThreadInfo represents thread information
type ThreadInfo struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	Type            string `json:"type"`
	Status          string `json:"status"`
	TopicID         string `json:"topic_id"`
	SourceMessageID string `json:"source_message_id"`
	ParentThreadID  string `json:"parent_thread_id,omitempty"`
	LastActiveAt    int64  `json:"last_active_at"`
	CreatedAt       int64  `json:"created_at"`
	UpdatedAt       int64  `json:"updated_at"`
}

// ListThreadsRequest represents a request to list threads
type ListThreadsRequest struct {
	TopicID string `json:"topic_id"`
}

// NewThreadManagementService creates a new thread management service
func NewThreadManagementService(app *application.App, db *database.Service) *ThreadManagementService {
	return &ThreadManagementService{
		app: app,
		db:  db,
	}
}

// CreateThread creates a new thread for conversation branching
func (s *ThreadManagementService) CreateThread(ctx context.Context, req CreateThreadRequest) (*CreateThreadResponse, error) {
	if req.SourceMessageID == "" {
		return &CreateThreadResponse{
			Error: "source_message_id is required",
		}, fmt.Errorf("source_message_id is required")
	}

	if req.TopicID == "" {
		return &CreateThreadResponse{
			Error: "topic_id is required",
		}, fmt.Errorf("topic_id is required")
	}

	// Validate thread type
	if req.Type != ThreadTypeContinuation && req.Type != ThreadTypeStandalone {
		req.Type = ThreadTypeContinuation // Default
	}

	// Generate thread ID
	threadID := uuid.New().String()
	now := time.Now().UnixMilli()

	// Generate title if not provided
	title := req.Title
	if title == "" {
		title = fmt.Sprintf("Thread from message %s", req.SourceMessageID[:8])
	}

	// Prepare parent thread ID
	var parentThreadID sql.NullString
	if req.ParentThreadID != "" {
		parentThreadID = sql.NullString{String: req.ParentThreadID, Valid: true}
	}

	// Create thread in database
	thread, err := s.db.Queries().CreateThread(ctx, db.CreateThreadParams{
		ID:              threadID,
		Title:           title,
		Type:            string(req.Type),
		Status:          sql.NullString{String: string(ThreadStatusActive), Valid: true},
		TopicID:         req.TopicID,
		SourceMessageID: req.SourceMessageID,
		ParentThreadID:  parentThreadID,
		LastActiveAt:    now,
		CreatedAt:       now,
		UpdatedAt:       now,
	})

	if err != nil {
		xlog.Warn("❌ Failed to create thread", "error", err)
		return &CreateThreadResponse{
			Error: fmt.Sprintf("Failed to create thread: %v", err),
		}, err
	}

	xlog.Info("🔀 Created thread", "id", threadID, "type", req.Type, "from_message", req.SourceMessageID)

	return &CreateThreadResponse{
		ThreadID:  thread.ID,
		Title:     thread.Title,
		Type:      thread.Type,
		Status:    thread.Status.String,
		CreatedAt: thread.CreatedAt,
	}, nil
}

// ListThreadsByTopic lists all threads for a given topic
func (s *ThreadManagementService) ListThreadsByTopic(ctx context.Context, req ListThreadsRequest) ([]ThreadInfo, error) {
	if req.TopicID == "" {
		return nil, fmt.Errorf("topic_id is required")
	}

	if req.TopicID == "" {
		return nil, fmt.Errorf("topic_id is required")
	}

	threads, err := s.db.Queries().ListThreadsByTopic(ctx, req.TopicID)

	if err != nil {
		return nil, fmt.Errorf("failed to list threads: %w", err)
	}

	result := make([]ThreadInfo, len(threads))
	for i, thread := range threads {
		result[i] = ThreadInfo{
			ID:              thread.ID,
			Title:           thread.Title,
			Type:            thread.Type,
			Status:          thread.Status.String,
			TopicID:         thread.TopicID,
			SourceMessageID: thread.SourceMessageID,
			ParentThreadID:  thread.ParentThreadID.String,
			LastActiveAt:    thread.LastActiveAt,
			CreatedAt:       thread.CreatedAt,
			UpdatedAt:       thread.UpdatedAt,
		}
	}

	return result, nil
}

// GetThread retrieves a single thread by ID
func (s *ThreadManagementService) GetThread(ctx context.Context, threadID string) (*ThreadInfo, error) {
	if threadID == "" {
		return nil, fmt.Errorf("thread_id is required")
	}

	thread, err := s.db.Queries().GetThread(ctx, threadID)

	if err != nil {
		return nil, fmt.Errorf("failed to get thread: %w", err)
	}

	return &ThreadInfo{
		ID:              thread.ID,
		Title:           thread.Title,
		Type:            thread.Type,
		Status:          thread.Status.String,
		TopicID:         thread.TopicID,
		SourceMessageID: thread.SourceMessageID,
		ParentThreadID:  thread.ParentThreadID.String,
		LastActiveAt:    thread.LastActiveAt,
		CreatedAt:       thread.CreatedAt,
		UpdatedAt:       thread.UpdatedAt,
	}, nil
}

// UpdateThreadStatus updates the status of a thread
func (s *ThreadManagementService) UpdateThreadStatus(ctx context.Context, threadID string, status ThreadStatus) error {
	if threadID == "" {
		return fmt.Errorf("thread_id is required")
	}

	// Get current thread
	thread, err := s.db.Queries().GetThread(ctx, threadID)

	if err != nil {
		return fmt.Errorf("failed to get thread: %w", err)
	}

	now := time.Now().UnixMilli()

	// Update thread
	_, err = s.db.Queries().UpdateThread(ctx, db.UpdateThreadParams{
		Title:        thread.Title,
		Status:       sql.NullString{String: string(status), Valid: true},
		LastActiveAt: now,
		UpdatedAt:    now,
		ID:           threadID,
	})

	if err != nil {
		return fmt.Errorf("failed to update thread status: %w", err)
	}

	xlog.Info("🔄 Updated thread status", "id", threadID, "status", status)

	return nil
}

// UpdateThreadLastActive updates the last active time of a thread
func (s *ThreadManagementService) UpdateThreadLastActive(ctx context.Context, threadID string) error {
	if threadID == "" {
		return fmt.Errorf("thread_id is required")
	}

	// Get current thread
	thread, err := s.db.Queries().GetThread(ctx, threadID)

	if err != nil {
		return fmt.Errorf("failed to get thread: %w", err)
	}

	now := time.Now().UnixMilli()

	// Update thread
	_, err = s.db.Queries().UpdateThread(ctx, db.UpdateThreadParams{
		Title:        thread.Title,
		Status:       thread.Status,
		LastActiveAt: now,
		UpdatedAt:    now,
		ID:           threadID,
	})

	if err != nil {
		return fmt.Errorf("failed to update thread last active: %w", err)
	}

	return nil
}

// DeleteThread deletes a thread
func (s *ThreadManagementService) DeleteThread(ctx context.Context, threadID string) error {
	if threadID == "" {
		return fmt.Errorf("thread_id is required")
	}

	err := s.db.Queries().DeleteThread(ctx, threadID)
	if err != nil {
		return fmt.Errorf("failed to delete thread: %w", err)
	}

	xlog.Info("🗑️  Deleted thread", "id", threadID)

	return nil
}

// GetThreadMessages retrieves all messages in a thread
// This returns messages that have the thread_id set
func (s *ThreadManagementService) GetThreadMessages(ctx context.Context, threadID string) ([]db.Message, error) {
	if threadID == "" {
		return nil, fmt.Errorf("thread_id is required")
	}

	messages, err := s.db.Queries().ListMessagesByThread(ctx, sql.NullString{String: threadID, Valid: true})

	if err != nil {
		return nil, fmt.Errorf("failed to get thread messages: %w", err)
	}

	return messages, nil
}
