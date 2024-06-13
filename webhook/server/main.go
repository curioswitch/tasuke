package main

import (
	"context"
	"embed"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"github.com/curioswitch/go-curiostack/server"

	"github.com/curioswitch/tasuke/webhook/server/internal/config"
	"github.com/curioswitch/tasuke/webhook/server/internal/handler"
)

var confFiles embed.FS // Currently empty

func main() {
	os.Exit(server.Main(&config.Config{}, confFiles, setupServer))
}

func setupServer(ctx context.Context, conf *config.Config, s *server.Server) error {
	fbApp, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: conf.Google.Project})
	if err != nil {
		return fmt.Errorf("main: create firebase app: %w", err)
	}

	firestore, err := fbApp.Firestore(ctx)
	if err != nil {
		return fmt.Errorf("main: create firestore client: %w", err)
	}
	defer firestore.Close()

	webhook, err := handler.New(conf, firestore)
	if err != nil {
		return fmt.Errorf("main: create handler: %w", err)
	}
	server.Mux(s).Handle("/github-webhook", webhook)

	return server.Start(ctx, s)
}
