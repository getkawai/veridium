# Changelog

All notable changes to this project will be documented in this file.

## [0.2.0] - 2025-11-09

### Added
- **GroupMessageFlatten Processor** - Complete implementation for flattening group messages into assistant + tool message sequences
- **Pipeline Management APIs**:
  - `AddProcessor()` - Add custom processors to the pipeline
  - `RemoveProcessor()` - Remove processors by name
  - `DisableProcessor()` / `EnableProcessor()` - Toggle built-in processors
  - `ClearCustomProcessors()` - Remove all custom processors
  - `GetProcessors()` - List all processor names
  - `Clone()` - Deep copy engine with same configuration
- **Pipeline Validation API**:
  - `Validate()` - Validate pipeline configuration
  - Returns `ValidationResult` with errors list
  - Checks for duplicate names, missing functions, invalid config
- **Pipeline Statistics API**:
  - `GetStats()` - Get comprehensive pipeline statistics
  - Returns processor counts, names, custom/disabled processors
- **Enhanced Message Types**:
  - Added `GroupChild`, `ToolCall`, `ToolResult` types
  - Extended `Message` with group message fields
  - Added ParentID, ThreadID, GroupID, AgentID, TargetID, TopicID fields
  - Added Reasoning field for reasoning models

### Changed
- MessageContentProcessor enhanced with better structure
- Graph orchestration now includes GroupMessageFlatten as first processor
- Engine struct now supports custom processors and disabled processors map

### Implementation Notes
- **Feature Parity**: Now matches TypeScript version's core functionality
- All 9 processors + 4 providers fully implemented
- Pipeline management matches TypeScript's ContextEngine API
- Ready for production use with complete feature set

## [0.1.0] - 2025-01-XX

### Added
- Initial implementation of Context Engine Go using Eino framework
- Phase 1: Foundation Setup
  - Go project structure
  - Eino dependencies
  - Core type definitions (Message, Context, Config)
  - Base graph builder interface
  - Testing framework setup

- Phase 2: Simple Processors Migration
  - HistoryTruncateProcessor → Lambda Node
  - MessageCleanupProcessor → Lambda Node
  - ToolMessageReorder → Lambda Node

- Phase 3: Providers Migration
  - SystemRoleInjector → Lambda Node
  - InboxGuideProvider → Lambda Node
  - ToolSystemRoleProvider → Lambda Node
  - HistorySummaryProvider → Lambda Node

- Phase 4: Complex Processors Migration
  - InputTemplateProcessor → Lambda Node
  - PlaceholderVariablesProcessor → Lambda Node
  - MessageContentProcessor → Lambda Node (basic implementation)
  - ToolCallProcessor → Lambda Node

- Phase 5: Graph Orchestration
  - Main context engineering graph builder
  - Node connection and workflow orchestration
  - MessageInput/MessageOutput wrappers for Eino compatibility

### Implementation Notes
- All processors and providers are implemented as Eino Lambda nodes
- Graph orchestration uses Eino Workflow with MessageInput/MessageOutput wrappers
- Basic functionality implemented; advanced features (streaming, parallelism) to be added in future phases

