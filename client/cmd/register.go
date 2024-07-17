// Package cmd для cli регистрации пользователя
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/Andromaril/Gopher-and-secrets/client/internal/grpc"
	pb "github.com/Andromaril/Gopher-and-secrets/server/proto"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// User запрос сервису
var (
	User pb.RegisterRequest
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "register new user",
	Long:  `register new user, use: client register and flag -l for login, -p for password`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Начата регистрация")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		c, err := grpc.Init()
		if err != nil {
			log.Fatalln(err)
			return
		}
		_, err = c.Register(ctx, &User)
		if err != nil {
			fmt.Println("Не удалось зарегистрировать, пожалуйста, попробуйте еще раз")
			return
		}
		fmt.Printf("Вы успешно зарегистрировались, залогиньтесь для продолжения работы")
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
	registerCmd.Flags().StringVarP(&User.Login, "login", "l", "", "new login")
	registerCmd.Flags().StringVarP(&User.Password, "password", "p", "", "new password")
	registerCmd.MarkFlagRequired("login")
	registerCmd.MarkFlagRequired("password")
}
