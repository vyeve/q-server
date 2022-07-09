package main

import (
	"github.com/vyeve/q-server/repository"
	"github.com/vyeve/q-server/server"
	"github.com/vyeve/q-server/utils"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		utils.ModuleLogger,
		utils.ModuleValidator,
		repository.Module,
		server.Module,
		fx.Invoke(func(srv server.Server) {
			srv.Init()
		}),
	)
	app.Run()
}
