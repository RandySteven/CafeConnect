package models

import "time"

type RoleUser struct {
	ID        uint64
	RoleID    uint64
	UserID    uint64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
