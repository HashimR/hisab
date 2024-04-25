package dbmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

type Account struct {
	ID                          int             `db:"id"`
	LeanAccountId               string          `db:"lean_account_id"`
	UserID                      int             `db:"user_id"`
	AccountType                 string          `db:"account_type"`
	IBAN                        string          `db:"iban"`
	Name                        string          `db:"name"`
	Number                      string          `db:"number"`
	CurrencyCode                string          `db:"currency_code"`
	IsCredit                    bool            `db:"is_credit"`
	CreditCardNumberLastFour    string          `db:"credit_card_number_last_four"`
	CreditCardLimit             decimal.Decimal `db:"credit_card_limit"`
	CreditCardNextPaymentDate   time.Time       `db:"credit_card_next_payment_date"`
	CreditCardNextPaymentAmount decimal.Decimal `db:"credit_card_next_payment_amount"`
	LastUpdatedTransactions     time.Time       `db:"last_updated_transactions"`
	LatestBalance               decimal.Decimal `db:"latest_balance"`
	LastUpdatedBalance          time.Time       `db:"last_updated_balance"`
}
