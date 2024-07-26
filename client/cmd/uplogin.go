// Package cmd cli обновления секретов формата логин/пароль
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

// OldLogin для обновляемого секрета
var OldLogin struct {
	Login    string
	Password string
}

// uploginCmd represents the uplogin command
var uploginCmd = &cobra.Command{
	Use:   "uplogin",
	Short: "update your login/password secret",
	Long: `update your login/password secret use: client uplogin and flags -l old login -p old password 
	-u new login -n new password`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Начат процесс обновления секрета формата логин/пароль")
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
		secretold := OldLogin.Login + " : " + OldLogin.Password
		secretnew := NewLogin.Login + " : " + NewLogin.Password
		_, err = c.UpdateSecret(ctxjwt, &pb.UpdateSecretRequest{Secret: secretold, SecretNew: secretnew})
		if err != nil {
			fmt.Println("Не удалось обновить секрет, пожалуйста, попробуйте еще раз")
			return
		}
		fmt.Println("Секрет успешно обновлен")
	},
}

func init() {
	rootCmd.AddCommand(uploginCmd)
	uploginCmd.Flags().StringVarP(&OldLogin.Login, "old login", "l", "", "old secret login")
	uploginCmd.Flags().StringVarP(&OldLogin.Password, "old password", "p", "", "old secret password")
	uploginCmd.Flags().StringVarP(&NewLogin.Login, "new login", "u", "", "new secret login")
	uploginCmd.Flags().StringVarP(&NewLogin.Password, "new password", "n", "", "new secret password")
	uploginCmd.MarkFlagRequired("old login")
	uploginCmd.MarkFlagRequired("old password")
	uploginCmd.MarkFlagRequired("new login")
	uploginCmd.MarkFlagRequired("new password")
}
