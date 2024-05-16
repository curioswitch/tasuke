package getuser

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

// Handler is the handler for the FrontendService.GetUser RPC.
type Handler struct {
	store ifirestore.Client
}

// GetUser implements FrontendService.GetUser.
func (h *Handler) GetUser(ctx context.Context, _ *frontendapi.GetUserRequest) (*frontendapi.GetUserResponse, error) {
	fbToken := firebaseauth.TokenFromContext(ctx)

	var u model.User
	if err := h.store.GetDocument(ctx, "users", fbToken.UID, &u); err != nil {
		return nil, fmt.Errorf("getuser: failed to get user document: %w", err)
	}

	return &frontendapi.GetUserResponse{
		User: u.ToProto(),
	}, nil
}
