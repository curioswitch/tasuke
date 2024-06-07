package ghapi

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v62/github"
	"github.com/stretchr/testify/require"

	"github.com/curioswitch/tasuke/webhook/server/internal/config"
)

func TestClientCreator(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	// Encode the private key to the PEM format
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	privateKeyBase64 := base64.StdEncoding.EncodeToString(privateKeyPEM)

	conf := &config.Config{
		GitHub: config.GitHub{
			AppID:            111,
			PrivateKeyBase64: privateKeyBase64,
		},
	}

	creator, err := NewClientCreator(conf)
	require.NoError(t, err)

	client, err := creator.NewClient(222)
	require.NoError(t, err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/app/installations/222/access_tokens":
			w.Write([]byte(`{"token":"test_token"}`))
		case "/repos/ghtest/repo/issues/3/comments":
			require.Equal(t, "token test_token", r.Header.Get("Authorization"))
		default:
			require.Failf(t, "unexpected request", "path: %s", r.URL.Path)
		}
	}))
	defer srv.Close()

	u, _ := url.Parse(srv.URL + "/")
	client.BaseURL = u

	client.Client().Transport.(*ghinstallation.Transport).BaseURL = srv.URL

	_, _, err = client.Issues.CreateComment(context.Background(), "ghtest", "repo", 3, &github.IssueComment{
		Body: github.String("test"),
	})
	require.NoError(t, err)
}
