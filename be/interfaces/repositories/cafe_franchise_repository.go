package repository_interfaces

import "github.com/RandySteven/CafeConnect/be/entities/models"

type CafeFranchiseRepository interface {
	Saver[models.CafeFranchise]
	Finder[models.CafeFranchise]
}
