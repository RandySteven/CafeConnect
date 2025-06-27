package repository_interfaces

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
)

type (
	TransactionHeaderRepository interface {
		Saver[models.TransactionHeader]
		FindByTransactionCode(ctx context.Context, transactionCode string) (result *models.TransactionHeader, err error)
		FindByUserID(ctx context.Context, userId uint64) (result []*models.TransactionHeader, err error)
		Updater[models.TransactionHeader]
		Index
		FindByTransactionStatus(ctx context.Context, status string) (result []*models.TransactionHeader, err error)
	}

	TransactionDetailRepository interface {
		Saver[models.TransactionDetail]
		Finder[models.TransactionDetail]
		FindByTransactionId(ctx context.Context, transactionId uint64) (result []*models.TransactionDetail, err error)
	}
)
