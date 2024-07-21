// Package cmd cli добавления текстовых секретов
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
	"google.golang.org/grpc/metadata"

	"github.com/spf13/cobra"
)

// NewSecret переменная запроса на добавление секрета
var (
	NewSecret pb.AddSecretRequest
)

// TypeText константа для текстовых секретов
const TypeText = "secret text"

// addtextCmd represents the addtext command
var addtextCmd = &cobra.Command{
	Use:   "addtext",
	Short: "add text secret",
	Long:  `add text secret, use: client addtext and flags -s secret -c comment`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Начат процесс добавления секрета")
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
		_, err = c.AddSecret(ctxjwt, &pb.AddSecretRequest{Secret: NewSecret.Secret, Meta: TypeText, Comment: NewSecret.Comment})
		if err != nil {
			fmt.Println("Не удалось добавить секрет, пожалуйста, попробуйте еще раз")
			fmt.Println(err)
			return
		}
		fmt.Printf("Текстовый секрет успешно сохранен")
	},
}

func init() {
	rootCmd.AddCommand(addtextCmd)
	addtextCmd.Flags().StringVarP(&NewSecret.Secret, "secret", "s", "", "new secret")
	addtextCmd.Flags().StringVarP(&NewSecret.Comment, "comment", "c", "", "new comment, optional")
	addtextCmd.MarkFlagRequired("secret")
}
