// Package models содержит модели, используемые сервисом
package model

// User описывает пользователя
type User struct {
	ID       int64
	Login    string
	PasswordHash []byte
}

