package repository_interfaces

import "github.com/RandySteven/CafeConnect/be/entities/models"

type CafeProductRepository interface {
	Saver[models.CafeProduct]
	Finder[models.CafeProduct]
}
