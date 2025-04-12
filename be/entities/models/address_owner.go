package models

import "time"

type AddressOwner struct {
	ID        uint64
	AddressID uint64
	OwnerID   uint64
	OwnerType string
	IsDefault bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
