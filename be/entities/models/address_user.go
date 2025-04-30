package models

import "time"

type AddressUser struct {
	ID        uint64
	AddressID uint64
	UserID    uint64
	IsDefault bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Address *Address
	User    *User
}
