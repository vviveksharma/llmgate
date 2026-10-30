package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/vviveksharma/llmgate/internal/config"
	"github.com/vviveksharma/llmgate/internal/observability"
	"github.com/vviveksharma/llmgate/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	log := observability.NewLogger()

	srv, err := server.New(cfg, log)
	if err != nil {
		log.Error("failed to create server", "err", err)
		os.Exit(1)
	}

	go func() {
		if err := srv.Start(); err != nil {
			log.Error("server stopped", "err", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
