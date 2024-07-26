// Package storage для взаимодействия с базой данных

package storage

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/Andromaril/Gopher-and-secrets/server/internal/model"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	url := "postgres://postgres:qwerty123@localhost:5432/server2"
	_, err := Init(url)
	assert.NoError(t, err)
}

func newMock() (Database, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal("Failed to create sql mock db")
	}
	return &Storage{DB: db}, mock
}

func newUser() *model.User {
	return &model.User{
		ID:           1,
		Login:        "test@mail.ru",
		PasswordHash: []byte("$2a$10$6LXF8d70LTX5Vn3VXTUCBOnmOjLL6dRTn8ha1L4uN9RZH/.eiYIru"),
	}
}

func newSecret() *model.Secret {
	return &model.Secret{
		ID:       1,
		SecretID: 1,
		Secret:   "test",
		Meta:     "test",
		Comment:  "test",
	}
}

func TestStorage_SaveUser(t *testing.T) {
	s, mock := newMock()
	user := newUser()

	t.Run("SaveUserOk", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO users").
			WithArgs(user.Login, user.PasswordHash).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(user.ID))

		userActual, err := s.SaveUser(context.Background(), user.Login, user.PasswordHash)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, userActual)
	})
	t.Run("UserExists", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO users").
			WithArgs(user.Login, user.PasswordHash).
			WillReturnError(sql.ErrNoRows)

		_, err := s.SaveUser(context.Background(), user.Login, user.PasswordHash)
		assert.Error(t, err)
	})
}

func TestGetUser(t *testing.T) {
	s, mock := newMock()
	user := newUser()

	t.Run("GetUserOk", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, login, password FROM users WHERE").
			WithArgs(user.Login).
			WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password"}).
				AddRow(user.ID, user.Login, user.PasswordHash))

		userActual, err := s.GetUser(context.Background(), user.Login)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, userActual.ID)
		assert.Equal(t, user.Login, userActual.Login)
		assert.Equal(t, user.PasswordHash, userActual.PasswordHash)

	})
	t.Run("UserNotFound", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, login, password FROM users WHERE").
			WithArgs(user.Login).
			WillReturnError(sql.ErrNoRows)

		_, err := s.GetUser(context.Background(), user.Login)
		assert.ErrorIs(t, err, ErrUserNotFound)
	})
}

func TestSaveSecret(t *testing.T) {
	s, mock := newMock()
	user := newUser()
	secret := newSecret()

	t.Run("SaveSecretOk", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO secrets").
			WithArgs(user.ID, secret.Secret, secret.Meta, secret.Comment).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(user.ID))

		userActual, err := s.SaveSecret(context.Background(), user.ID, secret.Secret, secret.Meta, secret.Comment)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, userActual)
	})
	t.Run("ErrorSaveSecret", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO secrets").
			WithArgs(user.ID, secret.Secret, secret.Meta, secret.Comment).
			WillReturnError(sql.ErrNoRows)

		_, err := s.SaveSecret(context.Background(), user.ID, secret.Secret, secret.Meta, secret.Comment)
		assert.Error(t, err)
	})
}

func TestGetSecret(t *testing.T) {
	s, mock := newMock()
	user := newUser()
	secret := newSecret()

	t.Run("GetSecretOk", func(t *testing.T) {
		mock.ExpectQuery("SELECT secret, comment FROM secrets WHERE").
			WithArgs(user.ID, secret.Meta).
			WillReturnRows(sqlmock.NewRows([]string{"secret", "comment"}).
				AddRow(secret.Secret, secret.Comment))

		secretActual, err := s.GetSecret(context.Background(), user.ID, secret.Meta)
		assert.NoError(t, err)
		for _, i := range secretActual {
			assert.Equal(t, secret.Secret, i.Secret)
			assert.Equal(t, secret.Comment, i.Comment)
		}

	})
	t.Run("SecretNotFound", func(t *testing.T) {
		mock.ExpectQuery("SELECT secret, comment FROM secrets WHERE").
			WithArgs(user.ID, secret.Meta).
			WillReturnError(sql.ErrNoRows)

		_, err := s.GetSecret(context.Background(), user.ID, secret.Meta)
		assert.ErrorIs(t, err, ErrSecretNotFound)
	})

}

func TestUpdateSecret(t *testing.T) {
	s, mock := newMock()
	user := newUser()
	secret := newSecret()
	newsecret := &model.Secret{
		ID:       1,
		SecretID: 1,
		Secret:   "test1",
		Meta:     "test1",
		Comment:  "test1",
	}

	t.Run("UpdateSecretOk", func(t *testing.T) {
		mock.ExpectExec("UPDATE secrets SET secret").
			WithArgs(newsecret.Secret, user.ID, secret.Secret).
			WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(nil)

		err := s.UpdateSecret(context.Background(), user.ID, secret.Secret, newsecret.Secret)
		assert.NoError(t, err)
	})
	t.Run("ErrorUpdateSecret", func(t *testing.T) {
		mock.ExpectExec("UPDATE secrets SET secret").
			WithArgs(newsecret.Secret, user.ID, secret.Secret).
			WillReturnError(fmt.Errorf("error insert"))

		err := s.UpdateSecret(context.Background(), user.ID, secret.Secret, newsecret.Secret)
		assert.Error(t, err)
	})
}

func TestDeleteSecret(t *testing.T) {
	s, mock := newMock()
	user := newUser()
	secret := newSecret()

	t.Run("DeleteSecretOk", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM secrets WHERE").
			WithArgs(secret.Secret, user.ID).
			WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(nil)

		err := s.DeleteSecret(context.Background(), user.ID, secret.Secret)
		assert.NoError(t, err)
	})
	t.Run("ErrorDeleteSecret", func(t *testing.T) {
		mock.ExpectExec("UPDATE secrets SET secret").
			WithArgs(secret.Secret, user.ID).
			WillReturnError(fmt.Errorf("error delete"))

		err := s.DeleteSecret(context.Background(), user.ID, secret.Secret)
		assert.Error(t, err)
	})
}

func TestStorage_GetAll(t *testing.T) {
	s, mock := newMock()
	user := newUser()
	secret := newSecret()

	t.Run("GetSecretAllOk", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, secret, comment, meta FROM secrets WHERE").
			WithArgs(user.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "secret", "comment", "meta"}).
				AddRow(user.ID, secret.Secret, secret.Comment, secret.Meta))

		secretActual, err := s.GetAll(context.Background(), user.ID)
		assert.NoError(t, err)
		for _, i := range secretActual {
			assert.Equal(t, secret.SecretID, i.SecretID)
			assert.Equal(t, secret.Secret, i.Secret)
			assert.Equal(t, secret.Comment, i.Comment)
			assert.Equal(t, secret.Meta, i.Meta)
		}

	})
	t.Run("SecretNotFound", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, secret, comment, meta FROM secrets WHERE").
			WithArgs(user.ID).
			WillReturnError(sql.ErrNoRows)

		_, err := s.GetAll(context.Background(), user.ID)
		assert.ErrorIs(t, err, ErrSecretNotFound)
	})

}
