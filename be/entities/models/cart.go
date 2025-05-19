package models

import "time"

type Cart struct {
	ID            uint64
	UserID        uint64
	CafeProductID uint64
	Qty           uint64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}
