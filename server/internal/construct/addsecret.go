package construct

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// Secret структура для секретов
type Secret struct {
	usrSave  UserSave
	usrGet   UserGet
	scrSave  SecretSave
	tokenTTL time.Duration
}

func NewSecret(userSave UserSave, userGet UserGet, tokenTTL time.Duration, scrSave SecretSave) *Secret {
	return &Secret{
		usrSave:  userSave,
		usrGet:   userGet,
		tokenTTL: tokenTTL,
		scrSave:  scrSave,
	}
}

// SecretSave для сохранения секретов
type SecretSave interface {
	SaveSecret(ctx context.Context, userID int64, secret []byte, meta string, comment []byte) (uid int64, err error)
}

// SaveNewSecret для сохранения секретов
func (s *Secret) SaveSecret(ctx context.Context, userID int64, secret []byte, meta string, comment []byte) (int64, error) {

	log.Info("save secret")

	scrHash, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate secret hash", err)

		return 0, fmt.Errorf("error generate secret: %w", err)
	}
	comHash, err := bcrypt.GenerateFromPassword([]byte(comment), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate comment hash", err)

		return 0, fmt.Errorf("error generate comment: %w", err)
	}

	id, err := s.scrSave.SaveSecret(ctx, userID, scrHash, meta, comHash)
	if err != nil {
		log.Error("failed to save secret", err)

		return 0, fmt.Errorf("error save secret: %w", err)
	}

	return id, nil
}

func (s *Secret) RegisterNewUser(ctx context.Context, email string, pass string) (int64, error) {

	return 0, nil
}

// Login функция для логина пользователя
func (s *Secret) Login(ctx context.Context, email string, password string) (token string, err error) {

	return "", nil
}
