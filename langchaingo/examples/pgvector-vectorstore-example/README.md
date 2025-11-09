# Chromem Vector Store with OpenAI Embeddings Example

This example demonstrates how to use Chromem, a local embedded vector database, with OpenAI embeddings in a Go application. It showcases the integration of langchain-go, OpenAI's API, and Chromem to create a powerful vector database for similarity searches with persistent storage.

## What This Example Does

1. **Sets up a Chromem Vector Store:**
   - Creates a local persistent vector database using Chromem
   - No external database server required

2. **Initializes OpenAI Embeddings:**
   - Creates an embeddings client using the OpenAI API
   - Requires an OpenAI API key to be set as an environment variable

3. **Creates a Chromem Store:**
   - Creates a persistent local vector store
   - Uses OpenAI embeddings for vector generation

4. **Adds Sample Documents:**
   - Inserts several documents (cities) with metadata into the vector store
   - Each document includes the city name, population, and area

5. **Performs Similarity Searches:**
   - Demonstrates various types of similarity searches:
     a. Basic search for documents similar to "japan"
     b. Search for South American cities with a score threshold
     c. Additional search with score threshold (note: metadata filters not yet supported in chromem)

## How to Run the Example

1. Set your OpenAI API key:
   ```
   export OPENAI_API_KEY=<your key>
   ```

2. Run the Go example:
   ```
   go run pgvector_vectorstore_example.go
   ```

The example will automatically:
- Create a temporary directory for the Chromem database
- Store vectors with persistent storage
- Clean up temporary files after completion

## Key Features

- **Local Storage**: No external database server required
- **Persistent Storage**: Data is stored locally and persists between runs
- **OpenAI Embeddings**: Uses OpenAI's embedding models for vector generation
- **Similarity Search**: Supports similarity search with score thresholds
- **Easy Setup**: No Docker or database configuration needed

This example provides a practical demonstration of using local vector databases for semantic search and similarity matching, which can be incredibly useful for various AI and machine learning applications.
