// Package app для старта сервера
package app

import (
	"time"

	grpcapp "github.com/Andromaril/Gopher-and-secrets/server/internal/app/grpc"
	"github.com/Andromaril/Gopher-and-secrets/server/internal/config"
	"github.com/Andromaril/Gopher-and-secrets/server/internal/construct"
	"github.com/Andromaril/Gopher-and-secrets/server/internal/storage"

	_ "github.com/jackc/pgx/v5/stdlib"
	log "github.com/sirupsen/logrus"
)

// App структура сервера
type App struct {
	GRPCServer *grpcapp.App
}

// New создает новый экземпляр структуры App
func New(grpcPort string, storagePath string, tokenTTL time.Duration) *App {
	db, err := storage.Init(config.Databaseflag)
	if err != nil {
		panic(err)
	}
	log.Infof("Init database")
	//defer db2.Close()
	authService := construct.New(db, db, tokenTTL, db)
	secretService := construct.NewSecret(db, db, tokenTTL, db)
	grpcApp := grpcapp.New(authService, grpcPort, secretService)

	return &App{
		GRPCServer: grpcApp,
	}
}
