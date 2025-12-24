package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kawai-network/veridium/pkg/api"
	"github.com/kawai-network/veridium/pkg/store"
)

func main() {
	// Shared Config (In prod, use env vars)
	accountID := "ceab218751d33cd804878196ad7bef74"
	apiToken := "OP8BZQhyeJxrovCPKt15eUOSC6i5LXTVECGRSMc1"
	namespaceID := "55bd9f26233940dabb65f2a1992bfae9"

	// Initialize Store
	kv, err := store.NewKVStore(apiToken, accountID, namespaceID)
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
