// Package cmd cli удаления секретов
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

// DeleteSecret структура запроса обновления секрета
var (
	DeleteSecret pb.DeleteSecretRequest
)

// deltextCmd represents the deltext command
var deltextCmd = &cobra.Command{
	Use:   "deltext",
	Short: "delete your secret text",
	Long:  `delete your secret text use: client deltext and flag -s secret`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Начат процесс удаления секрета")
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
		_, err = c.DeleteSecret(ctxjwt, &pb.DeleteSecretRequest{Secret: DeleteSecret.Secret})
		if err != nil {
			fmt.Println("Не удалось удалить секрет, пожалуйста, попробуйте еще раз")
			return
		}
		fmt.Println("Секрет успешно удален")
	},
}

func init() {
	rootCmd.AddCommand(deltextCmd)
	deltextCmd.Flags().StringVarP(&DeleteSecret.Secret, "secret", "s", "", "delete secret")
	deltextCmd.MarkFlagRequired("secret")
}
