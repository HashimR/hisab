package transaction

import "github.com/shopspring/decimal"

type Transactions struct {
	Transactions *[]Transaction `json:"transactions"`
}

type Transaction struct {
	ID       int             `json:"id"`
	Name     string          `json:"name"`
	Amount   decimal.Decimal `json:"amount"`
	DateTime string          `json:"date_time"`
	Category string          `json:"category"`
	ImageURL string          `json:"image_url"`
}
