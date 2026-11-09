package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ollama/ollama/api"
	"github.com/vviveksharma.in/llmgate/internal/env"
	"github.com/vviveksharma.in/llmgate/internal/ollama"
	"github.com/vviveksharma.in/llmgate/internal/server"
)

func main() {
	env.LoadFile(".env")

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "help":
			server.Help()
			return
		case "serve":
			runServer()
			return
		}
	}

	runDemoChat(os.Args[1:])
}

func runDemoChat(args []string) {
	model := strings.TrimSpace(os.Getenv("OLLAMA_MODEL"))
	if model == "" {
		log.Fatal("set OLLAMA_MODEL in .env or the environment (see .env.example)")
	}

	prompt := "Who are you?"
	if len(args) > 0 {
		prompt = strings.Join(args, " ")
	}

	client, err := ollama.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	fmt.Fprintf(os.Stderr, "model: %s\n→ %s\n\n", model, prompt)

	reply, err := client.Chat(ctx, model, []api.Message{{Role: "user", Content: prompt}})
	if err != nil {
		log.Fatalf("ollama: %v\n(start: ollama serve | pull: ollama pull %s)", err, model)
	}

	fmt.Println(reply)
}

func runServer() {
	client := server.MustOllamaClient()
	server.PingOllama(client)

	srv := server.New(client)
	if err := srv.ListenAndServe(server.AddrFromEnv()); err != nil {
		log.Fatal(err)
	}
}
