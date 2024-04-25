package auth

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/rs/zerolog/log"
	"main/store"
	"main/store/dbmodels"
	"time"
)

type RefreshTokenService struct {
	refreshTokenRepository *store.RefreshTokenRepository
}

func NewRefreshTokenService(rtr *store.RefreshTokenRepository) *RefreshTokenService {
	return &RefreshTokenService{
		refreshTokenRepository: rtr,
	}
}

func (rts *RefreshTokenService) ValidateRefreshToken(refreshTokenString string) (*dbmodels.RefreshToken, error) {
	// Get the refresh token from the repository - return err if not found
	token, err := rts.refreshTokenRepository.GetRefreshTokenByToken(refreshTokenString)
	if err != nil {
		return nil, err
	}

	token, err = rts.generateNewIfExpired(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (rts *RefreshTokenService) generateNewIfExpired(token *dbmodels.RefreshToken) (*dbmodels.RefreshToken, error) {
	if token.Expiration.After(time.Now()) {
		return token, nil
	}

	// Token is expired, so delete old and generate a new one
	err := rts.refreshTokenRepository.DeleteRefreshToken(token.Token)
	if err != nil {
		log.Error().Err(err).Msg("Could not delete refresh token")
		return nil, err
	}

	newToken, err := rts.GenerateNewRefreshToken(token.Username)
	if err != nil {
		log.Error().Err(err).Msg("Could not generate new refresh token")
		return nil, err
	}

	return newToken, nil
}

func (rts *RefreshTokenService) GenerateNewRefreshToken(username string) (*dbmodels.RefreshToken, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, err
	}
	token := hex.EncodeToString(tokenBytes)

	expiration := time.Now().Add(time.Hour * 24 * 30) // 30 days

	// Create a new RefreshToken struct.
	refreshToken := &dbmodels.RefreshToken{
		Username:   username,
		Token:      token,
		Expiration: expiration,
	}

	// Save the new refresh token in the repository.
	if err := rts.refreshTokenRepository.CreateRefreshToken(refreshToken); err != nil {
		return nil, err
	}

	return refreshToken, nil
}

func (rts *RefreshTokenService) GetRefreshToken(username string) (*dbmodels.RefreshToken, error) {
	token, err := rts.refreshTokenRepository.GetRefreshTokenByUsername(username)
	if err != nil {
		return nil, err
	}

	token, err = rts.generateNewIfExpired(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}
