package transaction

import (
	"encoding/json"
	"log"
	"os"
)

func (t *Service) GetTransactions() (*Transactions, error) {
	data, err := os.ReadFile("./services/transaction/hardcoded.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
		return nil, err
	}

	var transactions *Transactions

	if err := json.Unmarshal(data, &transactions); err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
		return nil, err
	}
	return transactions, nil
}

type Service struct {
}

func NewTransactionService() *Service {
	return &Service{}
}
