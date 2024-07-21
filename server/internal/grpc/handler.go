// Package handler содержит обработчики запросов grpc
package handler

import (
	"context"
	"errors"

	"github.com/Andromaril/Gopher-and-secrets/server/internal/construct"
	"github.com/Andromaril/Gopher-and-secrets/server/internal/model"
	"github.com/Andromaril/Gopher-and-secrets/server/internal/storage"
	pb "github.com/Andromaril/Gopher-and-secrets/server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var success = "success"

type serverAPI struct {
	pb.UnimplementedAuthServer
	auth   Auth
	secret Secret
}

// Auth интерфейс описывающий авторизацию
type Auth interface {
	Login(ctx context.Context, email string, password string) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (userID int64, err error)
}

// Secret интерфейс описывающий секрет
type Secret interface {
	SaveSecret(ctx context.Context, userID int64, secret string, meta string, comment string) (int64, error)
	GetNewSecret(ctx context.Context, userID int64, meta string) ([]model.Secret, error)
	UpdateSecret(ctx context.Context, userID int64, secret string, secretnew string) error
	DeleteSecret(ctx context.Context, userID int64, sec string) error
	GetAll(ctx context.Context, userID int64) ([]model.Secret, error)
}

// Register регистрация grpc сервиса
func Register(gRPCServer *grpc.Server, auth Auth, secret Secret) {
	pb.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth, secret: secret})
}

// Login обработчик запроса логина
func (s *serverAPI) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	if in.Login == "" || in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password is required")
	}

	token, err := s.auth.Login(ctx, in.GetLogin(), in.GetPassword())
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
	if in.Login == "" || in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password is required")
	}

	id, err := s.auth.RegisterNewUser(ctx, in.GetLogin(), in.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "failed to register user")
	}
	return &pb.RegisterResponse{UserId: id}, nil
}

func (s *serverAPI) AddSecret(ctx context.Context, in *pb.AddSecretRequest,
) (*pb.AddSecretResponse, error) {
	id := ctx.Value("id").(int64)
	if id == 0 || in.Secret == "" {
		return nil, status.Error(codes.InvalidArgument, "user id and secret is required")
	}

	id, err := s.secret.SaveSecret(ctx, id, in.GetSecret(), in.GetMeta(), in.GetComment())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to save secret")
	}
	return &pb.AddSecretResponse{SecretId: id}, nil
}

func (s *serverAPI) GetSecret(ctx context.Context, in *pb.GetSecretRequest,
) (*pb.GetSecretResponse, error) {
	id := ctx.Value("id").(int64)
	if id == 0 {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	secret, err := s.secret.GetNewSecret(ctx, id, in.GetMeta())
	if err != nil {
		if errors.Is(err, storage.ErrSecretNotFound) {
			return nil, status.Error(codes.NotFound, "secret not found")
		}

		return nil, status.Error(codes.Internal, "failed to get secret")
	}
	var pbSecrets []*pb.Secret
	for _, secret := range secret {
		pbSecrets = append(pbSecrets, &pb.Secret{
			SecretId: secret.SecretID,
			UserId:   secret.ID,
			Secret:   secret.Secret,
			Meta:     secret.Meta,
			Comment:  secret.Comment,
		})
	}
	return &pb.GetSecretResponse{Secret: pbSecrets}, nil
}

func (s *serverAPI) UpdateSecret(ctx context.Context, in *pb.UpdateSecretRequest,
) (*pb.UpdateSecretResponse, error) {
	id := ctx.Value("id").(int64)
	if id == 0 || in.Secret == "" {
		return nil, status.Error(codes.InvalidArgument, "user id and secret is required")
	}

	err := s.secret.UpdateSecret(ctx, id, in.GetSecret(), in.GetSecretNew())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update secret")
	}
	return &pb.UpdateSecretResponse{Status: success}, nil
}

func (s *serverAPI) DeleteSecret(ctx context.Context, in *pb.DeleteSecretRequest,
) (*pb.DeleteSecretResponse, error) {
	id := ctx.Value("id").(int64)
	if id == 0 || in.Secret == "" {
		return nil, status.Error(codes.InvalidArgument, "user id and secret is required")
	}

	err := s.secret.DeleteSecret(ctx, id, in.GetSecret())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update secret")
	}
	return &pb.DeleteSecretResponse{Status: success}, nil
}

func (s *serverAPI) GetAll(ctx context.Context, in *pb.GetAllRequest,
) (*pb.GetAllResponse, error) {
	if in.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	secret, err := s.secret.GetAll(ctx, in.GetUserId())
	if err != nil {
		if errors.Is(err, storage.ErrSecretNotFound) {
			return nil, status.Error(codes.NotFound, "secret not found")
		}

		return nil, status.Error(codes.Internal, "failed to get secret")
	}
	var pbSecrets []*pb.Secret
	for _, secret := range secret {
		pbSecrets = append(pbSecrets, &pb.Secret{
			SecretId: secret.SecretID,
			UserId:   secret.ID,
			Secret:   secret.Secret,
			Meta:     secret.Meta,
			Comment:  secret.Comment,
		})
	}
	return &pb.GetAllResponse{Secret: pbSecrets}, nil
}
