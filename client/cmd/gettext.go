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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
)

// gettextCmd represents the gettext command
var gettextCmd = &cobra.Command{
	Use:   "gettext",
	Short: "get your secret text",
	Long:  `get your secret text use: client gettext`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Начат процесс получения текстовых секретов")
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
		res, err := c.GetSecret(ctxjwt, &pb.GetSecretRequest{Meta: TypeText})
		if err != nil {
			fmt.Println("Не удалось получить секреты, пожалуйста, попробуйте еще раз")
			return
		}
		fmt.Println("Ваши текстовые секреты:")
		for _, i := range res.Secret {
			if i.Comment != "" {
				fmt.Printf("secret text: %s, comment: %s \n", i.Secret, i.Comment)
			} else {
				fmt.Printf("secret text: %s, comment: нет комментария \n", i.Secret)
			}
		}
		fmt.Println("Текстовые секреты успешно получены")
	},
}

func init() {
	rootCmd.AddCommand(gettextCmd)
}
