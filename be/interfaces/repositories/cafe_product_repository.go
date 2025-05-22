package repository_interfaces

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
)

type CafeProductRepository interface {
	Saver[models.CafeProduct]
	Finder[models.CafeProduct]
	FindByCafeID(ctx context.Context, cafeID uint64) (result []*models.CafeProduct, err error)
	FindByCafeIDs(ctx context.Context, cafeIDs []uint64) (result []*models.CafeProduct, err error)
}
