// Package cmd cli получение бинарных секретов
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

// getbinCmd represents the getbin command
var getbinCmd = &cobra.Command{
	Use:   "getbin",
	Short: "get your secret binary",
	Long:  `get your secret binary use: client getbin`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Начат процесс получения секретов бинарных данных")
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
		res, err := c.GetSecret(ctxjwt, &pb.GetSecretRequest{UserId: id, Meta: TypeBin})
		if err != nil {
			fmt.Println("Не удалось получить секреты, пожалуйста, попробуйте еще раз")
			return
		}
		fmt.Println("Ваши текстовые секреты:")
		for _, i := range res.Secret {
			if i.Comment != "" {
				fmt.Printf("secret bin: %s, comment: %s \n", i.Secret, i.Comment)
			} else {
				fmt.Printf("secret bin: %s, comment: нет комментария \n", i.Secret)
			}
		}
		fmt.Println("Секреты бинарных данных секреты успешно получены")
	},
}

func init() {
	rootCmd.AddCommand(getbinCmd)
}
