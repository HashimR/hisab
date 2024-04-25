package models

import "time"

type LeanAccountsResponse struct {
	Status       string            `json:"status" validate:"required"`
	Payload      LeanPayload       `json:"payload" validate:"required"`
	ResultsID    string            `json:"results_id" validate:"required"`
	Message      string            `json:"message" validate:"required"`
	Timestamp    time.Time         `json:"timestamp" validate:"required"`
	StatusDetail *LeanStatusDetail `json:"status_detail"`
}

type LeanPayload struct {
	Accounts []*LeanAccount `json:"accounts" validate:"required"`
	Type     string         `json:"type" validate:"required"`
}

type LeanAccount struct {
	AccountID     string          `json:"account_id" validate:"required"`
	Name          string          `json:"name" validate:"required"`
	CurrencyCode  string          `json:"currency_code" validate:"required"`
	Type          string          `json:"type" validate:"required"`
	IBAN          string          `json:"iban" validate:"required"`
	AccountNumber string          `json:"account_number" validate:"required"`
	Credit        *LeanCreditInfo `json:"credit"`
}

type LeanCreditInfo struct {
	CardNumberLastFour   string  `json:"card_number_last_four,omitempty"` // Use omitempty to handle null values
	Limit                float64 `json:"limit,omitempty"`
	NextPaymentDueDate   string  `json:"next_payment_due_date,omitempty"`
	NextPaymentDueAmount float64 `json:"next_payment_due_amount,omitempty"`
}

type LeanStatusDetail struct {
	GranularStatusCode   string `json:"granular_status_code,omitempty"`
	StatusAdditionalInfo string `json:"status_additional_info,omitempty"`
}
