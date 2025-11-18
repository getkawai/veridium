#!/bin/bash

# Test script for AgentChatService API
# Phase 4: Testing Topic & Thread integration

echo "🧪 Testing AgentChatService API - Phase 4"
echo "=========================================="
echo ""

# Test 1: First message (should auto-create session & topic)
echo "📝 Test 1: First message (auto-create session & topic)"
echo "-------------------------------------------------------"
SESSION_ID="test-session-$(date +%s)"
USER_ID="test-user-001"

curl -X POST http://localhost:8080/api/agent/chat \
  -H "Content-Type: application/json" \
  -d "{
    \"session_id\": \"$SESSION_ID\",
    \"user_id\": \"$USER_ID\",
    \"message\": \"Hello! What is CloudWeGo Eino?\",
    \"temperature\": 0.7,
    \"max_tokens\": 200
  }" | jq '.'

echo ""
echo ""

# Wait a bit for topic generation
sleep 2

# Test 2: Follow-up message (same session, should use existing topic)
echo "📝 Test 2: Follow-up message (same session)"
echo "-------------------------------------------"

curl -X POST http://localhost:8080/api/agent/chat \
  -H "Content-Type: application/json" \
  -d "{
    \"session_id\": \"$SESSION_ID\",
    \"user_id\": \"$USER_ID\",
    \"message\": \"Can you explain more about its features?\",
    \"temperature\": 0.7,
    \"max_tokens\": 200
  }" | jq '.'

echo ""
echo ""

# Test 3: New session (should create new topic)
echo "📝 Test 3: New session (new topic)"
echo "-----------------------------------"
NEW_SESSION_ID="test-session-$(date +%s)"

curl -X POST http://localhost:8080/api/agent/chat \
  -H "Content-Type: application/json" \
  -d "{
    \"session_id\": \"$NEW_SESSION_ID\",
    \"user_id\": \"$USER_ID\",
    \"message\": \"Tell me about artificial intelligence\",
    \"temperature\": 0.7,
    \"max_tokens\": 200
  }" | jq '.'

echo ""
echo ""

# Test 4: Message with explicit TopicID
echo "📝 Test 4: Message with explicit TopicID"
echo "----------------------------------------"

curl -X POST http://localhost:8080/api/agent/chat \
  -H "Content-Type: application/json" \
  -d "{
    \"session_id\": \"$SESSION_ID\",
    \"user_id\": \"$USER_ID\",
    \"topic_id\": \"explicit-topic-001\",
    \"message\": \"Continue from previous topic\",
    \"temperature\": 0.7,
    \"max_tokens\": 200
  }" | jq '.'

echo ""
echo ""

echo "✅ All tests completed!"
echo ""
echo "📊 Response should include:"
echo "   - message_id: Created message ID"
echo "   - session_id: Session ID"
echo "   - topic_id: Auto-created or provided topic ID"
echo "   - thread_id: Thread ID (if in thread)"
echo "   - message: Assistant response"
echo "   - created_at: Timestamp"

