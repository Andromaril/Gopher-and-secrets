// Package secret содержит функцию шифрования и дешифрования

package secret

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	t.Run("EncryptOk", func(t *testing.T) {
		s, err := Encrypt("test1", MySecret)
		assert.Equal(t, s, "HyVU91k=")
		assert.NoError(t, err)

	})
}

func TestDecrypt(t *testing.T) {
	t.Run("DeryptOk", func(t *testing.T) {
		s, err := Decrypt("HyVU91k=", MySecret)
		assert.Equal(t, s, "test1")
		assert.NoError(t, err)
	})
}
