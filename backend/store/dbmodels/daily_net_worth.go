package dbmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

type DailyNetWorth struct {
	UserID      int             `db:"user_id"`
	Date        time.Time       `db:"date"`
	NetWorth    decimal.Decimal `db:"net_worth"`
	LastUpdated time.Time       `db:"last_updated"`
}
