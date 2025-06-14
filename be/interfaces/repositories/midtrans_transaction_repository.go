package repository_interfaces

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
)

type MidtransTransactionRepository interface {
	Saver[models.MidtransTransaction]
	Finder[models.MidtransTransaction]
	FindByTransactionCode(ctx context.Context, transactionCode string) (result *models.MidtransTransaction, err error)
}
