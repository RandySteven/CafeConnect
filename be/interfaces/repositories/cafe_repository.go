package repository_interfaces

import "github.com/RandySteven/CafeConnect/be/entities/models"

type CafeRepository interface {
	Saver[models.Cafe]
	Finder[models.Cafe]
}
