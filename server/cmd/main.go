package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Andromaril/Gopher-and-secrets/server/internal/app"
	"github.com/Andromaril/Gopher-and-secrets/server/internal/config"
	log "github.com/sirupsen/logrus"
)

func main() {
	config.ParseFlags()
	application := app.New(config.GrpcPort, config.Databaseflag, 10000)

	go func() {
		application.GRPCServer.MustRun()
	}()

	// Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()
	log.Info("Gracefully stopped")
}
