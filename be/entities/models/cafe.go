package models

import "time"

type Cafe struct {
	ID        uint64
	Name      string
	LogoURL   string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
