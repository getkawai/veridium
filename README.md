# Veridium Startup Process

When `main.go` executes, the application follows a structured initialization sequence to set up the backend services, the AI engine, and the Wails frontend runtime.

## 1. Context & Core Service Initialization
- **Dev Mode Check**: Checks for development mode flags.
- **Context Creation**: `app.NewContext()` creates a centralized context aimed at being the "Single Source of Truth" for the application state.
- **Service Init**: `ctx.InitAll()` initializes all core backend services, including:
  - **SQLite**: The primary relational database used for storing application metadata, chat history, agent configurations, and session information.
  - **DuckDB**: An embedded analytical database utilized as the high-performance vector store for semantic retrieval and long-term memory.
  - **Vector Search Engine**: Provides semantic similarity lookup capabilities using high-dimensional embeddings (defaulting to 384 dimensions) to power RAG and memory features.
  - **Whisper Service**: A local implementation of the Whisper model for high-accuracy, privacy-preserving speech-to-text transcription.
  - **Knowledge Base (KB) Service**: Specifically manages the ingestion, chunking, and indexing of knowledge base documents (PDFs, text, etc.), coordinating with the RAG processor.
  - **Local LLM Library Service**: A centralized manager for running Large Language Models locally using a highly optimized `llama.cpp` wrapper, handling model lifecycle and inference.

## 2. Specialized Service Setup
- **File Processor**: A `FileProcessorService` is instantiated to bridge Wails with the file processing logic (PDFs, images, etc.), utilizing the database and vector search.
- **Stable Diffusion**: The image generation manager (`image.New`) is initialized and starts its background initialization process immediately.

## 3. Wails Application Configuration
- A Wails application (`application.App`) is created with specific configurations for Mac (lifecycle management) and assets.
- **Service Injection**: A comprehensive list of services is bound to the Wails application, making them callable from the frontend. This includes:
  - **Database Services**: Queries, DB access, Table Viewer.
  - **Core Features**: Search, TTS, Audio Recorder.
  - **AI/ML**: Vector Search, File Processor, KB Service.
  - **System Utilities**: File Service, Local File System, System Service, Machine ID.
  - **Native Wails Services**: Notifications, Logger, clean SQLite/KVStore access.

## 4. Agent & Logic wiring (`registerAgentServices`)
After the base app is created, complex dependencies are wired up:
- **Audio Recorder**: Connected to the Wails app instance.
- **Thread & Topic Management**: Services for handling chat threads and topics are initialized and registered.
- **Cascading Model Chain**: The AI utilizes a prioritized chain of providers for robustness: `OpenRouter (Free)` → `Pollinations AI` → `ZAI (GLM-4.6)` → `Local Llama`.
  - **Chat Model**: High-reasoning configuration with 100k context window and tool-calling/attachment support.
  - **Summary Model**: Specialized for long-context analysis with a 50k context window.
  - **Title Model**: Dedicated to semantic summarization for session naming.
- **MemGPT-style Memory Integration**: A persistent memory system that uses LLM enrichment to extract and store user facts and preferences in the vector store, enabling cross-session recall via specialized memory tools.
- **Thread & Topic Management**:
  - **Thread Service**: Handles the complete lifecycle and database persistence of chat sessions and messages.
  - **Topic Service**: Automates the generation and organization of session titles based on the Title Model's analysis of the interaction.
- **Cleanup Handlers**: Shutdown hooks are registered to properly close the LLM library and Stable Diffusion processes.

## 5. UI Initialization
- **Main Window**: The main application window is created with:
  - Custom macOS styling (hidden title bar, translucent backdrop).
  - Specific dimensions and background colors.
- **Drag & Drop**: A native drag-and-drop handler is registered to:
  - Intercept file drops.
  - Process files immediately via `FileProcessorService`.
  - Emit real-time events (`files:dropped`) to the frontend with file details.

## 6. Execution
- Finally, `wailsApp.Run()` is called to start the application event loop, serving the frontend and listening for backend calls.
