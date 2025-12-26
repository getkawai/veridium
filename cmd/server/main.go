package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kawai-network/veridium/pkg/api"
	"github.com/kawai-network/veridium/pkg/store"
)

func main() {
	// Initialize Store with multi-namespace configuration
	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		log.Fatalf("Failed to initialize KV store: %v", err)
	}

	claimHandler := api.NewClaimHandler(kv)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/claim/proof", claimHandler.GetProof)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
