// Package storage для взаимодействия с базой данных
package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Andromaril/Gopher-and-secrets/server/internal/model"
	log "github.com/sirupsen/logrus"
)

// переменные ошибок
var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

// Storage структура бд
type Storage struct {
	DB *sql.DB
}

// Init инициализация бд
func Init(storagePath string) (*Storage, error) {

	db, err := sql.Open("pgx", storagePath)
	if err != nil {
		return nil, fmt.Errorf("fatal open sql connection %w", err)
	}
	return &Storage{DB: db}, nil
}

// SaveUser добавление нового пользователя в бд
func (s *Storage) SaveUser(ctx context.Context, login string, password []byte) (int64, error) {
	var id int64
	err := s.DB.QueryRow(`
	INSERT INTO users (login, password)
	VALUES($1, $2) RETURNING id`, login, password).Scan(&id)
	// if err != nil {
	// 	return 0, fmt.Errorf("error insert %w", err)
	// }
	//id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error insert in save user: %w", err)
	}
	return id, nil
}

// GetUser получение пользователя
func (s *Storage) GetUser(ctx context.Context, login string) (model.User, error) {
	rows := s.DB.QueryRowContext(ctx, "SELECT id, login, password FROM users WHERE login=$1", login)
	var user model.User
	err := rows.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		log.Error("error in scan from user select", err)
		return model.User{}, fmt.Errorf("error in scan from user select: %w", err)
	}
	return user, nil
}
