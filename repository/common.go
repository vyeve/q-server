package repository

import (
	"errors"

	"github.com/vyeve/q-server/utils/logger"

	"go.uber.org/fx"
)

const (
	EnvPathToSQLITE = "APP_SQLITE_DB_PATH"
	EnvDBPoolSize   = "APP_SQLITE_DB_POOL_SIZE"

	defaultPathToSQLITE = "../assets/schema.sqlite"
	defaultPoolSize     = 10
)

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
