package net_worth

import "time"

type NetWorth struct {
	Amount string    `json:"amount"`
	Time   time.Time `json:"time"`
}

type NetWorthGraph struct {
	Points []*NetWorth `json:"points"`
}
