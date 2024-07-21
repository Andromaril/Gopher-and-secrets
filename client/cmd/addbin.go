// Package cmd cli добавления бинарных секретов
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

// NewBinSecret для секрета в бинарном формате
var (
	NewBinSecret []byte
)

// TypeBin константа для бинарного секрета
const TypeBin = "secret binary"

// addbinCmd represents the addbin command
var addbinCmd = &cobra.Command{
	Use:   "addbin",
	Short: "add bin secret",
	Long:  `add bin secret, use: client addbin and flags -s secret -c comment`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Начат процесс добавдения секрета в бинарном формате")
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
		_, err = c.AddSecret(ctxjwt, &pb.AddSecretRequest{UserId: id, Secret: NewSecret.Secret, Meta: TypeBin, Comment: NewSecret.Comment})
		if err != nil {
			fmt.Println("Не удалось добавить секрет, пожалуйста, попробуйте еще раз")
			fmt.Println(err)
			return
		}
		fmt.Printf("Бинарный секрет успешно сохранен")
	},
}

func init() {
	rootCmd.AddCommand(addbinCmd)
	addbinCmd.Flags().StringVarP(&NewSecret.Secret, "secret", "s", "", "Binary data to save.")
	addbinCmd.Flags().StringVarP(&NewSecret.Comment, "comment", "c", "", "new comment, optional")
	addbinCmd.MarkFlagRequired("secret")
}
