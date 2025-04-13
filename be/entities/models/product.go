package models

import "time"

type Product struct {
	ID                uint64
	Name              string
	PhotoURL          string
	ProductCategoryID uint64
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
}
