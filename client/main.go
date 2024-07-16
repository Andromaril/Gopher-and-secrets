package main

import (
	"github.com/Andromaril/Gopher-and-secrets/client/cmd"
	"github.com/Andromaril/Gopher-and-secrets/client/internal/config"
)

func main() {
	config.ParseFlags()
	cmd.Execute()
	//grpc.Init()
}
