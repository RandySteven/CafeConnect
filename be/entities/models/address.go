package models

import "time"

type Address struct {
	ID        uint64
	Address   string
	Latitude  float64
	Longitude float64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
