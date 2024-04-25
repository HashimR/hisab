package balances

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"main/config"
	"main/store"
	"main/store/dbmodels"
	"net/http"
	"time"
)

type Service struct {
	balanceRepository *store.BalanceRepository
}

func NewBalanceService(br *store.BalanceRepository) *Service {
	return &Service{
		balanceRepository: br,
	}
}

func (s *Service) GetLatestBalancesForUser(userId int) ([]*Balance, error) {
	dbBalances, err := s.balanceRepository.GetLatestBalancesForUser(userId)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get latest balances for user %d", userId)
		return nil, err
	}

	var balances []*Balance
	for _, dbBalance := range dbBalances {
		balances = append(balances, mapDBBalanceToModel(dbBalance))
	}

	return balances, nil
}

func (s *Service) FetchAndStoreBalanceForAccount(leanAccountId string, accountId int, leanEntityId string) error {
	// Fetch balance from Lean API
	balance, err := s.fetchBalanceFromAPI(leanAccountId, leanEntityId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch balance from Lean API")
		return err
	}

	// Convert LeanBalance to dbmodels.Balance
	dbBalance := convertToDBBalance(balance, accountId)

	// Store or update balance in the database
	err = s.balanceRepository.StoreOrUpdateBalance(dbBalance)
	if err != nil {
		log.Error().Err(err).Msg("Failed to store or update balance in the database")
		return err
	}

	log.Info().Msg("Successfully fetched and stored balance")
	return nil
}

func (s *Service) fetchBalanceFromAPI(leanAccountId string, leanEntityId string) (*LeanBalance, error) {
	url := "https://sandbox.leantech.me/data/v1/balance/"
	leanAppToken := config.GetConfig().GetString("lean-app-token")

	// Create the request body
	data := map[string]string{
		"account_id": leanAccountId,
		"entity_id":  leanEntityId,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("lean-app-token", leanAppToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request to Lean failed with status code: %d", resp.StatusCode)
	}

	// Read the response body
	var responseBody LeanBalancesResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody.Payload, nil
}

func convertToDBBalance(leanBalance *LeanBalance, accountId int) *dbmodels.Balance {
	// Convert LeanBalance to dbmodels.Balance
	dbBalance := &dbmodels.Balance{
		AccountID:   accountId,
		Amount:      decimal.NewFromFloat(leanBalance.Balance), // Assuming Balance is a float64
		Currency:    leanBalance.CurrencyCode,
		LastUpdated: time.Now(),
		Date:        time.Now(),
		Type:        leanBalance.AccountType,
	}

	return dbBalance
}

func mapDBBalanceToModel(dbBalance *dbmodels.Balance) *Balance {
	return &Balance{
		ID:          dbBalance.ID,
		AccountID:   dbBalance.AccountID,
		Amount:      dbBalance.Amount,
		Currency:    dbBalance.Currency,
		LastUpdated: dbBalance.LastUpdated,
		Date:        dbBalance.Date,
		Type:        dbBalance.Type,
	}
}
