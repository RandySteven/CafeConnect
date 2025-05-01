package repository_interfaces

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
)

type CafeRepository interface {
	Saver[models.Cafe]
	Finder[models.Cafe]

	FindByCafeFranchiseId(ctx context.Context, cafeFranchiseId uint64) (result []*models.Cafe, err error)
}
