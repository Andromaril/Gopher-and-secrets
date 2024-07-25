package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/Andromaril/Gopher-and-secrets/client/cmd"
	"github.com/Andromaril/Gopher-and-secrets/client/internal/config"
	"github.com/Andromaril/Gopher-and-secrets/client/internal/local"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	config.ParseFlags()
	err := local.InitStorage()
	if err != nil {
		log.Fatalln(err)
	}
	err = local.InitStorageTemp()
	if err != nil {
		log.Fatalln(err)
	}
	err = local.LoadUser()
	if err != nil {
		log.Fatalln(err)
	}
	log.Infof("Build Version %s, Build date %s, Build commit %s", buildVersion, buildDate, buildCommit)
	cmd.Execute()
	err = local.UpdateUser()
	if err != nil {
		log.Fatalln(err)
	}
	err = local.UpdateTemp()
	if err != nil {
		log.Fatalln(err)
	}

}
