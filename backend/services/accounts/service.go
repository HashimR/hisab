package accounts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"main/config"
	"main/models"
	"main/store"
	"main/store/dbmodels"
	"net/http"
	"time"
)

type Service struct {
	accountRepository *store.AccountRepository
}

func NewAccountService(ar *store.AccountRepository) *Service {
	return &Service{
		accountRepository: ar,
	}
}

func (s *Service) GetAccounts(userId int) ([]*Account, error) {
	//TODO: update response for required frontend data
	accountDataList, err := s.accountRepository.GetAccountsByUser(userId)
	if err != nil {
		return nil, err
	}

	// Initialize a slice to hold the mapped account data
	accounts := make([]*Account, len(accountDataList))

	// Map the database account data to the accounts.Account struct
	for i, accountData := range accountDataList {
		accounts[i] = &Account{
			ID:                          accountData.ID,
			LeanAccountId:               accountData.LeanAccountId,
			UserID:                      accountData.UserID,
			Name:                        accountData.Name,
			Number:                      accountData.Number,
			LatestBalance:               accountData.LatestBalance,
			AccountType:                 accountData.AccountType,
			IBAN:                        accountData.IBAN,
			CurrencyCode:                accountData.CurrencyCode,
			IsCredit:                    accountData.IsCredit,
			CreditCardNumberLastFour:    accountData.CreditCardNumberLastFour,
			CreditCardLimit:             accountData.CreditCardLimit,
			CreditCardNextPaymentDate:   accountData.CreditCardNextPaymentDate,
			CreditCardNextPaymentAmount: accountData.CreditCardNextPaymentAmount,
			LastUpdatedTransactions:     accountData.LastUpdatedTransactions,
			LastUpdatedBalance:          accountData.LastUpdatedBalance,
		}
	}

	return accounts, nil
}

//func (s *Service) FetchAndStoreAccounts(leanEntityId string, userId int) ([]string, error) {
//	accounts, err := fetchLeanAccounts(leanEntityId)
//	if err != nil {
//		log.Error().Err(err).Msgf("Failed to fetch lean accounts")
//		return nil, err
//	}
//
//	err = s.accountRepository.StoreLeanAccounts(accounts, userId)
//	if err != nil {
//		log.Error().Err(err).Msgf("Failed to store lean accounts")
//		return nil, err
//	}
//
//	return accounts
//}

func (s *Service) FetchAndStoreAccounts(leanEntityId string, userId int) ([]*dbmodels.Account, error) {
	// Fetch LeanAccounts from the external service
	leanAccounts, err := fetchLeanAccounts(leanEntityId)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to fetch lean accounts")
		return nil, err
	}

	// Iterate over the fetched LeanAccounts and map them to dbmodels.Account
	var accountsToStore []*dbmodels.Account
	for _, leanAcc := range leanAccounts {
		acc := mapLeanAccountToDBModel(leanAcc, userId)
		accountsToStore = append(accountsToStore, acc)
	}

	// Store the mapped accounts in the database
	if err := s.accountRepository.StoreAccounts(accountsToStore); err != nil {
		log.Error().Err(err).Msgf("Failed to store lean accounts for entity: %s", leanEntityId)
	}

	return accountsToStore, nil
}

func mapLeanAccountToDBModel(leanAcc *models.LeanAccount, userId int) *dbmodels.Account {
	acc := &dbmodels.Account{
		LeanAccountId:               leanAcc.AccountID,
		UserID:                      userId,
		AccountType:                 leanAcc.Type,
		IBAN:                        leanAcc.IBAN,
		Name:                        leanAcc.Name,
		Number:                      leanAcc.AccountNumber,
		CurrencyCode:                leanAcc.CurrencyCode,
		IsCredit:                    leanAcc.Credit != nil,
		CreditCardNumberLastFour:    "",
		CreditCardLimit:             decimal.NewFromFloat(0),
		CreditCardNextPaymentDate:   time.Time{},
		CreditCardNextPaymentAmount: decimal.NewFromFloat(0),
		LastUpdatedTransactions:     time.Time{},
		LatestBalance:               decimal.NewFromFloat(0),
		LastUpdatedBalance:          time.Time{},
	}

	if leanAcc.Credit != nil {
		acc.CreditCardNumberLastFour = leanAcc.Credit.CardNumberLastFour
		acc.CreditCardLimit = decimal.NewFromFloat(leanAcc.Credit.Limit)
		nextPaymentDate, _ := time.Parse(time.DateOnly, leanAcc.Credit.NextPaymentDueDate)
		acc.CreditCardNextPaymentDate = nextPaymentDate
		acc.CreditCardNextPaymentAmount = decimal.NewFromFloat(leanAcc.Credit.NextPaymentDueAmount)
	}

	return acc
}

func fetchLeanAccounts(entityID string) ([]*models.LeanAccount, error) {
	url := "https://sandbox.leantech.me/data/v1/accounts/"
	leanAppToken := config.GetConfig().GetString("lean-app-token")

	// Create the request body
	data := map[string]string{
		"entity_id": entityID,
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
		log.Error().Msgf("Request to Lean failed with status code: %d", resp.StatusCode)
		return nil, fmt.Errorf("request to Lean failed with status code: %d", resp.StatusCode)
	}

	// Read the response body
	var responseBody models.LeanAccountsResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode response body")
		return nil, err
	}

	return responseBody.Payload.Accounts, nil
}
