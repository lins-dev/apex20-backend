package http

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	v1 "github.com/apex20/contracts/proto/apex20/v1"
	"github.com/apex20/contracts/proto/apex20/v1/apex20v1connect"
	"github.com/apex20/backend/internal/application/port"
)

// AuthUseCases groups the use cases needed for the auth routes.
type AuthUseCases struct {
	SignUp port.UserSignUpper
	SignIn port.UserSignIner
}

type authServiceServer struct {
	uc AuthUseCases
}

var _ apex20v1connect.AuthServiceHandler = (*authServiceServer)(nil)

func RegisterAuthHandler(server *ChiServer, uc AuthUseCases) {
	path, handler := apex20v1connect.NewAuthServiceHandler(&authServiceServer{uc: uc})
	server.Mount(path, handler)
}

func (s *authServiceServer) SignUp(ctx context.Context, req *connect.Request[v1.SignUpRequest]) (*connect.Response[v1.SignUpResponse], error) {
	out, err := s.uc.SignUp.Execute(ctx, port.SignUpInput{
		Email:    req.Msg.Email,
		Password: req.Msg.Password,
		Name:     req.Msg.Name,
	})
	if err != nil {
		if errors.Is(err, port.ErrEmailAlreadyExists) {
			return nil, connect.NewError(connect.CodeAlreadyExists, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&v1.SignUpResponse{
		UserId:      out.User.ID.String(),
		AccessToken: out.AccessToken,
	}), nil
}

func (s *authServiceServer) SignIn(ctx context.Context, req *connect.Request[v1.SignInRequest]) (*connect.Response[v1.SignInResponse], error) {
	out, err := s.uc.SignIn.Execute(ctx, port.SignInInput{
		Email:    req.Msg.Email,
		Password: req.Msg.Password,
	})
	if err != nil {
		if errors.Is(err, port.ErrInvalidCredentials) {
			return nil, connect.NewError(connect.CodeUnauthenticated, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&v1.SignInResponse{
		UserId:      out.User.ID.String(),
		AccessToken: out.AccessToken,
	}), nil
}
