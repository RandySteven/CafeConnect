package repository_interfaces

import "github.com/RandySteven/CafeConnect/be/entities/models"

type PointRepository interface {
	Saver[models.Point]
	Finder[models.Point]
	Updater[models.Point]
}
