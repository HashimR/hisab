package dbmodels

type EntityData struct {
	ID          int    `db:"id"`
	LeanId      string `db:"lean_id"`
	LeanUserId  string `db:"lean_user_id"`
	BankId      int    `db:"bank_id"`
	Permissions string `db:"permissions"`
	UserId      int    `db:"user_id"`
}
