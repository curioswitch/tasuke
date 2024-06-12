package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-enry/go-enry/v2"
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

	if !event.Issue.IsPullRequest() {
		return nil
	}

	body := strings.TrimSpace(event.Comment.GetBody())

	if !strings.HasPrefix(body, "/tasuke") {
		return nil
	}

	client, err := h.clientCreator.NewClient(event.Installation.GetID())
	if err != nil {
		return fmt.Errorf("handler: create client: %w", err)
	}

	owner := event.Repo.GetOwner().GetLogin()
	repo := event.Repo.GetName()
	num := event.Issue.GetNumber()

	diff, _, err := client.PullRequests.GetRaw(ctx, owner, repo, num, github.RawOptions{
		Type: github.Diff,
	})
	if err != nil {
		return fmt.Errorf("handler: get pull request diff: %w", err)
	}

	langIDs := diffLanguages(diff)

	langs := make([]string, len(langIDs))
	for i, id := range langIDs {
		lang, _ := enry.GetLanguageInfoByID(id)
		langs[i] = lang.Name
	}

	if _, _, err := client.Issues.CreateComment(ctx, owner, repo, num, &github.IssueComment{
		Body: github.String(fmt.Sprintf(`Hi there! I'm currently under development but will be happy to help after I'm working.

Currently, I only detect the languages in the PR. The languages I detected for this PR are: %s.
`, strings.Join(langs, ", "))),
	}); err != nil {
		return fmt.Errorf("handler: create comment: %w", err)
	}

	return nil
}

func diffLanguages(diff string) []int {
	var langIDs []int

	skip := false

	var filename string
	var content bytes.Buffer
	for len(diff) > 0 {
		line, rest, _ := strings.Cut(diff, "\n")
		diff = rest

		switch {
		case len(line) == 0:
			// Note that this likely doesn't happen in practice, but prevent crashes on bad input.
		case strings.HasPrefix(line, "diff --git "):
			// Note that if the patch content contained diff --git, it would be following
			// a +/- or space, so this is surprisingly robust.

			langIDs = maybeAppendFileLanguageID(langIDs, filename, content.Bytes())
			skip = false
			filename = ""
			content.Reset()
		case skip:
		case len(filename) == 0:
			if add, ok := strings.CutPrefix(line, "+++ "); ok {
				if name, ok := strings.CutPrefix(add, "b/"); ok {
					filename = name
				} else {
					// File removal, we don't need to care about it for finding code reviewers.
					skip = true
				}
			}
		case line[0] == '+' || line[0] == ' ':
			content.WriteString(line[1:])
		}
	}
	langIDs = maybeAppendFileLanguageID(langIDs, filename, content.Bytes())

	return langIDs
}

func maybeAppendFileLanguageID(langIDs []int, filename string, content []byte) []int {
	if len(filename) == 0 {
		return langIDs
	}

	if langs := enry.GetLanguages(filename, content); len(langs) == 1 {
		langID, _ := enry.GetLanguageID(langs[0])
		return append(langIDs, langID)
	}

	return langIDs
}
