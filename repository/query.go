package repository

const (
	insertTransfersStmt = `
INSERT INTO transactions (
	bank_account_id,
    counterparty_name,
    counterparty_iban,
    counterparty_bic,
    amount_cents,
    amount_currency,
    description
    )
VALUES
`

	updateTotalBalanceStmt = `
UPDATE bank_accounts
SET (balance_cents) = (balance_cents - $1)
WHERE (iban, bic) = ($2, $3)
RETURNING 
    id,
    balance_cents;	
`
)
