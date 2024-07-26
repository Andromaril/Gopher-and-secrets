// Package cmd cli логин пользователя
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
)

// UserLogin запрос сервису
var (
	UserLogin pb.LoginRequest
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login new user",
	Long:  `login new user, use: client login and flag -l for login, -p for password`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Начат процесс логина")
		u, err := user.Current()
		if err != nil {
			log.Fatalln(err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		c, err := grpc.Init()
		if err != nil {
			log.Fatalln(err)
			return
		}
		res, err := c.Login(ctx, &UserLogin)
		if err != nil {
			fmt.Println("Не удалось залогинится, пожалуйста, попробуйте еще раз")
			return
		}
		local.User[u.Username] = res.GetToken()
		//metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+res.GetToken())
		//log.Info(local.User)
		fmt.Printf("Логин прошел успешно")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&UserLogin.Login, "login", "l", "", "login")
	loginCmd.Flags().StringVarP(&UserLogin.Password, "password", "p", "", "password")
	loginCmd.MarkFlagRequired("login")
	loginCmd.MarkFlagRequired("password")
}
