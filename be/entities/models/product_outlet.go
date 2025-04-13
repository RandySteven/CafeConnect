package models

import "time"

type ProductOutlet struct {
	ID        uint64
	OutletID  uint64
	ProductID uint64
	Price     uint64
	Stock     uint64
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
