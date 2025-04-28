package repository_interfaces

import "github.com/RandySteven/CafeConnect/be/entities/models"

type ProductRepository interface {
	Saver[models.Product]
	Finder[models.Product]
}
