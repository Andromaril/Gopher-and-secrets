// Package secret содержит функцию генерации токена
package secret

import (
	"time"

	"github.com/Andromaril/Gopher-and-secrets/server/internal/model"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
)

var tokenEncodeString = []byte("supersecrettoken")

// MyClaims для claims токена
type MyClaims struct {
	jwt.RegisteredClaims
	UID   int64  `json:"uid"`   // id пользователя
	Login string `json:"login"` // логин полтзователя
	Exp   int64  `json:"exp"`   // дата истечения токена
}

// NewToken функция для генерирования jwt token
func NewToken(user model.User, duration time.Duration) (string, error) {
	//token := jwt.New(jwt.SigningMethodHS256)
	// claims := token.Claims.(jwt.MapClaims)
	// claims["uid"] = user.ID
	// claims["login"] = user.Login
	// claims["exp"] = time.Now().Add(duration).Unix()
	claims := MyClaims{
		RegisteredClaims: jwt.RegisteredClaims{},
		UID:              user.ID,
		Login:            user.Login,
		Exp:              time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(tokenEncodeString)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// DecodeToken декодирование токена
func DecodeToken(t string) (int64, error) {
	// keyFunc := func(t *jwt.Token) (interface{}, error) {
	// 	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
	// 		return nil, fmt.Errorf("Unexpected signature method: %v", t.Header["alg"])
	// 	}
	// 	return tokenEncodeString, nil
	// }
	// claims := &MyClaims{}
	// parsedToken, err := jwt.ParseWithClaims(string(tokenEncodeString), claims, keyFunc)
	// if err != nil {
	// 	log.Fatalf("Parse error: %v", err)
	// }

	// if !parsedToken.Valid {
	// 	log.Fatalf("Invalid token")
	// }
	token, err := jwt.ParseWithClaims(t, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return tokenEncodeString, nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	claims, ok := token.Claims.(*MyClaims)
	if !ok {
		log.Fatal("unknown claims type, cannot proceed")
	}
	return claims.UID, nil
}
