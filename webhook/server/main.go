package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/curioswitch/go-curiostack/logging"
	"github.com/curioswitch/go-curiostack/server"

	"github.com/curioswitch/tasuke/webhook/server/internal/config"
	"github.com/curioswitch/tasuke/webhook/server/internal/handler"
)

func main() {
	ctx := context.Background()

	conf := config.Load()

	logging.Initialize(conf.Common)

	mux := server.NewMux()

	h := handler.New(conf)

	mux.Handle("/github-webhook", h)

	srv := server.NewServer(mux, conf.Common)

	slog.InfoContext(ctx, fmt.Sprintf("Starting server on address %v", srv.Addr))
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		slog.ErrorContext(ctx, fmt.Sprintf("Failed to start server: %v", err))
	}
}
