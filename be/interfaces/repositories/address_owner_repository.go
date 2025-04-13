package repository_interfaces

import "github.com/RandySteven/CafeConnect/be/entities/models"

type AddressOwnerRepository interface {
	Saver[models.AddressOwner]
	Finder[models.AddressOwner]
}
