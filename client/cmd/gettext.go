// Package cmd cli получение текстовых секретов
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

// gettextCmd represents the gettext command
var gettextCmd = &cobra.Command{
	Use:   "gettext",
	Short: "get your secret text",
	Long: `get your secret text use: client gettext`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gettext called")
		user, err := user.Current()
		if err != nil {
			log.Fatalln(err)
		}
		username, ok := local.User[user.Username]
		if !ok {
			fmt.Println("Вы не авторизированы, залогиньтесь или зарегистрируйтесь")
			return
		}
		log.Info(username)
		jwt, ok := local.User[user.Username]
		if !ok {
			fmt.Println("Вы не авторизированы, залогиньтесь или зарегистрируйтесь")
			return
		}
		id, err := secret.DecodeToken(jwt)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(id)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		c, err := grpc.Init()
		if err != nil {
			log.Fatalln(err)
			return
		}
		ctxjwt := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+jwt)
		res, err := c.GetSecret(ctxjwt, &pb.GetSecretRequest{UserId: id, Meta: "secret text"})
		if err != nil {
			fmt.Println("Не удалось получить секреты, пожалуйста, попробуйте еще раз")
			return
		}
		fmt.Println(res)
		fmt.Printf("Текстовые секреты успешно получены")
	},
}

func init() {
	rootCmd.AddCommand(gettextCmd)
}
