package store

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"main/store/dbmodels"
)

type RefreshTokenRepository struct {
	db *sqlx.DB
}

func NewRefreshTokenRepository(db *sqlx.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db}
}

func (rr *RefreshTokenRepository) CreateRefreshToken(token *dbmodels.RefreshToken) error {
	query := `
        INSERT INTO refresh_tokens (user_id, token, expiration)
        VALUES (?, ?, ?)
    `

	_, err := rr.db.Exec(query, token.Username, token.Token, token.Expiration)
	return err
}

func (rr *RefreshTokenRepository) GetRefreshTokenByToken(token string) (*dbmodels.RefreshToken, error) {
	var refreshToken dbmodels.RefreshToken
	query := `
        SELECT id, user_id, token, expiration
        FROM refresh_tokens
        WHERE token = ?
    `

	if err := rr.db.Get(&refreshToken, query, token); err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (rr *RefreshTokenRepository) DeleteRefreshToken(token string) error {
	query := `
        DELETE FROM refresh_tokens
        WHERE token = ?
    `

	_, err := rr.db.Exec(query, token)
	return err
}

func (rr *RefreshTokenRepository) GetRefreshTokenByUsername(username string) (*dbmodels.RefreshToken, error) {
	var refreshToken dbmodels.RefreshToken
	query := `
        SELECT id, user_id, token, expiration
        FROM refresh_tokens
        WHERE user_id = ?
    `

	if err := rr.db.Get(&refreshToken, query, username); err != nil {
		log.Error().Err(err).Msg("Could not get refresh token")
		return nil, err
	}

	return &refreshToken, nil
}
