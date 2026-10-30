package handler

import (
	"encoding/json"
	"net/http"
)

type Health struct {
	Ready func() error
}

func (h *Health) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	body := map[string]string{"status": "ok"}
	if h.Ready != nil {
		if err := h.Ready(); err != nil {
			status = http.StatusServiceUnavailable
			body["status"] = "degraded"
			body["error"] = err.Error()
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
