package models

import "time"

type (
	TransactionHeader struct {
		ID              uint64
		UserID          uint64
		CafeID          uint64
		TransactionCode string
		Status          string
		TransactionAt   time.Time
		CreatedAt       time.Time
		UpdatedAt       time.Time
		DeletedAt       *time.Time
	}

	TransactionDetail struct {
		ID            uint64
		TransactionID uint64
		CafeProductID uint64
		Qty           uint64
		CreatedAt     time.Time
		UpdatedAt     time.Time
		DeletedAt     *time.Time
	}
)
