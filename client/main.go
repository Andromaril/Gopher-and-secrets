package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/Andromaril/Gopher-and-secrets/client/cmd"
	"github.com/Andromaril/Gopher-and-secrets/client/internal/config"
	"github.com/Andromaril/Gopher-and-secrets/client/internal/local"
)

func main() {
	config.ParseFlags()
	err := local.InitStorage()
	if err != nil {
		log.Fatalln(err)
	}
	err = local.LoadUser()
	if err != nil {
		log.Fatalln(err)
	}
	cmd.Execute()
	err = local.UpdateUser()
	if err != nil {
		log.Fatalln(err)
	}

}
