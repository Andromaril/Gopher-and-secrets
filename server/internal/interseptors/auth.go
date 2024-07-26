package interceptors

import (
	"context"

	"github.com/Andromaril/Gopher-and-secrets/server/secret"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	skipmethods = map[string]struct{}{
		"/server.Auth/Register": {},
		"/server.Auth/Login":    {},
	}
)

// AuthCheck interceptor для верификации jwt-токена
func AuthCheck(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log.Println(info.FullMethod)
	_, ok := skipmethods[info.FullMethod]
	if ok {
		return handler(ctx, req)
	}
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	userID, err := secret.DecodeToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "неверный токен: %v", err)
	}

	return handler(context.WithValue(ctx, "id", userID), req)
}
