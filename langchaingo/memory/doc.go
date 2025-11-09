/*
Package memory provides an interface for managing conversational data and
a variety of implementations for storing and retrieving that data.

The main components of this package are:
  - ChatMessageHistory: a struct that stores chat messages.
  - ConversationBuffer: a simple form of memory that remembers previous conversational back and forth directly.
  - ConversationTokenBuffer: memory that prunes old messages when token limit is exceeded.
  - ConversationWindowBuffer: memory that keeps only the last N conversation pairs.
  - ConversationSummaryBuffer: memory that automatically summarizes old messages and keeps only recent pairs.
    This prevents LLM confusion from very long conversations by maintaining a summary + recent context.
  - SessionMemory: database-backed memory for persistent session storage.
*/
package memory
