package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ollama/ollama/api"
	"github.com/vviveksharma.in/llmgate/internal/ollama"
)

// Server is the HTTP front door for llmgate → Ollama.
type Server struct {
	ollama      *ollama.Client
	defaultModel string
}

func New(ollamaClient *ollama.Client) *Server {
	return &Server{
		ollama:       ollamaClient,
		defaultModel: strings.TrimSpace(os.Getenv("OLLAMA_MODEL")),
	}
}

func (s *Server) ListenAndServe(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /v1/models", s.handleModels)
	mux.HandleFunc("POST /v1/chat", s.handleChat)

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}
	log.Printf("llmgate listening on %s (OLLAMA_HOST controls upstream)", addr)
	return srv.ListenAndServe()
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleModels(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	names, err := s.ollama.ModelNames(ctx)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"models": names})
}

type chatRequest struct {
	Model    string         `json:"model"`
	Messages []api.Message  `json:"messages"`
}

type chatResponse struct {
	Model   string `json:"model"`
	Content string `json:"content"`
}

func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}

	model := strings.TrimSpace(req.Model)
	if model == "" {
		model = s.defaultModel
	}
	if model == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "set model in the request body or OLLAMA_MODEL",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
	defer cancel()

	content, err := s.ollama.Chat(ctx, model, req.Messages)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, chatResponse{Model: model, Content: content})
}

func writeError(w http.ResponseWriter, err error) {
	var statusErr api.StatusError
	if errors.As(err, &statusErr) {
		code := statusErr.StatusCode
		if code == 0 {
			code = http.StatusBadGateway
		}
		writeJSON(w, code, map[string]string{"error": statusErr.Error()})
		return
	}
	log.Printf("request error: %v", err)
	writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("encode response: %v", err)
	}
}

// AddrFromEnv returns the bind address for llmgate (default :8080).
func AddrFromEnv() string {
	if v := strings.TrimSpace(os.Getenv("LLMGATE_ADDR")); v != "" {
		return v
	}
	return ":8080"
}

func MustOllamaClient() *ollama.Client {
	c, err := ollama.NewClient()
	if err != nil {
		log.Fatalf("ollama: %v", err)
	}
	return c
}

// PingOllama logs a warning if the Ollama daemon is not reachable.
func PingOllama(c *ollama.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := c.ModelNames(ctx); err != nil {
		log.Printf("warning: ollama not reachable yet (%v). Start it with: ollama serve", err)
		log.Printf("tip: pull a model with: ollama pull %s", strings.TrimSpace(os.Getenv("OLLAMA_MODEL")))
		return
	}
	log.Print("connected to ollama")
}

// Help prints usage to stderr.
func Help() {
	fmt.Fprintln(os.Stderr, "llmgate — local Ollama client and optional HTTP gateway")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  go run .              chat with Ollama (default prompt: Who are you?)")
	fmt.Fprintln(os.Stderr, `  go run . "your prompt"  send a custom prompt (quote it in zsh if it has ? * [ etc.)`)
	fmt.Fprintln(os.Stderr, "  go run . serve        start the HTTP API")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Environment:")
	fmt.Fprintln(os.Stderr, "  OLLAMA_HOST    Ollama API URL (default http://127.0.0.1:11434)")
	fmt.Fprintln(os.Stderr, "  OLLAMA_MODEL   Default model name for POST /v1/chat")
	fmt.Fprintln(os.Stderr, "  LLMGATE_ADDR   Listen address (default :8080)")
}
