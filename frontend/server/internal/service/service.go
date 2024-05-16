package service

import (
	"context"

	"connectrpc.com/connect"

	frontendapi "github.com/curioswitch/tasuke/frontend/api/go"
	"github.com/curioswitch/tasuke/frontend/api/go/frontendapiconnect"
	"github.com/curioswitch/tasuke/frontend/server/internal/handler/getuser"
	"github.com/curioswitch/tasuke/frontend/server/internal/handler/saveuser"
)

// New returns a new service implementation for FrontendService.
func New(getUser *getuser.Handler, saveUser *saveuser.Handler) frontendapiconnect.FrontendServiceHandler {
	return &frontendService{
		getUser:  getUser,
		saveUser: saveUser,
	}
}

type frontendService struct {
	getUser  *getuser.Handler
	saveUser *saveuser.Handler

	frontendapiconnect.UnimplementedFrontendServiceHandler
}

// GetUser implements frontendapiconnect.FrontendServiceHandler.
func (s *frontendService) GetUser(ctx context.Context, req *connect.Request[frontendapi.GetUserRequest]) (*connect.Response[frontendapi.GetUserResponse], error) {
	res, err := s.getUser.GetUser(ctx, req.Msg)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(res), nil
}

// SaveUser implements frontendapiconnect.FrontendServiceHandler.
func (s *frontendService) SaveUser(ctx context.Context, req *connect.Request[frontendapi.SaveUserRequest]) (*connect.Response[frontendapi.SaveUserResponse], error) {
	res, err := s.saveUser.SaveUser(ctx, req.Msg)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(res), nil
}
