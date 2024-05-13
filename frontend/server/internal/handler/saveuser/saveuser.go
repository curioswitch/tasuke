package saveuser

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/curioswitch/go-usegcp/middleware/firebaseauth"

	frontendapi "github.com/curioswitch/tasuke/frontend/api/go"
	ifirestore "github.com/curioswitch/tasuke/frontend/server/internal/client/firestore"
	"github.com/curioswitch/tasuke/frontend/server/internal/model"
)

// NewHandler returns a Handler that uses the given Firestore client.
func NewHandler(client *firestore.Client) *Handler {
	return &Handler{
		store: ifirestore.NewClient(client),
	}
}

// Handler is the handler for the FrontendService.SaveUser RPC.
type Handler struct {
	store ifirestore.Client
}

// SaveUser implements FrontendService.SaveUser.
func (h *Handler) SaveUser(ctx context.Context, req *frontendapi.SaveUserRequest) (*frontendapi.SaveUserResponse, error) {
	fbToken := firebaseauth.TokenFromContext(ctx)

	u := model.User{
		ProgrammingLanguageIDs: req.GetUser().GetProgrammingLanguageIds(),
		MaxOpenReviews:         req.GetUser().GetMaxOpenReviews(),
	}

	if err := h.store.SetDocument(ctx, "users", fbToken.UID, u); err != nil {
		return nil, fmt.Errorf("saveuser: failed to save user document: %w", err)
	}

	return &frontendapi.SaveUserResponse{}, nil
}
