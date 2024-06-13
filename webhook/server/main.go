package main

import (
	"context"
	"embed"
	"fmt"
	"os"

	"github.com/curioswitch/go-curiostack/server"

	"github.com/curioswitch/tasuke/webhook/server/internal/config"
	"github.com/curioswitch/tasuke/webhook/server/internal/handler"
)

var confFiles embed.FS // Currently empty

func main() {
	os.Exit(server.Main(&config.Config{}, confFiles, setupServer))
}

func setupServer(ctx context.Context, conf *config.Config, s *server.Server) error {
	webhook, err := handler.New(conf)
	if err != nil {
		return fmt.Errorf("main: create handler: %w", err)
	}
	server.Mux(s).Handle("/github-webhook", webhook)

	return server.Start(ctx, s)
}
