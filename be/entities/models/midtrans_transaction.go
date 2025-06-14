package models

import "time"

type MidtransTransaction struct {
	ID              uint64
	TransactionCode string
	TotalAmt        int64
	Token           string
	RedirectURL     string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}
