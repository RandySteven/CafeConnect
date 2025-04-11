package repository_interfaces

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
)

type UserRepository interface {
	Saver[models.User]
	Finder[models.User]
	Updater[models.User]
	Deleter[models.User]

	FindByEmail(ctx context.Context, email string) (result *models.User, err error)
	FindByUsername(ctx context.Context, username string) (result *models.User, err error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (result *models.User, err error)
}
