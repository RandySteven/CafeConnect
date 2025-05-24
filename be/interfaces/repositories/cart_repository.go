package repository_interfaces

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
)

type CartRepository interface {
	Saver[models.Cart]
	Updater[models.Cart]
	FindByUserID(ctx context.Context, userId uint64) (result []*models.Cart, err error)
	FindByUserIDAndCafeProductID(ctx context.Context, userId uint64, cafeProductId uint64) (result *models.Cart, err error)
	DeleteByUserID(ctx context.Context, userId uint64) (err error)
}
