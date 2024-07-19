// Package construct содержит интерфейсы для определения взаимодействия с бд
package construct

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Andromaril/Gopher-and-secrets/server/internal/model"
	"github.com/Andromaril/Gopher-and-secrets/server/internal/storage"
	"github.com/Andromaril/Gopher-and-secrets/server/secret"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// Auth структура для регистрации пользователя
type Auth struct {
	usrSave  UserSave
	usrGet   UserGet
	tokenTTL time.Duration
}

// переменная ошибки невалидных данных
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// UserSave для регистрации(сохранения пользователя)
type UserSave interface {
	SaveUser(ctx context.Context, email string, passwordHash []byte) (uid int64, err error)
}

// UserGet для получение(логин) пользователя
type UserGet interface {
	GetUser(ctx context.Context, email string) (model.User, error)
}

// New для создания экземпляра структуры Auth
func New(userSave UserSave, userGet UserGet, tokenTTL time.Duration) *Auth {
	return &Auth{
		usrSave:  userSave,
		usrGet:   userGet,
		tokenTTL: tokenTTL,
	}
}

// RegisterNewUser функция для регистрации нового пользователя
func (a *Auth) RegisterNewUser(ctx context.Context, email string, pass string) (int64, error) {

	log.Info("register user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", err)

		return 0, fmt.Errorf("error generate password: %w", err)
	}

	id, err := a.usrSave.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user", err)

		return 0, fmt.Errorf("error save user: %w", err)
	}

	return id, nil
}

// Login функция для логина пользователя
func (a *Auth) Login(ctx context.Context, email string, password string) (token string, err error) {
	log.Info("login")
	user, err := a.usrGet.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("user not found", err)

			return "", fmt.Errorf("error get user: %w", ErrInvalidCredentials)
		}

		log.Error("failed to get user", err)
		return "", fmt.Errorf("error get user: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		log.Info(ErrInvalidCredentials, err)

		return "", fmt.Errorf("%s: %w", ErrInvalidCredentials, err)
	}
	log.Info("user logged in successfully")

	jwt, err := secret.NewToken(user, a.tokenTTL)
	if err != nil {
		log.Error("failed to generate token", err)

		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return jwt, nil
}
