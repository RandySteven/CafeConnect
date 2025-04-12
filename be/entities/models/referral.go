package models

import "time"

type Referral struct {
	ID          uint64
	Code        string
	UserID      uint64
	ExpiredTime time.Time
	Status      string
	NumbOfUsage uint64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
