// Package config считывает флаги сервера
package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

// Config структура конфига
type Config struct {
	GrpcPort     string `json:"port"`
	Databaseflag string `json:"database_dsn"`
}

var (
	GrpcPort     string // адрес запуска сервиса
	Databaseflag string // адрес бд
	ConfigKey    string // файл с конфигом в формате json
)

// ParseFlags для флагов либо переменных окружения
func ParseFlags() {
	flag.StringVar(&GrpcPort, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&Databaseflag, "d", "postgres://postgres:qwerty123@localhost:5432/server2", "database path")
	flag.StringVar(&ConfigKey, "c", "", "json-file flag")
	flag.Parse()
	if envGrpcPort := os.Getenv("ADDRESS"); envGrpcPort != "" {
		GrpcPort = envGrpcPort
	}
	if envDatabaseflag := os.Getenv("DATABASE_DSN"); envDatabaseflag != "" {
		Databaseflag = envDatabaseflag
	}
	if envConfigKey := os.Getenv("CONFIG"); envConfigKey != "" {
		ConfigKey = envConfigKey
	}

	if ConfigKey != "" {
		c, err := os.ReadFile(ConfigKey)
		if err != nil {
			panic(err)
		}
		var conf Config
		err = json.Unmarshal(c, &conf)
		if err != nil {
			log.Fatal(err)
		}
		if GrpcPort == "localhost:8080" {
			GrpcPort = conf.GrpcPort
		}
		if Databaseflag == "" {
			Databaseflag = conf.Databaseflag
		}
	}
}
