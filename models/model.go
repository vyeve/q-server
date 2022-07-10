package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type (
	// Fiat is an enumeration of supported currencies.
	Fiat string
	// CreditTransfer represents struct for credit-transfers
	CreditTransfer struct {
		AmountCent       int64  `json:"amount"`
		Currency         Fiat   `json:"currency,omitempty"`
		CounterpartyName string `json:"counterparty_name,omitempty"`
		CounterpartyBIC  string `json:"counterparty_bic,omitempty"`
		CounterpartyIBAN string `json:"counterparty_iban,omitempty"`
		Description      string `json:"description"`
	}
	// Receipt is a container for credit transfers and meta information about bank
	Receipt struct {
		OrganizationName string            `json:"organization_name,omitempty"`
		OrganizationBIC  string            `json:"organization_bic,omitempty"`
		OrganizationIBAN string            `json:"organization_iban,omitempty"`
		CreditTransfers  []*CreditTransfer `json:"credit_transfers,omitempty"`
	}
	// centWrapper is a helper to unmarshal amount cents
	centWrapper string
	// transfer is a helper to prevent recursion during unmarshal
	transfer CreditTransfer
	// transactionWrapper is a helper to unmarshal credit transfer
	transactionWrapper struct {
		transfer
		AmountCent centWrapper `json:"amount"`
	}
)

// supported fiats.
// EUR commonly used.
// USD for test and for future features.
const (
	FiatEUR Fiat = "EUR"
	FiatUSD Fiat = "USD"
)

// supportedFiats is a list fiats, needs for marshal/unmarshal methods
var supportedFiats = []Fiat{
	FiatEUR,
	FiatUSD,
}

// MarshalJSON implements json.Marshaller interface
func (f Fiat) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", f)), nil
}

// UnmarshalJSON implements json.Unmarshaler interface
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

// UnmarshalJSON implements json.Unmarshaler interface
func (c *centWrapper) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	*c = centWrapper(data)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler interface
func (t *CreditTransfer) UnmarshalJSON(data []byte) error {
	wr := new(transactionWrapper)
	err := json.Unmarshal(data, wr)
	if err != nil {
		return err
	}
	*t = CreditTransfer(wr.transfer)
	t.AmountCent, err = wr.AmountCent.convert()
	return err
}

// convert converts string representation amount to int64 cents
// without conversion to float. Money doesn't like float.
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
		return -1, ErrIncorrectAmount
	}
	return int64(out), nil
}
