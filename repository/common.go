package repository

import (
	"errors"

	"github.com/vyeve/q-server/utils/logger"

	"go.uber.org/fx"
)

// environments to configure repository
const (
	EnvPathToSQLITE = "APP_SQLITE_DB_PATH"
	EnvDBPoolSize   = "APP_SQLITE_DB_POOL_SIZE"
)

// common errors
var (
	ErrInsufficientFunds   = errors.New("insufficient funds")
	ErrUnknownOrganization = errors.New("unknown organization")
	ErrNoTransfers         = errors.New("no transfers to upload")
)

type Params struct {
	fx.In

	Logger    logger.Logger
	LifeCycle fx.Lifecycle
}
