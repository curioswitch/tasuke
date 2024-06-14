package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/go-enry/go-enry/v2"
	"github.com/google/go-github/v62/github"

	"github.com/curioswitch/tasuke/common/tasukedb"
	"github.com/curioswitch/tasuke/webhook/server/internal/config"
	"github.com/curioswitch/tasuke/webhook/server/internal/ghapi"
)

type githubEventType byte

const (
	githubEventTypeIssueComment githubEventType = iota
	githubEventTypePullRequest
)

func New(config *config.Config, fsClient *firestore.Client) (*Handler, error) {
	ghCreator, err := ghapi.NewClientCreator(config)
	if err != nil {
		return nil, fmt.Errorf("handler: create client creator: %w", err)
	}
	return &Handler{
		secret:    []byte(config.GitHub.Secret),
		ghCreator: ghCreator,

		store: fsClient,
		users: fsClient.Collection("users"),
	}, nil
}

type Handler struct {
	secret    []byte
	ghCreator *ghapi.ClientCreator

	store *firestore.Client
	users *firestore.CollectionRef
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

	payload, err := github.ValidatePayload(r, h.secret)
	if err != nil {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	switch eventType {
	case githubEventTypeIssueComment:
		var event github.IssueCommentEvent
		if err := json.Unmarshal(payload, &event); err != nil {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}
		err = h.handleIssueComment(ctx, &event)
	case githubEventTypePullRequest:
		var event github.PullRequestEvent
		if err := json.Unmarshal(payload, &event); err != nil {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}
		err = h.handlePullRequest(ctx, &event)
	}

	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Failed to handle event: %v", err))
		http.Error(w, "Failed to handle event", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) handleIssueComment(ctx context.Context, event *github.IssueCommentEvent) error {
	if !event.Issue.IsPullRequest() {
		return nil
	}

	body := strings.TrimSpace(event.Comment.GetBody())

	if !strings.HasPrefix(body, "/tasuke") {
		return nil
	}

	gh, err := h.ghCreator.NewClient(event.Installation.GetID())
	if err != nil {
		return fmt.Errorf("handler: create client: %w", err)
	}

	owner := event.Repo.GetOwner().GetLogin()
	repo := event.Repo.GetName()
	num := event.Issue.GetNumber()

	diff, _, err := gh.PullRequests.GetRaw(ctx, owner, repo, num, github.RawOptions{
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

	review := tasukedb.Review{
		Repo:        fmt.Sprintf("%s/%s", owner, repo),
		PullRequest: int64(num),
	}

	if existingReviews, err := h.getExistingReviews(ctx, review.Repo, review.PullRequest); err != nil {
		return fmt.Errorf("handler: check existing reviews: %w", err)
	} else if len(existingReviews) > 0 {
		if _, _, err := gh.Issues.CreateComment(ctx, owner, repo, num, &github.IssueComment{
			Body: github.String(`Sorry, this PR already has a reviewer assigned to it.`),
		}); err != nil {
			return fmt.Errorf("handler: create comment: %w", err)
		}
		return nil
	}

	// TODO: Actually match with a reviewer. First, we just find any.
	var user tasukedb.User
	existing := false

	reviewID := fmt.Sprintf("%s-%s-%d", owner, repo, num)

	if err := h.store.RunTransaction(ctx, func(ctx context.Context, t *firestore.Transaction) error {
		// Get candidate docs without using transaction. We don't want all candidates to be docs, but
		// do want to fetch fresh candidates on retries, at least for now.
		docs, err := h.users.WhereEntity(firestore.PropertyFilter{
			Path:     "remainingReviews",
			Operator: ">",
			Value:    0,
		}).Limit(100).Documents(ctx).GetAll()
		if err != nil {
			return fmt.Errorf("handler: get available users: %w", err)
		}
		if len(docs) == 0 {
			return nil
		}

		// For now just pick a random user.
		doc := docs[rand.Intn(len(docs))] //nolint:gosec // We don't need cryptographically secure randomness here.

		// Refetch the doc within the transaction.
		doc, err = t.Get(doc.Ref)
		if err != nil {
			return fmt.Errorf("handler: refetch user: %w", err)
		}

		if err := doc.DataTo(&user); err != nil {
			return fmt.Errorf("handler: parse user: %w", err)
		}

		user.RemainingReviews--
		if err := t.Update(doc.Ref, []firestore.Update{
			{
				Path:  "remainingReviews",
				Value: user.RemainingReviews,
			},
		}); err != nil {
			return fmt.Errorf("handler: update remaining reviews: %w", err)
		}

		if _, err := doc.Ref.Collection("reviews").Doc(reviewID).Set(ctx, review); err != nil {
			return fmt.Errorf("handler: create review: %w", err)
		}

		return nil
	}); err != nil {
		return err
	}

	if existing {
		if _, _, err := gh.Issues.CreateComment(ctx, owner, repo, num, &github.IssueComment{
			Body: github.String(`Sorry, this PR already has a reviewer assigned to it.`),
		}); err != nil {
			return fmt.Errorf("handler: create comment: %w", err)
		}
		return nil
	}

	ghUser, _, err := gh.Users.GetByID(ctx, user.GithubUserID)
	if err != nil {
		return fmt.Errorf("handler: get user: %w", err)
	}

	if _, _, err := gh.Issues.CreateComment(ctx, owner, repo, num, &github.IssueComment{
		Body: github.String(fmt.Sprintf(`Hi there! I'm currently under development but will be happy to help after I'm working.

Currently, I only detect the languages in the PR. The languages I detected for this PR are: %s.

I don't actually match against these languages yet. But I do still find an arbitrary user to nag.
@%s, could you please review this PR? Thanks!
`, strings.Join(langs, ", "), ghUser.GetLogin())),
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

func (h *Handler) handlePullRequest(ctx context.Context, event *github.PullRequestEvent) error {
	if *event.Action != "closed" {
		return nil
	}

	owner := event.Repo.GetOwner().GetLogin()
	repo := event.Repo.GetName()
	num := event.GetNumber()

	reviewDocs, err := h.getExistingReviews(ctx, fmt.Sprintf("%s/%s", owner, repo), int64(num))
	if err != nil {
		return fmt.Errorf("handler: get existing reviews: %w", err)
	}

	// Currently we only support one reviewer. Even if we supported multiple in the future,
	// making parallel is likely not important, we do iterate anyways.

	for _, doc := range reviewDocs {
		var review tasukedb.Review
		if err := doc.DataTo(&review); err != nil {
			return fmt.Errorf("handler: parse review: %w", err)
		}

		if err := h.store.RunTransaction(ctx, func(ctx context.Context, t *firestore.Transaction) error {
			if err := t.Update(doc.Ref, []firestore.Update{
				{
					Path:  "completed",
					Value: true,
				},
			}); err != nil {
				return fmt.Errorf("handler: mark review completed: %w", err)
			}

			userDoc, err := doc.Ref.Parent.Parent.Get(ctx)
			if err != nil {
				return fmt.Errorf("handler: get user for review completion: %w", err)
			}

			var user tasukedb.User
			if err := userDoc.DataTo(&user); err != nil {
				return fmt.Errorf("handler: parse user: %w", err)
			}

			user.RemainingReviews++

			if err := t.Update(userDoc.Ref, []firestore.Update{
				{
					Path:  "remainingReviews",
					Value: user.RemainingReviews,
				},
			}); err != nil {
				return fmt.Errorf("handler: increment user remaining reviews: %w", err)
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

func (h *Handler) getExistingReviews(ctx context.Context, repo string, num int64) ([]*firestore.DocumentSnapshot, error) {
	return h.store.CollectionGroup("reviews").WhereEntity(firestore.AndFilter{
		Filters: []firestore.EntityFilter{
			firestore.PropertyFilter{
				Path:     "repo",
				Operator: "==",
				Value:    repo,
			},
			firestore.PropertyFilter{
				Path:     "pullRequest",
				Operator: "==",
				Value:    num,
			},
			firestore.PropertyFilter{
				Path:     "completed",
				Operator: "==",
				Value:    false,
			},
		},
	}).Documents(ctx).GetAll()
}
