/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
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
	"google.golang.org/grpc/metadata"

	"github.com/spf13/cobra"
)

var (
	NewSecret pb.AddSecretRequest
)

// addtextCmd represents the addtext command
var addtextCmd = &cobra.Command{
	Use:   "addtext",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("addtext called")
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
		_, err = c.AddSecret(ctxjwt, &pb.AddSecretRequest{UserId: id, Secret: NewSecret.Secret, Meta: "secret text", Comment: NewSecret.Comment})
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
	addtextCmd.Flags().StringVarP(&NewSecret.Comment, "comment", "p", "", "new comment")
	addtextCmd.MarkFlagRequired("secret")
	addtextCmd.MarkFlagRequired("comment")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addtextCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addtextCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
