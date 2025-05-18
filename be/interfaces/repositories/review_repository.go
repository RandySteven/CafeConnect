package repository_interfaces

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
)

type ReviewRepository interface {
	Saver[models.Review]
	Finder[models.Review]
	FindByCafeID(ctx context.Context, cafeID uint64) (result []*models.Review, err error)
	AvgCafeRating(ctx context.Context, cafeID uint64) (result float64, err error)
}
