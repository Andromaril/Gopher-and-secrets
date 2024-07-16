// Package app содержит функции для старта сервера
package app

import (
	"fmt"
	"log/slog"
	"net"

	authgrpc "github.com/Andromaril/Gopher-and-secrets/server/internal/grpc"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// App структура для запуска сервера
type App struct {
	gRPCServer *grpc.Server
	port       string
}

// New создает новый экземпляр структуры App
func New(authService authgrpc.Auth, port string) *App {

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Error("Recovered from panic", slog.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
	))

	authgrpc.Register(gRPCServer, authService)

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
