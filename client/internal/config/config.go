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
	TempSecret   string `json:"temp_storage"`
}

// Переменные флагов
var (
	GrpcPort     string // адрес запуска сервиса
	Databaseflag string // адрес бд
	ConfigKey    string // файл с конфигом в формате json
	LocalStorage string // временный файл с хранением имени пользователя + токена
	TempSecret   string // временный файл для хранения секретов
)

// ParseFlags для флагов либо переменных окружения
func ParseFlags() {
	flag.StringVar(&GrpcPort, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&Databaseflag, "d", "", "database path")
	flag.StringVar(&ConfigKey, "c", "", "json-file flag")
	flag.StringVar(&LocalStorage, "s", "c:/Users/Мария/AppData/Local/Temp/user.json", "temp file")
	flag.StringVar(&TempSecret, "t", "c:/Users/Мария/AppData/Local/Temp/secrettemp.json", "temp secret file")
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

	if envTempSecret := os.Getenv("TEMP"); envTempSecret != "" {
		TempSecret = envTempSecret
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
		if TempSecret == "" {
			TempSecret = conf.TempSecret
		}
	}
}
