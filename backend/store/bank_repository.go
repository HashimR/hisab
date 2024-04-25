package store

import (
	"github.com/jmoiron/sqlx"
	"main/store/dbmodels"
)

type BankRepository struct {
	db *sqlx.DB
}

func NewBankRepository(db *sqlx.DB) *BankRepository {
	return &BankRepository{db}
}

func (br *BankRepository) CreateBankIfDoesNotExist(bank *dbmodels.BankData) (int, error) {
	var id int
	query := `
		SELECT id
		FROM banks
		WHERE lean_identifier = ?
	`

	err := br.db.Get(&id, query, bank.LeanIdentifier)
	if err == nil && id != 0 {
		return id, nil
	}

	storeQuery := `
        INSERT INTO banks (
		   name,
		   lean_identifier,
		   logo,
		   main_color,
		   background_color
        )
        VALUES (?, ?, ?, ?, ?)
    `

	_, err = br.db.Exec(storeQuery, bank.Name, bank.LeanIdentifier, bank.Logo, bank.MainColor, bank.BackgroundColor)
	if err != nil {
		return 0, err
	}

	//insertId, err := row.LastInsertId() // check last insert id
	//if err != nil {
	//	return 0, err
	//}

	return 1, nil
}
