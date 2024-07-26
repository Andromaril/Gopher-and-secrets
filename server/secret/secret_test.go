// Package secret содержит функцию генерации токена

package secret

import (
	"testing"

	"github.com/Andromaril/Gopher-and-secrets/server/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNewToken(t *testing.T) {
	t.Run("NewTokenOk", func(t *testing.T) {
		user := model.User{
			ID:    1,
			Login: "test@mail.ru",
		}
		_, err := NewToken(user, 1000)
		assert.NoError(t, err)
	})
}

func TestDecodeToken(t *testing.T) {
	t.Run("DecodeTokenOk", func(t *testing.T) {
		id, err := DecodeToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjEsImxvZ2luIjoidGVzdEBtYWlsLnJ1IiwiZXhwIjoxNzIxOTI1MDI2fQ.7eEl-oSGGFBk888ivsj34zsC0lkn3HtFF9zpBORxJIM")
		assert.Equal(t, id, int64(1))
		assert.NoError(t, err)
	})
}
