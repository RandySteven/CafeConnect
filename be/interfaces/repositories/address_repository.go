package repository_interfaces

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
)

type AddressRepository interface {
	Saver[models.Address]
	Finder[models.Address]
	Updater[models.Address]
	FindAddressBasedOnRadius(ctx context.Context, longitude, latitude float64, rangeKm uint64) (result []*models.Address, err error)
}
