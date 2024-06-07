package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/curioswitch/go-curiostack/logging"
	"github.com/curioswitch/go-curiostack/server"

	"github.com/curioswitch/tasuke/webhook/server/internal/config"
	"github.com/curioswitch/tasuke/webhook/server/internal/handler"
)

func main() {
	os.Exit(doMain())
}

func doMain() int {
	ctx := context.Background()

	conf := config.Load()

	logging.Initialize(conf.Common)

	mux := server.NewMux()

	h, err := handler.New(conf)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Failed to create handler: %v", err))
		return 1
	}

	mux.Handle("/github-webhook", h)

	srv := server.NewServer(mux, conf.Common)

	slog.InfoContext(ctx, fmt.Sprintf("Starting server on address %v", srv.Addr))
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		slog.ErrorContext(ctx, fmt.Sprintf("Failed to start server: %v", err))
		return 1
	}

	return 0
}
