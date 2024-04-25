package sync

import (
	"github.com/rs/zerolog/log"
	"main/services/accounts"
	"main/services/balances"
	"main/services/entity"
)

type Service struct {
	accountService *accounts.Service
	entityService  *entity.EntityService
	balanceService *balances.Service
}

func (s *Service) Sync(userId int) error {
	entities, err := s.entityService.GetEntitiesForUser(userId)
	if err != nil {
		log.Error().Err(err).Msgf("Could not get entities for user %d while syncing", userId)
		return err
	}

	// Update accounts
	// Fetch and update balances
	for _, e := range entities {
		// Fetch and store accounts for the entity
		accs, err := s.accountService.FetchAndStoreAccounts(e.LeanId, userId)
		if err != nil {
			log.Error().Err(err).Msgf("Could not fetch and store accounts for entityId: %s while syncing", e.LeanId)
			continue // Skip to the next entity if an error occurs
		}

		// Fetch and update balances for each account ID
		for _, acc := range accs {
			err := s.balanceService.FetchAndStoreBalanceForAccount(acc.LeanAccountId, acc.ID, e.LeanId)
			if err != nil {
				log.Error().Err(err).Msgf("Could not fetch and store balance for entityId: %s accountId: %d while syncing", e.LeanId, acc.ID)
				return err
			}
		}
	}
	return nil
}

func NewSyncService(
	as *accounts.Service,
	es *entity.EntityService,
	bs *balances.Service) *Service {
	return &Service{
		accountService: as,
		entityService:  es,
		balanceService: bs,
	}
}
