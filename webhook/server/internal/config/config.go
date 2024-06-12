package config

import (
	"embed"
	"log"

	"github.com/curioswitch/go-curiostack/config"
)

var confFiles embed.FS // Currently empty

// Github is the configuration for the Github webhook.
type GitHub struct {
	// AppID is the ID of the Github app.
	AppID int64 `koanf:"appid"`

	// Secret is the secret used to validate the webhook payload.
	Secret string `koanf:"secret"`

	// PrivateKeyBase64 is the private key for app authentication,
	// encoded as a Base64 string.
	PrivateKeyBase64 string `koanf:"privatekeybase64"`
}

// Config is the configuration for the webhook server.
type Config struct {
	config.Common

	// GitHub is the configuration for the Github webhook.
	GitHub *GitHub `koanf:"github"`
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
