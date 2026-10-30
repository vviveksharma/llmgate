package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/vviveksharma/llmgate/internal/middleware"
	"github.com/vviveksharma/llmgate/internal/provider"
	"github.com/vviveksharma/llmgate/internal/router"
	"github.com/vviveksharma/llmgate/internal/stream"
	"github.com/vviveksharma/llmgate/internal/usage"
)

type Chat struct {
	Router   *router.ModelRouter
	Recorder *usage.Recorder
}

func (h *Chat) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req provider.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.Stream {
		stream.SSEHeaders(w)
		ch := make(chan string, 8)
		go func() {
			p, _ := h.Router.Chat(r.Context(), req)
			if p != nil {
				ch <- p.Content
			}
			close(ch)
		}()
		for chunk := range ch {
			_ = stream.WriteSSE(w, "message", chunk)
		}
		return
	}
	resp, err := h.Router.Chat(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	if h.Recorder != nil {
		_ = h.Recorder.Record(r.Context(), usage.Record{
			TenantID: middleware.TenantIDFromContext(r.Context()),
			Model:    resp.Model,
			At:       time.Now(),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
