package repository_interfaces

import "github.com/RandySteven/CafeConnect/be/entities/models"

type ProductCategoryRepository interface {
	Finder[models.ProductCategory]
}
