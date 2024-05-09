package createuser

import (
	"context"
	"fmt"

	"github.com/curioswitch/go-usegcp/middleware/firebaseauth"

	frontendapi "github.com/curioswitch/tasuke/frontend/api"
	"github.com/curioswitch/tasuke/frontend/server/internal/client/firestore"
	"github.com/curioswitch/tasuke/frontend/server/internal/model"
)

// Handler is the handler for the FrontendService.CreateUser RPC.
type Handler struct {
	store firestore.Client
}

// CreateUser implements FrontendService.CreateUser.
func (h *Handler) CreateUser(ctx context.Context, req *frontendapi.CreateUserRequest) (*frontendapi.CreateUserResponse, error) {
	fbToken := firebaseauth.TokenFromContext(ctx)

	u := model.User{
		ProgrammingLanguageIDs: req.GetUser().GetProgrammingLanguageIds(),
		MaxOpenReviews:         req.GetUser().GetMaxOpenReviews(),
	}

	if err := h.store.CreateDocument(ctx, "users", fbToken.UID, u); err != nil {
		return nil, fmt.Errorf("createuser: failed to create user document: %w", err)
	}

	return &frontendapi.CreateUserResponse{}, nil
}
