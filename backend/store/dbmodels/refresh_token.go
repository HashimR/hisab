package dbmodels

import (
	"time"
)

type RefreshToken struct {
	ID         int       `db:"id"`
	Username   string    `db:"user_id"`
	Token      string    `db:"token"`
	Expiration time.Time `db:"expiration"`
}
