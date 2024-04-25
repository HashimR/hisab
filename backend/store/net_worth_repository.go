package store

import (
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"main/store/dbmodels"
	"time"
)

type NetWorthRepository struct {
	db *sqlx.DB
}

func NewNetWorthRepository(db *sqlx.DB) *NetWorthRepository {
	return &NetWorthRepository{db}
}

func (nwr *NetWorthRepository) StoreNetWorthForDate(userID int, date time.Time, netWorth decimal.Decimal, lastUpdated time.Time) error {
	// Key on user_id & date
	query := `
		INSERT INTO daily_net_worth (user_id, date, net_worth, last_updated)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		net_worth = VALUES(net_worth),
		last_updated = VALUES(last_updated)
	`

	// Execute the SQL statement to insert or update the net worth record
	_, err := nwr.db.Exec(query, userID, date, netWorth, lastUpdated)
	if err != nil {
		return err
	}

	return nil
}

func (nwr *NetWorthRepository) GetNetWorthForDate(userID int, date time.Time) (*dbmodels.DailyNetWorth, error) {
	query := `
        SELECT user_id, date, net_worth, last_updated
        FROM daily_net_worth
        WHERE user_id = :user_id AND date = :date
    `

	namedParams := map[string]interface{}{
		"user_id": userID,
		"date":    date,
	}

	var netWorth dbmodels.DailyNetWorth
	err := nwr.db.Get(&netWorth, query, namedParams)
	if err != nil {
		return nil, err
	}

	return &netWorth, nil
}

func (nwr *NetWorthRepository) GetNetWorthForDateRange(userID int, start time.Time, end time.Time) ([]*dbmodels.DailyNetWorth, error) {
	query := `
        SELECT user_id, date, net_worth, last_updated
        FROM daily_net_worth
        WHERE user_id = :user_id AND date >= :start_date AND date <= :end_date
    `

	namedParams := map[string]interface{}{
		"user_id":    userID,
		"start_date": start,
		"end_date":   end,
	}

	var netWorthList []*dbmodels.DailyNetWorth
	err := nwr.db.Select(&netWorthList, query, namedParams)
	if err != nil {
		return nil, err
	}

	return netWorthList, nil
}

func (nwr *NetWorthRepository) GetLastXNetWorths(userID, x int) ([]*dbmodels.DailyNetWorth, error) {
	// Prepare the SQL query to fetch the last X net worth records for the user
	query := `
		SELECT user_id, date, net_worth, last_updated
		FROM daily_net_worth
		WHERE user_id = :user_id
		ORDER BY date DESC
		LIMIT :limit
	`

	namedParams := map[string]interface{}{
		"user_id": userID,
		"limit":   x,
	}

	var netWorthList []*dbmodels.DailyNetWorth
	err := nwr.db.Select(&netWorthList, query, namedParams)
	if err != nil {
		return nil, err
	}

	return netWorthList, nil
}
