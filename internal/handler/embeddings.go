package handler

import (
	"encoding/json"
	"net/http"

	"github.com/vviveksharma/llmgate/internal/provider"
	"github.com/vviveksharma/llmgate/internal/router"
)

type Embeddings struct {
	Router *router.ModelRouter
}

func (h *Embeddings) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req provider.EmbeddingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	resp, err := h.Router.Embeddings(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
