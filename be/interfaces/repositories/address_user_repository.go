package repository_interfaces

import "github.com/RandySteven/CafeConnect/be/entities/models"

type AddressUserRepository interface {
	Saver[models.AddressUser]
	Finder[models.AddressUser]
}
