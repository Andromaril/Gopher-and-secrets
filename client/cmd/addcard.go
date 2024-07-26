// Package cmd cli добавление секретов банковских карт
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

// NewCard структура для секрета карты
var NewCard struct {
	Card string
	Data string
	Cvc  string
}

// TypeCard константа для секрета карты
const TypeCard = "card"

// addcardCmd represents the addcard command
var addcardCmd = &cobra.Command{
	Use:   "addcard",
	Short: "add bank card secret",
	Long:  `add bank card secret, use: client addcard and flags -n number -d data -v cvc`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Начат процесс добавления секрета банковской карты")
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
		secret := NewCard.Card + " : " + NewCard.Data + " : " + NewCard.Cvc
		_, err = c.AddSecret(ctxjwt, &pb.AddSecretRequest{Secret: secret, Meta: TypeCard, Comment: NewSecret.Comment})
		if err != nil {
			fmt.Println("Не удалось добавить секрет, пожалуйста, попробуйте еще раз")
			fmt.Println(err)
			return
		}
		fmt.Printf("Секрет банковской карты успешно сохранен")
	},
}

func init() {
	rootCmd.AddCommand(addcardCmd)
	addcardCmd.Flags().StringVarP(&NewCard.Card, "card", "n", "", "new secret card")
	addcardCmd.Flags().StringVarP(&NewCard.Data, "data", "d", "", "new secret data")
	addcardCmd.Flags().StringVarP(&NewCard.Cvc, "cvc", "v", "", "new secret cvc")
	addcardCmd.Flags().StringVarP(&NewSecret.Comment, "comment", "c", "", "new comment, optional")
	addcardCmd.MarkFlagRequired("card")
	addcardCmd.MarkFlagRequired("data")
	addcardCmd.MarkFlagRequired("cvc")

}
