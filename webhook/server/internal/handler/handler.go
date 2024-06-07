package handler

import (
	"log/slog"
	"net/http"

	"github.com/google/go-github/github"

	"github.com/curioswitch/tasuke/webhook/server/internal/config"
)

func New(config *config.Config) *Handler {
	return &Handler{
		secret: []byte(config.GitHub.Secret),
	}
}

type Handler struct {
	secret []byte
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, h.secret)
	if err != nil {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}
	slog.InfoContext(r.Context(), "Received payload"+string(payload))
}
