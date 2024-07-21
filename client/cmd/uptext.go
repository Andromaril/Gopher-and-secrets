// Package cmd cli обновления секретов
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

// UpdateSecret структура запроса обновления секрета
var (
	UpdateSecret pb.UpdateSecretRequest
)

// uptextCmd represents the uptext command
var uptextCmd = &cobra.Command{
	Use:   "uptext",
	Short: "update your secret text",
	Long:  `update your secret text use: client uptext and flags -o old secret -n new secret`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Начат процесс обновления секрета")
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
		_, err = c.UpdateSecret(ctxjwt, &pb.UpdateSecretRequest{UserId: id, Secret: UpdateSecret.Secret, SecretNew: UpdateSecret.SecretNew})
		if err != nil {
			fmt.Println("Не удалось обновить секрет, пожалуйста, попробуйте еще раз")
			return
		}
		fmt.Println("Секрет успешно обновлен")
	},
}

func init() {
	rootCmd.AddCommand(uptextCmd)
	uptextCmd.Flags().StringVarP(&UpdateSecret.Secret, "secret", "o", "", "old secret")
	uptextCmd.Flags().StringVarP(&UpdateSecret.SecretNew, "secret new", "n", "", "new secret")
	uptextCmd.MarkFlagRequired("secret")
	uptextCmd.MarkFlagRequired("secret new")
}
