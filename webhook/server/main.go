package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/curioswitch/go-curiostack/server"
	"github.com/go-chi/chi/v5"

	"github.com/curioswitch/tasuke/webhook/server/internal/config"
	"github.com/curioswitch/tasuke/webhook/server/internal/handler"
)

func main() {
	os.Exit(doMain())
}

func doMain() int {
	ctx := context.Background()

	conf := config.Load()

	webhook, err := handler.New(conf)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Failed to create handler: %v", err))
		return 1
	}

	return server.Run(ctx, &conf.Common,
		server.SetupMux(func(mux *chi.Mux) error {
			mux.Handle("/github-webhook", webhook)
			return nil
		}),
	)
}
