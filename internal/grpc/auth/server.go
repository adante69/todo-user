package auth

import (
	"context"
	"errors"
	ssov1 "github.com/adante69/todo-protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"todo-sso/internal/services/auth"
	"todo-sso/internal/storage"
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
}
type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRpc *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRpc, &serverAPI{auth: auth})
}

const (
	emptyvalue = 0
)

func (s *serverAPI) Login(
	ctx context.Context,
	req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {

	if err := validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))

	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.InvalidArgument, "failed to login")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, "failed to register new user")
	}
	return &ssov1.RegisterResponse{
		UserId: uint64(userID),
	}, nil
}

func validateLogin(req *ssov1.LoginRequest) error {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "empty email or password")
	}

	if req.GetAppId() == emptyvalue {
		return status.Error(codes.InvalidArgument, "app id is required")
	}
	return nil
}

func validateRegister(req *ssov1.RegisterRequest) error {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "empty email or password")
	}
	return nil
}
