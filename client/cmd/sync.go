// Package cmd cli синхронизация локального хранилища и хранища сервера
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

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync all your secret",
	Long:  `sync your secret use: client sync`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Процесс синхронизации с базой сервиса начат")
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
		res, err := c.GetAll(ctxjwt, &pb.GetAllRequest{UserId: id})
		if err != nil {
			fmt.Println("Не удалось получить секреты, пожалуйста, попробуйте еще раз")
			return
		}
		for _, i := range res.Secret {
			local.Secret = append(local.Secret, local.TempSecret{
				SecretID: i.SecretId,
				Secret:   i.Secret,
				Meta:     i.Meta,
				Comment:  i.Comment,
			})
		}
		fmt.Println("Cекреты успешно получены")
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
