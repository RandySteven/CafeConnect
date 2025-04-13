package repository_interfaces

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
)

type PointRepository interface {
	Saver[models.Point]
	Finder[models.Point]
	Updater[models.Point]
	FindByUserID(ctx context.Context, userId uint64) (result *models.Point, err error)
}
