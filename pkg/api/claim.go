package api

import (
	"encoding/json"
	"net/http"

	"github.com/kawai-network/veridium/pkg/store"
)

type ClaimHandler struct {
	Store store.Store
}

func NewClaimHandler(s store.Store) *ClaimHandler {
	return &ClaimHandler{Store: s}
}

func (h *ClaimHandler) GetProof(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address required", http.StatusBadRequest)
		return
	}

	proof, err := h.Store.GetMerkleProof(r.Context(), address)
	if err != nil {
		// Log error in real app
		http.Error(w, "Proof not found or error", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(proof); err != nil {
		http.Error(w, "Failed to encode proof", http.StatusInternalServerError)
		return
	}
}
