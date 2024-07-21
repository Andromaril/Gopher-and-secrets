// Package cmd cli удаления секретов бинарных данных
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

// delbinCmd represents the delbin command
var delbinCmd = &cobra.Command{
	Use:   "delbin",
	Short: "delete your secret binary",
	Long:  `delete your secret binary use: client delbin and flag -s secret`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Начат процесс удаления секретов бинарных данных")
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
		_, err = c.DeleteSecret(ctxjwt, &pb.DeleteSecretRequest{UserId: id, Secret: DeleteSecret.Secret})
		if err != nil {
			fmt.Println("Не удалось удалить секрет, пожалуйста, попробуйте еще раз")
			return
		}
		fmt.Println("Секрет успешно удален")
	},
}

func init() {
	rootCmd.AddCommand(delbinCmd)
	delbinCmd.Flags().StringVarP(&DeleteSecret.Secret, "secret", "s", "", "delete secret")
	delbinCmd.MarkFlagRequired("secret")
}
