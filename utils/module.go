package utils

import (
	"github.com/vyeve/q-server/utils/logger"
	"github.com/vyeve/q-server/utils/validator"

	"go.uber.org/fx"
)

var (
	ModuleValidator = fx.Provide(
		validator.New,
	)

	ModuleLogger = fx.Provide(
		logger.New,
	)
)
