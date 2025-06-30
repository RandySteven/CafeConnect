package jobs

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/enums"
	cache_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/caches"
	job_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/jobs"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
	"time"
)

type TransactionJob struct {
	transactionHeaderRepo   repository_interfaces.TransactionHeaderRepository
	transactionDetailRepo   repository_interfaces.TransactionDetailRepository
	midtransTransactionRepo repository_interfaces.MidtransTransactionRepository
	transactionCache        cache_interfaces.TransactionCache
	midtrans                midtrans_client.Midtrans
}

func (t *TransactionJob) CheckMidtransStatus(ctx context.Context) (err error) {
	pendingTransactions, err := t.transactionHeaderRepo.FindByTransactionStatus(ctx, enums.TransactionPENDING.String())
	if err != nil {
		return err
	}

	for _, transaction := range pendingTransactions {
		midtransTransactionStatusResponse, err := t.midtrans.CheckTransaction(ctx, transaction.TransactionCode)
		if err != nil {
			return err
		}

		switch midtransTransactionStatusResponse.TransactionStatus {
		case "settlement":
			transaction.Status = enums.TransactionSUCCESS.String()
			transaction.UpdatedAt = time.Now()
			_, err = t.transactionHeaderRepo.Update(ctx, transaction)
			if err != nil {
				return err
			}
		case "pending":

		case "cancel":
			transaction.Status = enums.TransactionFAILED.String()
			transaction.UpdatedAt = time.Now()
			_, err = t.transactionHeaderRepo.Update(ctx, transaction)
			if err != nil {
				return err
			}
		}

	}

	return
}

func (t *TransactionJob) HardFailedPendingTrx(ctx context.Context) (err error) {
	transactionHeaders, err := t.transactionHeaderRepo.FindByTransactionStatus(ctx, enums.TransactionPENDING.String())
	if err != nil {
		return err
	}

	for _, transactionHeader := range transactionHeaders {
		transactionHeader.Status = enums.TransactionFAILED.String()
		transactionHeader.UpdatedAt = time.Now()
		_, err = t.transactionHeaderRepo.Update(ctx, transactionHeader)
		if err != nil {
			return err
		}
	}

	return nil
}

var _ job_interfaces.TransactionJob = &TransactionJob{}

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
