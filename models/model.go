package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type (
	Fiat string

	Transaction struct {
		AmountCent       int64  `json:"amount"`
		Currency         Fiat   `json:"currency,omitempty"`
		CounterpartyName string `json:"counterparty_name,omitempty"`
		CounterpartyBIC  string `json:"counterparty_bic,omitempty"`
		CounterpartyIBAN string `json:"counterparty_iban,omitempty"`
		Description      string `json:"description"`
	}

	Receipt struct {
		OrganizationName string         `json:"organization_name,omitempty"`
		OrganizationBIC  string         `json:"organization_bic,omitempty"`
		OrganizationIBAN string         `json:"organization_iban,omitempty"`
		CreditTransfers  []*Transaction `json:"credit_transfers,omitempty"`
	}

	centWrapper string

	transaction Transaction

	transactionWrapper struct {
		transaction
		AmountCent centWrapper `json:"amount"`
	}
)

const (
	FiatEUR Fiat = "EUR"
	FiatUSD Fiat = "USD"
)

var supportedFiats = []Fiat{
	FiatEUR,
	FiatUSD,
}

func (f Fiat) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", f)), nil
}

func (f *Fiat) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	for _, v := range supportedFiats {
		if v == Fiat(data) {
			*f = v
			return nil
		}
	}
	return ErrUnsupportedFiat
}

func (c *centWrapper) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	*c = centWrapper(data)
	return nil
}

func (t *Transaction) UnmarshalJSON(data []byte) error {
	wr := new(transactionWrapper)
	err := json.Unmarshal(data, wr)
	if err != nil {
		return err
	}
	*t = Transaction(wr.transaction)
	t.AmountCent, err = wr.AmountCent.convert()
	return err
}

func (c centWrapper) convert() (int64, error) {
	if len(c) == 0 {
		return -1, ErrIncorrectAmount
	}
	var cents string
	i := strings.Index(string(c), ".")
	switch i {
	case -1: // 123
		cents = string(c) + "00"
	case len(c) - 3: // 123.45
		cents = string(c)[:i] + string(c)[i+1:]
	case len(c) - 2: // 123.4
		cents = string(c)[:i] + string(c)[i+1:] + "0"
	case len(c) - 1: // 123.
		cents = string(c)[:i] + "00"
	default:
		return -1, ErrIncorrectAmount
	}
	out, err := strconv.ParseUint(cents, 10, 64)
	if err != nil {
		return -1, err
	}
	return int64(out), nil
}
