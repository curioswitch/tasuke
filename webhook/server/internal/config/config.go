package config

import (
	"embed"
	"log"

	"github.com/curioswitch/go-curiostack/config"
)

var confFiles embed.FS // Currently empty

// Github is the configuration for the Github webhook.
type GitHub struct {
	// Secret is the secret used to validate the webhook payload.
	Secret string `koanf:"secret"`
}

// Config is the configuration for the webhook server.
type Config struct {
	config.Common

	// GitHub is the configuration for the Github webhook.
	GitHub GitHub `koanf:"github"`
}

// Load loads the configuration for the webhook server.
func Load() *Config {
	cfg := &Config{}
	if err := config.Load(cfg, confFiles); err != nil {
		// Should never happen in a proper setup, so just panic.
		log.Fatalf("Failed to load config: %v", err)
	}
	return cfg
}
