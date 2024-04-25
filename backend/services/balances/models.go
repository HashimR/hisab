package balances

import (
	"github.com/shopspring/decimal"
	"time"
)

type Balance struct {
	ID          int
	AccountID   int
	Amount      decimal.Decimal
	Currency    string
	LastUpdated time.Time
	Date        time.Time
	Type        string // CURRENT, SAVINGS, CREDIT
}

type LeanBalance struct {
	Balance      float64 `json:"balance"`
	CurrencyCode string  `json:"currency_code"`
	AccountID    string  `json:"account_id"`
	AccountName  string  `json:"account_name"`
	AccountType  string  `json:"account_type"`
	Type         string  `json:"type"`
	Status       string  `json:"status"`
	ResultsID    string  `json:"results_id"`
	Message      string  `json:"message"`
	Timestamp    string  `json:"timestamp"`
	Meta         struct {
		StatusDetail struct {
			GranularStatusCode   string `json:"granular_status_code"`
			StatusAdditionalInfo string `json:"status_additional_info"`
		} `json:"status_detail"`
	} `json:"meta"`
}

type LeanBalancesResponse struct {
	Payload   LeanBalance `json:"payload"`
	Status    string      `json:"status"`
	ResultsID string      `json:"results_id"`
	Message   string      `json:"message"`
	Timestamp string      `json:"timestamp"`
	Meta      struct {
		StatusDetail struct {
			GranularStatusCode   string `json:"granular_status_code"`
			StatusAdditionalInfo string `json:"status_additional_info"`
		} `json:"status_detail"`
	} `json:"meta"`
}
