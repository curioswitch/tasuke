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

	"github.com/curioswitch/tasuke/common/tasukedb"
	frontendapi "github.com/curioswitch/tasuke/frontend/api/go"
	ifirestore "github.com/curioswitch/tasuke/frontend/server/internal/client/firestore"
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

	var u tasukedb.User
	if err := h.store.GetDocument(ctx, "users", fbToken.UID, &u); err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, connect.NewError(connect.CodeNotFound, errUserNotFound)
		}
		return nil, fmt.Errorf("getuser: failed to get user document: %w", err)
	}

	return &frontendapi.GetUserResponse{
		User: userToProto(&u),
	}, nil
}

// ToProto converts a User to its API representation.
func userToProto(u *tasukedb.User) *frontendapi.User {
	return &frontendapi.User{
		ProgrammingLanguageIds: u.ProgrammingLanguageIDs,
		MaxOpenReviews:         u.MaxOpenReviews,
	}
}
