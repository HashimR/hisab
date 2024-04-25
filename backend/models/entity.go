package models

import "time"

type EntityResponse struct {
	Type      string        `json:"type" validate:"required"`
	Payload   EntityPayload `json:"payload" validate:"required"`
	Timestamp time.Time     `json:"timestamp" validate:"required"`
}

type EntityPayload struct {
	Id             string       `json:"id" validate:"required"`
	LeanCustomerId string       `json:"customer_id" validate:"required"`
	Permissions    []string     `json:"permissions" validate:"required"`
	BankDetails    BankResponse `json:"bank_details" validate:"required"`
}

type BankResponse struct {
	Name               string `json:"name" validate:"required"`
	LeanBankIdentifier string `json:"identifier" validate:"required"`
	Logo               string `json:"logo" validate:"required"`
	MainColor          string `json:"main_color" validate:"required"`
	BackgroundColor    string `json:"background_color" validate:"required"`
}

type Entity struct {
	Id             int
	LeanId         string
	LeanCustomerId string
	Permissions    []string
	BankId         int
}
