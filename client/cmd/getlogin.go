// Package cmd cli получение секретов формата логин/пароль
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

// getloginCmd represents the getlogin command
var getloginCmd = &cobra.Command{
	Use:   "getlogin",
	Short: "get your login/password",
	Long:  `get your login/password use: client getlogin`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Начат процесс получения секрета формата логин/пароль")
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
		res, err := c.GetSecret(ctxjwt, &pb.GetSecretRequest{UserId: id, Meta: TypeLogin})
		if err != nil {
			fmt.Println("Не удалось получить секреты, пожалуйста, попробуйте еще раз")
			return
		}
		fmt.Println("Ваши секреты формата логин/пароль")
		for _, i := range res.Secret {
			if i.Comment != "" {
				fmt.Printf("secret login/password: %s, comment: %s \n", i.Secret, i.Comment)
			} else {
				fmt.Printf("secret login/password: %s, comment: нет комментария \n", i.Secret)
			}
		}
		fmt.Println("Секреты формата логин/пароль успешно получены")
	},
}

func init() {
	rootCmd.AddCommand(getloginCmd)
}
