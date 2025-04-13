package repository_interfaces

import "github.com/RandySteven/CafeConnect/be/entities/models"

type AddressRepository interface {
	Saver[models.Address]
	Finder[models.Address]
	Updater[models.Address]
}
