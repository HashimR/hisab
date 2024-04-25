package dbmodels

type BankData struct {
	ID              int    `db:"id"`
	Name            string `db:"name"`
	LeanIdentifier  string `db:"lean_identifier"`
	Logo            string `db:"logo"`
	MainColor       string `db:"main_color"`
	BackgroundColor string `db:"background_color"`
}
