package config

import (
	"embed"
	"log"

	"github.com/curioswitch/go-curiostack/config"
)

//go:embed *.yaml
var confFiles embed.FS

// Google is configuration related to GCP.
type Google struct {
	// Project is the GCP project to target.
	Project string `koanf:"project"`
}

// Config is the configuration for the frontend server.
type Config struct {
	config.Common

	// Google is the GCP configuration.
	Google Google `koanf:"google"`
}

// Load loads the configuration for the frontend server.
func Load() *Config {
	cfg := &Config{}
	if err := config.Load(cfg, confFiles); err != nil {
		// Should never happen in a proper setup, so just panic.
		log.Fatalf("Failed to load config: %v", err)
	}
	return cfg
}
