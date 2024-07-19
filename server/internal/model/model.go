// Package models содержит модели, используемые сервисом
package model

// User описывает пользователя
type User struct {
	ID           int64
	Login        string
	PasswordHash []byte
}

// Secret описывает секрет
type Secret struct {
	SecretID int64
	ID       int64
	Secret   string
	Meta     string
	Comment  string
}
