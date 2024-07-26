// Package cmd cli удаления секретов формата логин/пароль
package cmd

import (
	"context"
	"fmt"
	"os/user"
	"time"

	"github.com/Andromaril/Gopher-and-secrets/client/internal/grpc"
	"github.com/Andromaril/Gopher-and-secrets/client/internal/local"
	pb "github.com/Andromaril/Gopher-and-secrets/server/proto"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
)

// delloginCmd represents the dellogin command
var delloginCmd = &cobra.Command{
	Use:   "dellogin",
	Short: "delete your login/password secret",
	Long:  `delete your secret text use: client deltext and flag -l login -p password`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Начат процесс удаления секрета формата логин.пароль")
		user, err := user.Current()
		if err != nil {
			log.Fatalln(err)
		}
		jwt, ok := local.User[user.Username]
		if !ok {
			fmt.Println("Вы не авторизированы, залогиньтесь или зарегистрируйтесь")
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		c, err := grpc.Init()
		if err != nil {
			log.Fatalln(err)
			return
		}
		ctxjwt := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+jwt)
		secret := OldLogin.Login + " : " + OldLogin.Password
		_, err = c.DeleteSecret(ctxjwt, &pb.DeleteSecretRequest{Secret: secret})
		if err != nil {
			fmt.Println("Не удалось удалить секрет, пожалуйста, попробуйте еще раз")
			return
		}
		fmt.Println("Секрет успешно удален")
	},
}

func init() {
	rootCmd.AddCommand(delloginCmd)
	delloginCmd.Flags().StringVarP(&OldLogin.Login, "login", "l", "", "secret login to delete")
	delloginCmd.Flags().StringVarP(&OldLogin.Password, "password", "p", "", "secret password to delete")
	delloginCmd.MarkFlagRequired("login")
	delloginCmd.MarkFlagRequired("password")
}
