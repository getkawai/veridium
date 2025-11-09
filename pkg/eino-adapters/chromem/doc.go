// Package chromem provides Eino-compatible adapters for chromem vector database.
//
// This package wraps chromem.Collection to implement Eino's Indexer and Retriever
// interfaces, allowing chromem to be used seamlessly in Eino workflows and graphs.
//
// # Features
//
//   - Local embeddings (no external API calls)
//   - Zero-dependency vector database
//   - In-memory with optional persistence
//   - Compatible with Eino graph orchestration
//
// # Basic Usage
//
//	// Create chromem DB and collection
//	db := chromem.NewDB()
//	collection, err := db.CreateCollection("my_docs", nil, chromem.NewEmbeddingFuncDefault())
//
//	// Wrap with Eino adapters
//	indexer := chromem.NewIndexer(collection)
//	retriever, err := chromem.NewRetriever(&chromem.RetrieverConfig{
//	    Collection: collection,
//	    TopK:       5,
//	})
//
//	// Use in Eino graph
//	graph := compose.NewGraph[Input, Output]()
//	graph.AddRetrieverNode("retriever", retriever)
//
// # Indexer
//
// The Indexer wraps chromem.Collection.AddDocuments() and implements the
// indexer.Indexer interface from Eino. It converts Eino documents to chromem
// documents and stores them in the collection.
//
//	docs := []*schema.Document{
//	    {Content: "Hello world", MetaData: map[string]any{"source": "test"}},
//	}
//	ids, err := indexer.Store(ctx, docs)
//
// # Retriever
//
// The Retriever wraps chromem.Collection.Query() and implements the
// retriever.Retriever interface from Eino. It performs similarity search
// and returns results as Eino documents.
//
//	docs, err := retriever.Retrieve(ctx, "hello",
//	    retriever.WithTopK(10),
//	    retriever.WithScoreThreshold(0.7),
//	)
//
// # Integration with Existing Code
//
// This adapter allows you to use your existing chromem collections in Eino
// workflows without any migration:
//
//	// Your existing chromem collection
//	collection := getExistingChromemCollection()
//
//	// Wrap it for Eino
//	einoIndexer := chromem.NewIndexer(collection)
//	einoRetriever, _ := chromem.NewRetriever(&chromem.RetrieverConfig{
//	    Collection: collection,
//	})
//
//	// Now use in Eino graph
//	graph.AddRetrieverNode("my_retriever", einoRetriever)
//
// # Advanced Features
//
// Access chromem-specific functionality through GetCollection():
//
//	// Get underlying collection for chromem-specific operations
//	collection := retriever.GetCollection()
//	results, err := collection.Query(ctx, chromem.QueryOptions{
//	    QueryText: "query",
//	    Negative: chromem.NegativeQueryOptions{
//	        Text: "exclude this",
//	        Mode: chromem.NEGATIVE_MODE_FILTER,
//	    },
//	})
package chromem

