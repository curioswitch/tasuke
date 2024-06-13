package saveuser

import (
	"context"
	"fmt"
	"strconv"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/v4/auth"
	"github.com/curioswitch/go-usegcp/middleware/firebaseauth"

	"github.com/curioswitch/tasuke/common/tasukedb"
	frontendapi "github.com/curioswitch/tasuke/frontend/api/go"
	ifirestore "github.com/curioswitch/tasuke/frontend/server/internal/client/firestore"
)

// NewHandler returns a Handler that uses the given Firestore client.
func NewHandler(client *firestore.Client) *Handler {
	return &Handler{
		store: ifirestore.NewClient[tasukedb.User](client, "users"),
	}
}

// Handler is the handler for the FrontendService.SaveUser RPC.
type Handler struct {
	store ifirestore.Client[tasukedb.User]
}

// SaveUser implements FrontendService.SaveUser.
func (h *Handler) SaveUser(ctx context.Context, req *frontendapi.SaveUserRequest) (*frontendapi.SaveUserResponse, error) {
	fbToken := firebaseauth.TokenFromContext(ctx)

	githubID, err := githubUserID(fbToken)
	if err != nil {
		return nil, err
	}

	u := tasukedb.User{
		GithubUserID:           int64(githubID),
		ProgrammingLanguageIDs: req.GetUser().GetProgrammingLanguageIds(),
		MaxOpenReviews:         req.GetUser().GetMaxOpenReviews(),
	}

	if err := h.store.SetDocument(ctx, fbToken.UID, &u); err != nil {
		return nil, fmt.Errorf("saveuser: failed to save user document: %w", err)
	}

	return &frontendapi.SaveUserResponse{}, nil
}

func githubUserID(fbToken *auth.Token) (int, error) {
	identity := fbToken.Firebase.Identities["github.com"]
	if identity == nil {
		// We only allow GitHub users so can't happen in practice.
		return 0, fmt.Errorf("saveuser: token not a github user: %v", fbToken.UID)
	}

	if idsAny, ok := identity.([]any); ok && len(idsAny) > 0 {
		if idStr, ok := idsAny[0].(string); ok {
			if id, err := strconv.Atoi(idStr); err == nil {
				return id, nil
			}
		}
	}

	return 0, fmt.Errorf("saveuser: malformed firebase token: %v", fbToken.UID)
}
