package repository_interfaces

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
)

type ReferralRepository interface {
	Saver[models.Referral]
	Finder[models.Referral]

	FindByUserID(ctx context.Context, userId uint64) (result *models.Referral, err error)
	FindByCode(ctx context.Context, code string) (result *models.Referral, err error)
}
