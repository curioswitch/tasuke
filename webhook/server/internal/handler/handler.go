package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/go-github/v62/github"

	"github.com/curioswitch/tasuke/webhook/server/internal/config"
	"github.com/curioswitch/tasuke/webhook/server/internal/ghapi"
)

type githubEventType byte

const (
	githubEventTypeIssueComment githubEventType = iota
	githubEventTypePullRequest
)

func New(config *config.Config) (*Handler, error) {
	clientCreator, err := ghapi.NewClientCreator(config)
	if err != nil {
		return nil, fmt.Errorf("handler: create client creator: %w", err)
	}
	return &Handler{
		secret:        []byte(config.GitHub.Secret),
		clientCreator: clientCreator,
	}, nil
}

type Handler struct {
	secret        []byte
	clientCreator *ghapi.ClientCreator
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	eventTypeStr := r.Header.Get("X-GitHub-Event")
	if eventTypeStr == "" {
		http.Error(w, "Missing X-GitHub-Event header", http.StatusBadRequest)
		return
	}

	var eventType githubEventType
	switch eventTypeStr {
	case "issue_comment":
		eventType = githubEventTypeIssueComment
	case "pull_request":
		eventType = githubEventTypePullRequest
	default:
		// Return success
		return
	}

	slog.InfoContext(ctx, "Received event "+eventTypeStr)

	payload, err := github.ValidatePayload(r, h.secret)
	if err != nil {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}
	slog.InfoContext(ctx, "Received payload"+string(payload))

	switch eventType {
	case githubEventTypeIssueComment:
		err = h.handleIssueComment(ctx, payload)
	case githubEventTypePullRequest:
		// TODO
	}

	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Failed to handle event: %v", err))
		http.Error(w, "Failed to handle event", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) handleIssueComment(ctx context.Context, payload []byte) error {
	var event github.IssueCommentEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("handler: unmarshal payload: %w", err)
	}

	body := strings.TrimSpace(event.Comment.GetBody())

	if !strings.HasPrefix(body, "/tasuke") {
		return nil
	}

	client, err := h.clientCreator.NewClient(event.Installation.GetID())
	if err != nil {
		return fmt.Errorf("handler: create client: %w", err)
	}

	if _, _, err := client.Issues.CreateComment(ctx, event.Repo.GetOwner().GetLogin(), event.Repo.GetName(), event.Issue.GetNumber(), &github.IssueComment{
		Body: github.String("Hi there! I'm currently under development but will be happy to help after I'm working."),
	}); err != nil {
		return fmt.Errorf("handler: create comment: %w", err)
	}

	return nil
}
