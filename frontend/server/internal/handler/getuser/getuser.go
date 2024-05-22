package getuser

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	"connectrpc.com/connect"
	"github.com/curioswitch/go-usegcp/middleware/firebaseauth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	frontendapi "github.com/curioswitch/tasuke/frontend/api/go"
	ifirestore "github.com/curioswitch/tasuke/frontend/server/internal/client/firestore"
	"github.com/curioswitch/tasuke/frontend/server/internal/model"
)

var errUserNotFound = errors.New("user not found")

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
		if status.Code(err) == codes.NotFound {
			return nil, connect.NewError(connect.CodeNotFound, errUserNotFound)
		}
		return nil, fmt.Errorf("getuser: failed to get user document: %w", err)
	}

	return &frontendapi.GetUserResponse{
		User: u.ToProto(),
	}, nil
}
