// Package handler содержит обработчики запросов grpc

package handler

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/Andromaril/Gopher-and-secrets/server/internal/construct"
	interceptors "github.com/Andromaril/Gopher-and-secrets/server/internal/interseptors"
	"github.com/Andromaril/Gopher-and-secrets/server/internal/model"
	"github.com/Andromaril/Gopher-and-secrets/server/internal/storage"
	"github.com/Andromaril/Gopher-and-secrets/server/internal/storage/mocks"
	pb "github.com/Andromaril/Gopher-and-secrets/server/proto"
	"github.com/Andromaril/Gopher-and-secrets/server/secret"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func checkErrorStatus(t *testing.T, err error, code codes.Code) {
	require.Error(t, err)

	errStatus, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, code, errStatus.Code())
}

//ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)

func StartTest(t *testing.T) (*mocks.MockDatabase, pb.AuthClient) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer(grpc.UnaryInterceptor(interceptors.AuthCheck))
	a2 := construct.NewSecret(mockDB, mockDB, mockDB, mockDB, mockDB)
	a1 := construct.New(mockDB, mockDB, 1000)

	pb.RegisterAuthServer(s, &serverAPI{auth: a1, secret: a2})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Test server exited with error: %v", err)
		}
	}()
	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	client := pb.NewAuthClient(conn)
	return mockDB, client
}

func Test_serverAPI_Register(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDB := mocks.NewMockDatabase(ctrl)

		lis = bufconn.Listen(bufSize)
		s := grpc.NewServer(grpc.UnaryInterceptor(interceptors.AuthCheck))
		a2 := construct.NewSecret(mockDB, mockDB, mockDB, mockDB, mockDB)
		a1 := construct.New(mockDB, mockDB, 1000)

		Register(s, a1, a2)
	})
}

func TestLogin(t *testing.T) {
	var err error
	mockDB, client := StartTest(t)
	ctx := context.Background()
	user := &model.User{
		ID:           0,
		Login:        "test@mail.ru",
		PasswordHash: []byte("hash"),
	}
	password := "password"
	t.Run("EmptyEmail", func(t *testing.T) {
		request := &pb.LoginRequest{Login: "", Password: "password"}
		_, err = client.Login(context.Background(), request)
		checkErrorStatus(t, err, codes.InvalidArgument)
	})
	t.Run("EmptyPasssword", func(t *testing.T) {
		request := &pb.LoginRequest{Login: "login", Password: ""}
		_, err = client.Login(context.Background(), request)
		checkErrorStatus(t, err, codes.InvalidArgument)
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		request := &pb.LoginRequest{Login: user.Login, Password: password}
		mockDB.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(model.User{}, fmt.Errorf("error in scan from user select: %w", storage.ErrUserNotFound))
		_, err := client.Login(ctx, request)
		checkErrorStatus(t, err, codes.InvalidArgument)
	})

	t.Run("InternalError", func(t *testing.T) {
		request := &pb.LoginRequest{Login: user.Login, Password: password}
		mockDB.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(model.User{Login: user.Login, PasswordHash: user.PasswordHash}, nil)
		_, err := client.Login(ctx, request)
		checkErrorStatus(t, err, codes.Internal)
	})

	t.Run("OK", func(t *testing.T) {
		user := &model.User{
			ID:           0,
			Login:        "test@mail.ru",
			PasswordHash: []byte("$2a$10$6LXF8d70LTX5Vn3VXTUCBOnmOjLL6dRTn8ha1L4uN9RZH/.eiYIru"),
		}
		request := &pb.LoginRequest{Login: user.Login, Password: password}
		mockDB.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(model.User{Login: user.Login, PasswordHash: user.PasswordHash}, nil)
		_, err := client.Login(ctx, request)
		assert.NoError(t, err)
	})
}

func TestRegister(t *testing.T) {
	var err error
	mockDB, client := StartTest(t)
	ctx := context.Background()
	password := "password"
	t.Run("EmptyEmail", func(t *testing.T) {
		request := &pb.RegisterRequest{Login: "", Password: "password"}
		_, err = client.Register(context.Background(), request)
		checkErrorStatus(t, err, codes.InvalidArgument)
	})
	t.Run("EmptyPasssword", func(t *testing.T) {
		request := &pb.RegisterRequest{Login: "login", Password: ""}
		_, err = client.Register(context.Background(), request)
		checkErrorStatus(t, err, codes.InvalidArgument)
	})

	t.Run("InternalError", func(t *testing.T) {
		user := &model.User{
			ID:           0,
			Login:        "test@mail.ru",
			PasswordHash: []byte("hash"),
		}
		request := &pb.RegisterRequest{Login: user.Login, Password: password}
		mockDB.EXPECT().SaveUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(0), fmt.Errorf("error generate password"))
		_, err := client.Register(ctx, request)
		checkErrorStatus(t, err, codes.Internal)
	})

	t.Run("OK", func(t *testing.T) {
		user := &model.User{
			ID:           0,
			Login:        "test@mail.ru",
			PasswordHash: []byte("$2a$10$6LXF8d70LTX5Vn3VXTUCBOnmOjLL6dRTn8ha1L4uN9RZH/.eiYIru"),
		}
		request := &pb.RegisterRequest{Login: user.Login, Password: password}
		mockDB.EXPECT().SaveUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(6), nil)
		_, err := client.Register(ctx, request)
		assert.NoError(t, err)
	})
}

func TestAddSecret(t *testing.T) {
	mockDB, client := StartTest(t)
	ctx := context.Background()
	user := model.User{
		ID: 6,
	}
	jwt, err := secret.NewToken(user, 1000)
	ctxjwt := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+jwt)
	secret1 := &model.Secret{
		ID:      6,
		Secret:  "test1",
		Meta:    "test",
		Comment: "test1",
	}

	t.Run("InternalError", func(t *testing.T) {
		request := &pb.AddSecretRequest{
			Secret:  secret1.Secret,
			Meta:    secret1.Meta,
			Comment: secret1.Comment,
		}
		mockDB.EXPECT().SaveSecret(gomock.Any(), secret1.ID, "HyVU91k=", "test", "HyVU91k=").Return(int64(0), fmt.Errorf("error insert in save secrets"))
		_, err = client.AddSecret(ctxjwt, request)
		checkErrorStatus(t, err, codes.Internal)
	})

	t.Run("OK", func(t *testing.T) {

		request := &pb.AddSecretRequest{
			Secret:  secret1.Secret,
			Meta:    secret1.Meta,
			Comment: secret1.Comment,
		}
		mockDB.EXPECT().SaveSecret(gomock.Any(), secret1.ID, "HyVU91k=", "test", "HyVU91k=").Return(int64(1), nil)
		_, err = client.AddSecret(ctxjwt, request)
		assert.NoError(t, err)
	})
}

func TestGetSecret(t *testing.T) {
	mockDB, client := StartTest(t)
	ctx := context.Background()
	user := model.User{
		ID: 6,
	}
	jwt, err := secret.NewToken(user, 1000)
	ctxjwt := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+jwt)
	secret1 := &model.Secret{
		ID:      6,
		Secret:  "test1",
		Meta:    "test",
		Comment: "test1",
	}

	secret2 := make([]model.Secret, 0)

	t.Run("InternalError", func(t *testing.T) {
		request := &pb.GetSecretRequest{
			Meta: secret1.Meta,
		}
		mockDB.EXPECT().GetSecret(gomock.Any(), secret1.ID, "test").Return(secret2, fmt.Errorf("failed to get secret"))
		_, err = client.GetSecret(ctxjwt, request)
		checkErrorStatus(t, err, codes.Internal)
	})

	t.Run("NotFound", func(t *testing.T) {
		request := &pb.GetSecretRequest{
			Meta: secret1.Meta,
		}
		mockDB.EXPECT().GetSecret(gomock.Any(), secret1.ID, "test").Return(secret2, fmt.Errorf("error in scan from secret select: %w", storage.ErrSecretNotFound))
		_, err = client.GetSecret(ctxjwt, request)
		checkErrorStatus(t, err, codes.NotFound)
	})

	t.Run("OK", func(t *testing.T) {
		request := &pb.GetSecretRequest{
			Meta: secret1.Meta,
		}
		mockDB.EXPECT().GetSecret(gomock.Any(), secret1.ID, "test").Return(secret2, nil)
		_, err = client.GetSecret(ctxjwt, request)
		assert.NoError(t, err)
	})
}

func TestUpdateSecret(t *testing.T) {
	mockDB, client := StartTest(t)
	ctx := context.Background()
	user := model.User{
		ID: 6,
	}
	jwt, err := secret.NewToken(user, 1000)
	ctxjwt := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+jwt)
	secret1 := &model.Secret{
		ID:      6,
		Secret:  "test1",
		Meta:    "test",
		Comment: "test1",
	}

	t.Run("InternalError", func(t *testing.T) {
		request := &pb.UpdateSecretRequest{
			Secret:    secret1.Secret,
			SecretNew: "test3",
		}
		mockDB.EXPECT().UpdateSecret(gomock.Any(), secret1.ID, "HyVU91k=", "HyVU91s=").Return(fmt.Errorf("error insert"))
		_, err = client.UpdateSecret(ctxjwt, request)
		checkErrorStatus(t, err, codes.Internal)
	})

	t.Run("OK", func(t *testing.T) {

		request := &pb.UpdateSecretRequest{
			Secret:    secret1.Secret,
			SecretNew: "test3",
		}
		mockDB.EXPECT().UpdateSecret(gomock.Any(), secret1.ID, "HyVU91k=", "HyVU91s=").Return(nil)
		_, err = client.UpdateSecret(ctxjwt, request)
		assert.NoError(t, err)
	})
}

func TestDeleteSecret(t *testing.T) {
	mockDB, client := StartTest(t)
	ctx := context.Background()
	user := model.User{
		ID: 6,
	}
	jwt, err := secret.NewToken(user, 1000)
	ctxjwt := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+jwt)
	secret1 := &model.Secret{
		ID:      6,
		Secret:  "test1",
		Meta:    "test",
		Comment: "test1",
	}

	t.Run("InternalError", func(t *testing.T) {
		request := &pb.DeleteSecretRequest{
			Secret: secret1.Secret,
		}
		mockDB.EXPECT().DeleteSecret(gomock.Any(), secret1.ID, "HyVU91k=").Return(fmt.Errorf("error delete"))
		_, err = client.DeleteSecret(ctxjwt, request)
		checkErrorStatus(t, err, codes.Internal)
	})

	t.Run("OK", func(t *testing.T) {

		request := &pb.DeleteSecretRequest{
			Secret: secret1.Secret,
		}
		mockDB.EXPECT().DeleteSecret(gomock.Any(), secret1.ID, "HyVU91k=").Return(nil)
		_, err = client.DeleteSecret(ctxjwt, request)
		assert.NoError(t, err)
	})
}

func TestGetAll(t *testing.T) {
	mockDB, client := StartTest(t)
	ctx := context.Background()
	user := model.User{
		ID: 6,
	}
	jwt, err := secret.NewToken(user, 1000)
	ctxjwt := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+jwt)
	secret1 := &model.Secret{
		ID:      6,
		Secret:  "test1",
		Meta:    "test",
		Comment: "test1",
	}

	secret2 := make([]model.Secret, 0)

	t.Run("UserIsRequired", func(t *testing.T) {
		request := &pb.GetAllRequest{
			UserId: 0,
		}
		_, err = client.GetAll(ctxjwt, request)
		checkErrorStatus(t, err, codes.InvalidArgument)
	})
	t.Run("InternalError", func(t *testing.T) {
		request := &pb.GetAllRequest{
			UserId: secret1.ID,
		}
		mockDB.EXPECT().GetAll(gomock.Any(), secret1.ID).Return(secret2, fmt.Errorf("failed to get secret"))
		_, err = client.GetAll(ctxjwt, request)
		checkErrorStatus(t, err, codes.Internal)
	})

	t.Run("NotFound", func(t *testing.T) {
		request := &pb.GetAllRequest{
			UserId: secret1.ID,
		}
		mockDB.EXPECT().GetAll(gomock.Any(), secret1.ID).Return(secret2, fmt.Errorf("error in scan from secret select: %w", storage.ErrSecretNotFound))
		_, err = client.GetAll(ctxjwt, request)
		checkErrorStatus(t, err, codes.NotFound)
	})

	t.Run("OK", func(t *testing.T) {
		request := &pb.GetAllRequest{
			UserId: secret1.ID,
		}
		mockDB.EXPECT().GetAll(gomock.Any(), secret1.ID).Return(secret2, nil)
		_, err = client.GetAll(ctxjwt, request)
		assert.NoError(t, err)
	})
}
