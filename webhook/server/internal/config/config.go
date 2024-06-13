package config

import (
	"github.com/curioswitch/go-curiostack/config"
)

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
