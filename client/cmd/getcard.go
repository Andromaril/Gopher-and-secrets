// Package cmd cli получение секретов банковских карт
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

// getcardCmd represents the getcard command
var getcardCmd = &cobra.Command{
	Use:   "getcard",
	Short: "get your bank card secret",
	Long:  `get your bank card use: client getcard`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Начат процесс получения секретов банковских карт")
		user, err := user.Current()
		if err != nil {
			log.Fatalln(err)
		}
		jwt, ok := local.User[user.Username]
		if !ok {
			fmt.Println("Вы не авторизированы, залогиньтесь или зарегистрируйтесь")
			return
		}
		//, err := secret.DecodeToken(jwt)
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
		res, err := c.GetSecret(ctxjwt, &pb.GetSecretRequest{Meta: TypeCard})
		if err != nil {
			fmt.Println("Не удалось получить секреты, пожалуйста, попробуйте еще раз")
			return
		}
		fmt.Println("Ваши секреты банковских карт")
		for _, i := range res.Secret {
			if i.Comment != "" {
				fmt.Printf("secret bank card: %s, comment: %s \n", i.Secret, i.Comment)
			} else {
				fmt.Printf("secret bank card: %s, comment: нет комментария \n", i.Secret)
			}
		}
		fmt.Println("Секреты банковских карт успешно получены")
	},
}

func init() {
	rootCmd.AddCommand(getcardCmd)
}
