// Package cmd cli добавление секретов логин/пароль
package cmd

import (
	"context"
	"fmt"
	"os/user"
	"time"

	"github.com/Andromaril/Gopher-and-secrets/client/internal/grpc"
	"github.com/Andromaril/Gopher-and-secrets/client/internal/local"
	pb "github.com/Andromaril/Gopher-and-secrets/server/proto"
	"github.com/Andromaril/Gopher-and-secrets/server/secret"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
)

// NewLogin для секрета
var NewLogin struct {
	Login    string
	Password string
}

// TypeLogin константа для типа секретов логин/пароль
const TypeLogin = "login/password"

// addloginCmd represents the addlogin command
var addloginCmd = &cobra.Command{
	Use:   "addlogin",
	Short: "add login/password secret",
	Long:  `add login/password secret, use: client addlogin and flags -l login -p password -c comment`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Начат процесс добавления секрета формата логин/пароль")
		user, err := user.Current()
		if err != nil {
			log.Fatalln(err)
		}
		jwt, ok := local.User[user.Username]
		if !ok {
			fmt.Println("Вы не авторизированы, залогиньтесь или зарегистрируйтесь")
			return
		}
		id, err := secret.DecodeToken(jwt)
		if err != nil {
			log.Fatalln(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		c, err := grpc.Init()
		if err != nil {
			log.Fatalln(err)
			return
		}
		ctxjwt := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+jwt)
		secret := NewLogin.Login + " : " + NewLogin.Password
		_, err = c.AddSecret(ctxjwt, &pb.AddSecretRequest{UserId: id, Secret: secret, Meta: TypeLogin, Comment: NewSecret.Comment})
		if err != nil {
			fmt.Println("Не удалось добавить секрет, пожалуйста, попробуйте еще раз")
			fmt.Println(err)
			return
		}
		fmt.Printf("Секрет формата логин/пароль успешно сохранен")
	},
}

func init() {
	rootCmd.AddCommand(addloginCmd)
	addloginCmd.Flags().StringVarP(&NewLogin.Login, "login", "l", "", "new secret login")
	addloginCmd.Flags().StringVarP(&NewLogin.Password, "password", "p", "", "new secret password")
	addloginCmd.Flags().StringVarP(&NewSecret.Comment, "comment", "c", "", "new comment, optional")
	addloginCmd.MarkFlagRequired("login")
	addloginCmd.MarkFlagRequired("password")
}
