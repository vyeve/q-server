package server

import (
	"github.com/vyeve/q-server/repository"
	"github.com/vyeve/q-server/utils/logger"
	"github.com/vyeve/q-server/utils/validator"

	"go.uber.org/fx"
)

const (
	EnvServerPort    = "APP_SERVER_PORT"
	defaultPort      = 8080
	transferEndpoint = "/transfer"
)

type Params struct {
	fx.In

	Logger    logger.Logger
	Repo      repository.Repository
	Validator validator.ValidatorJSON
	LifeCycle fx.Lifecycle
}
