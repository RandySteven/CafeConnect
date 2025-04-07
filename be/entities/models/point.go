package models

import "time"

type Point struct {
	ID        uint64
	Point     uint64
	UserID    uint64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
