package ghapi

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v62/github"

	"github.com/curioswitch/tasuke/webhook/server/internal/config"
)

// NewClientCreator returns a new ClientCreator for creating authenticated [github.Client]s.
func NewClientCreator(conf *config.Config) (*ClientCreator, error) {
	key, err := base64.StdEncoding.DecodeString(conf.GitHub.PrivateKeyBase64)
	if err != nil {
		return nil, fmt.Errorf("ghapi: decode base64 key: %w", err)
	}

	return &ClientCreator{
		appID: conf.GitHub.AppID,
		key:   key,
	}, nil
}

// ClientCreator is a factory for creating authenticated [github.Client]s.
type ClientCreator struct {
	appID int64
	key   []byte
}

// NewClient creates a new authenticated [github.Client] for the given installation.
func (c *ClientCreator) NewClient(installation int64) (*github.Client, error) {
	tr, err := ghinstallation.New(http.DefaultTransport, c.appID, installation, c.key)
	if err != nil {
		return nil, fmt.Errorf("ghapi: create transport: %w", err)
	}

	return github.NewClient(&http.Client{Transport: tr}), nil
}
