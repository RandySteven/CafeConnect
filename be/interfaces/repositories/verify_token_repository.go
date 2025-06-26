package repository_interfaces

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
)

type VerifyTokenRepository interface {
	Saver[models.VerifyToken]
	Updater[models.VerifyToken]
	FindByToken(ctx context.Context, token string) (result *models.VerifyToken, err error)
}
