package repository_interfaces

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
)

type AddressUserRepository interface {
	Saver[models.AddressUser]
	Finder[models.AddressUser]
	FindByAddressAndUserID(ctx context.Context, addressID uint64, userID uint64) (result *models.AddressUser, err error)
	FindByUserID(ctx context.Context, userID uint64) (results []*models.AddressUser, err error)
}
