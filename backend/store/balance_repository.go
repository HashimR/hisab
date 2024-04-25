package store

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"main/store/dbmodels"
)

type BalanceRepository struct {
	db *sqlx.DB
}

func NewBalanceRepository(db *sqlx.DB) *BalanceRepository {
	return &BalanceRepository{db}
}

func (br *BalanceRepository) GetLatestBalancesForUser(userID int) ([]*dbmodels.Balance, error) {
	query := `
		SELECT b.account_id, b.user_id, b.amount, b.currency, MAX(b.date) AS date, b.type
		FROM balances b
		WHERE b.user_id = ?
		GROUP BY b.account_id
	`

	var balances []*dbmodels.Balance
	if err := br.db.Select(&balances, query, userID); err != nil {
		return nil, err
	}

	return balances, nil
}

func (br *BalanceRepository) StoreOrUpdateBalance(balance *dbmodels.Balance) error {
	query := `
		INSERT INTO balances (account_id, user_id, amount, currency, last_updated, date, type)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		amount = VALUES(amount),
		last_updated = VALUES(last_updated)
	`

	_, err := br.db.Exec(query, balance.AccountID, balance.UserID, balance.Amount, balance.Currency, balance.LastUpdated, balance.Date, balance.Type)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to store or update balance for user %d", balance.UserID)
		return err
	}

	return nil
}
