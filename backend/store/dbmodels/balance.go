package dbmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

type Balance struct {
	ID          int             `db:"id"`
	AccountID   int             `db:"account_id"`
	UserID      int             `db:"user_id"`
	Amount      decimal.Decimal `db:"amount"`
	Currency    string          `db:"currency"`
	LastUpdated time.Time       `db:"last_updated"`
	Date        time.Time       `db:"date"` // Max one balance per account & date
	Type        string          `db:"type"` // CURRENT, SAVINGS, CREDIT
}
