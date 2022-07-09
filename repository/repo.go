package repository

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"syscall"

	"github.com/vyeve/q-server/models"
	"github.com/vyeve/q-server/utils/logger"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/fx"
)

type repoImpl struct {
	db     *sql.DB
	logger logger.Logger
}

// newConnection extract parameters through environment and returns connection to DB
func newConnection() (*sql.DB, error) {
	path, found := syscall.Getenv(EnvPathToSQLITE)
	if !found {
		return nil, fmt.Errorf("not provided path to database. use %s environment", EnvPathToSQLITE)
	}

	var poolSize int
	poolSizeEnv, found := syscall.Getenv(EnvDBPoolSize)
	if found {
		var err error
		poolSize, err = strconv.Atoi(poolSizeEnv)
		if err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	if poolSize > 0 {
		db.SetMaxIdleConns(poolSize)
		db.SetMaxOpenConns(poolSize)
	}
	return db, nil
}

// NewRepository returns Repository instance
func NewRepository(params Params) (Repository, error) {
	db, err := newConnection()
	if err != nil {
		return nil, err
	}
	r := &repoImpl{
		db:     db,
		logger: params.Logger,
	}
	params.LifeCycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})
	return r, nil
}

// UploadTransfers implements business logic of transfers.
// In case it cannot find organization by IBAN+BIC the ErrUnknownOrganization will be returned.
// In case total transfers will be more than balance, the ErrInsufficientFunds will be returned.
func (r *repoImpl) UploadTransfers(ctx context.Context, receipt *models.Receipt) (err error) {
	if receipt == nil || len(receipt.CreditTransfers) == 0 {
		return ErrNoTransfers
	}
	r.logger.Debugf("try to upload %d transactions", len(receipt.CreditTransfers))
	var (
		buf       bytes.Buffer
		numFields = 6
		data      = make([]interface{}, 1, numFields*len(receipt.CreditTransfers)+1)
	)
	// open transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		// commit transaction if all is OK
		if err == nil {
			err = tx.Commit()
			return
		}
		// rollback transaction on failed
		if err2 := tx.Rollback(); err2 != nil {
			r.logger.Debugf("failed to rollback transaction. err: %v", err)
		}
	}()

	var (
		bankID         int64
		resultBalance  int64
		totalTransfers int64
	)
	// calculate totalTransfers and prepare data to upload
	for i, tr := range receipt.CreditTransfers {
		if i != 0 {
			buf.WriteRune(',')
		}
		totalTransfers += tr.AmountCent
		buf.WriteString("($1,") // placeholder for bank_account_id
		for j := 0; j < numFields; j++ {
			if j != 0 {
				buf.WriteRune(',')
			}
			buf.WriteString("$" + strconv.Itoa(i*numFields+j+2))
		}
		buf.WriteRune(')')
		data = append(data,
			tr.CounterpartyName,
			tr.CounterpartyIBAN,
			tr.CounterpartyBIC,
			tr.AmountCent,
			tr.Currency,
			tr.Description,
		)
	}
	// extract bank_account_id and total balance after transfers
	if err = tx.QueryRowContext(ctx, updateTotalBalanceStmt, totalTransfers, receipt.OrganizationIBAN, receipt.OrganizationBIC).
		Scan(&bankID, &resultBalance); err != nil {
		if err == sql.ErrNoRows {
			return ErrUnknownOrganization // no records for IBAN-BIC in bank_account_table
		}
		return err // something happen wrong
	}
	if resultBalance < 0 {
		return ErrInsufficientFunds // not enough funds
	}
	data[0] = bankID // set the first placeholder
	_, err = tx.ExecContext(ctx, insertTransfersStmt+buf.String(), data...)
	if err != nil {
		return err
	}

	return nil
}
