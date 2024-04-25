package store

import (
	"github.com/jmoiron/sqlx"
	"main/store/dbmodels"
)

type EntityRepository struct {
	db *sqlx.DB
}

func NewEntityRepository(db *sqlx.DB) *EntityRepository {
	return &EntityRepository{db}
}

func (er *EntityRepository) CreateEntity(entity *dbmodels.EntityData) error {
	storeQuery := `
        INSERT INTO entity (
            lean_id,
            lean_user_id,
            bank_id,
            permissions,
            user_id
        )
        VALUES (
            ?,
            ?,
            ?,
            ?,
            ?
        )
        ON DUPLICATE KEY UPDATE
            lean_user_id = VALUES(lean_user_id),
            bank_id = VALUES(bank_id),
            permissions = VALUES(permissions),
            user_id = VALUES(user_id)
    `
	_, err := er.db.Exec(storeQuery, entity.LeanId, entity.LeanUserId, entity.BankId, entity.Permissions, entity.UserId)
	if err != nil {
		return err
	}
	return nil
}

func (er *EntityRepository) GetEntitiesForUser(userID int) ([]*dbmodels.EntityData, error) {
	query := `
        SELECT id, lean_id, lean_user_id, bank_id, permissions, user_id
        FROM entity
        WHERE user_id = ?
    `

	var entities []*dbmodels.EntityData
	if err := er.db.Select(&entities, query, userID); err != nil {
		return nil, err
	}

	return entities, nil
}
