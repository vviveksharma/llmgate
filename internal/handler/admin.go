package handler

import (
	"encoding/json"
	"net/http"

	"github.com/vviveksharma/llmgate/internal/auth"
	"github.com/vviveksharma/llmgate/internal/tenant"
)

type Admin struct {
	Tenants *tenant.Service
	JWT     *auth.JWTValidator
}

func (h *Admin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	list, err := h.Tenants.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(list)
}
