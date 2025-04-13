package repository_interfaces

import "github.com/RandySteven/CafeConnect/be/entities/models"

type RoleUserRepository interface {
	Saver[models.RoleUser]
	Finder[models.RoleUser]
}
