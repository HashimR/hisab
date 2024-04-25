package accounts

import (
	"github.com/shopspring/decimal"
	"time"
)

type Account struct {
	ID                          int             `json:"id"`
	LeanAccountId               string          `json:"lean_account_id"`
	UserID                      int             `json:"user_id"`
	AccountType                 string          `json:"account_type"`
	IBAN                        string          `json:"iban"`
	Name                        string          `json:"name"`
	Number                      string          `json:"number"`
	CurrencyCode                string          `json:"currency_code"`
	IsCredit                    bool            `json:"is_credit"`
	CreditCardNumberLastFour    string          `json:"credit_card_number_last_four"`
	CreditCardLimit             decimal.Decimal `json:"credit_card_limit"`
	CreditCardNextPaymentDate   time.Time       `json:"credit_card_next_payment_date"`
	CreditCardNextPaymentAmount decimal.Decimal `json:"credit_card_next_payment_amount"`
	LastUpdatedTransactions     time.Time       `json:"last_updated_transactions"`
	LatestBalance               decimal.Decimal `json:"latest_balance"`
	LastUpdatedBalance          time.Time       `json:"last_updated_balance"`
}
