package dbmodels

type UserData struct {
	ID             int    `db:"id"`
	Email          string `db:"email"`
	Password       string `db:"password"`
	FirstName      string `db:"first_name"`
	LastName       string `db:"last_name"`
	PhoneNumber    string `db:"phone_number"`
	Country        string `db:"country"`
	OpenBankingId  string `db:"open_banking_id"`
	LeanCustomerId string `db:"lean_customer_id"`
	ConnectedState bool   `db:"connected_state"`
}
