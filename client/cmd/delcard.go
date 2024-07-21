// Package cmd cli удаления банковских карт
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

// delcardCmd represents the delcard command
var delcardCmd = &cobra.Command{
	Use:   "delcard",
	Short: "delete your login/password secret",
	Long:  `delete your secret text use: client delcard and flag -k old nubmer -t old data -s old cvc`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Начат процесс удаления банковской карты")
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
		secret := OldCard.Card + " : " + OldCard.Data + " : " + OldCard.Cvc
		_, err = c.DeleteSecret(ctxjwt, &pb.DeleteSecretRequest{UserId: id, Secret: secret})
		if err != nil {
			fmt.Println("Не удалось удалить секрет, пожалуйста, попробуйте еще раз")
			return
		}
		fmt.Println("Секрет успешно удален")
	},
}

func init() {
	rootCmd.AddCommand(delcardCmd)
	delcardCmd.Flags().StringVarP(&OldCard.Card, "old number", "k", "", "old secret card to delete")
	delcardCmd.Flags().StringVarP(&OldCard.Data, "old data", "t", "", "old secret data to delete")
	delcardCmd.Flags().StringVarP(&OldCard.Cvc, "old cvc", "s", "", "old secret cvc to delete")
	delcardCmd.MarkFlagRequired("old number")
	delcardCmd.MarkFlagRequired("old data")
	delcardCmd.MarkFlagRequired("old cvc")
}
