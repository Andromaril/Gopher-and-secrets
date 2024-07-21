// Package app содержит функции для старта сервера
package app

import (
	"fmt"
	"log/slog"
	"net"

	authgrpc "github.com/Andromaril/Gopher-and-secrets/server/internal/grpc"
	interceptors "github.com/Andromaril/Gopher-and-secrets/server/internal/interseptors"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// App структура для запуска сервера
type App struct {
	gRPCServer *grpc.Server
	port       string
}

// New создает новый экземпляр структуры App
func New(authService authgrpc.Auth, port string, secretService authgrpc.Secret) *App {

	// recoveryOpts := []recovery.Option{
	// 	recovery.WithRecoveryHandler(func(p interface{}) (err error) {
	// 		log.Error("Recovered from panic", slog.Any("panic", p))

	// 		return status.Errorf(codes.Internal, "internal error")
	// 	}),
	// }

	gRPCServer := grpc.NewServer(grpc.UnaryInterceptor(interceptors.AuthCheck))

	authgrpc.Register(gRPCServer, authService, secretService)

	return &App{
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// MustRun запускает grpc сервер
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run запускает gRPC сервер
func (a *App) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s", a.port))
	if err != nil {
		return fmt.Errorf("error start server: %w", err)
	}

	log.Info("grpc server started ", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("error in serve: %w", err)
	}

	return nil
}

// Stop останавливает gRPC сервер
func (a *App) Stop() {
	const op = "grpcapp.Stop"

	log.Info("stopping gRPC server")

	a.gRPCServer.GracefulStop()
}
