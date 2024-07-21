// Package cmd cli обновления секретов банковских карт
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

// OldCard для старого секрета карты
var OldCard struct {
	Card string
	Data string
	Cvc  string
}

// upcardCmd represents the upcard command
var upcardCmd = &cobra.Command{
	Use:   "upcard",
	Short: "update your bank card",
	Long: `update your bank card use: client upcard and flags -n old number -d old data -c old cvc 
	-k new number -t new data -s new cvc`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Начат процесс обновления банковской карты")
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
		secretold := OldCard.Card + " : " + OldCard.Data + " : " + OldCard.Cvc
		secretnew := NewCard.Card + " : " + NewCard.Data + " : " + NewCard.Cvc
		_, err = c.UpdateSecret(ctxjwt, &pb.UpdateSecretRequest{Secret: secretold, SecretNew: secretnew})
		if err != nil {
			fmt.Println("Не удалось обновить секрет, пожалуйста, попробуйте еще раз")
			return
		}
		fmt.Println("Секрет успешно обновлен")
	},
}

func init() {
	rootCmd.AddCommand(upcardCmd)
	upcardCmd.Flags().StringVarP(&OldCard.Card, "old number", "n", "", "old secret card")
	upcardCmd.Flags().StringVarP(&OldCard.Data, "old data", "d", "", "old secret data")
	upcardCmd.Flags().StringVarP(&OldCard.Cvc, "old cvc", "c", "", "old secret cvc")
	upcardCmd.Flags().StringVarP(&NewCard.Card, "new number", "k", "", "new secret card")
	upcardCmd.Flags().StringVarP(&NewCard.Data, "new data", "t", "", "new secret data")
	upcardCmd.Flags().StringVarP(&NewCard.Cvc, "new cvc", "s", "", "new secret cvc")
	upcardCmd.MarkFlagRequired("old number")
	upcardCmd.MarkFlagRequired("old data")
	upcardCmd.MarkFlagRequired("old cvc")
	upcardCmd.MarkFlagRequired("new number")
	upcardCmd.MarkFlagRequired("new data")
	upcardCmd.MarkFlagRequired("new cvc")
}
