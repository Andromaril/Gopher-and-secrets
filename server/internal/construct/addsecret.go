package construct

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// Secret структура для секретов
type Secret struct {
	scrSave SecretSave
}

func NewSecret(scrSave SecretSave) *Secret {
	return &Secret{
		scrSave: scrSave,
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
