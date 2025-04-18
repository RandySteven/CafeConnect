package models

import "time"

type User struct {
	ID             uint64
	Name           string
	Username       string
	Email          string
	Password       string
	PhoneNumber    string
	ProfilePicture string
	DoB            string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}
