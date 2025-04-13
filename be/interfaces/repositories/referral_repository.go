package repository_interfaces

import "github.com/RandySteven/CafeConnect/be/entities/models"

type ReferralRepository interface {
	Saver[models.Referral]
	Finder[models.Referral]
}
