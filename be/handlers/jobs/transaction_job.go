package jobs

import (
	cache_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/caches"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
)

type TransactionJob struct {
	transactionHeaderRepo   repository_interfaces.TransactionHeaderRepository
	transactionDetailRepo   repository_interfaces.TransactionDetailRepository
	midtransTransactionRepo repository_interfaces.MidtransTransactionRepository
	transactionCache        cache_interfaces.TransactionCache
	midtrans                midtrans_client.Midtrans
}

func NewTransactionJob(
	transactionHeaderRepo repository_interfaces.TransactionHeaderRepository,
	transactionDetailRepo repository_interfaces.TransactionDetailRepository,
	midtransTransactionRepo repository_interfaces.MidtransTransactionRepository,
	transactionCache cache_interfaces.TransactionCache,
	midtrans midtrans_client.Midtrans,
) *TransactionJob {
	return &TransactionJob{
		transactionCache:        transactionCache,
		transactionDetailRepo:   transactionDetailRepo,
		transactionHeaderRepo:   transactionHeaderRepo,
		midtransTransactionRepo: midtransTransactionRepo,
		midtrans:                midtrans,
	}
}
