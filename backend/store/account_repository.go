package store

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"main/store/dbmodels"
)

type AccountRepository struct {
	db *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) *AccountRepository {
	return &AccountRepository{db}
}

func (ar *AccountRepository) GetAccountsByUser(userId int) ([]*dbmodels.Account, error) {
	var accounts []*dbmodels.Account
	query := `
        SELECT id, lean_account_id, user_id, account_type, iban, name, number, currency_code,
               is_credit, credit_card_number_last_four, credit_card_limit,
               credit_card_next_payment_date, credit_card_next_payment_amount,
               last_updated_transactions, latest_balance, last_updated_balance
        FROM accounts
        WHERE user_id = ?
        AND is_active = true
    `

	if err := ar.db.Select(&accounts, query, userId); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (ar *AccountRepository) StoreAccounts(accounts []*dbmodels.Account) error {
	// Prepare the SQL statement for inserting or updating the Account table
	query := `
        INSERT INTO accounts (
            lean_account_id, user_id, account_type, iban, name, 
            number, currency_code, is_credit, credit_card_number_last_four, 
            credit_card_limit, credit_card_next_payment_date, credit_card_next_payment_amount,
            last_updated_transactions, latest_balance, last_updated_balance
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
            user_id = VALUES(user_id),
            account_type = VALUES(account_type),
            iban = VALUES(iban),
            name = VALUES(name),
            number = VALUES(number),
            currency_code = VALUES(currency_code),
            is_credit = VALUES(is_credit),
            credit_card_number_last_four = VALUES(credit_card_number_last_four),
            credit_card_limit = VALUES(credit_card_limit),
            credit_card_next_payment_date = VALUES(credit_card_next_payment_date),
            credit_card_next_payment_amount = VALUES(credit_card_next_payment_amount),
            last_updated_transactions = VALUES(last_updated_transactions),
            latest_balance = VALUES(latest_balance),
            last_updated_balance = VALUES(last_updated_balance)
    `

	// Prepare the statement once to be reused for each account
	stmt, err := ar.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Iterate over the list of accounts
	for _, acc := range accounts {
		// Execute the SQL statement to insert or update the account in the database
		_, err := stmt.Exec(
			acc.LeanAccountId, acc.UserID, acc.AccountType, acc.IBAN, acc.Name,
			acc.Number, acc.CurrencyCode, acc.IsCredit, acc.CreditCardNumberLastFour,
			acc.CreditCardLimit, acc.CreditCardNextPaymentDate, acc.CreditCardNextPaymentAmount,
			acc.LastUpdatedTransactions, acc.LatestBalance, acc.LastUpdatedBalance,
		)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to insert/update account into database, id: %s", acc.LeanAccountId)
			continue
		}
	}

	return nil
}
