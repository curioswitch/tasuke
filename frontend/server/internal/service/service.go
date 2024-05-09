package service

import (
	"context"

	"connectrpc.com/connect"

	frontendapi "github.com/curioswitch/tasuke/frontend/api"
	"github.com/curioswitch/tasuke/frontend/api/frontendapiconnect"
	"github.com/curioswitch/tasuke/frontend/server/internal/handler/saveuser"
)

// New returns a new service implementation for FrontendService.
func New(saveUser *saveuser.Handler) frontendapiconnect.FrontendServiceHandler {
	return &frontendService{
		saveUser: saveUser,
	}
}

type frontendService struct {
	saveUser *saveuser.Handler

	frontendapiconnect.UnimplementedFrontendServiceHandler
}

// SaveUser implements frontendapiconnect.FrontendServiceHandler.
func (s *frontendService) SaveUser(ctx context.Context, req *connect.Request[frontendapi.SaveUserRequest]) (*connect.Response[frontendapi.SaveUserResponse], error) {
	res, err := s.saveUser.SaveUser(ctx, req.Msg)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(res), nil
}
