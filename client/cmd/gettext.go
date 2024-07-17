/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os/user"

	"github.com/Andromaril/Gopher-and-secrets/client/internal/local"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// gettextCmd represents the gettext command
var gettextCmd = &cobra.Command{
	Use:   "gettext",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gettext called")
		user, err := user.Current()
		if err != nil {
			log.Fatalln(err)
		}
		jwt, ok := local.User[user.Username]
		if !ok {
			fmt.Println("User not authenticated.")
			return
		}
		log.Info(jwt)
	},
}

func init() {
	rootCmd.AddCommand(gettextCmd)
	gettextCmd.Flags().String("title", "t", "Text title to search for.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gettextCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gettextCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
