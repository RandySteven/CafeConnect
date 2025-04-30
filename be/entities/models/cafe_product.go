package models

import "time"

type CafeProduct struct {
	ID        uint64
	CafeID    uint64
	ProductID uint64
	Price     uint64
	Stock     uint64
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Cafe    *Cafe
	Product *Product
}
