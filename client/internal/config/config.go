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
	LocalStorage string `json:"local_storage"`
}

// Переменные флагов
var (
	GrpcPort     string // адрес запуска сервиса
	Databaseflag string // адрес бд
	ConfigKey    string // файл с конфигом в формате json
	LocalStorage string // временный файл с хранением имени пользователя + токена
)

// ParseFlags для флагов либо переменных окружения
func ParseFlags() {
	flag.StringVar(&GrpcPort, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&Databaseflag, "d", "", "database path")
	flag.StringVar(&ConfigKey, "c", "", "json-file flag")
	flag.StringVar(&LocalStorage, "s", "c:/Users/Мария/AppData/Local/Temp/user.json", "temp file")
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

	if envLocalStorage := os.Getenv("LOCAL_STORAGE"); envLocalStorage != "" {
		LocalStorage = envLocalStorage
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
		if LocalStorage == "" {
			LocalStorage = conf.LocalStorage
		}
	}
}
