package models

import "time"

type Cafe struct {
	ID              uint64
	AddressID       uint64
	CafeFranchiseID uint64
	CafeType        string
	PhotoURLs       string
	OpenHour        string
	CloseHour       string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}
