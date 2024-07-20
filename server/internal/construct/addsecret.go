package construct

import (
	"context"
	"errors"
	"fmt"

	"github.com/Andromaril/Gopher-and-secrets/server/internal/model"
	"github.com/Andromaril/Gopher-and-secrets/server/internal/storage"
	"github.com/Andromaril/Gopher-and-secrets/server/secret"
	log "github.com/sirupsen/logrus"
)

var super = "supersupersecret"
var supercom = "supersupercom"

// Secret структура для секретов
type Secret struct {
	scrSave   SecretSave
	scrGet    SecretGet
	scrUpdate SecretUpdate
	scrDelete SecretDelete
}

// NewSecret создание нового экземпляра Secret
func NewSecret(scrSave SecretSave, scrGet SecretGet, scrUpdate SecretUpdate, scrDelete SecretDelete) *Secret {
	return &Secret{
		scrSave:   scrSave,
		scrGet:    scrGet,
		scrUpdate: scrUpdate,
		scrDelete: scrDelete,
	}
}

// SecretSave для сохранения секретов
type SecretSave interface {
	SaveSecret(ctx context.Context, userID int64, secret string, meta string, comment string) (uid int64, err error)
}

// SecretGet для получения секретов
type SecretGet interface {
	GetSecret(ctx context.Context, userID int64, meta string) ([]model.Secret, error)
}

// SecretUpdate для получения секретов
type SecretUpdate interface {
	UpdateSecret(ctx context.Context, userID int64, secret string, secretnew string) error
}

// SecretDelete для удаления секретов
type SecretDelete interface {
	DeleteSecret(ctx context.Context, userID int64, sec string) error
}

// SaveSecret для сохранения секретов
func (s *Secret) SaveSecret(ctx context.Context, userID int64, sec string, meta string, comment string) (int64, error) {

	log.Info("save secret")

	scr, err := secret.Encrypt(sec, secret.MySecret)
	if err != nil {
		log.Error("failed to generate secret hash", err)

		return 0, fmt.Errorf("error generate secret: %w", err)
	}
	com, err := secret.Encrypt(comment, secret.MySecret)

	if err != nil {
		log.Error("failed to generate comment hash", err)

		return 0, fmt.Errorf("error generate comment: %w", err)
	}

	id, err := s.scrSave.SaveSecret(ctx, userID, scr, meta, com)
	if err != nil {
		log.Error("failed to save secret", err)

		return 0, fmt.Errorf("error save secret: %w", err)
	}

	return id, nil
}

// GetNewSecret для получения секрета
func (s *Secret) GetNewSecret(ctx context.Context, userID int64, meta string) ([]model.Secret, error) {

	log.Info("get secret")
	sec := make([]model.Secret, 0)
	scr, err := s.scrGet.GetSecret(ctx, userID, meta)
	if err != nil {
		if errors.Is(err, storage.ErrSecretNotFound) {
			log.Error("secret not found ", err)

			return scr, fmt.Errorf("error get secret: %w", err)
		}
		log.Error("failed to get secret ", err)
		return scr, fmt.Errorf("error get secret: %w", err)
	}
	for _, value := range scr {
		s, err := secret.Decrypt(value.Secret, secret.MySecret)
		if err != nil {
			return sec, fmt.Errorf("error decrypt secret: %w", err)
		}
		c, err := secret.Decrypt(value.Comment, secret.MySecret)
		if err != nil {
			return sec, fmt.Errorf("error decrypt secret: %w", err)
		}
		sec = append(sec, model.Secret{Secret: s, Comment: c})

	}
	log.Info("get secret successfully")

	return sec, nil
}

// UpdateSecret для обновления секрета
func (s *Secret) UpdateSecret(ctx context.Context, userID int64, sec string, secretnew string) error {
	log.Info("update secret")
	scrold, err := secret.Encrypt(sec, secret.MySecret)
	if err != nil {
		log.Error("failed to generate secret hash", err)

		return fmt.Errorf("error generate secret: %w", err)
	}

	scr, err := secret.Encrypt(secretnew, secret.MySecret)
	if err != nil {
		log.Error("failed to generate secret hash", err)

		return fmt.Errorf("error generate secret: %w", err)
	}

	err = s.scrUpdate.UpdateSecret(ctx, userID, scrold, scr)
	if err != nil {
		log.Error("failed to update secret", err)

		return fmt.Errorf("error update secret: %w", err)
	}
	return nil
}

// DeleteSecret для удаления секрета
func (s *Secret) DeleteSecret(ctx context.Context, userID int64, sec string) error {
	log.Info("delete secret")

	scr, err := secret.Encrypt(sec, secret.MySecret)
	if err != nil {
		log.Error("failed to generate secret hash", err)

		return fmt.Errorf("error generate secret: %w", err)
	}

	err = s.scrDelete.DeleteSecret(ctx, userID, scr)
	if err != nil {
		log.Error("failed to delete secret", err)

		return fmt.Errorf("error delete secret: %w", err)
	}
	return nil
}
