package models

import "time"

type Cafe struct {
	ID              uint64
	AddressID       uint64
	CafeFranchiseID uint64
	CafeType        string
	PhotoURL        string
	OpenHour        time.Time
	CloseHour       time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}
