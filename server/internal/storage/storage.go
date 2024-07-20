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
	ErrUserExists     = errors.New("user already exists")
	ErrUserNotFound   = errors.New("user not found")
	ErrSecretNotFound = errors.New("secret not found")
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
	if err != nil {
		return 0, fmt.Errorf("error insert %w", err)
	}
	// if err != nil {
	// 	var pgErr *pgconn.PgError
	// 	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
	// 		return 0, fmt.Errorf("error insert %w", ErrUserExists)
	// 	}
	// }

	return id, nil
}

// GetUser получение пользователя
func (s *Storage) GetUser(ctx context.Context, login string) (model.User, error) {
	rows := s.DB.QueryRowContext(ctx, "SELECT id, login, password FROM users WHERE login=$1", login)
	var user model.User
	err := rows.Scan(&user.ID, &user.Login, &user.PasswordHash)
	if err != nil {
		log.Error("error in scan from user select ", err)
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, fmt.Errorf("error in scan from user select: %w", ErrUserNotFound)
		}
	}
	return user, nil
}

// SaveSecret добавление нового секрета в бд
func (s *Storage) SaveSecret(ctx context.Context, userID int64, secret string, meta string, comment string) (int64, error) {
	var id int64
	err := s.DB.QueryRow(`
	INSERT INTO secrets (user_id, secret, meta, comment)
	VALUES($1, $2, $3, $4) RETURNING id`, userID, secret, meta, comment).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error insert in save secrets: %w", err)
	}
	return id, nil
}

// GetSecret получает секрет из бд
func (s *Storage) GetSecret(ctx context.Context, userID int64, meta string) ([]model.Secret, error) {
	rows, err := s.DB.QueryContext(ctx, "SELECT secret, comment FROM secrets WHERE user_id=$1 AND meta=$2", userID, meta)
	secret := make([]model.Secret, 0)
	if err != nil {
		log.Error("error in scan from secret select ", err)
		if errors.Is(err, sql.ErrNoRows) {
			return secret, fmt.Errorf("error in scan from secret select: %w", ErrSecretNotFound)
		}
	}
	defer rows.Close()
	for rows.Next() {
		var result model.Secret
		err = rows.Scan(&result.Secret, &result.Comment)
		if err != nil {
			return secret, fmt.Errorf("invalid scan when get secrets %w", err)
		}
		secret = append(secret, model.Secret{Secret: result.Secret, Comment: result.Comment})
	}
	err = rows.Err()
	if err != nil {
		return secret, fmt.Errorf("error select %w", err)
	}
	return secret, nil
}

// UpdateSecret для обновления секрета
func (s *Storage) UpdateSecret(ctx context.Context, userID int64, secret string, secretnew string) error {
	_, err := s.DB.ExecContext(ctx, `
	UPDATE secrets SET secret=$1 WHERE user_id=$2 AND secret=$3`, secretnew, userID, secret)
	if err != nil {
		return fmt.Errorf("error insert %w", err)
	}
	return nil
}

// DeleteSecret для удаления секрета
func (s *Storage) DeleteSecret(ctx context.Context, userID int64, secret string) error {
	_, err := s.DB.ExecContext(ctx, `
	DELETE FROM secrets WHERE secret=$1 AND user_id=$2`, secret, userID)
	if err != nil {
		return fmt.Errorf("error delete %w", err)
	}
	return nil
}
