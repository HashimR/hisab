package store

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"main/store/dbmodels"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db}
}

func (ur *UserRepository) CreateUser(user *dbmodels.UserData) error {
	query := `
        INSERT INTO users (
            email,
            password,
            first_name,
            last_name,
            phone_number,
            country,
            open_banking_id,
            lean_customer_id
        )
        VALUES (
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?
        )
    `

	_, err := ur.db.Exec(query, user.Email, user.Password, user.FirstName, user.LastName, user.PhoneNumber, user.Country, user.OpenBankingId, user.LeanCustomerId)
	return err
}

func (ur *UserRepository) GetUserByEmail(email string) (*dbmodels.UserData, error) {
	var user dbmodels.UserData
	query := `
		SELECT id, email, password, first_name, last_name, phone_number, country, lean_customer_id, connected_state
		FROM users
		WHERE email = ?
	`

	if err := ur.db.Get(&user, query, email); err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) GetUserByLeanCustomerId(leanCustomerId string) (*dbmodels.UserData, error) {
	var user dbmodels.UserData
	query := `
		SELECT id, email, password, first_name, last_name, phone_number, country, lean_customer_id, connected_state
		FROM users
		WHERE lean_customer_id = ?
	`

	if err := ur.db.Get(&user, query, leanCustomerId); err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) StoreConnectedState(leanCustomerId string) error {
	query := `
        UPDATE users
        SET connected_state = true
        WHERE lean_customer_id = ?
    `

	_, err := ur.db.Exec(query, leanCustomerId)
	if err != nil {
		log.Error().Err(err).Msg("Could not store connected state")
		return err
	}

	return nil
}
