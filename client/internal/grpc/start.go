// Package grpc для инициализации grpc клиента
package grpc

import (
	"github.com/Andromaril/Gopher-and-secrets/client/internal/config"
	pb "github.com/Andromaril/Gopher-and-secrets/server/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Init инициализация клиента
func Init() (pb.AuthClient, error) {
	conn, err := grpc.NewClient(config.GrpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	log.Infof("grpc client started adress=%s", config.ConfigKey)
	c := pb.NewAuthClient(conn)
	return c, nil
}
