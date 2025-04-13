package repository_interfaces

import "github.com/RandySteven/CafeConnect/be/entities/models"

type RoleRepository interface {
	Saver[models.Role]
	Finder[models.Role]
}
