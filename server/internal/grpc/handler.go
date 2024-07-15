// Package handler содержит обработчики запросов grpc
package handler

import (
	"context"
	"errors"

	"github.com/Andromaril/Gopher-and-secrets/server/internal/construct"
	"github.com/Andromaril/Gopher-and-secrets/server/internal/storage"
	pb "github.com/Andromaril/Gopher-and-secrets/server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	pb.UnimplementedAuthServer
	auth Auth
}

// Auth интерфейс описывающий авторизацию
type Auth interface {
	Login(ctx context.Context, email string, password string) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (userID int64, err error)
}

// Register регистрация grpc сервиса
func Register(gRPCServer *grpc.Server, auth Auth) {
	pb.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

// Login обработчик запроса логина
func (s *serverAPI) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	if in.Email == "" || in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password is required")
	}

	token, err := s.auth.Login(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		if errors.Is(err, construct.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}
		return nil, status.Error(codes.Internal, "failed to login")
	}
	return &pb.LoginResponse{Token: token}, nil
}

// Register обработчик запроса регистрации
func (s *serverAPI) Register(ctx context.Context, in *pb.RegisterRequest,
) (*pb.RegisterResponse, error) {
	if in.Email == "" || in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password is required")
	}

	id, err := s.auth.RegisterNewUser(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "failed to register user")
	}
	return &pb.RegisterResponse{UserId: id}, nil
}
