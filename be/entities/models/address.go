package models

import "time"

type Address struct {
	ID        uint64
	Address   string
	Latitude  float32
	Longitude float32
	IsDefault bool
	UserID    uint64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
