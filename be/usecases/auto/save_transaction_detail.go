package auto_transfer_usecases

import (
	"context"
	"fmt"

	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
)

func (t *autoTransferWorkflow) saveTransactionDetail(ctx context.Context, transactionID uint64, checkoutList []*requests.CheckoutList) (err error) {
	for _, checkout := range checkoutList {
		transactionDetail := &models.TransactionDetail{
			TransactionID: transactionID,
			CafeProductID: checkout.CafeProductID,
			Qty:           checkout.Qty,
		}
		_, err = t.transactionDetailRepository.Save(ctx, transactionDetail)
		if err != nil {
			return fmt.Errorf("failed to save transaction detail: %w", err)
		}
	}
	return nil
}
