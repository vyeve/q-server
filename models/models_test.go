package models

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestConvertFiatToCents(t *testing.T) {
	testCases := []struct {
		name      string
		cents     centWrapper
		expResult int64
		needErr   bool
	}{
		{
			name:      "test integer value",
			cents:     "123",
			expResult: 12300,
			needErr:   false,
		},
		{
			name:      "test with dot at the end",
			cents:     "123.",
			expResult: 12300,
			needErr:   false,
		},
		{
			name:      "test with one digit on the right",
			cents:     "123.4",
			expResult: 12340,
			needErr:   false,
		},
		{
			name:      "test with two digits on the right",
			cents:     "123.45",
			expResult: 12345,
			needErr:   false,
		},
		{
			name:      "test with three digits on the right",
			cents:     "123.456",
			expResult: -1,
			needErr:   true,
		},
		{
			name:      "test with four digits on the right",
			cents:     "123.4567",
			expResult: -1,
			needErr:   true,
		},
		{
			name:      "test with two decimals points",
			cents:     "123.45.67",
			expResult: -1,
			needErr:   true,
		},
		{
			name:      "test with string value",
			cents:     "foobar",
			expResult: -1,
			needErr:   true,
		},
		{
			name:      "test with string value and a dot",
			cents:     "foo.bar",
			expResult: -1,
			needErr:   true,
		},
		{
			name:      "test with negative value",
			cents:     "-123.45",
			expResult: -1,
			needErr:   true,
		},
		{
			name:      "test with empty string",
			cents:     "",
			expResult: -1,
			needErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := tc.cents.convert()
			if tc.needErr && err == nil {
				t.Error("expected not <nil> error")
			}
			if !tc.needErr && err != nil {
				t.Error(err)
			}
			if out != tc.expResult {
				t.Errorf("expected: %d, got: %d", tc.expResult, out)
			}
		})
	}
}

func TestReceipt_Unmarshal(t *testing.T) {
	// test with correct json
	p, err := ioutil.ReadFile("../assets/sample1.json")
	if err != nil {
		t.Fatal(err)
	}
	r := new(Receipt)
	err = json.Unmarshal(p, r)
	if err != nil {
		t.Fatal(err)
	}
	if cf := len(r.CreditTransfers); cf != 3 {
		t.Fatalf("expected 3 credit transfers, got: %d", cf)
	}
	exp := Transfer{
		AmountCent:       1450,
		Currency:         FiatEUR,
		CounterpartyName: "Bip Bip",
		CounterpartyBIC:  "CRLYFRPPTOU",
		CounterpartyIBAN: "EE383680981021245685",
		Description:      "Wonderland/4410",
	}
	if !reflect.DeepEqual(*r.CreditTransfers[0], exp) {
		t.Errorf("expected: %+v, got: %+v", exp, *r.CreditTransfers[0])
	}
	// test with incorrect json
	p, err = ioutil.ReadFile("../assets/sample2.json")
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(p, r)
	if err == nil {
		t.Error("expected not <nil> error")
	}
}

func TestTransaction_Unmarshal(t *testing.T) {
	b := `
	{
		"amount": "14.5",
		"currency": "EUR",
		"counterparty_name": "Bip Bip",
		"counterparty_bic": "CRLYFRPPTOU",
		"counterparty_iban": "EE383680981021245685",
		"description": "Wonderland/4410"
	}
	`
	tr := new(Transfer)
	err := json.Unmarshal([]byte(b), tr)
	if err != nil {
		t.Fatal(err)
	}
	exp := Transfer{
		AmountCent:       1450,
		Currency:         FiatEUR,
		CounterpartyName: "Bip Bip",
		CounterpartyBIC:  "CRLYFRPPTOU",
		CounterpartyIBAN: "EE383680981021245685",
		Description:      "Wonderland/4410",
	}
	if !reflect.DeepEqual(*tr, exp) {
		t.Errorf("expected: %+v, got: %+v", exp, *tr)
	}
}
