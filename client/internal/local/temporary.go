// Package local содержит все для локальной работы с метаинформацией
package local

import (
	"encoding/json"
	"os"

	"github.com/Andromaril/Gopher-and-secrets/client/internal/config"
	log "github.com/sirupsen/logrus"
)

// User для локального хранения метаинформации о юзере
var User = make(map[string]string)

// Storage создание файла
func Storage() error {
	_, err := os.OpenFile(config.LocalStorage, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}

	return nil
}

// NewUser добавление в файл нового пользователя:токен
func NewUser() error {
	err := Storage()
	if err != nil {
		log.Error(err)
		return err
	}
	err = LoadUser()
	if err != nil {
		return err
	}
	return nil
}

// LoadUser чтение файла
func LoadUser() error {
	data, err := os.ReadFile(config.LocalStorage)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if len(data) != 0 {
		return json.Unmarshal(data, &User)
	}

	return nil
}

// UpdateUser запись в информации о новом юзере
func UpdateUser() error {
	data, err := json.Marshal(User)
	if err != nil {
		return err
	}
	err = os.WriteFile(config.LocalStorage, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

// InitStorage для инициализации локального хранилища
func InitStorage() (err error) {
	if err = NewUser(); err != nil {
		return err
	}
	return nil
}
