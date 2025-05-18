package models

import "time"

type Review struct {
	ID        uint64
	UserID    uint64
	CafeID    uint64
	Score     float64
	Comment   string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
