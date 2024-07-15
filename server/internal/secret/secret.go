// Package secret содержит функцию генерации токена
package secret

import (
	"time"

	"github.com/Andromaril/Gopher-and-secrets/server/internal/model"
	"github.com/golang-jwt/jwt"
)

var tokenEncodeString string = "supersecrettoken"

// NewToken функция для генерирования jwt token
func NewToken(user model.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(tokenEncodeString))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
