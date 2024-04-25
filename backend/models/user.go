package models

type UserRequest struct {
	Username       string `json:"username" validate:"required"`
	Password       string `json:"password" validate:"required"`
	FirstName      string `json:"first_name" validate:"required"`
	LastName       string `json:"last_name" validate:"required"`
	PhoneNumber    string `json:"phone_number" validate:"required"`
	Country        string `json:"country" validate:"required"`
	LeanCustomerId string `json:"lean_customer_id"`
}

type User struct {
	ID             int
	Email          string
	FirstName      string
	LastName       string
	PhoneNumber    string
	Country        string
	OpenBankingId  string
	LeanCustomerId string
	ConnectedState bool
}
