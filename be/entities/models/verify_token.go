package models

import "time"

type VerifyToken struct {
	ID          uint64
	Token       string
	UserID      uint64
	IsClicked   bool
	ExpiredTime time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
