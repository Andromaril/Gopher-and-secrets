// Package local содержит все для локальной работы с метаинформацией
package local

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

// TempSecret для хранения информации о секретах
type TempSecret struct {
	SecretID int64  `json:"secret_id"`
	Secret   string `json:"secret"`
	Meta     string `json:"meta"`
	Comment  string `json:"comment"`
}

func dir() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	d := dirname + "\\Desktop\\secrettemp.json"
	return d
}

func dirlocal() string {
	s, err := os.UserCacheDir()
	if err != nil {
		log.Fatal(err)
	}
	dir := s + "\\Temp\\user.json"
	return dir
}

// User для локального хранения метаинформации о юзере
var User = make(map[string]string)

// Secret для хранения информации о секрете
var Secret = make([]TempSecret, 0)

// Storage создание файла
func Storage() error {
	d := dirlocal()
	_, err := os.OpenFile(d, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	return nil
}

// StorageTemp создание файла с секретам
func StorageTemp() error {
	d := dir()
	_, err := os.OpenFile(d, os.O_RDONLY|os.O_CREATE, 0777)
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

// NewTemp добавление в файл нового секрета
func NewTemp() error {
	err := StorageTemp()
	if err != nil {
		log.Error(err)
		return err
	}
	err = LoadTemp()
	if err != nil {
		return err
	}
	return nil
}

// LoadUser чтение файла
func LoadUser() error {
	d := dirlocal()
	data, err := os.ReadFile(d)
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

// LoadTemp загрузка файла с секретами
func LoadTemp() error {
	d := dir()
	datatemp, err := os.ReadFile(d)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if len(datatemp) != 0 {
		return json.Unmarshal(datatemp, &Secret)
	}

	return nil
}

// UpdateUser запись в информации о новом юзере
func UpdateUser() error {
	data, err := json.Marshal(User)
	if err != nil {
		return err
	}
	d := dirlocal()
	err = os.WriteFile(d, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

// UpdateTemp обновления файла с секретами
func UpdateTemp() error {
	d := dir()
	datatemp, err := json.Marshal(Secret)
	if err != nil {
		return err
	}
	err = os.WriteFile(d, datatemp, 0644)
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

// InitStorageTemp для инициализации локального хранилица с секретами
func InitStorageTemp() (err error) {

	if err = NewTemp(); err != nil {
		return err
	}
	return nil
}
