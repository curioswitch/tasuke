package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"slices"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/go-enry/go-enry/v2"
	"github.com/google/go-github/v62/github"
	"google.golang.org/api/iterator"

	"github.com/curioswitch/tasuke/common/languages"
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

	requesterDoc, err := h.users.WhereEntity(firestore.PropertyFilter{
		Path:     "githubUserId",
		Operator: "==",
		Value:    event.Sender.GetID(),
	}).Limit(1).Documents(ctx).Next()
	if err != nil && !errors.Is(err, iterator.Done) {
		return fmt.Errorf("handler: get sender: %w", err)
	}
	var requester tasukedb.User
	if requesterDoc != nil {
		if err := requesterDoc.DataTo(&requester); err != nil {
			return fmt.Errorf("handler: parse sender: %w", err)
		}
	}
	if requester.RemainingReviews <= 0 {
		if _, _, err := gh.Issues.CreateComment(ctx, owner, repo, num, &github.IssueComment{
			Body: github.String(`Sorry, currently we require review requests to come from registered reviewers (this may change in the future).
Please create an account at https://tasuke.dev and set maximum open reviews to at least 1.`),
		}); err != nil {
			return fmt.Errorf("handler: create comment: %w", err)
		}
		return nil
	}

	diff, _, err := gh.PullRequests.GetRaw(ctx, owner, repo, num, github.RawOptions{
		Type: github.Diff,
	})
	if err != nil {
		return fmt.Errorf("handler: get pull request diff: %w", err)
	}

	mainLangID := diffMainLanguage(diff)
	var mainLang string
	if mainLangID >= 0 {
		lang, _ := enry.GetLanguageInfoByID(mainLangID)
		mainLang = lang.Name
	} else {
		if _, _, err := gh.Issues.CreateComment(ctx, owner, repo, num, &github.IssueComment{
			Body: github.String(`Sorry, I couldn't detect a supported programming language for this PR.
Please make sure it is one of our [supported languages](https://github.com/github-linguist/linguist/blob/master/lib/linguist/popular.yml).`),
		}); err != nil {
			return fmt.Errorf("handler: create comment: %w", err)
		}
		return nil
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

	var user *tasukedb.User
	existing := false

	reviewID := fmt.Sprintf("%s-%s-%d", owner, repo, num)

	if err := h.store.RunTransaction(ctx, func(ctx context.Context, t *firestore.Transaction) error {
		// Get candidate docs without using transaction. We don't want all candidates to be locks, but
		// do want to fetch fresh candidates on retries, at least for now.
		docs, err := h.users.WhereEntity(firestore.AndFilter{
			Filters: []firestore.EntityFilter{
				firestore.PropertyFilter{
					Path:     "programmingLanguageIds",
					Operator: "array-contains",
					Value:    mainLangID,
				},
				firestore.PropertyFilter{
					Path:     "remainingReviews",
					Operator: ">",
					Value:    0,
				},
			},
		}).OrderBy("remainingReviews", firestore.Desc).Limit(100).Documents(ctx).GetAll()
		if err != nil {
			return fmt.Errorf("handler: get available users: %w", err)
		}
		if len(docs) == 0 {
			return nil
		}

		docs = slices.DeleteFunc(docs, func(doc *firestore.DocumentSnapshot) bool {
			return doc.Ref.ID == requesterDoc.Ref.ID
		})

		// For now just pick a random user.
		doc := docs[rand.Intn(len(docs))] //nolint:gosec // We don't need cryptographically secure randomness here.

		// Refetch the doc within the transaction.
		doc, err = t.Get(doc.Ref)
		if err != nil {
			return fmt.Errorf("handler: refetch user: %w", err)
		}

		var u tasukedb.User
		if err := doc.DataTo(&u); err != nil {
			return fmt.Errorf("handler: parse user: %w", err)
		}

		u.RemainingReviews--
		if err := t.Update(doc.Ref, []firestore.Update{
			{
				Path:  "remainingReviews",
				Value: u.RemainingReviews,
			},
		}); err != nil {
			return fmt.Errorf("handler: update remaining reviews: %w", err)
		}

		if err := t.Set(doc.Ref.Collection("reviews").Doc(reviewID), review); err != nil {
			return fmt.Errorf("handler: create review: %w", err)
		}

		user = &u

		return nil
	}); err != nil {
		return err
	}

	if user == nil {
		if _, _, err := gh.Issues.CreateComment(ctx, owner, repo, num, &github.IssueComment{
			Body: github.String(fmt.Sprintf(`Sorry, I could not find an available reviewer for %s. Please try again later.`, mainLang)),
		}); err != nil {
			return fmt.Errorf("handler: create comment: %w", err)
		}
		return nil
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
		Body: github.String(fmt.Sprintf(`Hi there! @%s, could you please review this PR for %s? Thanks!
`, ghUser.GetLogin(), mainLang)),
	}); err != nil {
		return fmt.Errorf("handler: create comment: %w", err)
	}

	return nil
}

func diffMainLanguage(diff string) int {
	langLines := map[int]int{}

	skip := false

	filename := ""
	lines := 0
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

			addLanguageLines(langLines, filename, content.Bytes(), lines)
			skip = false
			filename = ""
			lines = 0
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
			lines++
		}
	}
	addLanguageLines(langLines, filename, content.Bytes(), lines)

	mainLang := -1
	mainLangLines := -1

	for langID, lines := range langLines {
		if lines > mainLangLines {
			mainLang = langID
			mainLangLines = lines
		}
	}

	return mainLang
}

func addLanguageLines(langLines map[int]int, filename string, content []byte, lines int) {
	if len(filename) == 0 {
		return
	}

	if langs := enry.GetLanguages(filename, content); len(langs) == 1 {
		langID, _ := enry.GetLanguageID(langs[0])
		if languages.IsSupported(langID) {
			langLines[langID] += lines
		}
	}
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

	if (len(reviewDocs)) == 0 {
		return nil
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
